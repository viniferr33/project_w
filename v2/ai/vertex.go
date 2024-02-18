package ai

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"project_w/v2/config"
)

func PredictText(prompt string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, config.PROJECT_ID, config.PROJECT_LOCATION)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			fmt.Println(part)
		}
	}
	return "", nil
}
