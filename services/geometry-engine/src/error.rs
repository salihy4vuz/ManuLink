use thiserror::Error;

#[derive(Error, Debug)]
pub enum EngineError {
    #[error("Invalid Input: {0}")]
    InvalidInput(String),
    
    #[error("Calculation Error: {0}")]
    CalculationError(String),
    
    #[error("System Error: {0}")]
    SystemError(String),
    
    #[error("IO Error: {0}")]
    IoError(#[from] std::io::Error),
    
    #[error("RabbitMQ Error: {0}")]
    LapinError(#[from] lapin::Error),
}
