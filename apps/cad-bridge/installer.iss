; Inno Setup Script for Sheet Metal Automation CAD Add-ins
; This script packages the SolidWorks, Inventor, and SolidEdge add-ins.

[Setup]
AppName=Sheet Metal Automation CAD Bridge
AppVersion=1.0.0
DefaultDirName={pf}\SheetMetalAutomation
DefaultGroupName=Sheet Metal Automation
UninstallDisplayIcon={app}\cad_bridge_logo.ico
Compression=lzma2
SolidCompression=yes
OutputDir=Output
OutputBaseFilename=SheetMetalAutomation_CAD_Setup

[Files]
; Assuming the DLLs are built in bin/Release
Source: "bin\Release\SolidWorksBridge.dll"; DestDir: "{app}"; Flags: ignoreversion regserver
Source: "bin\Release\SolidEdgeAddin.dll"; DestDir: "{app}"; Flags: ignoreversion regserver
Source: "bin\Release\AutodeskInventorAddin.dll"; DestDir: "{app}"; Flags: ignoreversion regserver
Source: "bin\Release\Newtonsoft.Json.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\Release\WebSocket4Net.dll"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
; SolidWorks Add-in Registration
Root: HKLM; Subkey: "SOFTWARE\SolidWorks\Addins\{{YOUR-GUID-HERE}}"; ValueType: string; ValueName: "Title"; ValueData: "Sheet Metal Automation"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\SolidWorks\Addins\{{YOUR-GUID-HERE}}"; ValueType: string; ValueName: "Description"; ValueData: "B2B SaaS CAD Bridge for Quote Generation"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\SolidWorks\Addins\{{YOUR-GUID-HERE}}"; ValueType: dword; ValueName: "Default"; ValueData: "1"; Flags: uninsdeletekey

; SolidEdge Add-in Registration
Root: HKLM; Subkey: "SOFTWARE\Unigraphics Solutions\Solid Edge\Version 222\AddIns\{{YOUR-SE-GUID-HERE}}"; ValueType: string; ValueName: "Summary"; ValueData: "Sheet Metal Automation"; Flags: uninsdeletekey

; Inventor Add-in Registration (Usually done via .addin file, so we copy it to the correct AppData folder)
[Dirs]
Name: "{userappdata}\Autodesk\ApplicationPlugins\SheetMetalAutomation.bundle\Contents"

[Files]
Source: "AutodeskInventorAddin.addin"; DestDir: "{userappdata}\Autodesk\ApplicationPlugins\SheetMetalAutomation.bundle"; Flags: ignoreversion

[Run]
; Re-register COM objects explicitly using RegAsm if needed
Filename: "{dotnet4064}\regasm.exe"; Parameters: "/codebase ""{app}\SolidWorksBridge.dll"""; StatusMsg: "Registering SolidWorks Addin..."; Flags: runhidden
Filename: "{dotnet4064}\regasm.exe"; Parameters: "/codebase ""{app}\SolidEdgeAddin.dll"""; StatusMsg: "Registering SolidEdge Addin..."; Flags: runhidden

[UninstallRun]
Filename: "{dotnet4064}\regasm.exe"; Parameters: "/unregister ""{app}\SolidWorksBridge.dll"""; Flags: runhidden
Filename: "{dotnet4064}\regasm.exe"; Parameters: "/unregister ""{app}\SolidEdgeAddin.dll"""; Flags: runhidden
