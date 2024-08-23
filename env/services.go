package env

import "os/exec"

func GenerateThumbnail(videoPath string, thumbnailPath string) error {
	timePosition := "00:01:01"
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", timePosition, "-vframes", "1", thumbnailPath)
	return cmd.Run()
}
