#!/bin/bash
cd hack-analyzer
go mod init github.com/haghfizzuddin/hack-analyzer 2>/dev/null
go get github.com/fatih/color github.com/spf13/cobra
go build -o ../hack-analyzer cmd/analyzer/main.go
cd ..
chmod +x hack-analyzer
echo "✅ Build complete! Run: ./hack-analyzer --help"
