use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use crate::cam::dxf::LineSegment;
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct MicroJointRequest {
    pub outer_contour: Vec<LineSegment>,
    pub part_area: f64,
    pub thickness: f64,
    pub material_density_g_cm3: f64,
    pub thin_sheet_joint_width_mm: f64,
    pub med_sheet_joint_width_mm: f64,
    pub thick_sheet_joint_width_mm: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct MicroJointResponse {
    pub joint_locations: Vec<(f64, f64)>,
    pub joint_width: f64,
    pub recommended_count: u32,
    pub is_hazardous_placement: bool,
}

#[instrument(skip(req))]
pub fn generate_micro_joints(req: &MicroJointRequest) -> Result<MicroJointResponse, EngineError> {
    info!("Generating micro joints for contour of length {}", req.outer_contour.len());

    if req.outer_contour.is_empty() {
        return Err(EngineError::InvalidInput("Outer contour cannot be empty".to_string()));
    }

    let weight_kg = (req.part_area / 1000000.0) * req.thickness * req.material_density_g_cm3;

    let joint_width = if req.thickness < 2.0 {
        req.thin_sheet_joint_width_mm 
    } else if req.thickness < 5.0 {
        req.med_sheet_joint_width_mm
    } else {
        req.thick_sheet_joint_width_mm
    };

    let mut recommended_count = 1;
    if weight_kg > 2.0 {
        recommended_count = 2;
    }
    if weight_kg > 10.0 {
        recommended_count = 4;
    }

    let mut sorted_segments = req.outer_contour.clone();
    sorted_segments.sort_by(|a, b| {
        let len_a = ((a.x2 - a.x1).powi(2) + (a.y2 - a.y1).powi(2)).sqrt();
        let len_b = ((b.x2 - b.x1).powi(2) + (b.y2 - b.y1).powi(2)).sqrt();
        len_b.partial_cmp(&len_a).unwrap()
    });

    let mut locations = Vec::new();
    for i in 0..std::cmp::min(recommended_count as usize, sorted_segments.len()) {
        let seg = &sorted_segments[i];
        let mid_x = (seg.x1 + seg.x2) / 2.0;
        let mid_y = (seg.y1 + seg.y2) / 2.0;
        locations.push((mid_x, mid_y));
    }

    let is_hazardous = recommended_count > 2 && req.thickness > 3.0;
    if is_hazardous {
        warn!("Micro joint placement may be hazardous for thickness {} and count {}", req.thickness, recommended_count);
    }

    Ok(MicroJointResponse {
        joint_locations: locations,
        joint_width,
        recommended_count,
        is_hazardous_placement: is_hazardous,
    })
}
