pub mod gcode;

pub mod bending;
pub mod deburring;
pub mod collision;
pub mod quality_control;
pub mod packaging;
pub mod painting;
pub mod laser_gas;
pub mod tool_wear;
pub mod assembly;
pub mod flange;
pub mod exporter;
pub mod dxf;
pub mod punching;
pub mod bending_sequence;
pub mod nesting;
pub mod holes;
pub mod micro_joints;
pub mod step_parser;
pub mod welding;
pub mod corner_relief;
pub mod tooling_ai;
pub mod hemming;
pub mod grain_direction;
pub mod forming;
pub mod beveling;
pub mod haz_analytics;
pub mod etching;
pub mod surface_protection;
pub mod hardware;
pub mod machining;
pub mod sheet_metal;

// === FAZ 19: Profesyonel CAD Format ve CAM Entegrasyonu ===
pub mod dxf_reader; // DXF R12-2018 okuma motoru (FAZ 19.A.2)
pub mod iges_reader; // IGES 5.3 okuma motoru (FAZ 19.A.3)
pub mod lantek;      // Lantek Expert CAM entegrasyonu (FAZ 19.B)
