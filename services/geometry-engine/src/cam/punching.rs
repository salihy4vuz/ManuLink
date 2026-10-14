use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub enum ToolShape {
    Round,
    Rectangle,
    Square,
    Obround,
    SingleD,
    DoubleD,
    Special(String),
}

#[derive(Debug, Serialize, Deserialize)]
pub enum TurretStation {
    A, 
    B, 
    C, 
    D, 
    E, 
}

#[derive(Debug, Serialize, Deserialize)]
pub enum HitStrategy {
    SingleHit,
    Nibbling, 
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PunchRequest {
    pub desired_shape: ToolShape,
    pub dimension_x: f64, 
    pub dimension_y: Option<f64>,
    pub station_a_max_mm: f64,
    pub station_b_max_mm: f64,
    pub station_c_max_mm: f64,
    pub station_d_max_mm: f64,
    pub station_e_max_mm: f64,
    pub clearance_factor: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PunchResponse {
    pub recommended_station: TurretStation,
    pub strategy: HitStrategy,
    pub tool_clearance: f64, 
}

#[instrument(skip(req))]
pub fn calculate_punch_tooling(req: &PunchRequest, material_thickness: f64) -> Result<PunchResponse, EngineError> {
    info!("Calculating punch tooling for {:?}", req.desired_shape);
    if req.dimension_x <= 0.0 || material_thickness <= 0.0 {
        return Err(EngineError::InvalidInput("Dimensions and thickness must be > 0".to_string()));
    }

    let max_dim = match req.desired_shape {
        ToolShape::Round | ToolShape::Square => req.dimension_x,
        _ => {
            let y = req.dimension_y.unwrap_or(req.dimension_x);
            (req.dimension_x.powi(2) + y.powi(2)).sqrt()
        }
    };

    let (station, strategy) = if max_dim <= req.station_a_max_mm {
        (TurretStation::A, HitStrategy::SingleHit)
    } else if max_dim <= req.station_b_max_mm {
        (TurretStation::B, HitStrategy::SingleHit)
    } else if max_dim <= req.station_c_max_mm {
        (TurretStation::C, HitStrategy::SingleHit)
    } else if max_dim <= req.station_d_max_mm {
        (TurretStation::D, HitStrategy::SingleHit)
    } else if max_dim <= req.station_e_max_mm {
        (TurretStation::E, HitStrategy::SingleHit)
    } else {
        (TurretStation::D, HitStrategy::Nibbling)
    };

    let tool_clearance = material_thickness * req.clearance_factor;

    Ok(PunchResponse {
        recommended_station: station,
        strategy,
        tool_clearance,
    })
}
