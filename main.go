// Cloud Clip
// Copyright (C) 2020  Marcel Joachim Kloubert <marcel.kloubert@gmx.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"

	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "test",
				Value:    "wurst",
				Required: false,
				Aliases:  []string{"t"},
				Usage:    "a test flag",
			},
		},
		Name:   "Cloud Clip",
		Usage:  "Cloud Clip server and command line interface",
		Action: runApp,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp(c *cli.Context) error {
	fmt.Println("Hello, cclip! There are", c.Args().Len(), "argument(s)")

	return nil
}
