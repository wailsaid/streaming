package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/saidwail/streaming/database"
	"github.com/saidwail/streaming/models"
)

func ScanVideoForAdultContent(videoPath string) ([]float64, error) {
	// This is a placeholder function. In a real-world scenario, you'd use a more sophisticated
	// video processing library and AI service to detect adult content.

	// For demonstration purposes, we'll use ffprobe to get video duration and simulate finding adult content
	duration, err := getVideoDuration(videoPath)
	if err != nil {
		return nil, err
	}

	// Simulate finding adult content at 25%, 50%, and 75% of the video duration
	timestamps := []float64{
		duration * 0.25,
		duration * 0.50,
		duration * 0.75,
	}

	return timestamps, nil
}

func getVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func RemoveAdultContent(videoID string, timestamps []string) error {
	// This is a placeholder function. In a real-world scenario, you'd implement video editing logic here.
	db := database.DB
	db.Model(&models.Video{}).Where("id = ?", videoID).Update("description", "This video has been flagged for adult content and is not suitable for all viewers.")

	fmt.Printf("Removing adult content from video %s at timestamps: %v\n", videoID, timestamps)
	return nil
}
