use serde::{Deserialize, Serialize};
use tracing::{info, warn, instrument};

#[derive(Serialize, Deserialize, Debug)]
pub struct HoleData {
    pub diameter: f64,
    pub count: u32,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct HoleCuttingRequest {
    pub thickness: f64,
    pub holes: Vec<HoleData>,
    pub piercing_time_multiplier: f64,
    pub circular_cut_speed_mm_per_min: f64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct HoleCuttingResult {
    pub piercing_time_sec: f64,
    pub cutting_time_sec: f64,
}

#[instrument(skip(req))]
pub fn calculate_holes(req: &HoleCuttingRequest) -> HoleCuttingResult {
    info!("Calculating hole cut times for {} types of holes", req.holes.len());
    
    let mut total_pierce = 0.0;
    let mut total_cut = 0.0;

    let pierce_per_hole = req.piercing_time_multiplier * req.thickness.powf(1.5);

    for hole in &req.holes {
        let count = hole.count as f64;
        total_pierce += pierce_per_hole * count;
        
        let circumference = std::f64::consts::PI * hole.diameter;
        let time_minutes = circumference / req.circular_cut_speed_mm_per_min;
        let cut_sec = time_minutes * 60.0;
        
        total_cut += cut_sec * count;
    }

    HoleCuttingResult {
        piercing_time_sec: total_pierce,
        cutting_time_sec: total_cut,
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct HoleProximityRequest {
    pub thickness: f64,
    pub hole_diameter: f64,
    pub distance_to_bend_line: f64,
    pub proximity_factor: f64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct HoleProximityResult {
    pub is_safe: bool,
    pub minimum_safe_distance: f64,
    pub warning: Option<String>,
}

#[instrument(skip(req))]
pub fn check_hole_proximity(req: &HoleProximityRequest) -> HoleProximityResult {
    info!("Checking hole proximity (distance: {}mm)", req.distance_to_bend_line);

    let minimum_safe_distance = (req.thickness * req.proximity_factor) + (req.hole_diameter / 2.0);
    
    let is_safe = req.distance_to_bend_line >= minimum_safe_distance;
    
    let warning = if !is_safe {
        let w = format!(
            "DANGER: Hole is too close to bend line. It will deform (ovalize). Min safe distance: {:.1}mm, Actual: {:.1}mm",
            minimum_safe_distance, req.distance_to_bend_line
        );
        warn!("{}", w);
        Some(w)
    } else {
        None
    };

    HoleProximityResult {
        is_safe,
        minimum_safe_distance,
        warning,
    }
}
