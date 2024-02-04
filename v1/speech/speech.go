package speech

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"fmt"
	"project_w/v1/audio"
)

func GetTranscription(a audio.Audio) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	op, err := client.LongRunningRecognize(ctx, &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:              speechpb.RecognitionConfig_FLAC,
			LanguageCode:          "pt-BR",
			AudioChannelCount:     2,
			EnableWordTimeOffsets: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: "gs://project_w_audio/ffmpeg.flac",
			},
		},
	})
	if err != nil {
		return err
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
			for _, w := range alt.Words {
				fmt.Printf(
					"Word: \"%v\" (startTime=%3f, endTime=%3f)\n",
					w.Word,
					float64(w.StartTime.Seconds)+float64(w.StartTime.Nanos)*1e-9,
					float64(w.EndTime.Seconds)+float64(w.EndTime.Nanos)*1e-9,
				)
			}
		}
	}

	return nil
}
