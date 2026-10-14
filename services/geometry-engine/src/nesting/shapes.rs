use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, Serialize, Deserialize)]
pub struct Rect {
    pub x: f64,
    pub y: f64,
    pub width: f64,
    pub height: f64,
}

impl Rect {
    pub fn new(x: f64, y: f64, width: f64, height: f64) -> Self {
        Self { x, y, width, height }
    }

    pub fn area(&self) -> f64 {
        self.width * self.height
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Part {
    pub id: String,
    pub width: f64,
    pub height: f64,
    // Final placement details
    pub placed_x: f64,
    pub placed_y: f64,
    pub rotated: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Sheet {
    pub width: f64,
    pub height: f64,
    pub placed_parts: Vec<Part>,
    pub free_rects: Vec<Rect>,
}

impl Sheet {
    pub fn new(width: f64, height: f64) -> Self {
        Self {
            width,
            height,
            placed_parts: Vec::new(),
            free_rects: vec![Rect::new(0.0, 0.0, width, height)],
        }
    }

    pub fn used_area(&self) -> f64 {
        self.placed_parts.iter().map(|p| p.width * p.height).sum()
    }

    pub fn efficiency(&self) -> f64 {
        let total_area = self.width * self.height;
        if total_area == 0.0 {
            return 0.0;
        }
        (self.used_area() / total_area) * 100.0
    }
}
