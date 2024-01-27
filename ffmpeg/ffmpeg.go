package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"project_w/video"
	"strconv"
	"strings"
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

	startTime, err := strconv.ParseFloat(videoInfo.Format.StartTime, 32)
	if err != nil {
		return nil, err
	}

	duration, err := strconv.ParseFloat(videoInfo.Format.Duration, 32)
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(videoInfo.Format.Size)
	if err != nil {
		return nil, err
	}

	splitFilename := strings.Split(videoInfo.Format.Filename, ".")

	return &video.Video{
		Id:         splitFilename[0],
		Filename:   videoInfo.Format.Filename,
		FormatName: videoInfo.Format.FormatName,
		StartTime:  startTime,
		Duration:   duration,
		Size:       size,
	}, nil
}

func SliceVideo(v video.Video, videoSlices *[]video.Video) error {
	outputPattern := v.Id + "_seg"
	cmd := exec.Command("ffmpeg", "-i", v.Filename, "-c", "copy", "-map", "0", "-f", "segment", "-segment_time", "600", "-reset_timestamps", "1", "-segment_format", "mp4", outputPattern+"_%03d.mp4")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	segments, err := filepath.Glob(fmt.Sprintf("%s*.mp4", outputPattern))
	if err != nil {
		return err
	}

	for _, segment := range segments {
		newVideo, err := GetVideo(segment)
		if err != nil {
			return err
		}

		*videoSlices = append(*videoSlices, *newVideo)
	}

	return nil
}
