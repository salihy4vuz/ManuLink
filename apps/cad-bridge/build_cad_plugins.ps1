<#
.SYNOPSIS
    SolidWorks ve Autodesk Inventor eklentilerini (Add-ins) otomatik derleyen otomasyon betiği.
.DESCRIPTION
    Bu script, çalıştırıldığı mühendislik bilgisayarında:
    1. Yüklü olan SolidWorks ve Inventor programlarının DLL yollarını Kayıt Defteri'nden (Registry) otomatik bulur.
    2. .csproj dosyalarındaki DLL referans yollarını (HintPath) dinamik olarak günceller.
    3. Bilgisayardaki MSBuild veya dotnet araçlarını bularak C# projelerini otomatik derler (Build).
    4. Üretilen Eklenti (.dll) dosyalarını 'dist' klasörüne kopyalar.
#>

param (
    [switch]$BuildSolidWorks = $true,
    [switch]$BuildInventor = $true
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$distDir = Join-Path $scriptDir "dist"
if (!(Test-Path $distDir)) { New-Item -ItemType Directory -Path $distDir | Out-Null }

Write-Host "=============================================" -ForegroundColor Cyan
Write-Host " SHEET METAL AUTOMATION - CAD BUILDER SCRIPT " -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan

# 1. MSBuild veya dotnet yolunu bulma
$msbuildPath = ""
if (Get-Command "dotnet" -ErrorAction SilentlyContinue) {
    $msbuildPath = "dotnet"
    Write-Host "[OK] 'dotnet' CLI bulundu." -ForegroundColor Green
} else {
    $vswherePath = "${env:ProgramFiles(x86)}\Microsoft Visual Studio\Installer\vswhere.exe"
    if (Test-Path $vswherePath) {
        $msbuildPath = & $vswherePath -latest -requires Microsoft.Component.MSBuild -find MSBuild\Current\Bin\MSBuild.exe
        if ($msbuildPath) {
            Write-Host "[OK] MSBuild bulundu: $msbuildPath" -ForegroundColor Green
        }
    }
}

if (-not $msbuildPath) {
    Write-Host "[HATA] Sistemde 'dotnet' veya 'Visual Studio MSBuild' bulunamadı." -ForegroundColor Red
    Write-Host "Lütfen .NET SDK veya Visual Studio yükleyin." -ForegroundColor Red
    exit 1
}

# 2. SolidWorks Otomatik Derleme
if ($BuildSolidWorks) {
    Write-Host "`n--- SolidWorks Eklentisi Derleniyor ---" -ForegroundColor Yellow
    $swPath = ""
    # Registry üzerinden SolidWorks arama
    $swRegPath = "HKLM:\SOFTWARE\SolidWorks\SolidWorks*"
    if (Test-Path $swRegPath) {
        $swKeys = Get-ChildItem -Path "HKLM:\SOFTWARE\SolidWorks" | Where-Object { $_.Name -match "SolidWorks \d+" }
        if ($swKeys) {
            $swPath = (Get-ItemProperty $swKeys[0].PSPath)."SolidWorks Executable"
            if ($swPath) { $swPath = Split-Path $swPath -Parent }
        }
    }
    
    if (-not $swPath) {
        # Varsayılan yolları kontrol et
        $defaultSwPaths = @("C:\Program Files\SolidWorks Corp\SolidWorks", "C:\Program Files\SOLIDWORKS Corp\SOLIDWORKS")
        foreach ($p in $defaultSwPaths) {
            if (Test-Path $p) { $swPath = $p; break }
        }
    }

    if ($swPath) {
        $dllPath = Join-Path $swPath "api\redist\SolidWorks.Interop.sldworks.dll"
        if (Test-Path $dllPath) {
            Write-Host "[OK] SolidWorks DLL bulundu: $dllPath" -ForegroundColor Green
            $csprojPath = Join-Path $scriptDir "SolidWorksBridge.csproj"
            
            # csproj HintPath guncelleme (Regex ile)
            $content = Get-Content $csprojPath -Raw
            $content = $content -replace "(<HintPath>).*?(</HintPath>)", "`$1$dllPath`$2"
            Set-Content -Path $csprojPath -Value $content

            # Derleme
            Write-Host "Derleniyor..." -ForegroundColor Cyan
            if ($msbuildPath -eq "dotnet") {
                & dotnet build $csprojPath -c Release
            } else {
                & $msbuildPath $csprojPath /p:Configuration=Release
            }
            
            # Kopyalama
            $outputDll = Join-Path $scriptDir "bin\Release\net48\SolidWorksAddin.dll"
            if (Test-Path $outputDll) {
                Copy-Item $outputDll $distDir -Force
                Write-Host "[BAŞARILI] SolidWorks eklentisi dist/ klasörüne kopyalandı!" -ForegroundColor Green
            }
        } else {
            Write-Host "[HATA] SolidWorks yüklü ama 'SolidWorks.Interop.sldworks.dll' bulunamadı." -ForegroundColor Red
        }
    } else {
        Write-Host "[HATA] Bu bilgisayarda SolidWorks kurulumu bulunamadı." -ForegroundColor Red
    }
}

# 3. Autodesk Inventor Otomatik Derleme
if ($BuildInventor) {
    Write-Host "`n--- Autodesk Inventor Eklentisi Derleniyor ---" -ForegroundColor Yellow
    $invPath = ""
    # Registry üzerinden Inventor arama
    $invRegPath = "HKLM:\SOFTWARE\Autodesk\Inventor\RegistryVersion*"
    if (Test-Path $invRegPath) {
        $invKeys = Get-ChildItem -Path "HKLM:\SOFTWARE\Autodesk\Inventor" | Where-Object { $_.Name -match "RegistryVersion" }
        if ($invKeys) {
            $invPath = (Get-ItemProperty $invKeys[0].PSPath).InstallLocation
        }
    }

    if (-not $invPath) {
        # Varsayılan yolları kontrol et
        $defaultInvPaths = @("C:\Program Files\Autodesk\Inventor 2026", "C:\Program Files\Autodesk\Inventor 2025", "C:\Program Files\Autodesk\Inventor 2024")
        foreach ($p in $defaultInvPaths) {
            if (Test-Path $p) { $invPath = $p; break }
        }
    }

    if ($invPath) {
        $dllPath = Join-Path $invPath "Bin\Public Assemblies\autodesk.inventor.interop.dll"
        if (Test-Path $dllPath) {
            Write-Host "[OK] Inventor DLL bulundu: $dllPath" -ForegroundColor Green
            $csprojPath = Join-Path $scriptDir "InventorBridge.csproj"
            
            $content = Get-Content $csprojPath -Raw
            $content = $content -replace "(<HintPath>).*?(</HintPath>)", "`$1$dllPath`$2"
            Set-Content -Path $csprojPath -Value $content

            Write-Host "Derleniyor..." -ForegroundColor Cyan
            if ($msbuildPath -eq "dotnet") {
                & dotnet build $csprojPath -c Release
            } else {
                & $msbuildPath $csprojPath /p:Configuration=Release
            }
            
            $outputDll = Join-Path $scriptDir "bin\Release\net48\InventorAddin.dll"
            if (Test-Path $outputDll) {
                Copy-Item $outputDll $distDir -Force
                Write-Host "[BAŞARILI] Inventor eklentisi dist/ klasörüne kopyalandı!" -ForegroundColor Green
            }
        } else {
            Write-Host "[HATA] Inventor yüklü ama 'autodesk.inventor.interop.dll' bulunamadı." -ForegroundColor Red
        }
    } else {
        Write-Host "[HATA] Bu bilgisayarda Autodesk Inventor kurulumu bulunamadı." -ForegroundColor Red
    }
}

Write-Host "`nİşlem Tamamlandı." -ForegroundColor Cyan
