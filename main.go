package main

import (
	"Mc-Lang-GPT-translator/translator"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "MC-LANG-GPT-Translator",
		Commands: []*cli.Command{
			{
				Name:    "translator",
				Aliases: []string{"tran"},
				Usage:   "translate lang file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "specify input file",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "specify output file",
					},
					&cli.IntFlag{
						Name:    "maxConcurrent",
						Aliases: []string{"c"},
						Usage:   "specify max concurrent",
						Value:   8,
					},
					&cli.IntFlag{
						Name:    "maxLine",
						Aliases: []string{"l"},
						Usage:   "specify max line",
						Value:   20,
					},
					&cli.StringFlag{
						Name:    "Model",
						Aliases: []string{"m"},
						Usage:   "gpt-4-32k-0314 | gpt-4-32k | gpt-4-0314 | gpt-4 | gpt-3.5-turbo-0301 | gpt-3.5-turbo",
						Value:   "gpt-3.5-turbo-0301",
					},
				},
				Action: func(c *cli.Context) error {
					input := c.String("input")
					output := c.String("output")
					translator.MaxConcurrent = c.Int("maxConcurrent")
					translator.MaxLine = c.Int("maxLine")
					translator.Model = c.String("Model")
					if input == "" {
						log.Fatal("input file not specified")
					}
					if output == "" {
						output = "out-" + input
					}
					translator.StartTranslator(input, output)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
