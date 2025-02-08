package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os/exec"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandomString(n int) string {
	var str string

	for i := 0; i < n; i++ {
		str += string(alphabet[rand.Intn(len(alphabet))])
	}
	return str
}

func getVideoAspectRatio(path string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", path)
	var buff bytes.Buffer
	cmd.Stdout = &buff
	cmd.Run()
	result := struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}{}
	err := json.Unmarshal(buff.Bytes(), &result)
	if err != nil {
		return "", err
	}
	ratio := float64(result.Streams[0].Width) / float64(result.Streams[0].Height)
	if ratio > 1.75 && ratio < 1.8 {
		return "landscape", nil
	}
	if ratio > 0.55 && ratio < 0.57 {
		return "portrait", nil
	}
	return "other", nil
}

func processVideoForFastStart(path string) (string, error) {
	outputPath := path + ".processing"

	cmd := exec.Command("ffmpeg", "-i", path, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return outputPath, nil
}
