@echo off
echo === Building Photo Video Organizer for multiple platforms ===

if not exist build mkdir build

echo   [WIN] Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o build/photo-organizer-windows-amd64.exe cmd/organizer/main.go
if %errorlevel% == 0 echo     [OK] Success: photo-organizer-windows-amd64.exe

echo   [MAC] Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w" -o build/photo-organizer-darwin-amd64 cmd/organizer/main.go
if %errorlevel% == 0 echo     [OK] Success: photo-organizer-darwin-amd64

echo   [MAC] Building for macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-s -w" -o build/photo-organizer-darwin-arm64 cmd/organizer/main.go
if %errorlevel% == 0 echo     [OK] Success: photo-organizer-darwin-arm64

echo   [LNX] Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o build/photo-organizer-linux-amd64 cmd/organizer/main.go
if %errorlevel% == 0 echo     [OK] Success: photo-organizer-linux-amd64

echo   [LNX] Building for Linux (arm64)...
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-s -w" -o build/photo-organizer-linux-arm64 cmd/organizer/main.go
if %errorlevel% == 0 echo     [OK] Success: photo-organizer-linux-arm64

echo.
echo [DONE] Cross-platform builds complete!
echo.
echo Built files:
dir build\photo-organizer-* /b 2>nul