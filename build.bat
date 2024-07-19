@echo off

REM Set the desired architecture
SET ARCH=amd64

REM Run go generate
go generate ./...

REM Build for the specified architecture
SET GOOS=windows
SET GOARCH=%ARCH%
go build -o build\windows_%ARCH%\moba-converter-go.exe

echo Build complete for Windows (%ARCH%)