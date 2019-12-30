package main

import (
	"io"
	"os"

	"github.com/m-mizutani/badman"
	"github.com/m-mizutani/badman/source"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var logger = logrus.New()

func main() {
	if err := handler(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func handler(args []string) error {
	var output string

	app := &cli.App{
		Name:  "badman",
		Usage: "CLI utility for badman",
		Commands: []*cli.Command{
			{
				Name:    "dump",
				Aliases: []string{"d"},
				Usage:   "Download sources and output serialized data",
				Action: func(c *cli.Context) error {
					man := badman.New()
					if err := man.Download(source.DefaultSet); err != nil {
						return errors.Wrapf(err, "Fail to download blacklists")
					}

					var out io.Writer
					if output == "-" {
						out = os.Stdout
					} else {
						fd, err := os.Create(output)
						if err != nil {
							return errors.Wrapf(err, "Fail to create output file: %s", output)
						}
						defer fd.Close()
						out = fd
					}

					if err := man.Dump(out); err != nil {
						return errors.Wrapf(err, "Fail to output blacklists")
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Usage:       "Output file name, '-' means stdout",
						Aliases:     []string{"o"},
						Value:       "-",
						Destination: &output,
					},
				},
			},
		},
	}

	return app.Run(args)
}
