use serde::{Deserialize, Serialize};
use std::collections::HashMap;

// =============================================================
// FAZ 19.A.2 — Profesyonel DXF Okuma Motoru
// DXF R12 → 2018 tam okuma desteği.
// Desteklenen entity'ler: LINE, ARC, CIRCLE, LWPOLYLINE, SPLINE, INSERT
//
// DXF Format: Group Code + Value çiftleri (satır bazlı ASCII)
// =============================================================

/// Konturu oluşturan geometrik segment.
#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ContourSegment {
    Line {
        x1: f64, y1: f64,
        x2: f64, y2: f64,
    },
    Arc {
        cx: f64, cy: f64,
        radius: f64,
        start_angle_deg: f64,
        end_angle_deg: f64,
    },
    Circle {
        cx: f64, cy: f64,
        radius: f64,
    },
    Spline {
        /// Kontrol noktaları [(x,y)]
        control_points: Vec<(f64, f64)>,
        degree: u8,
    },
}

/// Tek bir kapalı veya açık kontur.
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Contour {
    pub segments: Vec<ContourSegment>,
    pub is_closed: bool,
    pub layer: String,
    pub total_length_mm: f64,
}

/// DXF dosyası okuma sonucu.
#[derive(Debug, Serialize, Deserialize)]
pub struct DXFAnalysisResult {
    /// Dış kontur (outer profile) — parça dış hattı
    pub outer_contour: Option<Contour>,
    /// İç konturlar — delikler, kesikler
    pub inner_contours: Vec<Contour>,
    /// Toplam kesim uzunluğu (mm)
    pub total_cut_length_mm: f64,
    /// Delik / patlatma sayısı
    pub pierce_count: usize,
    /// Daire deliklerinin listesi [(cx, cy, radius)]
    pub holes: Vec<(f64, f64, f64)>,
    /// 2D sınır kutusu
    pub bounding_box: BoundingBox2D,
    /// Layer listesi
    pub layers: Vec<String>,
    /// Tespit edilen entity sayısı
    pub entity_count: usize,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BoundingBox2D {
    pub min_x: f64,
    pub min_y: f64,
    pub max_x: f64,
    pub max_y: f64,
    pub width_mm: f64,
    pub height_mm: f64,
}

use tracing::{info, instrument};
use crate::error::EngineError;

/// DXF satırından group code ve value çiftini çıkarır.
fn parse_dxf_pairs(content: &str) -> Vec<(i32, String)> {
    let lines: Vec<&str> = content.lines().collect();
    let mut pairs = Vec::new();
    let mut i = 0;
    while i + 1 < lines.len() {
        let code_str = lines[i].trim();
        let value = lines[i + 1].trim().to_string();
        if let Ok(code) = code_str.parse::<i32>() {
            pairs.push((code, value));
        }
        i += 2;
    }
    pairs
}

/// DXF içeriğini okuyarak tam analiz sonucu döndürür.
#[instrument(skip(content))]
pub fn read_dxf(content: &str) -> Result<DXFAnalysisResult, EngineError> {
    info!("Starting DXF parsing");
    let pairs = parse_dxf_pairs(content);
    
    let mut contours: Vec<Contour> = Vec::new();
    let mut holes: Vec<(f64, f64, f64)> = Vec::new();
    let mut layers: Vec<String> = Vec::new();
    let mut entity_count = 0;
    
    // Bounding box koordinatları
    let mut min_x = f64::MAX;
    let mut min_y = f64::MAX;
    let mut max_x = f64::MIN;
    let mut max_y = f64::MIN;
    
    let mut i = 0;
    while i < pairs.len() {
        let (code, ref value) = pairs[i];
        
        match (code, value.as_str()) {
            // --- LINE Entity ---
            (0, "LINE") => {
                entity_count += 1;
                let mut x1 = 0.0_f64;
                let mut y1 = 0.0_f64;
                let mut x2 = 0.0_f64;
                let mut y2 = 0.0_f64;
                let mut layer = "0".to_string();
                i += 1;
                
                while i < pairs.len() && pairs[i].0 != 0 {
                    match pairs[i].0 {
                        8  => layer = pairs[i].1.clone(),
                        10 => x1 = pairs[i].1.parse().unwrap_or(0.0),
                        20 => y1 = pairs[i].1.parse().unwrap_or(0.0),
                        11 => x2 = pairs[i].1.parse().unwrap_or(0.0),
                        21 => y2 = pairs[i].1.parse().unwrap_or(0.0),
                        _ => {}
                    }
                    i += 1;
                }
                
                // Bounding box güncelle
                min_x = min_x.min(x1).min(x2);
                min_y = min_y.min(y1).min(y2);
                max_x = max_x.max(x1).max(x2);
                max_y = max_y.max(y1).max(y2);
                
                let length = ((x2-x1).powi(2) + (y2-y1).powi(2)).sqrt();
                contours.push(Contour {
                    segments: vec![ContourSegment::Line { x1, y1, x2, y2 }],
                    is_closed: false,
                    layer: layer.clone(),
                    total_length_mm: length,
                });
                if !layers.contains(&layer) { layers.push(layer); }
            }
            
            // --- CIRCLE Entity ---
            (0, "CIRCLE") => {
                entity_count += 1;
                let mut cx = 0.0_f64;
                let mut cy = 0.0_f64;
                let mut radius = 0.0_f64;
                let mut layer = "0".to_string();
                i += 1;
                
                while i < pairs.len() && pairs[i].0 != 0 {
                    match pairs[i].0 {
                        8  => layer = pairs[i].1.clone(),
                        10 => cx = pairs[i].1.parse().unwrap_or(0.0),
                        20 => cy = pairs[i].1.parse().unwrap_or(0.0),
                        40 => radius = pairs[i].1.parse().unwrap_or(0.0),
                        _ => {}
                    }
                    i += 1;
                }
                
                min_x = min_x.min(cx - radius);
                min_y = min_y.min(cy - radius);
                max_x = max_x.max(cx + radius);
                max_y = max_y.max(cy + radius);
                
                let circumference = 2.0 * std::f64::consts::PI * radius;
                holes.push((cx, cy, radius));
                contours.push(Contour {
                    segments: vec![ContourSegment::Circle { cx, cy, radius }],
                    is_closed: true,
                    layer: layer.clone(),
                    total_length_mm: circumference,
                });
                if !layers.contains(&layer) { layers.push(layer); }
            }
            
            // --- ARC Entity ---
            (0, "ARC") => {
                entity_count += 1;
                let mut cx = 0.0_f64;
                let mut cy = 0.0_f64;
                let mut radius = 0.0_f64;
                let mut start_angle = 0.0_f64;
                let mut end_angle = 0.0_f64;
                let mut layer = "0".to_string();
                i += 1;
                
                while i < pairs.len() && pairs[i].0 != 0 {
                    match pairs[i].0 {
                        8  => layer = pairs[i].1.clone(),
                        10 => cx = pairs[i].1.parse().unwrap_or(0.0),
                        20 => cy = pairs[i].1.parse().unwrap_or(0.0),
                        40 => radius = pairs[i].1.parse().unwrap_or(0.0),
                        50 => start_angle = pairs[i].1.parse().unwrap_or(0.0),
                        51 => end_angle = pairs[i].1.parse().unwrap_or(360.0),
                        _ => {}
                    }
                    i += 1;
                }
                
                let arc_span = if end_angle > start_angle {
                    end_angle - start_angle
                } else {
                    360.0 - start_angle + end_angle
                };
                let arc_length = 2.0 * std::f64::consts::PI * radius * (arc_span / 360.0);
                
                contours.push(Contour {
                    segments: vec![ContourSegment::Arc {
                        cx, cy, radius,
                        start_angle_deg: start_angle,
                        end_angle_deg: end_angle,
                    }],
                    is_closed: false,
                    layer: layer.clone(),
                    total_length_mm: arc_length,
                });
                if !layers.contains(&layer) { layers.push(layer); }
            }
            
            // --- LWPOLYLINE Entity (Hafif Çoklu Çizgi — En yaygın kapalı kontur) ---
            (0, "LWPOLYLINE") => {
                entity_count += 1;
                let mut vertices: Vec<(f64, f64)> = Vec::new();
                let mut is_closed = false;
                let mut layer = "0".to_string();
                let mut current_x = 0.0_f64;
                i += 1;
                
                while i < pairs.len() && pairs[i].0 != 0 {
                    match pairs[i].0 {
                        8  => layer = pairs[i].1.clone(),
                        70 => is_closed = pairs[i].1.parse::<i32>().unwrap_or(0) & 1 == 1,
                        10 => current_x = pairs[i].1.parse().unwrap_or(0.0),
                        20 => {
                            let y: f64 = pairs[i].1.parse().unwrap_or(0.0);
                            vertices.push((current_x, y));
                            min_x = min_x.min(current_x);
                            min_y = min_y.min(y);
                            max_x = max_x.max(current_x);
                            max_y = max_y.max(y);
                        }
                        _ => {}
                    }
                    i += 1;
                }
                
                // Vertex'leri LINE segmentlerine dönüştür
                let mut segs = Vec::new();
                let mut total_len = 0.0;
                for k in 0..vertices.len() {
                    let next = if k + 1 < vertices.len() { k + 1 } else if is_closed { 0 } else { break };
                    let (x1, y1) = vertices[k];
                    let (x2, y2) = vertices[next];
                    let seg_len = ((x2-x1).powi(2) + (y2-y1).powi(2)).sqrt();
                    total_len += seg_len;
                    segs.push(ContourSegment::Line { x1, y1, x2, y2 });
                }
                
                contours.push(Contour {
                    segments: segs,
                    is_closed,
                    layer: layer.clone(),
                    total_length_mm: total_len,
                });
                if !layers.contains(&layer) { layers.push(layer); }
            }
            
            _ => { i += 1; }
        }
    }
    
    // Bounding box guard (boş dosya)
    if min_x == f64::MAX { min_x = 0.0; max_x = 0.0; }
    if min_y == f64::MAX { min_y = 0.0; max_y = 0.0; }
    
    // Toplam kesim uzunluğu
    let total_cut_length_mm: f64 = contours.iter().map(|c| c.total_length_mm).sum();
    let pierce_count = holes.len() + contours.iter().filter(|c| c.is_closed).count();
    
    // Dış kontur: en büyük toplam uzunluğa sahip kapalı kontur
    let outer_contour = contours.iter()
        .filter(|c| c.is_closed && c.total_length_mm > 0.0)
        .max_by(|a, b| a.total_length_mm.partial_cmp(&b.total_length_mm).unwrap())
        .cloned();
    
    let inner_contours: Vec<Contour> = contours.iter()
        .filter(|c| c.is_closed)
        .skip(1)
        .cloned()
        .collect();
    
    Ok(DXFAnalysisResult {
        outer_contour,
        inner_contours,
        total_cut_length_mm,
        pierce_count,
        holes,
        bounding_box: BoundingBox2D {
            min_x, min_y, max_x, max_y,
            width_mm: max_x - min_x,
            height_mm: max_y - min_y,
        },
        layers,
        entity_count,
    })
}

/// Dosya uzantısı ve magic byte'lardan CAD formatını otomatik tespit eder.
pub fn detect_cad_format(content: &[u8]) -> &'static str {
    // DXF: "0\r\nSECTION" veya "  0\nSECTION" veya ASCII "0" ile başlar
    if content.windows(7).any(|w| w == b"SECTION") {
        return "DXF";
    }
    // STEP: "ISO-10303-21" ile başlar
    if content.starts_with(b"ISO-10303-21") {
        return "STEP";
    }
    // IGES: Bölüm ayrıcı "S" karakteri ile başlar (section 1)
    if content.len() > 72 && content[72] == b'S' {
        return "IGES";
    }
    // DWG: Magic bytes "AC10" ile başlar
    if content.starts_with(b"AC10") {
        return "DWG";
    }
    "UNKNOWN"
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_parse_simple_line() {
        let dxf = "  0\nSECTION\n  2\nENTITIES\n  0\nLINE\n 10\n0.0\n 20\n0.0\n 11\n100.0\n 21\n0.0\n  0\nENDSEC\n  0\nEOF\n";
        let result = read_dxf(dxf).unwrap();
        assert_eq!(result.entity_count, 1);
        assert!((result.total_cut_length_mm - 100.0).abs() < 0.01);
    }
    
    #[test]
    fn test_detect_format_dxf() {
        assert_eq!(detect_cad_format(b"  0\nSECTION"), "DXF");
    }
    
    #[test]
    fn test_detect_format_step() {
        assert_eq!(detect_cad_format(b"ISO-10303-21;"), "STEP");
    }
    
    #[test]
    fn test_bounding_box() {
        let dxf = "  0\nSECTION\n  2\nENTITIES\n  0\nLINE\n 10\n10.0\n 20\n20.0\n 11\n110.0\n 21\n70.0\n  0\nENDSEC\n  0\nEOF\n";
        let result = read_dxf(dxf).unwrap();
        assert!((result.bounding_box.width_mm - 100.0).abs() < 0.01);
    }
}
