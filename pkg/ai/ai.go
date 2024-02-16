package ai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
}

func New(token string) *Client {
	return &Client{
		client: openai.NewClient(token),
	}
}

func (c *Client) Correct(ctx context.Context, input string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: `You're helpful german tutor. Correct user message to standard German`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
