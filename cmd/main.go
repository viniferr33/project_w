package main

import (
	"fmt"
	"project_w/ffmpeg"
)

func main() {
	_, err := ffmpeg.GetVideo("tmp/video.mp4")
	if err != nil {
		fmt.Println(err)
	}
}
