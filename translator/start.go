package translator

import (
	"Mc-Lang-GPT-translator/langparser"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	openai "github.com/sashabaranov/go-openai"
)

var Gentry []langparser.LangEntry
var Gstr string

var MaxConcurrent = 8
var MaxLine = 20
var Model = "gpt-3.5-turbo-0301"

func StartTranslator(input, output string) {
	client := openai.NewClient(os.Getenv("GPT_API_KEY"))
	Gentry, _ = langparser.ParseLangFile(input)
	langparser.WriteLangFile(input, Gentry)

	data, _ := os.ReadFile(input)
	Gstr = string(data)

	stringSlice := splitFile(input)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, MaxConcurrent)
	var entries []langparser.LangEntry
	for i, value := range stringSlice {
		ch <- struct{}{}
		wg.Add(1)
		go func(value string, i int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			log.Printf("Translating %d/%d\n", i+1, len(stringSlice))
			test := translate(value, client, i)
			entry, _ := langparser.ReadLangByString(test)
			entries = append(entries, entry...)
			log.Println("Translated successfully", i+1)
		}(value, i)
	}
	wg.Wait()
	langparser.WriteLangFile(output, langparser.ReplaceLangEntry(Gentry, entries))
}

func splitFile(input string) []string {
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0
	var buffer string
	var stringSlice []string
	for scanner.Scan() {
		buffer += scanner.Text() + "\n"
		counter++

		if counter%MaxLine == 0 {
			stringSlice = append(stringSlice, buffer)
			buffer = ""
		}
	}

	if buffer != "" {
		stringSlice = append(stringSlice, buffer)
	}

	return stringSlice
}

func translate(value string, client *openai.Client, i int) string {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: Model,

			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Minecraft language file translation assistant",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Please translate it into Chinese according to the semantics of Minecraft,The text is a Minecraft lang file in the format key=value, I just need you to translate the value and keep the key=.Only return the translate result, don’t interpret it, don't return the delimited character.",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "OK, please send me the language file and I will translate it into Chinese for you as soon as possible.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: value,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		fmt.Println("Retrying...", i+1)
		return translate(value, client, i)
	}
	message := resp.Choices[0].Message.Content

	if strings.Contains(value, message) || strings.Contains(message, value) || strings.Contains(Gstr, message) {
		fmt.Println("Translated result is the same as the original text.")
		fmt.Println("Retrying...", i+1)
		return translate(value, client, i)
	}
	return message
}
