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
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

type serverInfo struct {
	IP   string `json:"ip"`
	Time string `json:"time"`
}

func getServerInfo(w http.ResponseWriter, req *http.Request) {
	var info serverInfo
	info.IP = req.RemoteAddr
	info.Time = time.Now().Format(time.RFC3339)

	bytes, err := json.Marshal(info)
	if err != nil {
		SendError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(bytes)
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

	// initialize routes
	router.HandleFunc("/api/v1", getServerInfo).Methods("GET")

	log.Println("Server will run on port", port, "...")

	// try start server
	err = http.ListenAndServe(":"+strconv.Itoa(port), router)
	if err != nil {
		// failed
		log.Fatalln(err.Error())
	}

	return nil
}
