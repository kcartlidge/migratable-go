# Linux script to build for supported platforms.
# AGPL license. Copyright 2024 K Cartlidge.

echo Building Windows edition
env GOOS=windows GOARCH=amd64 go build -o ../builds/windows/migratable.exe

echo Building Mac edition - Intel X64
env GOOS=darwin GOARCH=amd64 go build -o ../builds/macos-x64/migratable

echo Building Mac edition - Apple Silicon M1 ARM64
env GOOS=darwin GOARCH=arm64 go build -o ../builds/macos/migratable

echo Building Linux edition
env GOOS=linux GOARCH=amd64 go build -o ../builds/linux/migratable
