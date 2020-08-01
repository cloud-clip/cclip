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

	"github.com/urfave/cli/v2"
)

// AppCommands - all known app commands
var AppCommands = []*cli.Command{
	{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "a test command",
		Action:  test,
	},
}

func test(c *cli.Context) error {
	fmt.Println("completed task: ", c.Args().First())
	return nil
}
