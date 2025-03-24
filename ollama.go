package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ollama/ollama/api"
)

func AskOllama(msg string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", fmt.Errorf("failed to create Ollama client: %w", err)
	}

	var prompt = `You are an expert in Markdown and HTML. You are well-versed in Hugo for creating GitHub Pages.

I have some files that are not 100% compatible with the syntax of Markdown files for Hugo. They need to be fixed following these guidelines:

- I would like to remove the HTML but preserve the formatting that the HTML has, making the text compatible with Markdown.
- I would like to change the script embeds, such as the following example line where xxxxxxxx is a placeholder simulating numbers:

      <script src="https://gist.github.com/vicendominguez/xxxxxxxx.js"></script>

  and preserve the original url but changing it to this format:

      [Gist](https://gist.github.com/vicendominguez/xxxxxxxx.js)


- I would like to achieve good English writing. A bit informal but professional enough for a blog post.
- I would like it to be compatible with a Hugo template.
- Your response should only and exclusively be the raw code of the file so I can pipe it to a file. 
- No extra comments. No summary. Nothing else.
- Only in English. 

The code is as follows:
` + "\n" + msg

	messages := []api.Message{
		{Role: "system", Content: "You are an expert code analyzer."},
		{Role: "user", Content: prompt},
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "qwen2.5-coder:7b",
		Messages: messages,
		Stream:   new(bool),
	}

	var commitMessage string

	err = client.Chat(ctx, req, func(resp api.ChatResponse) error {
		commitMessage += resp.Message.Content
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	commitMessage = SanitizeResponse(commitMessage)
	return commitMessage, nil
}

func SanitizeResponse(response string) string {
	// Compile the regular expression to find content between <think> and </think>
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)

	// Replace all matches with an empty string
	sanitizedResponse := re.ReplaceAllString(response, "")
	sanitizedResponse = removeQwenShits(sanitizedResponse)

	return sanitizedResponse
}
func removeQwenShits(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return s // Si hay menos de 2 líneas, devolvemos el string original
	}

	// Verificamos si la primera línea contiene "`markdown"
	if strings.Contains(lines[0], "``markdown") {
		// Eliminamos la primera y la última línea
		if len(lines) > 2 {
			return strings.Join(lines[1:len(lines)-1], "\n")
		}
		return "" // Si solo quedaba una línea, devolvemos vacío
	}

	return s // Si no cumple la condición, devolvemos el string original
}
