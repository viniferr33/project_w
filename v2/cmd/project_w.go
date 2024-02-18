package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"project_w/v2/config"
	_ "project_w/v2/config"
	"project_w/v2/ffmpeg"
	"project_w/v2/filehandler"
	"project_w/v2/gcs"
	"project_w/v2/speech"
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

	gcsAudioList := make([]string, 0)
	for _, audio := range segmentAudioList {
		uploadFile, err := gcs.UploadFile(*audio.File, config.GCS_BUCKET)
		if err != nil {
			log.Panic(err)
		}

		slog.Info(fmt.Sprintf("Uploaded -> %s", uploadFile))
		gcsAudioList = append(gcsAudioList, uploadFile)
	}

	//gcsAudioList = append(gcsAudioList, "gs://project_w_audio/2b3a3779-e479-4af2-b59c-2e7a26df7272_000.flac")
	//gcsAudioList = append(gcsAudioList, "gs://project_w_audio/2b3a3779-e479-4af2-b59c-2e7a26df7272_001.flac")

	transcriptFiles := make([]*filehandler.File, 0)
	for _, s := range gcsAudioList {
		f, err := speech.SpeechToText(s, config.TMP_DIR)
		if err != nil {
			log.Panic(err)
		}
		slog.Info(fmt.Sprintf("Transcript -> %s", s))
		transcriptFiles = append(transcriptFiles, f)
	}

}
