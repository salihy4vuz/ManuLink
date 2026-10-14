use serde::{Deserialize, Serialize};
use anyhow::{Context, Result};
use tracing::{info, warn, error, instrument};
use regex::Regex;
use base64::{Engine as _, engine::general_purpose::STANDARD};

#[derive(Deserialize, Debug)]
pub struct MachiningRequest {
    pub step_file_base64: String,
    pub material_density: f64,
    pub machining_type: String,
    // Database/Config provided parameters to avoid hardcoding
    pub time_per_kg_removal_min: f64,
    pub machine_topology_factor: Option<f64>,
    pub tool_diameter_mm: f64,
    pub cutting_speed_m_min: f64,
}

#[derive(Serialize, Debug)]
pub struct MachiningResponse {
    pub raw_bounding_box_volume_mm3: f64,
    pub final_part_volume_mm3: f64,
    pub removed_volume_mm3: f64,
    pub removed_weight_kg: f64,
    pub estimated_machining_time_min: f64,
    pub suggested_rpm: f64,
    pub suggested_feed_rate_mm_min: f64,
}

pub struct MachiningAnalyzer;

impl MachiningAnalyzer {
    #[instrument(skip(req), fields(machining_type = %req.machining_type))]
    pub fn calculate_machining_volume(req: &MachiningRequest) -> Result<MachiningResponse> {
        info!("Starting Enterprise 3D Bounding Box Calculation for {}", req.machining_type);

        if req.step_file_base64.is_empty() {
            let err_msg = "STEP file payload cannot be empty";
            error!("{}", err_msg);
            return Err(anyhow::anyhow!(err_msg));
        }

        let decoded_bytes = STANDARD.decode(&req.step_file_base64)
            .context("Failed to decode Base64 STEP file from RabbitMQ payload")?;
        
        let step_content = String::from_utf8_lossy(&decoded_bytes);
        
        let mut min_x = f64::MAX;
        let mut min_y = f64::MAX;
        let mut min_z = f64::MAX;
        
        let mut max_x = f64::MIN;
        let mut max_y = f64::MIN;
        let mut max_z = f64::MIN;

        let re_point = Regex::new(r"CARTESIAN_POINT\s*\([^,]*,\s*\(\s*([+-]?\d*\.?\d+(?:[eE][+-]?\d+)?)\s*,\s*([+-]?\d*\.?\d+(?:[eE][+-]?\d+)?)\s*,\s*([+-]?\d*\.?\d+(?:[eE][+-]?\d+)?)\s*\)\s*\)")
            .context("Failed to compile REGEX for CARTESIAN_POINT")?;
        
        for cap in re_point.captures_iter(&step_content) {
            if let (Ok(x), Ok(y), Ok(z)) = (
                cap[1].parse::<f64>(),
                cap[2].parse::<f64>(),
                cap[3].parse::<f64>()
            ) {
                if x < min_x { min_x = x; }
                if y < min_y { min_y = y; }
                if z < min_z { min_z = z; }
                
                if x > max_x { max_x = x; }
                if y > max_y { max_y = y; }
                if z > max_z { max_z = z; }
            }
        }

        let re_face = Regex::new(r"ADVANCED_FACE")
            .context("Failed to compile REGEX for ADVANCED_FACE")?;
        let parsed_faces = re_face.find_iter(&step_content).count();

        let (bounding_box_x, bounding_box_y, bounding_box_z) = if min_x == f64::MAX {
            warn!("No valid points found in STEP file; using fallback 100x100x100mm box.");
            (100.0, 100.0, 100.0)
        } else {
            (max_x - min_x, max_y - min_y, max_z - min_z)
        };

        let raw_volume = bounding_box_x * bounding_box_y * bounding_box_z;
        
        let topology_factor = req.machine_topology_factor.unwrap_or_else(|| {
            match parsed_faces {
                0..=6   => 0.85,
                7..=20  => 0.70,
                21..=50 => 0.55,
                _       => 0.40,
            }
        });

        let machine_factor = match req.machining_type.as_str() {
            "lathe" | "swiss_lathe" => 0.30,
            "drilling"              => 0.05,
            "milling" | "milling_3axis" => topology_factor,
            "press"                 => 0.02,
            "welding"               => 0.00,
            _                       => topology_factor,
        };

        let final_volume = raw_volume * (1.0 - machine_factor);
        let removed_volume = raw_volume - final_volume;
        let removed_weight = (removed_volume / 1_000_000.0) * req.material_density;

        let num_teeth = 4.0;
        let feed_per_tooth = 0.05;

        let suggested_rpm = (req.cutting_speed_m_min * 1000.0) / (std::f64::consts::PI * req.tool_diameter_mm);
        let suggested_feed_rate = suggested_rpm * num_teeth * feed_per_tooth;

        info!("Deep Topology analyzed. Bounding Box: {}x{}x{}. Faces detected: {}. Suggested RPM: {:.0}, Feed: {:.0} mm/min", 
            bounding_box_x, bounding_box_y, bounding_box_z, parsed_faces, suggested_rpm, suggested_feed_rate);

        Ok(MachiningResponse {
            raw_bounding_box_volume_mm3: raw_volume,
            final_part_volume_mm3: final_volume,
            removed_volume_mm3: removed_volume,
            removed_weight_kg: removed_weight,
            estimated_machining_time_min: removed_weight * req.time_per_kg_removal_min,
            suggested_rpm: suggested_rpm.round(),
            suggested_feed_rate_mm_min: suggested_feed_rate.round(),
        })
    }
}
