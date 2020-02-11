package main

//go:generate esc -o static.go -pkg main -modtime 1500000000 -prefix assets assets

import (
	"encoding/json"
	"os"

	"github.com/ghetzel/cli"
	"github.com/ghetzel/diecast"
	"github.com/ghetzel/go-stockutil/log"
)

func main() {
	app := cli.NewApp()
	app.Name = `owndoc`
	app.Usage = `Generate a static site documenting a Golang package and all subpackages`
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `debug`,
			EnvVar: `LOGLEVEL`,
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetLevelString(c.String(`log-level`))
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  `generate`,
			Usage: `Generate a JSON manifest decribing the current package and all subpackages.`,
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				if mod, err := ScanDir(c.Args().First()); err == nil {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent(``, `    `)
					enc.Encode(mod)
				} else {
					log.Fatal(err)
				}
			},
		}, {
			Name:  `render`,
			Usage: `Render a module's documentation as a standalone static site.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   `output-dir, o`,
					Usage:  `The output directory where generated files will be placed.`,
					Value:  `docs`,
					EnvVar: `OWNDOC_DIR`,
				},
			},
			Action: func(c *cli.Context) {
				if mod, err := ScanDir(c.Args().First()); err == nil {
					log.FatalIf(
						RenderHTML(c.String(`output-dir`), mod),
					)
				} else {
					log.Fatal(err)
				}
			},
		}, {
			Name:  `serve`,
			Usage: `Render a module's documentation as a standalone static site.`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   `output-dir, o`,
					Usage:  `The output directory where generated files will be placed.`,
					Value:  `docs`,
					EnvVar: `OWNDOC_DIR`,
				},
				cli.StringFlag{
					Name:   `address, a`,
					Usage:  `The address that the webserver should listen on.`,
					Value:  `localhost:16060`,
					EnvVar: `OWNDOC_ADDRESS`,
				},
			},
			Action: func(c *cli.Context) {
				log.SetLevelString(`error`)

				if mod, err := ScanDir(c.Args().First()); err == nil {
					log.FatalIf(
						RenderHTML(c.String(`output-dir`), mod),
					)
				} else {
					log.Fatal(err)
				}

				log.SetLevelString(c.GlobalString(`log-level`))
				log.Infof("Listening at http://%s", c.String(`address`))

				log.FatalIf(
					diecast.NewServer(c.String(`output-dir`)).ListenAndServe(c.String(`address`)),
				)
			},
		},
	}

	app.Run(os.Args)
}
