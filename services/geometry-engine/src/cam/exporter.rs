use serde::{Deserialize, Serialize};
use base64::{Engine as _, engine::general_purpose};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct GLTFExportRequest {
    pub part_id: String,
    pub vertices: Vec<[f32; 3]>,
    pub indices: Vec<u16>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct GLTFExportResponse {
    pub gltf_json: String,
}

#[instrument(skip(req))]
pub fn export_to_gltf(req: &GLTFExportRequest) -> GLTFExportResponse {
    info!("Exporting GLTF for part {} with {} vertices", req.part_id, req.vertices.len());
    
    let mut buffer_bytes = Vec::new();
    
    // Write vertices
    for v in &req.vertices {
        for &coord in v {
            buffer_bytes.extend_from_slice(&coord.to_le_bytes());
        }
    }
    let vertices_byte_length = buffer_bytes.len();
    
    // Write indices
    for &idx in &req.indices {
        buffer_bytes.extend_from_slice(&idx.to_le_bytes());
    }
    let indices_byte_length = buffer_bytes.len() - vertices_byte_length;
    
    let b64 = general_purpose::STANDARD.encode(&buffer_bytes);
    let uri = format!("data:application/octet-stream;base64,{}", b64);
    
    let gltf = serde_json::json!({
        "asset": { "version": "2.0", "generator": "SheetMetal CAM Engine" },
        "scene": 0,
        "scenes": [ { "nodes": [0] } ],
        "nodes": [ { "mesh": 0, "name": req.part_id } ],
        "meshes": [ {
            "primitives": [ {
                "attributes": { "POSITION": 1 },
                "indices": 0,
                "material": 0
            } ]
        } ],
        "materials": [ {
            "pbrMetallicRoughness": {
                "baseColorFactor": [0.2, 0.7, 0.9, 1.0],
                "metallicFactor": 0.8,
                "roughnessFactor": 0.2
            }
        } ],
        "buffers": [ { "uri": uri, "byteLength": buffer_bytes.len() } ],
        "bufferViews": [
            { "buffer": 0, "byteOffset": vertices_byte_length, "byteLength": indices_byte_length, "target": 34963 },
            { "buffer": 0, "byteOffset": 0, "byteLength": vertices_byte_length, "target": 34962 }
        ],
        "accessors": [
            { "bufferView": 0, "byteOffset": 0, "componentType": 5123, "count": req.indices.len(), "type": "SCALAR" },
            { "bufferView": 1, "byteOffset": 0, "componentType": 5126, "count": req.vertices.len(), "type": "VEC3", 
              "min": [-10.0, -10.0, -10.0], "max": [10.0, 10.0, 10.0] }
        ]
    });

    GLTFExportResponse {
        gltf_json: gltf.to_string(),
    }
}
