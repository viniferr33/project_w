package ffmpeg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"project_w/v2/filehandler"
	"strconv"
	"strings"
)

type Video struct {
	File      *filehandler.File
	Duration  float64
	StartTime float64
}

type Audio struct {
	File  *filehandler.File
	Video *Video
}

var InvalidExtensionError = errors.New("file extension should be .mp4")

func NewVideoFromFile(f filehandler.File) (*Video, error) {
	type format struct {
		Filename   string `json:"filename"`
		FormatName string `json:"format_name"`
		StartTime  string `json:"start_time"`
		Duration   string `json:"duration"`
		Size       string `json:"size"`
	}

	type videoInfo struct {
		Format format `json:"format"`
	}

	if f.Extension != ".mp4" {
		return nil, InvalidExtensionError
	}

	cmd := exec.Command("ffprobe", "-hide_banner", "-loglevel", "fatal", "-show_error", "-show_format", "-print_format", "json", f.Filepath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var vi videoInfo
	err = json.Unmarshal(output, &vi)
	if err != nil {
		return nil, err
	}

	startTime, err := strconv.ParseFloat(vi.Format.StartTime, 32)
	if err != nil {
		return nil, err
	}

	duration, err := strconv.ParseFloat(vi.Format.Duration, 32)
	if err != nil {
		return nil, err
	}

	return &Video{
		File:      &f,
		Duration:  duration,
		StartTime: startTime,
	}, nil
}

func (v *Video) SegmentVideo(segmentTime string) ([]Video, error) {
	splitFilepath := strings.Split(v.File.Filepath, v.File.Extension)
	outputPattern := fmt.Sprintf("%s_%s%s", splitFilepath[0], "%03d", v.File.Extension)
	extWithoutDot := v.File.Extension[1:]

	cmd := exec.Command("ffmpeg", "-i", v.File.Filepath, "-c", "copy", "-map", "0", "-f", "segment", "-segment_time", segmentTime, "-reset_timestamps", "1", "-segment_format", extWithoutDot, outputPattern)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	segmentsFilepath, err := filepath.Glob(fmt.Sprintf("%s_*%s", splitFilepath[0], v.File.Extension))
	if err != nil {
		return nil, err
	}

	segmentVideoList := make([]Video, 0)
	for _, segment := range segmentsFilepath {
		segmentFile, err := filehandler.GetFileFromFilepath(segment)
		if err != nil {
			return nil, err
		}

		segmentVideo, err := NewVideoFromFile(*segmentFile)
		if err != nil {
			return nil, err
		}

		segmentVideoList = append(segmentVideoList, *segmentVideo)
	}

	return segmentVideoList, nil
}

func (v *Video) ConvertToAudio(format string) (*Audio, error) {
	splitFilepath := strings.Split(v.File.Filepath, v.File.Extension)
	outputPattern := fmt.Sprintf("%s.%s", splitFilepath[0], format)

	cmd := exec.Command("ffmpeg", "-i", v.File.Filepath, "-vn", "-c:a", format, outputPattern)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	audioFile, err := filehandler.GetFileFromFilepath(outputPattern)
	if err != nil {
		return nil, err
	}

	return &Audio{
		File:  audioFile,
		Video: v,
	}, nil
}
