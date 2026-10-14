use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct LaserGasRequest {
    pub material: String,
    pub thickness_mm: f64,
    pub cutting_time_sec: f64,
    pub high_pressure_base_bar: f64,
    pub high_pressure_multiplier: f64,
    pub low_pressure_base_bar: f64,
    pub low_pressure_multiplier: f64,
    pub n2_flow_multiplier: f64,
    pub o2_flow_multiplier: f64,
    pub n2_cost_per_liter: f64,
    pub o2_cost_per_liter: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LaserGasResult {
    pub gas_type: String,
    pub pressure_bar: f64,
    pub consumption_liters: f64,
    pub estimated_cost_eur: f64,
}

pub struct LaserGasAnalyzer;

impl LaserGasAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &LaserGasRequest) -> LaserGasResult {
        info!("Analyzing laser gas for material {} at {}mm", req.material, req.thickness_mm);
        
        let mat = req.material.to_lowercase();
        
        let (gas_type, mut pressure, flow_rate_liters_per_min) = if mat.contains("stainless") || mat.contains("aluminum") {
            let p = req.high_pressure_base_bar + (req.thickness_mm * req.high_pressure_multiplier); 
            ("Nitrogen (N2)".to_string(), p, p * req.n2_flow_multiplier) 
        } else {
            let p = req.low_pressure_base_bar + (req.thickness_mm * req.low_pressure_multiplier); 
            ("Oxygen (O2)".to_string(), p, p * req.o2_flow_multiplier)
        };

        if pressure > 25.0 {
            pressure = 25.0; 
        }

        let consumption = (req.cutting_time_sec / 60.0) * flow_rate_liters_per_min;
        let cost_per_liter = if gas_type.contains("N2") { req.n2_cost_per_liter } else { req.o2_cost_per_liter };
        let cost_eur = consumption * cost_per_liter;

        LaserGasResult {
            gas_type,
            pressure_bar: pressure,
            consumption_liters: consumption,
            estimated_cost_eur: cost_eur,
        }
    }
}
