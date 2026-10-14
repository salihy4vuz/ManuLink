use serde::{Deserialize, Serialize};
use tracing::{info, instrument};

#[derive(Serialize, Deserialize, Debug)]
pub struct FlangeData {
    pub length: f64,
    pub angle_deg: f64,
    pub inner_radius: f64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct FlangeResult {
    pub bend_allowance: f64,
    pub bend_deduction: f64,
    pub flat_pattern_length: f64, // Açınım boyu
}

#[instrument]
pub fn calculate_bend_allowance(angle_deg: f64, inner_radius: f64, k_factor: f64, thickness: f64) -> f64 {
    let angle_rad = angle_deg.to_radians();
    // BA = A * (R + K * T)
    angle_rad * (inner_radius + k_factor * thickness)
}

#[instrument(skip(flanges))]
pub fn calculate_flange_flat_length(k_factor: f64, thickness: f64, flanges: &[FlangeData]) -> f64 {
    info!("Calculating flat length for {} flanges", flanges.len());
    let mut total_flat_length = 0.0;

    for flange in flanges {
        let ba = calculate_bend_allowance(flange.angle_deg, flange.inner_radius, k_factor, thickness);
        total_flat_length += flange.length + ba;
    }

    total_flat_length
}
