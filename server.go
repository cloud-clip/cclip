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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

type clipItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	MIME             string `json:"mime"`
	CreationTime     int64  `json:"ctime"`
	ModificationTime int64  `json:"mtime"`
	Size             int64  `json:"size"`
	ResourceLink     string `json:"resource"`
	ShareLink        string `json:"share"`
}

type clipMetaData struct {
	Name string `json:"name"`
	MIME string `json:"mime"`
}

type serverInfo struct {
	IP   string `json:"ip"`
	Time string `json:"time"`
}

type uploadFileResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	MIME             string `json:"mime"`
	CreationTime     int64  `json:"ctime"`
	ModificationTime int64  `json:"mtime"`
	Size             int64  `json:"size"`
	ResourceLink     string `json:"resource"`
	ShareLink        string `json:"share"`
}

// Password - The API password
var Password string

func checkPassword(next http.Handler) http.Handler {
	authorization := "Bearer " + Password

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != authorization {
			w.WriteHeader(401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func deleteAllClips(w http.ResponseWriter, req *http.Request) {
	clips, err := ScanClipDirectory()
	if err == nil {
		for _, c := range clips {
			err := c.Delete()
			if err != nil {
				SendError(w, err)
				return
			}
		}

		w.Header().Set("Content-Length", "0")
		w.WriteHeader(204)
	} else {
		SendError(w, err)
	}
}

func deleteClip(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	clip, err := GetClipByID(vars["id"])
	if err == nil {
		err = clip.Delete()
		if err == nil {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(204)

			return
		}
	}

	SendError(w, err)
}

func getClipData(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	clipFileName := path.Join(ClipDirectory, vars["id"])
	clipFileStat, err := os.Stat(clipFileName)
	if err != nil {
		SendError(w, err)
		return
	}

	if clipFileStat.IsDir() {
		w.WriteHeader(404)
		return
	}

	clipMetaFileName := path.Join(ClipDirectory, vars["id"]+".meta")
	clipMetaFileStat, err := os.Stat(clipMetaFileName)
	if err != nil {
		SendError(w, err)
		return
	}

	if clipMetaFileStat.IsDir() {
		w.WriteHeader(404)
		return
	}

	// get mime type
	var clipMime string
	{
		clipMetaBytes, err := ioutil.ReadFile(clipMetaFileName)
		if err == nil {
			var clipMeta clipMetaData

			err := json.Unmarshal(clipMetaBytes, &clipMeta)
			if err == nil {
				clipMime = clipMeta.MIME
			}
		}
	}
	clipMime = strings.TrimSpace(strings.ToLower(clipMime))

	file, err := os.Open(clipFileName)
	if err != nil {
		SendError(w, err)
		return
	}

	defer file.Close()

	if clipMime != "" {
		w.Header().Set("Content-Type", clipMime)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(clipFileStat.Size(), 10))
	w.Header().Set("Date", clipFileStat.ModTime().Format(http.TimeFormat))
	io.Copy(w, file)
}

func getClips(w http.ResponseWriter, req *http.Request) {
	// try scan directory for ".meta" files
	clips, err := ScanClipDirectory()
	if err != nil {
		SendError(w, err)
		return
	}

	items := make([]clipItem, 0)

	for _, c := range clips {
		clipMetaBytes, err := ioutil.ReadFile(c.metaFile)
		if err != nil {
			continue
		}

		var clipMeta clipMetaData
		err = json.Unmarshal(clipMetaBytes, &clipMeta)
		if err != nil {
			continue
		}

		// clip meta from JSON
		err = json.Unmarshal(clipMetaBytes, &clipMeta)
		if err != nil {
			continue
		}

		// create a new clip item for the result list
		var newItem clipItem
		newItem.ID = c.id
		newItem.MIME = clipMeta.MIME
		newItem.Name = clipMeta.Name
		newItem.ModificationTime = c.fileInfo.ModTime().Unix()
		newItem.CreationTime = newItem.ModificationTime
		newItem.Size = c.fileInfo.Size()
		newItem.ResourceLink = "/api/v1/clips/" + url.PathEscape(newItem.ID)
		newItem.ShareLink = "/api/v1/shares/" + url.PathEscape(newItem.ID)

		items = append(items, newItem)
	}

	// serialize list to JSON
	bytes, err := json.Marshal(items)
	if err != nil {
		SendError(w, err)
		return
	}

	if len(clips) > 0 {
		w.Header().Set("Date", clips[0].fileInfo.ModTime().Format(http.TimeFormat))
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Write(bytes)
}

func getClipsHead(w http.ResponseWriter, req *http.Request) {
	clips, err := ScanClipDirectory()
	if err == nil {
		if len(clips) > 0 {
			w.Header().Set("Date", clips[0].fileInfo.ModTime().Format(http.TimeFormat))
		}

		w.Header().Set("Content-Length", "0")
		w.Header().Set("X-Cclip-Count", strconv.Itoa(len(clips)))

		w.WriteHeader(204)
	} else {
		w.WriteHeader(500)
	}
}

func getServerInfo(w http.ResponseWriter, req *http.Request) {
	// collect data
	var info serverInfo
	info.IP = req.RemoteAddr
	info.Time = time.Now().Format(time.RFC3339)

	// serialize to JSON
	bytes, err := json.Marshal(info)
	if err != nil {
		SendError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(bytes)
}

func uploadClip(w http.ResponseWriter, req *http.Request) {
	tmpFile, err := ioutil.TempFile("", "cclip")
	if err != nil {
		SendError(w, err)
		return
	}

	// try delete, when leave function
	defer os.Remove(tmpFile.Name())

	if MaxClipSize > 0 {
		// has a maximum size
		req.Body = http.MaxBytesReader(w, req.Body, MaxClipSize)
	}

	defer req.Body.Close()

	_, err = io.Copy(tmpFile, req.Body)
	if err != nil {
		SendError(w, err)
		return
	}

	ctime := time.Now().Unix()
	id := strings.ReplaceAll(uuid.New().String(), "-", "")

	clipFileName := path.Join(ClipDirectory, id)
	clipMetaFileName := path.Join(ClipDirectory, id+".meta")

	err = os.Rename(tmpFile.Name(), clipFileName)
	if err != nil {
		SendError(w, err)
		return
	}

	clipFile, err := os.Open(clipFileName)
	if err != nil {
		os.Remove(clipFileName)

		SendError(w, err)
		return
	}

	defer clipFile.Close()

	clipMime := strings.TrimSpace(strings.ToLower(req.Header.Get("Content-Type")))
	if clipMime == "" {
		clipMime, err = GetFileContentType(clipFile)
		if err != nil {
			os.Remove(clipFileName)

			SendError(w, err)
			return
		}
	}

	// create clip meta
	var clipMeta clipMetaData
	clipMeta.MIME = clipMime
	clipMeta.Name = strings.TrimSpace(req.Header.Get("X-Cclip-Name"))

	// serialize meta to JSON
	bytes, err := json.Marshal(clipMeta)
	if err != nil {
		os.Remove(clipFileName)

		SendError(w, err)
		return
	}

	// try write to .meta file
	err = ioutil.WriteFile(clipMetaFileName, bytes, 0644)
	if err != nil {
		os.Remove(clipFileName)

		SendError(w, err)
		return
	}

	// create response object
	var response uploadFileResponse
	response.ID = id
	response.MIME = clipMeta.MIME
	response.Name = clipMeta.Name
	response.ResourceLink = "/api/v1/clips/" + url.PathEscape(id)
	response.ShareLink = "/api/v1/shares/" + url.PathEscape(id)
	response.CreationTime = ctime
	response.ModificationTime = -1
	response.Size = -1

	clipFileStat, err := os.Stat(clipFileName)
	if err == nil {
		response.ModificationTime = clipFileStat.ModTime().Unix()
		response.Size = clipFileStat.Size()
	}

	// serialize response
	bytes, err = json.Marshal(response)
	if err != nil {
		SendError(w, err)
		return
	}

	// send data
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.WriteHeader(201)
	w.Write(bytes)
}

// RunServer - Runs the server component
func RunServer(c *cli.Context) error {
	// CCLIP_PORT
	envPort := strings.TrimSpace(os.Getenv("CCLIP_PORT"))
	if envPort == "" {
		// default port
		envPort = "50979"
	}

	// CCLIP_DIR
	envDir := strings.TrimSpace(os.Getenv("CCLIP_DIR"))
	if envDir == "" {
		// default clip dir
		envDir = "clips"
	}
	if !path.IsAbs(envDir) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalln("Could not be current working directory", err.Error())
		}

		envDir = path.Join(cwd, envDir)
	}

	// CCLIP_MAX_SIZE
	envMaxSize := strings.TrimSpace(os.Getenv("CCLIP_MAX_SIZE"))
	if envMaxSize == "" {
		// default port
		envMaxSize = "134217728"
	}

	// convert CCLIP_PORT to integer
	port, err := strconv.Atoi(envPort)
	if err != nil {
		// invalid integer string
		log.Fatalln("Invalid TCP port", envPort, err.Error())
	}

	// check if valid TCP port value
	if port < 0 || port > 65535 {
		log.Fatalln("Invalid TCP port", port)
	}

	// convert CCLIP_MAX_SIZE to integer
	maxSize, err := strconv.ParseInt(envMaxSize, 10, 64)
	if err != nil {
		// invalid integer string
		log.Fatalln("Invalid value for maximum clip size", maxSize, err.Error())
	}

	// create output directory, if needed
	clipDirStat, err := os.Stat(envDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(envDir, 0755)
		if err != nil {
			log.Fatalln("Creating", envDir, "directory failed")
		}
	} else if !clipDirStat.IsDir() {
		log.Fatalln(envDir, "is no directory")
	}

	ClipDirectory = envDir // set clip directory
	log.Println("Use clip directory", ClipDirectory)

	if maxSize > 0 {
		MaxClipSize = maxSize

		log.Println("Using maximum clip size of", MaxClipSize, "bytes ...")
	} else {
		log.Println("[WARN] You have no maximum clip size defined")
	}

	Password = strings.TrimSpace(os.Getenv("CCLIP_PASSWORD"))

	router := mux.NewRouter()

	if Password != "" {
		router.Use(checkPassword)
	} else {
		log.Println("[WARN] You have no password defined! Use CCLIP_PASSWORD to set one")
	}

	// initialize routes
	AddHTTPAction(router, "", getServerInfo, "GET")
	AddHTTPAction(router, "/clips", deleteAllClips, "DELETE")
	AddHTTPAction(router, "/clips", getClips, "GET")
	AddHTTPAction(router, "/clips", getClipsHead, "HEAD")
	AddHTTPAction(router, "/clips", uploadClip, "POST")
	AddHTTPAction(router, "/clips/{id:[0-9a-f]{32}}", getClipData, "GET")
	AddHTTPAction(router, "/clips/{id:[0-9a-f]{32}}", deleteClip, "DELETE")

	log.Println("Server will run on port", port, "...")

	// try start server
	err = http.ListenAndServe(":"+strconv.Itoa(port), router)
	if err != nil {
		// failed
		log.Fatalln(err.Error())
	}

	return nil
}
