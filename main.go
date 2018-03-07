package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codenaut/imgtool/processor"
	"github.com/pelletier/go-toml"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1"
	app.Name = "imgtool"
	app.Usage = "Perform a number for operations on given images"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "Load configuration from FILE. Default to std input",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Output to FILE. Default to std output",
		},
	}
	app.Action = func(c *cli.Context) error {
		config, err := loadConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}
		output := os.Stdout
		outputfile := c.GlobalString("output")
		if outputfile != "" {
			output, err = os.Create(outputfile)
			if err != nil {
				return err
			}
		}
		processor := processor.New(*config)
		return processor.Process(output, c.Args())
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
	}

}

func loadConfig(filename string) (*processor.PageConfig, error) {
	unmarshall := func(content []byte) (*processor.PageConfig, error) {
		var config processor.PageConfig
		if err := toml.Unmarshal(content, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}
	if filename == "" {
		if content, err := ioutil.ReadAll(os.Stdin); err != nil {
			return nil, err
		} else {
			return unmarshall(content)
		}
	} else if content, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return unmarshall(content)
	}

}
