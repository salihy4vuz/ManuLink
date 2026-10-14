mod error;
mod math;
mod cam;
mod worker;
pub mod nesting;

use tracing::info;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::INFO)
        .init();
    
    info!("Starting Geometry & CAM Engine...");
    
    // Start RabbitMQ worker block
    if let Err(e) = worker::start_worker().await {
        tracing::error!("Worker failed: {:?}", e);
    }

    Ok(())
}
