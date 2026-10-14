use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct WeldingJoint {
    pub material: String,
    pub process: WeldingProcess,
    pub length_mm: f64,
    pub thickness_mm: f64,
    pub joint_type: JointType,
    pub mig_base_speed: f64,
    pub tig_base_speed: f64,
    pub laser_base_speed: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum WeldingProcess {
    MIG,
    TIG,
    Laser,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum JointType {
    Butt,  
    TJoint, 
    Corner, 
    Lap,    
}

#[derive(Debug, Serialize, Deserialize)]
pub struct WeldingAnalysisResult {
    pub required_wire_length_mm: f64,
    pub estimated_welding_time_sec: f64,
    pub recommended_gap_mm: f64,
    pub tab_and_slot_count: usize,
}

pub struct WeldingAnalyzer;

impl WeldingAnalyzer {
    #[instrument(skip(joint))]
    pub fn calculate_welding_params(joint: &WeldingJoint) -> WeldingAnalysisResult {
        info!("Calculating welding params for {:?} joint of length {}mm", joint.joint_type, joint.length_mm);

        let base_speed = match joint.process {
            WeldingProcess::MIG => joint.mig_base_speed, 
            WeldingProcess::TIG => joint.tig_base_speed,
            WeldingProcess::Laser => joint.laser_base_speed,
        };

        let mat_factor = match joint.material.to_lowercase().as_str() {
            "aluminum" => 0.8, 
            "stainless steel" | "stainless" => 0.9,
            _ => 1.0, 
        };

        let thickness_factor = (2.0 / joint.thickness_mm.max(0.5)).powf(0.85);

        let speed_mm_per_min = base_speed * mat_factor * thickness_factor;
        let time_sec = (joint.length_mm / speed_mm_per_min) * 60.0;
        
        let gap = match joint.joint_type {
            JointType::Butt => joint.thickness_mm * 0.5,
            JointType::TJoint => 0.0,
            JointType::Corner => joint.thickness_mm * 0.3,
            JointType::Lap => 0.0,
        };

        let tab_count = (joint.length_mm / 100.0).floor() as usize;

        WeldingAnalysisResult {
            required_wire_length_mm: joint.length_mm * 1.2, 
            estimated_welding_time_sec: time_sec,
            recommended_gap_mm: gap,
            tab_and_slot_count: tab_count,
        }
    }
}
