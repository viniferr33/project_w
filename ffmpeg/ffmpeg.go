package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"project_w/video"
)

type videoInfo struct {
	Format format `json:"format"`
}

type format struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	StartTime  string `json:"start_time"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
}

func init() {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("ffmpeg not found in path")
	}
}

func GetVideo(filepath string) (*video.Video, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}

	cmd := exec.Command("ffprobe", "-hide_banner", "-loglevel", "fatal", "-show_error", "-show_format", "-print_format", "json", filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var videoInfo videoInfo
	err = json.Unmarshal(output, &videoInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(output))
	fmt.Println(videoInfo)
	return &video.Video{}, nil
}
