package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("missing API_KEY")
	}

	ctx := context.Background()
	client := openai.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with chatGPT in console",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Println("Say something ('quit' to end): ")
				if !scanner.Scan() {
					break
				}
				question := scanner.Text()
				switch question {
				case "quit":
					quit = true
				default:
					GetResponse(client, ctx, question)
				}
			}
		},
	}
	rootCmd.Execute()
}

func GetResponse(client *openai.Client, ctx context.Context, question string) {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 500,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	})

	if err != nil {
		fmt.Println("line 32: ", err)
		os.Exit(13)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

/*
CLI OUTPUT:
$ go run .
Say something ('quit' to end):
Tell me about Go lang
Go, also known as Golang, is a statically typed, compiled programming language that was created by Google in 2007
Say something ('quit' to end):
1 puls 2. how much?
The price for 1 plus 2 would be 3
Say something ('quit' to end):
quit
*/
