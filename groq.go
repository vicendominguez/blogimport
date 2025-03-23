package main

import (
	"context"

	"github.com/conneroisu/groq-go"
)

func AskGroq(g groq.Client, prompt string) (string, error) {

	IAMessage := groq.ChatCompletionMessage{
		Role:    groq.RoleUser,
		Content: prompt,
	}

	var IAMessages []groq.ChatCompletionMessage

	IAMessages = append(IAMessages, IAMessage)

	ctx := context.Background()

	IAResponse, err := g.ChatCompletion(ctx, groq.ChatCompletionRequest{
		//		Model:     groq.ModelLlama370B8192,
		Model:     groq.ModelLlama3370BVersatile,
		Messages:  IAMessages,
		MaxTokens: 2000,
	},
	)

	if err != nil {
		return "", err
	}

	return IAResponse.Choices[0].Message.Content, nil

}
