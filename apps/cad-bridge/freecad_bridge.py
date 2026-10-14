#!/usr/bin/env python3
"""FreeCAD Sheet Metal to SMA CAD Bridge"""
import json
import base64
import os

import FreeCAD as App
import FreeCADGui as Gui
import json
import os
import tempfile
import importDXF

# WebSocket-client kütüphanesini gerektirir: pip install websocket-client
try:
    import websocket
except ImportError:
    App.Console.PrintError("websocket-client modülü eksik! Lütfen 'pip install websocket-client' komutunu çalıştırın.\n")

def extract_sheet_metal_data():
    """FreeCAD aktif belgesinden sac metal verilerini çıkarır"""
    doc = App.ActiveDocument
    if not doc:
        App.Console.PrintError("Aktif bir FreeCAD belgesi bulunamadı!\n")
        return None

    thickness = 2.0
    material = "Unknown Steel"
    bends = []
    mass_kg = 0.0
    surface_area_m2 = 0.0
    flat_pattern_obj = None

    # Parçaları Tara
    for obj in doc.Objects:
        # Sheet Metal eklentisindeki ana parça kontrolü
        if hasattr(obj, "Thickness"):
            thickness = obj.Thickness.Value
            
        # Büküm (Bend) kontrolü
        if obj.TypeId == "PartDesign::Bend" or "Bend" in obj.Name:
            angle = 90.0
            radius = 2.0
            if hasattr(obj, "angle"):
                angle = obj.angle.Value
            if hasattr(obj, "radius"):
                radius = obj.radius.Value
            bends.append({"angle": angle, "radius": radius, "direction": "UP"})
            
        # Kütle/Hacim analizi (Density = 7.85 kg/dm^3 for Steel as default)
        if hasattr(obj, "Shape") and obj.Shape.Volume > 0:
            volume_mm3 = obj.Shape.Volume
            mass_kg = (volume_mm3 / 1e6) * 7.85 # Hacimden kiloya geçiş
            surface_area_m2 = obj.Shape.Area / 1e6 # mm^2 to m^2
            
        # Açınım (Unfold/Flat Pattern) nesnesini bul
        if "Unfold" in obj.Name or "FlatPattern" in obj.Name:
            flat_pattern_obj = obj

    if not flat_pattern_obj:
        App.Console.PrintError("Açınım (Unfold) nesnesi bulunamadı. Lütfen SheetMetal eklentisi ile parçayı açın.\n")
        return None

    # Bounding Box (Açınım Ölçüleri)
    bbox = flat_pattern_obj.Shape.BoundBox
    length_x = bbox.XLength
    width_y = bbox.YLength

    # DXF Export
    temp_dir = tempfile.gettempdir()
    dxf_path = os.path.join(temp_dir, f"SMA_FreeCAD_{doc.Name}.dxf")
    importDXF.export([flat_pattern_obj], dxf_path)
    
    App.Console.PrintMessage(f"DXF başarıyla çıkarıldı: {dxf_path}\n")

    return {
        "action": "upload_cad",
        "partName": doc.Name,
        "material": material,
        "thickness": thickness,
        "bounding_box": {
            "x": length_x,
            "y": width_y
        },
        "bends": bends,
        "mass_kg": mass_kg,
        "surface_area_m2": surface_area_m2,
        "dxf_path": dxf_path.replace("\\", "\\\\"),
        "cad_source": "FreeCAD"
    }

def send_to_sma_api(payload: dict):
    """Çıkarılan veriyi SMA Go API'ye WebSocket üzerinden gönderir"""
    ws_url = os.getenv("SMA_API_URL", "ws://localhost:3000/ws")
    
    if "websocket" not in globals():
        return
        
    try:
        ws = websocket.WebSocket()
        ws.connect(ws_url)
        ws.send(json.dumps(payload))
        ws.close()
        App.Console.PrintMessage("Veri başarıyla Sheet Metal Automation sistemine gönderildi!\n")
    except Exception as e:
        App.Console.PrintError(f"WebSocket bağlantı hatası: {str(e)}\n")

def run():
    payload = extract_sheet_metal_data()
    if payload:
        send_to_sma_api(payload)

if __name__ == "__main__":
    run()
