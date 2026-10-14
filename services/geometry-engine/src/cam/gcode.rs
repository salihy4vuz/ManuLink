use serde::{Deserialize, Serialize};
use tracing::{info, instrument};
use crate::error::EngineError;

#[derive(Debug, Serialize, Deserialize)]
pub struct GCodeRequest {
    pub contours: Vec<Vec<(f64, f64)>>, // List of paths, each path is a list of points
    pub feed_rate: f64,
    pub pierce_delay_ms: Option<u32>,
    pub lead_in_mm: Option<f64>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct GCodeResponse {
    pub gcode_string: String,
    pub estimated_time_sec: f64,
}

#[instrument(skip(req))]
pub fn generate_gcode(req: &GCodeRequest) -> Result<GCodeResponse, EngineError> {
    info!("Generating GCode for {} contours", req.contours.len());
    let mut gcode = String::new();
    let mut estimated_time = 0.0;

    gcode.push_str("%\n");
    gcode.push_str("G21 ; Set units to millimeters\n");
    gcode.push_str("G90 ; Absolute positioning\n");
    gcode.push_str(&format!("F{:.1} ; Set cutting feed rate\n", req.feed_rate));

    let delay_sec = req.pierce_delay_ms.unwrap_or(0) as f64 / 1000.0;
    let lead_in = req.lead_in_mm.unwrap_or(0.0);

    for contour in &req.contours {
        if let Some(first) = contour.first() {
            gcode.push_str("M5 ; Laser/Plasma OFF\n");
            
            let start_x = first.0 - lead_in;
            let start_y = first.1;

            gcode.push_str(&format!("G0 X{:.3} Y{:.3}\n", start_x, start_y));
            
            gcode.push_str("M3 ; Laser/Plasma ON\n");
            if delay_sec > 0.0 {
                gcode.push_str(&format!("G4 P{:.2} ; Pierce delay\n", delay_sec));
                estimated_time += delay_sec;
            }

            if lead_in > 0.0 {
                gcode.push_str(&format!("G1 X{:.3} Y{:.3}\n", first.0, first.1));
                estimated_time += lead_in / (req.feed_rate / 60.0); 
            }

            let mut current_point = (first.0, first.1);
            for point in contour.iter().skip(1) {
                gcode.push_str(&format!("G1 X{:.3} Y{:.3}\n", point.0, point.1));
                
                let dx = point.0 - current_point.0;
                let dy = point.1 - current_point.1;
                let dist = (dx*dx + dy*dy).sqrt();
                estimated_time += dist / (req.feed_rate / 60.0);
                
                current_point = *point;
            }

            if lead_in > 0.0 {
                gcode.push_str(&format!("G1 X{:.3} Y{:.3} ; Lead-out\n", current_point.0 + lead_in, current_point.1));
            }
        }
    }

    gcode.push_str("M5 ; Laser/Plasma OFF\n");
    gcode.push_str("M30 ; End of program\n");
    gcode.push_str("%\n");

    Ok(GCodeResponse { 
        gcode_string: gcode,
        estimated_time_sec: estimated_time,
    })
}
