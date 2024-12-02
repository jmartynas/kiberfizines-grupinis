@echo off
echo Building for Windows 32-bit...
set GOOS=windows
set GOARCH=386
go build -o windows_386.exe

echo Building for Windows 64-bit...
set GOOS=windows
set GOARCH=amd64
go build -o windows_amd64.exe