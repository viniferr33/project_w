package main

import (
	"fmt"
	"project_w/ffmpeg"
)

func main() {
	v, err := ffmpeg.GetVideo("tmp/video.mp4")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)

	a, err := ffmpeg.ConvertMp4ToFlac(*v)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(a)
}
