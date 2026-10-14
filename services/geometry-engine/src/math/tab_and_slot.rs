use serde::{Deserialize, Serialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct TabSlotRequest {
    pub material_thickness: f64,
    pub joint_length: f64,
    pub tab_width_multiplier: f64,
    pub clearance_mm: f64,
    pub distance_between_tabs_mm: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TabSlotResponse {
    pub tab_width: f64,
    pub slot_width: f64,
    pub slot_length: f64,
    pub recommended_count: u32,
}

#[instrument(skip(req))]
pub fn calculate_tab_and_slot(req: &TabSlotRequest) -> TabSlotResponse {
    info!("Calculating tab and slot parameters for joint length {}", req.joint_length);

    let tab_width = req.material_thickness * req.tab_width_multiplier;
    
    let slot_width = req.material_thickness + req.clearance_mm;
    let slot_length = tab_width + req.clearance_mm;

    let recommended_count = (req.joint_length / req.distance_between_tabs_mm).floor() as u32;
    let recommended_count = if recommended_count == 0 { 1 } else { recommended_count };

    TabSlotResponse {
        tab_width,
        slot_width,
        slot_length,
        recommended_count,
    }
}
