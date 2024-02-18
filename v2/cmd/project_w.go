package main

import (
	"fmt"
	"log"
	"log/slog"
	_ "project_w/v2/config"
	"project_w/v2/filehandler"
	"project_w/v2/speech"
	"project_w/v2/transcript"
)

func main() {
	slog.Info("Started Project W!")

	/*
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

		transcriptFiles := make([]*filehandler.File, 0)
		for _, s := range gcsAudioList {
			f, err := speech.SpeechToText(s, config.TMP_DIR)
			if err != nil {
				log.Panic(err)
			}
			slog.Info(fmt.Sprintf("Transcript -> %s", s))
			transcriptFiles = append(transcriptFiles, f)
		}
	*/

	/*
		copyFile, _ := filehandler.GetFileFromFilepath("tmp/video.mp4")
		v, _ := ffmpeg.NewVideoFromFile(*copyFile)
		segmentVideoList := make([]ffmpeg.Video, 0)
		seg1_f, _ := filehandler.GetFileFromFilepath("tmp/video_0.mp4")
		seg1_v, _ := ffmpeg.NewVideoFromFile(*seg1_f)
		segmentVideoList = append(segmentVideoList, *seg1_v)

		seg2_f, _ := filehandler.GetFileFromFilepath("tmp/video_1.mp4")
		seg2_v, _ := ffmpeg.NewVideoFromFile(*seg2_f)
		segmentVideoList = append(segmentVideoList, *seg2_v)

		segmentAudioList := make([]ffmpeg.Audio, 0)
		seg1_f, _ = filehandler.GetFileFromFilepath("tmp/audio_0.flac")
		segmentAudioList = append(segmentAudioList, ffmpeg.Audio{
			File:  seg1_f,
			Video: seg1_v,
		})
		seg2_f, _ = filehandler.GetFileFromFilepath("tmp/audio_1.flac")
		segmentAudioList = append(segmentAudioList, ffmpeg.Audio{
			File:  seg2_f,
			Video: seg2_v,
		})
	*/

	transcriptFiles := make([]*filehandler.File, 0)
	text1, _ := filehandler.GetFileFromFilepath("tmp/text_0.json")
	text2, _ := filehandler.GetFileFromFilepath("tmp/text_1.json")
	transcriptFiles = append(transcriptFiles, text1)
	transcriptFiles = append(transcriptFiles, text2)

	bestAlternatives := make([]*speech.TranscriptResult, 0)
	for _, file := range transcriptFiles {
		transcriptResult, err := speech.LoadResultsFromFile(file)
		if err != nil {
			log.Panic(err)
		}
		for _, results := range transcriptResult {
			bestAlternative := speech.FindBestAlternative(results)
			bestAlternatives = append(bestAlternatives, bestAlternative)
		}
	}

	_, err := transcript.ConvertResultToTranscript(bestAlternatives, "tmp")
	if err != nil {
		log.Panic(err)
	}
	slog.Info(fmt.Sprintf("Converted to Transcript"))
}
