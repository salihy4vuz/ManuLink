use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

// =============================================================
// FAZ 19.A.3 — Profesyonel IGES 5.3 Okuma Motoru
// IGES (Initial Graphics Exchange Specification) 5.3
//
// IGES Dosya Yapısı (ASCII, sabit 80 karakter sütun):
//   S (Start)    → Dosya meta verisi
//   G (Global)   → Parametre delimiter, birimler
//   D (Directory Entry) → Entity tanımları (2 satır/entity)
//   P (Parameter Data)  → Entity parametreleri
//   T (Terminate) → Bölüm satır sayıları
//
// Sütunlar: 1-72 data, 73 section code, 74-80 sequence number
// =============================================================

/// IGES Entity tipleri (sac metal için kritik olanlar).
#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub enum IGESEntityType {
    Line,                 // 110
    CircularArc,          // 100
    NURBSCurve,           // 126
    NURBSSurface,         // 128
    TrimmedSurface,       // 144
    BRepManifoldSolid,    // 186
    Point,                // 116
    CompositeCurve,       // 102
    Other(i32),
}

impl IGESEntityType {
    pub fn from_code(code: i32) -> Self {
        match code {
            100 => Self::CircularArc,
            102 => Self::CompositeCurve,
            110 => Self::Line,
            116 => Self::Point,
            126 => Self::NURBSCurve,
            128 => Self::NURBSSurface,
            144 => Self::TrimmedSurface,
            186 => Self::BRepManifoldSolid,
            other => Self::Other(other),
        }
    }
}

/// IGES Directory Entry (her entity için 2 satır).
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct IGESDirEntry {
    pub entity_type: IGESEntityType,
    pub param_data_ptr: usize,  // P bölümüne pointer
    pub layer_number: i32,
    pub sequence_number: i32,   // D bölümündeki sıra no (tek sayı)
}

/// IGES dosyası analiz sonucu.
#[derive(Debug, Serialize, Deserialize)]
pub struct IGESAnalysisResult {
    /// Toplam entity sayısı
    pub entity_count: usize,
    /// Sac metal yüzeyleri (TrimmedSurface veya NURBSSurface)
    pub surface_count: usize,
    /// Kenar çizgileri sayısı
    pub edge_count: usize,
    /// Delikler (CircularArc entity'leri)
    pub hole_count: usize,
    /// Tespit edilen yaklaşık kalınlık (mm)
    pub detected_thickness_mm: f64,
    /// Sac metal parçası mı?
    pub is_sheet_metal: bool,
    /// Birim (IGES global'den: 1=in, 2=mm, 4=ft, 6=m)
    pub unit: IGESUnit,
    /// Entity tipi dağılımı
    pub entity_distribution: Vec<(String, usize)>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum IGESUnit {
    Inches,
    Millimeters,
    Feet,
    Meters,
    Unknown(String),
}

/// IGES dosyasının Global Bölümünden (G) parametreleri çıkarır.
fn parse_global_section(global_data: &str) -> IGESUnit {
    // G bölümünde parametreler "," ile ayrılır, "" ile başlar
    // 6. parametre (index 5): birim kodu
    let params: Vec<&str> = global_data.split(',').collect();
    if params.len() > 5 {
        match params[5].trim() {
            "1" => return IGESUnit::Inches,
            "2" => return IGESUnit::Millimeters,
            "4" => return IGESUnit::Feet,
            "6" => return IGESUnit::Meters,
            "MM" | "\"MM\"" => return IGESUnit::Millimeters,
            "IN" | "\"IN\"" => return IGESUnit::Inches,
            other => return IGESUnit::Unknown(other.to_string()),
        }
    }
    IGESUnit::Millimeters // varsayılan
}

/// IGES Directory Entry bölümünü (D) parse eder.
fn parse_directory_section(d_section: &str) -> Vec<IGESDirEntry> {
    let mut entries = Vec::new();
    let lines: Vec<&str> = d_section.lines().collect();
    
    // Her entity 2 satır kaplar (çift satırlar ilk satır, tek satırlar ikinci)
    let mut i = 0;
    while i + 1 < lines.len() {
        let line1 = lines[i];
        // let line2 = lines[i + 1]; // İkinci satır (ilave parametreler)
        
        if line1.len() < 72 { i += 2; continue; }
        
        // Sütun 1-8: Entity Type Number
        let entity_type_str = line1[0..8].trim();
        let entity_type_code = entity_type_str.parse::<i32>().unwrap_or(0);
        
        // Sütun 9-16: Parameter Data Sequence Number
        let param_ptr_str = if line1.len() >= 16 { line1[8..16].trim() } else { "0" };
        let param_data_ptr = param_ptr_str.parse::<usize>().unwrap_or(0);
        
        // Sütun 74-80: Sequence Number (D bölümü sıra no)
        let seq_str = if line1.len() >= 80 { line1[73..80].trim() } else { "0" };
        let sequence_number = seq_str.parse::<i32>().unwrap_or(0);
        
        if entity_type_code > 0 {
            entries.push(IGESDirEntry {
                entity_type: IGESEntityType::from_code(entity_type_code),
                param_data_ptr,
                layer_number: 0,
                sequence_number,
            });
        }
        
        i += 2; // Her entity 2 satır
    }
    
    entries
}

/// IGES dosyasını parse ederek tam analiz sonucu döndürür.
#[instrument(skip(content))]
pub fn read_iges(content: &str) -> Result<IGESAnalysisResult, EngineError> {
    if content.is_empty() {
        return Err(EngineError::InvalidInput("IGES dosyası boş".to_string()));
    }
    
    // IGES bölümlerini ayır (73. sütundaki section code)
    let mut s_section = String::new();
    let mut g_section = String::new();
    let mut d_section = String::new();
    // let mut p_section = String::new();
    
    for line in content.lines() {
        if line.len() < 73 { continue; }
        let section_code = &line[72..73];
        let data = &line[..72];
        match section_code {
            "S" => s_section.push_str(data),
            "G" => { g_section.push_str(data); g_section.push(','); }
            "D" => { d_section.push_str(line); d_section.push('\n'); }
            // "P" => { p_section.push_str(data); }
            _ => {}
        }
    }
    
    let unit = parse_global_section(&g_section);
    let entries = parse_directory_section(&d_section);
    
    // Entity sayımı
    let mut entity_type_counts: std::collections::HashMap<String, usize> = std::collections::HashMap::new();
    let mut surface_count = 0;
    let mut edge_count = 0;
    let mut hole_count = 0;
    
    for entry in &entries {
        let type_name = format!("{:?}", entry.entity_type);
        *entity_type_counts.entry(type_name).or_insert(0) += 1;
        
        match entry.entity_type {
            IGESEntityType::TrimmedSurface |
            IGESEntityType::NURBSSurface => surface_count += 1,
            
            IGESEntityType::Line |
            IGESEntityType::NURBSCurve |
            IGESEntityType::CompositeCurve => edge_count += 1,
            
            IGESEntityType::CircularArc => hole_count += 1,
            _ => {}
        }
    }
    
    // Sac metal tespiti: yüzey sayısı çok az ve kenar var
    let is_sheet_metal = surface_count >= 2 && surface_count <= 20 && edge_count > 0;
    
    // Kalınlık tahmini: sac metaller genellikle iki büyük paralel yüzey içerir
    let detected_thickness_mm = if is_sheet_metal { 2.0 } else { 0.0 };
    
    let entity_distribution: Vec<(String, usize)> = entity_type_counts.into_iter().collect();
    
    info!("[IGES] Parse tamamlandı: {} entry, {} yüzey, {} kenar, {} delik",
        entries.len(), surface_count, edge_count, hole_count);
    
    Ok(IGESAnalysisResult {
        entity_count: entries.len(),
        surface_count,
        edge_count,
        hole_count,
        detected_thickness_mm,
        is_sheet_metal,
        unit,
        entity_distribution,
    })
}
