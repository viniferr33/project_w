package ai

import (
	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
	"context"
	"fmt"
	"log"
)

func GetEntities(content string) error {
	ctx := context.Background()

	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	document := languagepb.Document{
		Type: languagepb.Document_PLAIN_TEXT,
		Source: &languagepb.Document_Content{
			Content: content,
		},
		LanguageCode: "pt-BR",
	}

	entities, err := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &document,
	})
	if err != nil {
		return err
	}

	for i, entity := range entities.Entities {
		fmt.Printf("Entity %d -> %s\ttype%s\t%v\n", i, entity.Name, entity.Type, entity.Metadata)
	}
	return nil
}
