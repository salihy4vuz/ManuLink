use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub enum FlangeType {
    Base,
    Edge,
    HemOpen,
    HemClosed,
    HemTeardrop,
    Joggle,
    Miter,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FlangeRequest {
    pub flange_type: FlangeType,
    pub thickness: f64,
    pub length: f64,
    pub internal_radius: Option<f64>,
    pub miter_angle_degrees: Option<f64>, 
    
    // Config params
    pub edge_deduction_multiplier: f64,
    pub hem_allowance_multiplier: f64,
    pub joggle_offset_multiplier: f64,
    pub joggle_allowance_multiplier: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FlangeResponse {
    pub flat_length: f64,
    pub deformation_allowance: f64, 
}

#[instrument(skip(req))]
pub fn calculate_flange(req: &FlangeRequest) -> Result<FlangeResponse, EngineError> {
    info!("Calculating flange parameters for {:?}", req.flange_type);

    if req.thickness <= 0.0 || req.length <= 0.0 {
        return Err(EngineError::InvalidInput("Thickness and Length must be positive".to_string()));
    }

    let ir = req.internal_radius.unwrap_or(req.thickness);

    let (flat_len, def_allowance) = match req.flange_type {
        FlangeType::Base => {
            (req.length, 0.0)
        },
        FlangeType::Edge => {
            let deduction = req.edge_deduction_multiplier * req.thickness;
            (req.length - deduction, deduction)
        },
        FlangeType::HemOpen => {
            let allowance = std::f64::consts::PI * (ir + req.hem_allowance_multiplier * req.thickness);
            (req.length + allowance, allowance)
        },
        FlangeType::HemClosed => {
            let allowance = req.hem_allowance_multiplier * req.thickness;
            (req.length + allowance, allowance)
        },
        FlangeType::HemTeardrop => {
            let allowance = std::f64::consts::PI * (ir + req.hem_allowance_multiplier * req.thickness) + req.thickness;
            (req.length + allowance, allowance)
        },
        FlangeType::Joggle => {
            let offset_distance = req.joggle_offset_multiplier * req.thickness;
            let allowance = offset_distance * req.joggle_allowance_multiplier; 
            (req.length + allowance, allowance)
        },
        FlangeType::Miter => {
            let angle = req.miter_angle_degrees.unwrap_or(45.0);
            if angle <= 0.0 || angle >= 90.0 {
                return Err(EngineError::InvalidInput("Miter angle must be between 0 and 90 degrees".to_string()));
            }
            let deduction = req.length * (angle * std::f64::consts::PI / 180.0).tan();
            (req.length - deduction, deduction)
        }
    };

    Ok(FlangeResponse {
        flat_length: flat_len,
        deformation_allowance: def_allowance,
    })
}
