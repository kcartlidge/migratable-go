@echo off
@echo Building Linux edition
@set GOOS=linux
@set GOARCH=amd64
@go build -o ../builds/linux/migratable

@echo Building Mac edition - Intel X64
@set GOOS=darwin
@set GOARCH=amd64
@go build -o ../builds/macos-x64/migratable

@echo Building Mac edition - Apple Silicon M1 ARM64
@set GOOS=darwin
@set GOARCH=arm64
@go build -o ../builds/macos/migratable

@echo Building Windows edition
@set GOOS=windows
@set GOARCH=amd64
@go build -o ../builds/windows/migratable.exe
