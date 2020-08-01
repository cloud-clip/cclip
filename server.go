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
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

func hello(w http.ResponseWriter, req *http.Request) {
	b := []byte("Hello!")

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(b)
}

// RunServer - Runs the server component
func RunServer(c *cli.Context) error {
	envPort := strings.TrimSpace(os.Getenv("CCLIP_PORT"))
	if envPort == "" {
		// default port
		envPort = "50979"
	}

	port, err := strconv.Atoi(envPort)
	if err != nil {
		// invalid integer string
		log.Fatalln(err.Error())
	}

	if port < 0 || port > 65535 {
		// port value is out of range
		log.Fatalln("Invalid TCP port")
	}

	router := mux.NewRouter()

	router.HandleFunc("/hello", hello).Methods("GET")

	log.Println("Server will run on port", port, "...")

	// try start server
	err = http.ListenAndServe(":"+strconv.Itoa(port), router)
	if err != nil {
		// failed
		log.Fatalln(err.Error())
	}

	return nil
}
