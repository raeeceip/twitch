package twitchstreamer

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/kbinani/screenshot"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format/rtmp"
	"github.com/nicklaw5/helix"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 1. Screen Capture Module
type ScreenCapture struct {
	frames chan *image.RGBA
}

func (sc *ScreenCapture) Capture() {
	for {
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			log.Println("Failed to capture screen:", err)
			continue
		}
		sc.frames <- img
		time.Sleep(time.Second / 30) // 30 FPS
	}
}

// 2. Audio Capture Module
type AudioCapture struct {
	samples chan []float32
}

func (ac *AudioCapture) Capture() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, 4096, ac.handleAudio)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	stream.Start()
	select {} // Run forever
}

func (ac *AudioCapture) handleAudio(in []float32) {
	ac.samples <- in
}

// 3 & 4. Video and Audio Encoder
type Encoder struct {
	videoTrack *webrtc.TrackLocalStaticSample
	audioTrack *webrtc.TrackLocalStaticSample
}

func (e *Encoder) EncodeVideo(frame *image.RGBA) {
	// Simplified: In reality, you'd need to implement actual H.264 encoding
	e.videoTrack.WriteSample(media.Sample{Data: frame.Pix, Duration: time.Second / 30})
}

func (e *Encoder) EncodeAudio(samples []float32) {
	// Simplified: In reality, you'd need to implement actual AAC encoding
	e.audioTrack.WriteSample(media.Sample{Data: []byte(fmt.Sprintf("%v", samples)), Duration: time.Second / 100})
}

// 5. RTMP Packager
type RTMPPackager struct {
	url string
}

func (rp *RTMPPackager) Stream(videoTrack, audioTrack *webrtc.TrackLocalStaticSample) {
	rtmpConn, err := rtmp.Dial(rp.url)
	if err != nil {
		log.Fatal(err)
	}
	defer rtmpConn.Close()

	// Simplified: In reality, you'd need to properly mux video and audio
	avutil.CopyFile(rtmpConn, rtmpConn)
}

// 6. Configuration Manager
type Config struct {
	TwitchStreamKey string
	FrameRate       int
	Resolution      string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	return &Config{
		TwitchStreamKey: viper.GetString("twitch_stream_key"),
		FrameRate:       viper.GetInt("frame_rate"),
		Resolution:      viper.GetString("resolution"),
	}
}

// 7. User Interface (simplified console UI for this example)
type UI struct {
	controller *StreamController
}

func (ui *UI) Run() {
	fmt.Println("1. Start Stream")
	fmt.Println("2. Stop Stream")
	fmt.Println("3. Exit")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		ui.controller.StartStream()
	case 2:
		ui.controller.StopStream()
	case 3:
		ui.controller.Exit()
	}
}

// 8. Stream Controller
type StreamController struct {
	isStreaming bool
	screenCap   *ScreenCapture
	audioCap    *AudioCapture
	encoder     *Encoder
	packager    *RTMPPackager
}

func (sc *StreamController) StartStream() {
	sc.isStreaming = true
	go sc.screenCap.Capture()
	go sc.audioCap.Capture()
	// Additional logic to start encoding and packaging
}

func (sc *StreamController) StopStream() {
	sc.isStreaming = false
	// Additional logic to stop all processes
}

func (sc *StreamController) Exit() {
	sc.StopStream()
	// Additional cleanup logic
}

// 9. Twitch API Client
type TwitchClient struct {
	client *helix.Client
}

func NewTwitchClient(clientID, clientSecret string) (*TwitchClient, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return nil, err
	}
	return &TwitchClient{client: client}, nil
}

func (tc *TwitchClient) GetStreamKey(userID string) (string, error) {
	resp, err := tc.client.GetStreamKey(&helix.StreamKeyParams{
		BroadcasterID: userID,
	})
	if err != nil {
		return "", err
	}
	return resp.Data.StreamKey, nil
}

// 10. Error Handler and Logging System
var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Cannot initialize logger")
	}
}

func LogInfo(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func LogError(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Main function to tie everything together
func RunTwitchStreamer() {
	config := LoadConfig()

	screenCap := &ScreenCapture{frames: make(chan *image.RGBA, 100)}
	audioCap := &AudioCapture{samples: make(chan []float32, 100)}
	encoder := &Encoder{} // Initialize with proper tracks
	packager := &RTMPPackager{url: "rtmp://live.twitch.tv/app/" + config.TwitchStreamKey}

	controller := &StreamController{
		screenCap: screenCap,
		audioCap:  audioCap,
		encoder:   encoder,
		packager:  packager,
	}

	ui := &UI{controller: controller}

	LogInfo("Twitch Streamer initialized")

	ui.Run()
}
