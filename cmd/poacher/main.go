package main

import (
	"fmt"
	"log"
	"os"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
)

func main() {
	poacher := cli.NewApp()
	poacher.Name = "poacher"
	poacher.Usage = "poach data"
	poacher.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}
	poacher.Commands = []cli.Command{
		{
			Name: "coinmarketcap",
			Subcommands: []cli.Command{
				{
					Name: "ticker",
					Action: func(c *cli.Context) error {
						// Connect to influx database
						influxClient, err := client.NewHTTPClient(client.HTTPConfig{
							Addr:     c.String("influx-url"),
							Username: c.String("influx-user"),
							Password: c.String("influx-pass"),
						})
						if err != nil {
							log.Fatalf("Could not connect to database, err: %v", err)
						}
						defer influxClient.Close()

						for {
							log.Println("Saving tickers")
							saveTickers(influxClient, c.String("currency"))
							log.Println("Done")

							log.Println("Sleeping", c.Duration("interval"))
							<-time.After(c.Duration("interval"))
						}
						return nil
					},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "influx-url",
							Value: "http://localhost:8086",
						},
						cli.StringFlag{
							Name:  "influx-user",
							Value: "",
						},
						cli.StringFlag{
							Name:  "influx-pass",
							Value: "",
						},
						cli.DurationFlag{
							Name:  "interval",
							Value: time.Duration(time.Minute * 5),
							Usage: "scrape data every `DURATION`",
						},
						cli.StringFlag{
							Name:  "currency",
							Value: "USD",
							Usage: "request data in this currency",
						},
					},
				},
			},
			Category: "Coinmarketcap",
		},
	}

	err := poacher.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
