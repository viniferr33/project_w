package speech

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"project_w/v2/filehandler"
)

type Word struct {
	StartTime       int64  `json:"start_time,omitempty"`
	EndTime         int64  `json:"end_time,omitempty"`
	Word            string `json:"word,omitempty"`
	Speaker         string `json:"speaker,omitempty"`
	ConfidenceScore float32
}

type TranscriptResult struct {
	Source          string  `json:"source,omitempty"`
	EndTime         int64   `json:"end_time,omitempty"`
	LanguageCode    string  `json:"language_code,omitempty"`
	Text            string  `json:"text,omitempty"`
	Words           []*Word `json:"words,omitempty"`
	ConfidenceScore float32
}

/**
OUTPUT FILE EXAMPLE

@Source: uri

@ResultListStart
@Result: N
@resultEndTime: NN.N
@languageCode: ss-SS

@transcriptStart
LOREM IPSUM DOLOR ...
@transcriptEnd

@WordListStart
WordIndex WordStartTime WordEndTime @Word @SpeakerLabel
@WordListEnd
*/

func SpeechToText(uri string, fileDestination string) (*filehandler.File, error) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	op, err := client.LongRunningRecognize(ctx, &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_FLAC,
			AudioChannelCount:          2,
			LanguageCode:               "pt-BR",
			EnableWordTimeOffsets:      true,
			EnableAutomaticPunctuation: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: uri,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	/**
	Instead of Wait, should loop until op.Done() and print the op.Metadata().ProgressPercent
	*/
	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, err
	}

	transcriptResultList := make([][]TranscriptResult, 0)
	for _, result := range resp.Results {

		transcriptAlternatives := make([]TranscriptResult, 0)
		for _, alternative := range result.Alternatives {
			words := make([]*Word, 0)
			for _, word := range alternative.Words {
				w := Word{
					StartTime:       word.StartTime.Seconds,
					EndTime:         word.EndTime.Seconds,
					Word:            word.Word,
					Speaker:         word.SpeakerLabel,
					ConfidenceScore: word.Confidence,
				}

				words = append(words, &w)
			}

			transcriptResult := TranscriptResult{
				Source:          uri,
				EndTime:         result.ResultEndTime.Seconds,
				LanguageCode:    result.LanguageCode,
				Text:            alternative.Transcript,
				Words:           words,
				ConfidenceScore: alternative.Confidence,
			}

			transcriptAlternatives = append(transcriptAlternatives, transcriptResult)
		}

		transcriptResultList = append(transcriptResultList, transcriptAlternatives)
	}

	resultTranscriptFileName := fmt.Sprintf("%s/%s.json", fileDestination, uuid.New().String())
	file, err := os.OpenFile(resultTranscriptFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(transcriptResultList)
	if err != nil {
		return nil, err
	}

	fileH, err := filehandler.GetFileFromFilepath(resultTranscriptFileName)
	if err != nil {
		return nil, err
	}

	return fileH, nil
}

func LoadResultsFromFile(f *filehandler.File) ([][]TranscriptResult, error) {
	file, err := os.Open(f.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var transcriptResults [][]TranscriptResult
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&transcriptResults)
	if err != nil {
		return nil, err
	}

	return transcriptResults, nil
}

func FindBestAlternative(alternatives []TranscriptResult) *TranscriptResult {
	var bestScore float32 = 0.0
	var bestAlternative TranscriptResult

	for _, alternative := range alternatives {
		if alternative.ConfidenceScore > bestScore {
			bestAlternative = alternative
		}
	}

	return &bestAlternative
}
