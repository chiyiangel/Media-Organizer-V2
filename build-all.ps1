# Photo Video Organizer - Cross-platform Build Script

Write-Host "=== Building Photo Video Organizer for multiple platforms ===" -ForegroundColor Green

# Create build directory
if (!(Test-Path "build")) {
    New-Item -ItemType Directory -Name "build"
}

# Define build configurations
$platforms = @(
    @{OS="windows"; Arch="amd64"; Ext=".exe"; Icon="[WIN]"}
    @{OS="darwin"; Arch="amd64"; Ext=""; Icon="[MAC]"}
    @{OS="darwin"; Arch="arm64"; Ext=""; Icon="[MAC]"}
    @{OS="linux"; Arch="amd64"; Ext=""; Icon="[LNX]"}
    @{OS="linux"; Arch="arm64"; Ext=""; Icon="[LNX]"}
)

# Build for each platform
foreach ($platform in $platforms) {
    $filename = "media-organizer-$($platform.OS)-$($platform.Arch)$($platform.Ext)"
    
    Write-Host "  $($platform.Icon) Building for $($platform.OS) ($($platform.Arch))..." -ForegroundColor Cyan
    
    $env:GOOS = $platform.OS
    $env:GOARCH = $platform.Arch
    
    go build -ldflags "-s -w" -o "build/$filename" "./cmd/organizer"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    [OK] Success: $filename" -ForegroundColor Green
    } else {
        Write-Host "    [FAIL] Failed: $filename" -ForegroundColor Red
    }
}

# Reset environment variables
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host "`n[DONE] Cross-platform builds complete!" -ForegroundColor Green

# Show built files
Write-Host "`nBuilt files:" -ForegroundColor Yellow
Get-ChildItem "build/media-organizer-*" | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  $($_.Name) ($size MB)" -ForegroundColor White
}
