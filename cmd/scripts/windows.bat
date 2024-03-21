@rem Windows script to build for supported platforms.
@rem AGPL license. Copyright 2024 K Cartlidge.

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

@echo.
@echo Windows does not use executable file attributes:
@echo * Setting builds to executable in git
@git update-index --chmod=+x ../builds/linux/migratable
@git update-index --chmod=+x ../builds/macos-x64/migratable
@git update-index --chmod=+x ../builds/macos/migratable
@git update-index --chmod=+x ../builds/windows/migratable.exe

@echo * Ensuring builds scripts remain executable in git
@git update-index --chmod=+x ./scripts/linux.sh
@git update-index --chmod=+x ./scripts/mac.sh
@git update-index --chmod=+x ./scripts/windows.bat
