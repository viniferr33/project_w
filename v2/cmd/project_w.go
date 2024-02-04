package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"project_w/v2/config"
	"project_w/v2/ffmpeg"
	"project_w/v2/filehandler"

	_ "project_w/v2/config"
)

func main() {
	slog.Info("Started Project W!")

	input := flag.String("i", "none", "Input File")
	flag.Parse()

	if *input == "none" {
		log.Panic("Invalid input file!")
	}

	slog.Info("Input File -> " + *input)

	file, err := filehandler.GetFileFromFilepath(*input)
	if err != nil {
		log.Panic(err)
	}

	slog.Info(fmt.Sprintf("%v", *file))

	copyFile, err := file.Copy(config.TMP_DIR)
	if err != nil {
		log.Panic(err)
	}
	slog.Info(fmt.Sprintf("%v", copyFile))

	v, err := ffmpeg.NewVideoFromFile(*copyFile)
	if err != nil {
		log.Panic(err)
	}

	slog.Info(fmt.Sprintf("%v", *v))

	segmentVideoList, err := v.SegmentVideo("600")
	if err != nil {
		log.Panic(err)
	}

	slog.Info("Sliced video")
	for i, video := range segmentVideoList {
		slog.Info(fmt.Sprintf("%d -> %v", i, video))
	}

	segmentAudioList := make([]ffmpeg.Audio, 0)
	for _, video := range segmentVideoList {
		audio, err := video.ConvertToAudio("flac")
		if err != nil {
			log.Panic(err)
		}

		segmentAudioList = append(segmentAudioList, *audio)
		slog.Info(fmt.Sprintf("%v", *audio))
	}

}
