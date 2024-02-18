package transcript

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"project_w/v2/filehandler"
	"project_w/v2/speech"
)

type Transcript struct {
	File      *filehandler.File
	Sentences []*Sentence
	Abstract  string
}

type Sentence struct {
	StartTime float64
	EndTime   float64
	RawText   string
	Abstract  string
}

func ConvertResultToTranscript(result []*speech.TranscriptResult, destination string) (*Transcript, error) {
	filepath := fmt.Sprintf("%s/%s.txt", destination, uuid.New().String())
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	transcriptFile, err := filehandler.GetFileFromFilepath(filepath)
	if err != nil {
		return nil, err
	}

	transcript := Transcript{
		File:      transcriptFile,
		Sentences: make([]*Sentence, 0),
		Abstract:  "",
	}

	var currentTime = 0.0
	for _, transcriptResult := range result {
		_, err := file.WriteString(transcriptResult.Text + "\n")
		if err != nil {
			return nil, err
		}

		sentence := Sentence{
			StartTime: currentTime,
			EndTime:   float64(transcriptResult.EndTime),
			RawText:   transcriptResult.Text,
		}

		currentTime += sentence.EndTime
		transcript.Sentences = append(transcript.Sentences, &sentence)
	}

	return &transcript, nil
}
