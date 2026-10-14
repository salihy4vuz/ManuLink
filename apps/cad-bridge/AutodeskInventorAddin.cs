using System;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
// Note: In a real project, reference "autodesk.inventor.interop.dll"
// using Inventor;

namespace SheetMetalApex.AutodeskBridge
{
    /// <summary>
    /// Faz 9.9: Autodesk Inventor Enterprise CAD Bridge
    /// This module hooks into Autodesk Inventor's COM API to extract Sheet Metal Flat Patterns, 
    /// Mass properties, Bending sequences, and ERP metadata, and sends it to the Go WebSocket Server.
    /// </summary>
    public class AutodeskInventorAddin
    {
        // WebSocket connection to our Go API
        private static readonly string wsUrl = "ws://localhost:3000/ws";

        public void ExportActiveSheetMetalPart()
        {
            try
            {
                // 1. Connect to Autodesk Inventor Application (Simulated COM Interop)
                // Inventor.Application invApp = (Inventor.Application)System.Runtime.InteropServices.Marshal.GetActiveObject("Inventor.Application");
                // PartDocument invDoc = (PartDocument)invApp.ActiveDocument;
                // SheetMetalComponentDefinition smDef = (SheetMetalComponentDefinition)invDoc.ComponentDefinition;
                
                bool isSheetMetal = true; // smDef.HasFlatPattern
                if (!isSheetMetal) 
                {
                    Console.WriteLine("Uyarı: Bu parça bir Sac Levha (Sheet Metal) değil!");
                    return;
                }

                // 2. Extract Basic Properties (Material, Thickness)
                string material = "Aluminum 6061"; // invDoc.PropertySets["Design Tracking Properties"]["Material"].Value
                double thicknessMm = 3.0; // (double)smDef.Thickness.Value * 10.0; // cm to mm

                // 3. ERP Metadata Extraction (Custom Properties)
                string partNo = "INV-PRT-001"; // invDoc.PropertySets["Design Tracking Properties"]["Part Number"].Value
                string revision = "B";         // invDoc.PropertySets["Summary Information"]["Revision Number"].Value
                string description = "Autodesk Lazer Kesim Braket";

                // 4. Extract Mass Properties & Bounding Box
                // MassProperties massProps = smDef.MassProperties;
                double massKg = 2.15; // massProps.Mass
                double surfaceAreaM2 = 0.45; // massProps.Area / 10000.0
                double lengthX = 500.0; // FlatPattern Extents
                double widthY = 250.0;
                
                // 5. Cut-List / Laser Metrics
                double cutLengthMm = 1850.0; // Extracted from FlatPattern edges
                int pierceCount = 6;

                // 6. Secondary Operations (Hole Wizard / Tapped Holes)
                int countersinkCount = 2; // Iterated through HoleFeatures
                int tappedHoleCount = 4;

                // 7. Bending Sequence (Bend Table)
                string bendDetailsJson = @"[
                    { ""angle"": 90.0, ""radius"": 3.0, ""direction"": ""UP"" },
                    { ""angle"": 45.0, ""radius"": 3.0, ""direction"": ""DOWN"" }
                ]";

                // 8. Generate Flat Pattern DXF (DataIO)
                // In Autodesk, DXF export of flat pattern is done via DataIO string formatting:
                // string sOut = "FLAT PATTERN DXF?AcadVersion=2012&OuterProfileLayer=Outer";
                // smDef.FlatPattern.DataIO.WriteDataToFile(sOut, tempDxfPath);
                string dxfBase64 = "QXV0b2Rlc2sgSW52ZW50b3IgRFhGIE1vY2sgRGF0YQ=="; // Mock base64 of generated DXF

                // 9. Generate STEP and Thumbnail
                // invDoc.SaveAs(tempStepPath, true);
                string stepBase64 = "QXV0b2Rlc2sgU1RFUCBNb2NrIERhdGE=";
                string thumbBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="; // 1x1 mock JPEG

                // 10. Construct Enterprise CAD Payload for Go API
                string payload = $@"{{
                    ""action"": ""cad_export"",
                    ""data"": {{
                        ""part_number"": ""{partNo}"",
                        ""revision"": ""{revision}"",
                        ""description"": ""{description}"",
                        ""material"": ""{material}"",
                        ""thickness_mm"": {thicknessMm},
                        ""flat_area_mm2"": {lengthX * widthY},
                        ""bounding_box_x"": {lengthX},
                        ""bounding_box_y"": {widthY},
                        ""mass_kg"": {massKg.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                        ""surface_area_m2"": {surfaceAreaM2.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                        ""cut_length_mm"": {cutLengthMm.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                        ""pierce_count"": {pierceCount},
                        ""countersink_count"": {countersinkCount},
                        ""tapped_hole_count"": {tappedHoleCount},
                        ""bend_count"": 2,
                        ""source"": ""Autodesk_Inventor_Addin"",
                        ""dxf_base64"": ""{dxfBase64}"",
                        ""step_base64"": ""{stepBase64}"",
                        ""thumbnail_base64"": ""{thumbBase64}"",
                        ""bend_details"": {bendDetailsJson}
                    }}
                }}";

                // 11. Send to WebSocket Server
                Console.WriteLine("Autodesk Inventor'dan Go sunucusuna veri aktarılıyor...");
                SendToBackendAsync(payload).Wait();
                Console.WriteLine("Aktarım Başarılı!");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Autodesk Export Hatası: {ex.Message}");
            }
        }

        private async Task SendToBackendAsync(string payload)
        {
            using (ClientWebSocket wsClient = new ClientWebSocket())
            {
                Uri serverUri = new Uri(wsUrl);
                await wsClient.ConnectAsync(serverUri, CancellationToken.None);
                
                byte[] bytes = Encoding.UTF8.GetBytes(payload);
                await wsClient.SendAsync(new ArraySegment<byte>(bytes), WebSocketMessageType.Text, true, CancellationToken.None);
                
                await wsClient.CloseAsync(WebSocketCloseStatus.NormalClosure, "Exported Autodesk Data", CancellationToken.None);
            }
        }
    }
}
