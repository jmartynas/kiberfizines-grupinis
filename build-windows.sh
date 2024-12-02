#!/bin/bash
echo "Building for Windows 32-bit..."
GOOS=windows GOARCH=386 go build -o windows_386.exe

echo "Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 go build -o windows_amd64.exe