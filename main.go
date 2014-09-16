package main

import (
	"github.com/codegangsta/cli"
	"github.com/cloudfoundry/dropsonde/autowire/metrics"
	"github.com/cloudfoundry/dropsonde/autowire/logs"
	"github.com/cloudfoundry/loggregator_consumer/dropsonde_consumer"
	"fmt"
	"time"
	"crypto/tls"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "sounder"
	app.Usage = "acceptance tool for the metric system"

	app.Commands = []cli.Command{
		{
			Name:      "emit",
			ShortName: "e",
			Usage:     "emit metrics",
			Action: func(*cli.Context) {
				metrics.IncrementCounter("counter")
				metrics.SendValue("value", 42.0, "unknown")

				fmt.Println("metrics sent")
			},
		},

		{
			Name:      "log",
			ShortName: "l",
			Usage:     "emit logs",
			Action: func(c *cli.Context) {
				logs.SendAppLog(c.Args().First(), "This is a log message", "sounder", "0")
				logs.SendAppErrorLog(c.Args().First(), "This is a error message", "sounder", "0")
				fmt.Println("logs sent")
			},
		},

		{
			Name:      "stream",
			ShortName: "s",
			Usage:     "stream dropsonde messages",
			Action: func(c *cli.Context) {
				consumer := dropsonde_consumer.NewDropsondeConsumer(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.Stream(c.Args().Get(1), c.Args().Get(2))
				if err != nil {
					panic(err)
				}
				for message := range messages {
					ts := message.GetTimestamp()
					timestamp := time.Unix(0, ts)
					fmt.Printf("%s - Type: %s\n", timestamp.String(), message.GetEventType().String())
				}
			},
		},

		{
			Name:      "recent",
			ShortName: "r",
			Usage:     "recent dropsonde log messages",
			Action: func(c *cli.Context) {
				consumer := dropsonde_consumer.NewDropsondeConsumer(c.Args().First(), &tls.Config{InsecureSkipVerify: true}, nil)
				messages, err := consumer.RecentLogs(c.Args().Get(1), c.Args().Get(2))
				if err != nil {
					panic(err)
				}
				for _, message := range messages {
					timestamp := time.Unix(0, message.GetTimestamp())
					fmt.Printf("%s - Type: %s\n", timestamp.String(), message.GetEventType().String())
				}
			},
		},
	}
	app.Run(os.Args)
}
