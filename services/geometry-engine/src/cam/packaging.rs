use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct PackagingRequest {
    pub part_length_x: f64,
    pub part_width_y: f64,
    pub part_height_z: f64,
    pub part_weight_kg: f64,
    pub total_quantity: u32,
    
    // Pallet properties
    pub pallet_length_mm: f64,
    pub pallet_width_mm: f64,
    pub pallet_max_height_mm: f64,
    pub pallet_max_weight_kg: f64,
    pub pallet_empty_weight_kg: f64,
    pub cost_per_pallet_eur: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PackagingResult {
    pub pallets_required: u32,
    pub total_gross_weight_kg: f64,
    pub desi_volumetric_weight: f64,
    pub packaging_cost_estimate: f64,
}

pub struct PackagingAnalyzer;

impl PackagingAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &PackagingRequest) -> PackagingResult {
        info!("Analyzing packaging for {} parts", req.total_quantity);

        let total_weight = req.part_weight_kg * req.total_quantity as f64;
        
        let parts_per_layer = ((req.pallet_length_mm / req.part_length_x).floor() * (req.pallet_width_mm / req.part_width_y).floor()).max(1.0);
        let layers_per_pallet = (req.pallet_max_height_mm / req.part_height_z).floor().max(1.0);
        
        let parts_per_pallet = (parts_per_layer * layers_per_pallet) as u32;
        
        let mut pallets_required = 1;
        if parts_per_pallet > 0 {
            pallets_required = (req.total_quantity as f64 / parts_per_pallet as f64).ceil() as u32;
        }

        let pallets_by_weight = (total_weight / req.pallet_max_weight_kg).ceil() as u32;
        let final_pallets = pallets_required.max(pallets_by_weight).max(1);

        let desi_per_pallet = ((req.pallet_length_mm / 10.0) * (req.pallet_width_mm / 10.0) * (req.pallet_max_height_mm / 10.0)) / 3000.0;
        let total_desi = desi_per_pallet * final_pallets as f64;

        let cost = final_pallets as f64 * req.cost_per_pallet_eur;

        PackagingResult {
            pallets_required: final_pallets,
            total_gross_weight_kg: total_weight + (final_pallets as f64 * req.pallet_empty_weight_kg), 
            desi_volumetric_weight: total_desi,
            packaging_cost_estimate: cost,
        }
    }
}
