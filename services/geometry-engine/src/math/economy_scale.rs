use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct EconomyScaleRequest {
    pub setup_time_hours: f64,
    pub run_time_per_part_sec: f64,
    pub hourly_machine_rate: f64,
    pub target_quantities: Vec<u32>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EconomyScaleResult {
    pub unit_costs_by_quantity: Vec<(u32, f64)>, 
}

pub struct EconomyScaleAnalyzer;

impl EconomyScaleAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &EconomyScaleRequest) -> EconomyScaleResult {
        info!("Analyzing economy of scale for {} quantities", req.target_quantities.len());
        let mut results = Vec::new();
        let rate_per_sec = req.hourly_machine_rate / 3600.0;
        let setup_cost = req.setup_time_hours * req.hourly_machine_rate;

        for &qty in &req.target_quantities {
            if qty == 0 { continue; }
            
            let run_cost_total = req.run_time_per_part_sec * rate_per_sec * qty as f64;
            let total_cost = setup_cost + run_cost_total;
            let unit_cost = total_cost / qty as f64;
            
            results.push((qty, unit_cost));
        }

        EconomyScaleResult {
            unit_costs_by_quantity: results,
        }
    }
}
