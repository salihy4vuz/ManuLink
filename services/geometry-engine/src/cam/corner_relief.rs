use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub enum ReliefShape {
    Circular,
    Rectangular,
    Tear,
    Obround,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReliefRequest {
    pub thickness_mm: f64,
    pub bend_radius_mm: f64,
    pub requested_shape: ReliefShape,
    // Database params
    pub safety_margin_mm: f64,
    pub min_laser_cut_width_mm: f64,
    pub tear_slit_width_mm: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ReliefResult {
    pub relief_width_mm: f64,
    pub relief_depth_mm: f64,
    pub is_tearing_prevented: bool,
    pub shape_applied: ReliefShape,
}

pub struct CornerReliefEngine;

impl CornerReliefEngine {
    #[instrument(skip(req))]
    pub fn calculate_relief(req: &ReliefRequest) -> ReliefResult {
        info!("Calculating corner relief for thickness: {}mm", req.thickness_mm);

        let min_depth = req.thickness_mm + req.bend_radius_mm + req.safety_margin_mm;
        let min_width = req.thickness_mm.max(req.min_laser_cut_width_mm); 

        let (final_width, final_depth) = match req.requested_shape {
            ReliefShape::Circular => {
                let diameter = min_depth.max(min_width);
                (diameter, diameter)
            },
            ReliefShape::Rectangular | ReliefShape::Obround => (min_width, min_depth),
            ReliefShape::Tear => (req.tear_slit_width_mm, min_depth), 
        };

        ReliefResult {
            relief_width_mm: final_width,
            relief_depth_mm: final_depth,
            is_tearing_prevented: true,
            shape_applied: req.requested_shape.clone(),
        }
    }
}
