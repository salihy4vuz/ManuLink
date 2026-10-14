use serde::{Deserialize, Serialize};
use crate::cam::dxf_reader::BoundingBox2D;
use crate::error::EngineError;
use tracing::{info, instrument};

// =============================================================
// FAZ 19.B.1-19.B.3 — Lantek Expert CAM Entegrasyonu
//
// Lantek Expert, Avrupa ve Türkiye'de Trumpf, Bystronic, Prima
// Power lazer makinelerinin standart CAM yazılımıdır.
//
// Bu modül:
//   - LTK (Lantek Part Description) dosyası yazar
//   - Lantek Teknoloji Tablosu parametrelerini yönetir
//   - Lantek Nesting Report (LNR) dosyasını okur
// =============================================================

/// Desteklenen Lantek makine profilleri.
#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum LantekMachineProfile {
    TrumpfTruLaser3030_6kW,
    TrumpfTruLaser5030_10kW,
    BystronicByStar4020_4kW,
    BystronicByStar6225_8kW,
    PrimaPowerPlatino_4kW,
    MazakOptonixFabri_6kW,
    Custom(String),
}

impl LantekMachineProfile {
    pub fn machine_code(&self) -> &str {
        match self {
            Self::TrumpfTruLaser3030_6kW => "TRUMPF_TRULASER_3030_6KW",
            Self::TrumpfTruLaser5030_10kW => "TRUMPF_TRULASER_5030_10KW",
            Self::BystronicByStar4020_4kW => "BYSTRONIC_BYSTAR_4020_4KW",
            Self::BystronicByStar6225_8kW => "BYSTRONIC_BYSTAR_6225_8KW",
            Self::PrimaPowerPlatino_4kW => "PRIMA_POWER_PLATINO_4KW",
            Self::MazakOptonixFabri_6kW => "MAZAK_OPTONIX_FABRI_6KW",
            Self::Custom(s) => s.as_str(),
        }
    }
}

/// Lantek Teknoloji Tablosu parametreleri.
/// Lazer kesimde malzeme+kalınlık+makine kombinasyonuna göre değişir.
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LantekTechParams {
    pub material_code: String,       // "MILD_STEEL", "STAINLESS_304", "AL_5052" vb.
    pub thickness_mm: f64,
    pub machine_code: String,
    
    // Kesim Parametreleri
    pub cutting_speed_mm_min: f64,   // Örn: 8000 mm/dk
    pub piercing_time_sec: f64,      // Örn: 1.2 sn
    pub pierce_pressure_bar: f64,    // Örn: 16 bar
    pub cutting_pressure_bar: f64,   // Örn: 6 bar
    pub focus_position_mm: f64,      // Örn: -1.5 mm (odak noktası)
    pub nozzle_diameter_mm: f64,     // Örn: 1.5 mm
    pub gas_type: String,            // "N2" (azot), "O2" (oksijen), "AIR"
    pub laser_power_w: f64,          // Örn: 4000 W
    
    // Kalite Parametreleri
    pub roughness_um: f64,           // Yüzey pürüzlülüğü (Ra, mikron)
    pub taper_angle_deg: f64,        // Koniklik açısı (°)
    pub kerf_width_mm: f64,          // Kesim izi genişliği (mm)
    
    // Maliyet Faktörleri
    pub gas_consumption_l_min: f64,  // Gaz tüketimi (litre/dk)
}

impl LantekTechParams {
    /// Malzeme ve kalınlık kombinasyonu için endüstri standardı parametreleri döndürür.
    pub fn lookup(material: &str, thickness_mm: f64, machine: &LantekMachineProfile) -> Self {
        let material_upper = material.to_uppercase();
        
        // Kesim hızı (mm/dk): Fiber lazer, N2 gazı, referans değerler
        let cutting_speed = match (material_upper.as_str(), thickness_mm as u32) {
            ("DKP" | "MILD_STEEL" | "ST37", 1) => 14000.0,
            ("DKP" | "MILD_STEEL" | "ST37", 2) => 8000.0,
            ("DKP" | "MILD_STEEL" | "ST37", 3) => 5500.0,
            ("DKP" | "MILD_STEEL" | "ST37", 4) => 4000.0,
            ("DKP" | "MILD_STEEL" | "ST37", 5) => 2800.0,
            ("DKP" | "MILD_STEEL" | "ST37", 6) => 2000.0,
            ("SS304" | "STAINLESS_304" | "PASLANMAZ", 1) => 12000.0,
            ("SS304" | "STAINLESS_304" | "PASLANMAZ", 2) => 6500.0,
            ("SS304" | "STAINLESS_304" | "PASLANMAZ", 3) => 4000.0,
            ("AL5052" | "ALUMINYUM" | "ALUMINUM", 1) => 18000.0,
            ("AL5052" | "ALUMINYUM" | "ALUMINUM", 2) => 10000.0,
            ("AL5052" | "ALUMINYUM" | "ALUMINUM", 3) => 7000.0,
            _ => 5000.0, // Güvenli varsayılan
        };
        
        // Pierce süresi kalınlıkla artar
        let piercing_time = 0.3 + (thickness_mm * 0.25);
        
        // Odak pozisyonu (negatif = lens altında)
        let focus = -(thickness_mm * 0.5).min(3.0);
        
        // Gaz seçimi: St37 ince→N2, kalın→O2; Paslanmaz her zaman N2; Alüminyum N2
        let gas = if material_upper.contains("ST37") || material_upper.contains("DKP") {
            if thickness_mm > 4.0 { "O2".to_string() } else { "N2".to_string() }
        } else {
            "N2".to_string()
        };
        
        let kerf = 0.1 + (thickness_mm * 0.02); // ~0.1-0.25mm
        
        LantekTechParams {
            material_code: material_upper,
            thickness_mm,
            machine_code: machine.machine_code().to_string(),
            cutting_speed_mm_min: cutting_speed,
            piercing_time_sec: piercing_time,
            pierce_pressure_bar: 12.0 + thickness_mm,
            cutting_pressure_bar: 6.0,
            focus_position_mm: focus,
            nozzle_diameter_mm: if thickness_mm <= 3.0 { 1.5 } else { 2.0 },
            gas_type: gas,
            laser_power_w: 4000.0,
            roughness_um: 3.2,
            taper_angle_deg: 0.5,
            kerf_width_mm: kerf,
            gas_consumption_l_min: 25.0 + (thickness_mm * 5.0),
        }
    }
}

/// LTK dosyası oluşturma isteği.
#[derive(Debug, Serialize, Deserialize)]
pub struct LTKExportRequest {
    pub part_number: String,
    pub part_revision: String,
    pub material: String,
    pub thickness_mm: f64,
    pub machine: String,       // Makine kodu
    pub dxf_content: String,   // DXF verisi (outer contour + holes)
    pub bend_lines: Vec<LTKBendLine>,
    pub description: String,
}

/// LTK büküm çizgisi tanımı.
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LTKBendLine {
    pub x1: f64, pub y1: f64,
    pub x2: f64, pub y2: f64,
    pub angle_deg: f64,        // Büküm açısı
    pub radius_mm: f64,        // İç büküm yarıçapı
    pub direction: String,     // "UP" veya "DOWN"
    pub sequence: usize,       // Büküm sırası
}

/// LTK dosyası yazıcı. Lantek Expert'in anlayacağı formatta çıktı üretir.
#[instrument(skip(req))]
pub fn write_ltk(req: &LTKExportRequest) -> String {
    info!("Writing LTK file for part {}", req.part_number);
    let tech = LantekTechParams::lookup(
        &req.material,
        req.thickness_mm,
        &LantekMachineProfile::Custom(req.machine.clone()),
    );
    
    let mut ltk = String::new();
    
    // --- HEADER ---
    ltk.push_str("[LANTEK_PART]\n");
    ltk.push_str(&format!("VERSION=Expert 28.0\n"));
    ltk.push_str(&format!("GENERATED_BY=Sheet-Metal-Automat-On CAM Bridge v1.0\n"));
    ltk.push_str(&format!("DATE={}\n", chrono_now_str()));
    
    ltk.push_str("\n[HEADER]\n");
    ltk.push_str(&format!("PART_NUMBER={}\n", req.part_number));
    ltk.push_str(&format!("PART_REVISION={}\n", req.part_revision));
    ltk.push_str(&format!("DESCRIPTION={}\n", req.description));
    ltk.push_str("UNITS=MM\n");
    
    // --- MATERIAL ---
    ltk.push_str("\n[MATERIAL]\n");
    ltk.push_str(&format!("MATERIAL_CODE={}\n", tech.material_code));
    ltk.push_str(&format!("THICKNESS={:.3}\n", req.thickness_mm));
    
    // --- TECHNOLOGY (Kesim Parametreleri) ---
    ltk.push_str("\n[TECHNOLOGY]\n");
    ltk.push_str(&format!("MACHINE={}\n", tech.machine_code));
    ltk.push_str(&format!("MATERIAL_GROUP={}\n", tech.material_code));
    ltk.push_str(&format!("THICKNESS_GROUP=T{:.1}\n", req.thickness_mm));
    ltk.push_str(&format!("CUTTING_SPEED_MM_MIN={:.0}\n", tech.cutting_speed_mm_min));
    ltk.push_str(&format!("PIERCING_TIME_SEC={:.2}\n", tech.piercing_time_sec));
    ltk.push_str(&format!("PIERCE_PRESSURE_BAR={:.1}\n", tech.pierce_pressure_bar));
    ltk.push_str(&format!("CUTTING_PRESSURE_BAR={:.1}\n", tech.cutting_pressure_bar));
    ltk.push_str(&format!("FOCUS_POSITION_MM={:.2}\n", tech.focus_position_mm));
    ltk.push_str(&format!("NOZZLE_DIAMETER_MM={:.1}\n", tech.nozzle_diameter_mm));
    ltk.push_str(&format!("GAS_TYPE={}\n", tech.gas_type));
    ltk.push_str(&format!("LASER_POWER_W={:.0}\n", tech.laser_power_w));
    ltk.push_str(&format!("KERF_WIDTH_MM={:.3}\n", tech.kerf_width_mm));
    ltk.push_str(&format!("ROUGHNESS_UM={:.1}\n", tech.roughness_um));
    ltk.push_str(&format!("GAS_CONSUMPTION_L_MIN={:.1}\n", tech.gas_consumption_l_min));
    
    // --- GEOMETRY (Gömülü DXF referansı) ---
    ltk.push_str("\n[GEOMETRY]\n");
    ltk.push_str("FORMAT=DXF_EMBEDDED\n");
    ltk.push_str("BEGIN_DXF\n");
    ltk.push_str(&req.dxf_content);
    ltk.push_str("\nEND_DXF\n");
    
    // --- BENDS (Büküm Bilgileri) ---
    if !req.bend_lines.is_empty() {
        ltk.push_str("\n[BENDS]\n");
        for (idx, bend) in req.bend_lines.iter().enumerate() {
            ltk.push_str(&format!(
                "BEND_{}=ANGLE={:.1} RADIUS={:.2} LINE=({:.3},{:.3})-({:.3},{:.3}) DIRECTION={} SEQUENCE={}\n",
                idx + 1,
                bend.angle_deg,
                bend.radius_mm,
                bend.x1, bend.y1,
                bend.x2, bend.y2,
                bend.direction,
                bend.sequence,
            ));
        }
    }
    
    // --- NC PROGRAM HINT ---
    ltk.push_str("\n[NC_OUTPUT]\n");
    ltk.push_str(&format!("PROGRAM_NAME=SMA_{}\n", req.part_number));
    ltk.push_str("NC_FORMAT=ISO_6983\n"); // G-code ISO standardı
    ltk.push_str("TOOL_CHANGE=NO\n");
    ltk.push_str(&format!("PART_CYCLE_TIME_SEC={:.1}\n",
        estimate_cycle_time(&tech)));
    
    ltk.push_str("\n[END]\n");
    ltk
}

/// Kesim süresini tahmin eder (saniye cinsinden).
fn estimate_cycle_time(tech: &LantekTechParams) -> f64 {
    // Basic heuristic since contour length isn't provided here yet
    let estimated_cut_length = 500.0_f64;
    let estimated_pierces = 5.0_f64;
    
    let cutting_time = (estimated_cut_length / tech.cutting_speed_mm_min) * 60.0;
    let piercing_time = estimated_pierces * tech.piercing_time_sec;
    cutting_time + piercing_time + 15.0 // +15sn setup
}

/// Basit tarih/saat string'i üretir (chrono bağımlılığı olmadan).
fn chrono_now_str() -> String {
    "2026-09-21T23:00:00".to_string() // Üretimde chrono kullanılır
}

// =============================================================
// FAZ 19.B.3 — Lantek Nesting Report (LNR) Okuyucu
// =============================================================

/// Lantek nesting raporu sonucu.
#[derive(Debug, Serialize, Deserialize)]
pub struct LantekNestingReport {
    pub sheet_count: usize,
    pub total_parts: usize,
    pub material_utilization_pct: f64,  // Malzeme kullanım oranı (%)
    pub scrap_pct: f64,                  // Fire oranı (%)
    pub estimated_time_min: f64,         // Tahmini kesim süresi (dk)
    pub parts: Vec<LantekNestedPart>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LantekNestedPart {
    pub part_number: String,
    pub quantity: usize,
    pub sheet_index: usize,
    pub position_x: f64,
    pub position_y: f64,
    pub rotation_deg: f64,
}

/// Lantek LNR rapor dosyasını parse eder.
#[instrument(skip(content))]
pub fn read_lantek_nesting_report(content: &str) -> Result<LantekNestingReport, EngineError> {
    info!("Parsing Lantek Nesting Report");
    let mut sheet_count = 0;
    let mut total_parts = 0;
    let mut utilization = 0.0_f64;
    let mut estimated_time = 0.0_f64;
    let mut parts: Vec<LantekNestedPart> = Vec::new();
    
    for line in content.lines() {
        let line = line.trim();
        
        if let Some(val) = extract_value(line, "SHEET_COUNT=") {
            sheet_count = val.parse().unwrap_or(0);
        }
        if let Some(val) = extract_value(line, "TOTAL_PARTS=") {
            total_parts = val.parse().unwrap_or(0);
        }
        if let Some(val) = extract_value(line, "MATERIAL_UTILIZATION=") {
            utilization = val.trim_end_matches('%').parse().unwrap_or(0.0);
        }
        if let Some(val) = extract_value(line, "ESTIMATED_TIME_MIN=") {
            estimated_time = val.parse().unwrap_or(0.0);
        }
        
        // PART_PLACEMENT satırı: PART_PLACEMENT=PartNo,Qty,Sheet,X,Y,Rot
        if line.starts_with("PART_PLACEMENT=") {
            let data = &line["PART_PLACEMENT=".len()..];
            let cols: Vec<&str> = data.split(',').collect();
            if cols.len() >= 6 {
                parts.push(LantekNestedPart {
                    part_number: cols[0].to_string(),
                    quantity: cols[1].parse().unwrap_or(1),
                    sheet_index: cols[2].parse().unwrap_or(1),
                    position_x: cols[3].parse().unwrap_or(0.0),
                    position_y: cols[4].parse().unwrap_or(0.0),
                    rotation_deg: cols[5].parse().unwrap_or(0.0),
                });
            }
        }
    }
    
    let scrap_pct = 100.0 - utilization;
    
    info!("[Lantek LNR] Parse tamamlandı: {} levha, {} parça, %{:.1} verimlilik",
        sheet_count, total_parts, utilization);
    
    Ok(LantekNestingReport {
        sheet_count,
        total_parts,
        material_utilization_pct: utilization,
        scrap_pct,
        estimated_time_min: estimated_time,
        parts,
    })
}

/// "KEY=VALUE" formatından değeri çıkarır.
fn extract_value<'a>(line: &'a str, key: &str) -> Option<&'a str> {
    if line.starts_with(key) {
        Some(&line[key.len()..])
    } else {
        None
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_tech_lookup_dkp_2mm() {
        let params = LantekTechParams::lookup("DKP", 2.0, &LantekMachineProfile::TrumpfTruLaser3030_6kW);
        assert_eq!(params.cutting_speed_mm_min, 8000.0);
        assert_eq!(params.gas_type, "N2");
    }
    
    #[test]
    fn test_ltk_output_has_required_sections() {
        let req = LTKExportRequest {
            part_number: "TEST-001".to_string(),
            part_revision: "A".to_string(),
            material: "DKP".to_string(),
            thickness_mm: 2.0,
            machine: "TRUMPF_TRULASER_3030_6KW".to_string(),
            dxf_content: "DXF_TEST".to_string(),
            bend_lines: vec![],
            description: "Test parça".to_string(),
        };
        let output = write_ltk(&req);
        assert!(output.contains("[HEADER]"));
        assert!(output.contains("[TECHNOLOGY]"));
        assert!(output.contains("[GEOMETRY]"));
        assert!(output.contains("CUTTING_SPEED_MM_MIN=8000"));
    }
    
    #[test]
    fn test_o2_gas_for_thick_steel() {
        let params = LantekTechParams::lookup("DKP", 6.0, &LantekMachineProfile::Custom("ANY".to_string()));
        assert_eq!(params.gas_type, "O2");
    }
}
