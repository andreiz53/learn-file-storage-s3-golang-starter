package main

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
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

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	client := s3.NewPresignClient(s3Client)
	request, err := client.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}
	return request.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return database.Video{}, nil
	}
	url := strings.Split(*(video.VideoURL), ",")
	if len(url) != 2 {
		return database.Video{}, nil
	}
	bucket := url[0]
	key := url[1]

	newURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Hour)
	if err != nil {
		return database.Video{}, err
	}
	video.VideoURL = &newURL
	return video, nil
}
