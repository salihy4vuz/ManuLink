use lapin::{
    options::*, types::FieldTable, BasicProperties, Connection, ConnectionProperties,
};
use deadpool_lapin::{Config, Manager, Pool};
use tracing::{instrument, info, error};
use serde_json::Value;
use futures_lite::stream::StreamExt;
use std::sync::Arc;

use crate::error::EngineError;
use crate::math::{
    flat_pattern::{calculate_flat_pattern, FlatPatternRequest},
    k_factor::{DynamicKFactorCalculator, KFactorCalculator, KFactorRequest},
    flange::{calculate_flange, FlangeRequest},
    kerf::{calculate_kerf, KerfRequest},
    economy_scale::{EconomyScaleAnalyzer, EconomyScaleRequest},
    sustainability::{SustainabilityAnalyzer, SustainabilityRequest},
    scrap_value::{ScrapAnalyzer, ScrapValueRequest},
};
use crate::cam::{
    nesting::{calculate_nesting, NestingRequest},
    bending::{calculate_bending_force, BendingRequest},
    exporter::{export_to_gltf, GLTFExportRequest},
    dxf::{export_dxf, DXFRequest},
    punching::{calculate_punch_tooling, PunchRequest},
    bending_sequence::{optimize_bending_sequence, SequenceRequest},
    collision::{detect_collision, CollisionRequest},
    micro_joints::{generate_micro_joints, MicroJointRequest},
    step_parser::StepParser,
    welding::{WeldingAnalyzer, WeldingJoint},
    corner_relief::{CornerReliefEngine, ReliefRequest},
    tooling_ai::{ToolingAI, ToolingRequest},
    hemming::{HemmingAnalyzer, HemmingRequest},
    holes::{HoleProximityRequest, check_hole_proximity},
    grain_direction::{GrainDirectionAnalyzer, GrainDirectionRequest},
    forming::{FormingAnalyzer, FormingRequest},
    beveling::{BevelAnalyzer, BevelRequest},
    haz_analytics::{HazAnalyzer, HazRequest},
    etching::{EtchingAnalyzer, EtchingRequest},
    surface_protection::{SurfaceAnalyzer, SurfaceProtectionRequest},
    hardware::{HardwareAnalyzer, HardwareRequest},
    deburring::{DeburringAnalyzer, DeburringRequest},
    quality_control::{QCAnalyzer, QualityControlRequest},
    packaging::{PackagingAnalyzer, PackagingRequest},
    painting::{PaintingAnalyzer, PaintingRequest},
    laser_gas::{LaserGasAnalyzer, LaserGasRequest},
    tool_wear::{ToolWearAnalyzer, ToolWearRequest},
    assembly::{AssemblyAnalyzer, AssemblyRequest},
    machining::{MachiningAnalyzer, MachiningRequest},
    sheet_metal::{SheetMetalAnalyzer, SheetMetalRequest},
};

pub async fn route_command(task_type: &str, payload: Value) -> Result<String, EngineError> {
    match task_type {
        "k_factor" => {
            let req: KFactorRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let calc = DynamicKFactorCalculator;
            let res = calc.calculate(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "flat_pattern" => {
            let req: FlatPatternRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_flat_pattern(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "nesting" => {
            let req: NestingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_nesting(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "bending_force" => {
            let req: BendingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_bending_force(&req); // Doesn't fail yet, but could be refactored
            Ok(serde_json::to_string(&res).unwrap())
        },
        "export_gltf" => {
            let req: GLTFExportRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = export_to_gltf(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "generate_gcode" => {
            let req: crate::cam::gcode::GCodeRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = crate::cam::gcode::generate_gcode(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "export_dxf" => {
            let req: DXFRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = export_dxf(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "flange" => {
            let req: FlangeRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_flange(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "punching" => {
            let req: PunchRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_punch_tooling(&req, 2.0)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "bending_sequence" => {
            let req: SequenceRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = optimize_bending_sequence(&req).await?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "collision" => {
            let req: CollisionRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = detect_collision(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "micro_joints" => {
            let req: MicroJointRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = generate_micro_joints(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "kerf" => {
            let req: KerfRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = calculate_kerf(&req)?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_step" => {
            let req: String = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let parser = StepParser::new();
            let res = parser.analyze_topology().map_err(|e| EngineError::SystemError(e))?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_welding" => {
            let req: WeldingJoint = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = WeldingAnalyzer::calculate_welding_params(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_relief" => {
            let req: ReliefRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = CornerReliefEngine::calculate_relief(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "select_tooling" => {
            let req: ToolingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = ToolingAI::select_tooling(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_hemming" => {
            let req: HemmingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = HemmingAnalyzer::calculate(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "check_hole_proximity" => {
            let req: HoleProximityRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = check_hole_proximity(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_grain_direction" => {
            let req: GrainDirectionRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = GrainDirectionAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_forming" => {
            let req: FormingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = FormingAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_beveling" => {
            let req: BevelRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = BevelAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_haz" => {
            let req: HazRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = HazAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_etching" => {
            let req: EtchingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = EtchingAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_surface_protection" => {
            let req: SurfaceProtectionRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = SurfaceAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_hardware" => {
            let req: HardwareRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = HardwareAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_deburring" => {
            let req: DeburringRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = DeburringAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_economy_scale" => {
            let req: EconomyScaleRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = EconomyScaleAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_sustainability" => {
            let req: SustainabilityRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = SustainabilityAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_scrap_value" => {
            let req: ScrapValueRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = ScrapAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_quality_control" => {
            let req: QualityControlRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = QCAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_packaging" => {
            let req: PackagingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = PackagingAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_painting" => {
            let req: PaintingRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = PaintingAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_laser_gas" => {
            let req: LaserGasRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = LaserGasAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_tool_wear" => {
            let req: ToolWearRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = ToolWearAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "calculate_assembly" => {
            let req: AssemblyRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = AssemblyAnalyzer::analyze(&req);
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_machining" => {
            let req: MachiningRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let req = if req.machining_type.is_empty() {
                MachiningRequest { machining_type: "milling_3axis".to_string(), ..req }
            } else {
                req
            };
            let res = MachiningAnalyzer::calculate_machining_volume(&req).map_err(|e| EngineError::SystemError(e.to_string()))?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        "analyze_sheet_metal" => {
            let req: SheetMetalRequest = serde_json::from_value(payload).map_err(|e| EngineError::InvalidInput(e.to_string()))?;
            let res = SheetMetalAnalyzer::analyze(&req).map_err(|e| EngineError::SystemError(e.to_string()))?;
            Ok(serde_json::to_string(&res).unwrap())
        },
        _ => Err(EngineError::SystemError(format!("Unknown task type: {}", task_type)))
    }
}



#[instrument]
pub async fn start_worker() -> Result<(), EngineError> {
    let addr = std::env::var("RABBITMQ_URL").unwrap_or_else(|_| "amqp://guest:guest@localhost:5672/".into());
    let mut cfg = Config::default();
    cfg.url = Some(addr);
    
    let pool = cfg.create_pool(Some(deadpool_lapin::Runtime::Tokio1))
        .map_err(|e| EngineError::SystemError(format!("Pool creation error: {}", e)))?;

    let conn = pool.get().await.map_err(|e| EngineError::SystemError(format!("RabbitMQ connection failed: {}", e)))?;
    let channel = conn.create_channel().await?;

    // DLX Setup for dead letters
    let mut args = FieldTable::default();
    args.insert(
        "x-dead-letter-exchange".into(),
        lapin::types::AMQPValue::LongString("geometry_dlx".into()),
    );
    args.insert(
        "x-dead-letter-routing-key".into(),
        lapin::types::AMQPValue::LongString("geometry_dlq".into()),
    );

    // Ensure exchange and queue for DLX exist (simple setup)
    channel.exchange_declare("geometry_dlx", lapin::ExchangeKind::Direct, ExchangeDeclareOptions::default(), FieldTable::default()).await?;
    channel.queue_declare("geometry_dlq", QueueDeclareOptions::default(), FieldTable::default()).await?;
    channel.queue_bind("geometry_dlq", "geometry_dlx", "geometry_dlq", QueueBindOptions::default(), FieldTable::default()).await?;

    channel.queue_declare(
        "geometry_tasks",
        QueueDeclareOptions::default(),
        args,
    ).await?;

    let mut consumer = channel.basic_consume(
        "geometry_tasks",
        "geometry_engine_enterprise",
        BasicConsumeOptions::default(),
        FieldTable::default(),
    ).await?;

    info!("Enterprise Geometry engine listening for tasks on 'geometry_tasks' queue with Deadpool");

    let amqp_channel = Arc::new(channel);

    while let Some(delivery) = consumer.next().await {
        if let Ok(delivery) = delivery {
            let channel_clone = Arc::clone(&amqp_channel);

            tokio::spawn(async move {
                let payload_str = String::from_utf8_lossy(&delivery.data);
                let mut response_payload = String::new();

                if let Ok(json) = serde_json::from_str::<Value>(&payload_str) {
                    let task_type = json["task"].as_str().unwrap_or("unknown");
                    info!("Processing task: {}", task_type);
                    
                    match route_command(task_type, json["payload"].clone()).await {
                        Ok(res) => response_payload = res,
                        Err(e) => {
                            error!("Command Error ({}): {}", task_type, e);
                            response_payload = format!(r#"{{"error": "{}"}}"#, e);
                        }
                    }
                } else {
                    // Send to DLX by rejecting with requeue=false
                    error!("Invalid JSON Payload. Sending to DLX.");
                    let _ = delivery.reject(lapin::options::BasicRejectOptions { requeue: false }).await;
                    return;
                }
                
                // RPC Reply
                if let Some(reply_to) = delivery.properties.reply_to().as_ref() {
                    if let Some(correlation_id) = delivery.properties.correlation_id().as_ref() {
                        let props = BasicProperties::default().with_correlation_id(correlation_id.clone());
                        
                        let _ = channel_clone.basic_publish(
                            "",
                            reply_to.as_str(),
                            BasicPublishOptions::default(),
                            response_payload.as_bytes(),
                            props,
                        ).await;
                    }
                }

                let _ = delivery.ack(BasicAckOptions::default()).await;
            });
        }
    }

    Ok(())
}
