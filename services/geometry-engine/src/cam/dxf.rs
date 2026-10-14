use serde::{Deserialize, Serialize};
use tracing::{info, instrument};
use crate::error::EngineError;

#[derive(Debug, Serialize, Deserialize)]
pub struct DXFRequest {
    pub layer_name: String,
    pub lines: Vec<LineSegment>,
    pub circles: Vec<CircleSegment>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct LineSegment {
    pub x1: f64,
    pub y1: f64,
    pub x2: f64,
    pub y2: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CircleSegment {
    pub cx: f64,
    pub cy: f64,
    pub radius: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct DXFResponse {
    pub dxf_content: String,
}

#[instrument(skip(req))]
pub fn export_dxf(req: &DXFRequest) -> Result<DXFResponse, EngineError> {
    info!("Exporting DXF for layer '{}' with {} lines, {} circles", req.layer_name, req.lines.len(), req.circles.len());
    
    let mut dxf = String::new();
    
    // Header
    dxf.push_str("0\nSECTION\n2\nHEADER\n0\nENDSEC\n");
    // Tables
    dxf.push_str("0\nSECTION\n2\nTABLES\n0\nENDSEC\n");
    // Blocks
    dxf.push_str("0\nSECTION\n2\nBLOCKS\n0\nENDSEC\n");
    // Entities
    dxf.push_str("0\nSECTION\n2\nENTITIES\n");

    for line in &req.lines {
        dxf.push_str("0\nLINE\n");
        dxf.push_str("8\n"); // Layer
        dxf.push_str(&format!("{}\n", req.layer_name));
        dxf.push_str(&format!("10\n{:.4}\n", line.x1));
        dxf.push_str(&format!("20\n{:.4}\n", line.y1));
        dxf.push_str(&format!("11\n{:.4}\n", line.x2));
        dxf.push_str(&format!("21\n{:.4}\n", line.y2));
    }

    for circle in &req.circles {
        dxf.push_str("0\nCIRCLE\n");
        dxf.push_str("8\n"); // Layer
        dxf.push_str(&format!("{}\n", req.layer_name));
        dxf.push_str(&format!("10\n{:.4}\n", circle.cx));
        dxf.push_str(&format!("20\n{:.4}\n", circle.cy));
        dxf.push_str(&format!("40\n{:.4}\n", circle.radius));
    }

    dxf.push_str("0\nENDSEC\n");
    dxf.push_str("0\nEOF\n");

    Ok(DXFResponse { dxf_content: dxf })
}
