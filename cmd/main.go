package main

import (
	"fmt"
	"project_w/ffmpeg"
	"project_w/video"
)

func main() {
	v, err := ffmpeg.GetVideo("tmp/video.mp4")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)

	vList := make([]video.Video, 0)
	err = ffmpeg.SliceVideo(*v, &vList)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vList)
}
