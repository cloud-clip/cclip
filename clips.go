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
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ByNewestClipFileList - describes a ClipFile list, which can be sorted by newest item descending
type ByNewestClipFileList []ClipFile

func (a ByNewestClipFileList) Len() int { return len(a) }
func (a ByNewestClipFileList) Less(i, j int) bool {
	// order descending
	return a[i].fileInfo.ModTime().Unix() > a[j].fileInfo.ModTime().Unix()
}
func (a ByNewestClipFileList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// ClipFile - A clip file
type ClipFile struct {
	file         string
	fileInfo     os.FileInfo
	id           string
	metaFile     string
	metaFileInfo os.FileInfo
}

// ClipDirectory - The clip / output directory
var ClipDirectory string

// MaxClipSize - Maximum size for a clip, in bytes
var MaxClipSize int64 = 0

// GetClipByID - Returns a clip file by its ID
func GetClipByID(id string) (ClipFile, error) {
	var clipFile ClipFile

	clipFileName := path.Join(ClipDirectory, id)
	clipFileStat, err := os.Stat(clipFileName)
	if err == nil {
		if !clipFileStat.IsDir() {
			clipMetaFileName := path.Join(ClipDirectory, id+".meta")
			clipMetaFileStat, err := os.Stat(clipMetaFileName)
			if err == nil {
				if !clipMetaFileStat.IsDir() {
					clipFile.file = clipFileName
					clipFile.fileInfo = clipFileStat
					clipFile.id = id
					clipFile.metaFile = clipFileName
					clipFile.metaFileInfo = clipMetaFileStat

					err = nil
				} else {
					err = errors.New("Clip meta is not file")
				}
			}
		} else {
			err = errors.New("Clip is not file")
		}
	}

	return clipFile, err
}

// ScanClipDirectory - Scans clip directory for clip files
func ScanClipDirectory() ([]ClipFile, error) {
	var files []ClipFile

	// try scan directory for ".meta" files
	err := filepath.Walk(ClipDirectory, func(metaFilePath string, metaFileStat os.FileInfo, err error) error {
		if err == nil {
			if !metaFileStat.IsDir() {
				if strings.HasSuffix(metaFilePath, ".meta") {
					metaFileName := path.Base(metaFilePath)

					fileName := metaFileName
					fileName = fileName[0 : len(fileName)-5]

					doesMatch, _ := regexp.MatchString("^([0-9a-f]{32})$", fileName)
					if doesMatch {
						filePath := metaFilePath[0 : len(metaFilePath)-5]

						fileStat, err := os.Stat(filePath)
						if err == nil {
							if !fileStat.IsDir() {
								var newFileItem ClipFile
								newFileItem.file = filePath
								newFileItem.fileInfo = fileStat
								newFileItem.id = fileName
								newFileItem.metaFile = metaFilePath
								newFileItem.metaFileInfo = metaFileStat

								files = append(files, newFileItem)
							}
						}
					}
				}
			}
		}

		return nil
	})

	sort.Sort(ByNewestClipFileList(files))

	return files, err
}
