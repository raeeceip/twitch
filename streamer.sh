#!/bin/bash

# Create main project directory
mkdir -p twitch-streamer
cd twitch-streamer

# Create subdirectories
mkdir -p cmd/twitch-streamer
mkdir -p internal/{capture,encode,rtmp,config,ui,api,logger}

# Create main application file
touch cmd/twitch-streamer/main.go

# Create internal package files
touch internal/capture/{screen.go,audio.go}
touch internal/encode/{video.go,audio.go}
touch internal/rtmp/packager.go
touch internal/config/manager.go
touch internal/ui/controller.go
touch internal/api/twitch.go
touch internal/logger/logger.go

# Create go.mod file
go mod init github.com/yourusername/twitch-streamer

# Create README.md
touch README.md

echo "Project structure created successfully!"