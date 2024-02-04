package main

import (
	"fmt"
	"project_w/v1/ffmpeg"
	"project_w/v1/speech"
)

func main() {
	v, err := ffmpeg.GetVideo("tmp/ffmpeg.mp4")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)

	a, err := ffmpeg.ConvertMp4ToFlac(*v)
	if err != nil {
		fmt.Println(err)
	}

	err = speech.GetTranscription(*a)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(a)
}
