package main

import (
	"github.com/codegangsta/cli"
	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/noaa/events"
	"fmt"
	"crypto/tls"
	"encoding/json"
	"os"
	"io/ioutil"
	"path/filepath"
)

type config struct {
	AccessToken string
}

func main() {
	app := cli.NewApp()
	app.Name = "sounder"
	app.Usage = "acceptance tool for the metric system"

	app.Commands = []cli.Command{
		{
			Name:      "stream",
			ShortName: "s",
			Usage:     "stream messages",
			Action: func(c *cli.Context) {
				consumer := noaa.NewNoaa(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.Stream(c.Args().Get(1), authToken())
				if err != nil {
					panic(err)
				}
				for message := range messages {
					displayMessage(message)
				}
			},
		},
		{
			Name:      "recent",
			ShortName: "r",
			Usage:     "recent log messages",
			Action: func(c *cli.Context) {
				consumer := noaa.NewNoaa(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.RecentLogs(c.Args().Get(1), authToken())
				if err != nil {
					panic(err)
				}
				for _, message := range messages {
					displayMessage(message)
				}
			},
		},
		{
			Name:      "tail logs",
			ShortName: "t",
			Usage:     "tail log messages",
			Action: func(c *cli.Context) {
				consumer := noaa.NewNoaa(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.TailingLogs(c.Args().Get(1), authToken())
				if err != nil {
					panic(err)
				}
				for message := range messages {
					displayMessage(message)
				}
			},
		},
		{
			Name:      "firehose",
			ShortName: "f",
			Usage:     "all data",
			Action: func(c *cli.Context) {
				consumer := noaa.NewNoaa(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.Firehose(authToken())
				if err != nil {
					panic(err)
				}
				for message := range messages {
					displayMessage(message)
				}
			},
		},
	}
	app.Run(os.Args)
}

func displayMessage(m *events.Envelope) {
	fmt.Printf("%v \n", m)
}

func authToken() string {
	var c config

	configDir := filepath.Join(os.Getenv("HOME"), ".cf")
	file, err := os.Open(filepath.Join(configDir, "config.json"))
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}

	fmt.Println("TOKEN:", c.AccessToken)

	return c.AccessToken
}
