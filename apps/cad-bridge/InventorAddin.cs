using System;
using System.Runtime.InteropServices;
using Inventor;
using System.Diagnostics;
using System.Threading.Tasks;

namespace SheetMetalApex.Bridge
{
    [ComVisible(true)]
    [Guid("B2C3D4E5-F6A7-8901-2345-67890ABCDEF1")]
    public class InventorAddin : ApplicationAddInServer
    {
        private Application invApp;
        private string clientId;
        private WebSocketClient wsClient;

        public void Activate(ApplicationAddInSite addInSiteObject, bool firstTime)
        {
            invApp = addInSiteObject.Application;
            clientId = addInSiteObject.Parent.ClientId;

            // Connect to Golang Gateway
            wsClient = new WebSocketClient("ws://localhost:3000/ws");
            _ = ConnectWebSocketAsync();
        }

        private async Task ConnectWebSocketAsync()
        {
            try {
                await wsClient.ConnectAsync();
            } catch (Exception ex) {
                Debug.WriteLine("WebSocket connection failed: " + ex.Message);
            }
        }

        public void Deactivate()
        {
            if (wsClient != null) {
                wsClient.DisconnectAsync().Wait();
            }
            Marshal.ReleaseComObject(invApp);
            invApp = null;
        }

        public void ExecuteCommand(int commandID)
        {
            // Placeholder for command execution
        }

        public object Automation => null;

        // Called when the user clicks "Send to B2B Quote" button in Inventor
        public void ExportActiveSheetMetalPart()
        {
            Document invDoc = invApp.ActiveDocument;
            if (invDoc == null || invDoc.DocumentType != DocumentTypeEnum.kPartDocumentObject) {
                // In a real add-in, you would use MessageBox or invApp.StatusBarText
                Debug.WriteLine("Sadece sac metal (Part) dosyaları gönderilebilir.");
                return;
            }

            PartDocument invPart = (PartDocument)invDoc;

            // Check if it's a Sheet Metal Component
            if (!(invPart.ComponentDefinition is SheetMetalComponentDefinition))
            {
                Debug.WriteLine("Bu parça sac metal değil!");
                return;
            }

            SheetMetalComponentDefinition smDef = (SheetMetalComponentDefinition)invPart.ComponentDefinition;

            // 1. Get Material
            string material = "Unknown Material";
            try {
                material = invPart.ComponentDefinition.Material.Name;
            } catch { }

            // 2. Extract Thickness
            double thickness = 0.0;
            try {
                // Inventor uses cm internally. Convert to mm.
                thickness = (double)smDef.Thickness.Value * 10.0; 
            } catch {
                thickness = 2.0; // Fallback
            }

            // 3. Ensure Flat Pattern exists
            if (!smDef.HasFlatPattern)
            {
                try {
                    smDef.Unfold();
                } catch {
                    Debug.WriteLine("Açınım (Flat Pattern) oluşturulamadı.");
                    return;
                }
            }

            FlatPattern flatPattern = smDef.FlatPattern;

            // 4. Bounding Box (Extents)
            double lengthX = flatPattern.Length * 10.0; // cm to mm
            double widthY = flatPattern.Width * 10.0; // cm to mm

            // 5. Extract Bend Information
            System.Collections.Generic.List<string> bendDetailsList = new System.Collections.Generic.List<string>();
            foreach (Edge oEdge in flatPattern.BendLinesBody.Edges)
            {
                try {
                    object isBendUpObj;
                    object bendAngleObj;
                    flatPattern.GetBendInfo(oEdge, out isBendUpObj, out bendAngleObj);
                    
                    bool isBendUp = (bool)isBendUpObj;
                    double bendAngleRad = (double)bendAngleObj;
                    double bendAngleDeg = bendAngleRad * (180.0 / Math.PI);
                    
                    string direction = isBendUp ? "UP" : "DOWN";
                    
                    // We assume standard radius based on thickness if API doesn't provide it directly here
                    // In deeper API calls, BendFeatures could be iterated directly.
                    double bendRadius = thickness; 

                    bendDetailsList.Add($@"{{ ""angle"": {bendAngleDeg.ToString(System.Globalization.CultureInfo.InvariantCulture)}, ""radius"": {bendRadius.ToString(System.Globalization.CultureInfo.InvariantCulture)}, ""direction"": ""{direction}"" }}");
                } catch (Exception ex) {
                    Debug.WriteLine("Bend extract error: " + ex.Message);
                }
            }
            
            string bendDetailsJson = "[" + string.Join(",", bendDetailsList) + "]";

            // 6. Mass Properties Extraction
            double massKg = 0.0;
            double surfaceAreaM2 = 0.0;
            try {
                MassProperties massProps = invPart.ComponentDefinition.MassProperties;
                massKg = massProps.Mass; // Mass in kg
                surfaceAreaM2 = massProps.Area / 10000.0; // Area from cm^2 to m^2
            } catch { }

            // 7. Export Flat Pattern to DXF
            string tempDxfPath = System.IO.Path.Combine(System.IO.Path.GetTempPath(), "SheetMetalApex_Inventor_FlatPattern.dxf");
            try {
                // Format parameter string for DXF export in Inventor
                string sOut = "FLAT PATTERN DXF?AcadVersion=2018";
                flatPattern.DataIO.WriteDataToFile(sOut, tempDxfPath);
            } catch (Exception ex) {
                Debug.WriteLine("DXF Export failed: " + ex.Message);
            }

            // 8. Build JSON Payload and Send via WebSocket
            string jsonPayload = $@"
            {{
                ""action"": ""upload_cad"",
                ""partName"": ""{invPart.DisplayName}"",
                ""material"": ""{material}"",
                ""thickness"": {thickness.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                ""bounding_box"": {{
                    ""x"": {lengthX.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                    ""y"": {widthY.ToString(System.Globalization.CultureInfo.InvariantCulture)}
                }},
                ""bends"": {bendDetailsJson},
                ""mass_kg"": {massKg.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                ""surface_area_m2"": {surfaceAreaM2.ToString(System.Globalization.CultureInfo.InvariantCulture)},
                ""dxf_path"": ""{tempDxfPath.Replace("\\", "\\\\")}""
            }}";

            wsClient.Send(jsonPayload);
        }
    }
}
