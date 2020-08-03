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
	"io"
	"net/http"
	"os"
)

// GetFileContentType - Returns the MIME type of a file
//
// https://golangcode.com/get-the-content-type-of-file/
func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

// MoveFile - Moves a file, which works also in Docker containers with mounted volumns
func MoveFile(src string, dest string) error {
	// open source for read
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// open destination for writing
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// copy from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// remove source file
	err = os.Remove(srcFile.Name())
	if err != nil {
		os.Remove(destFile.Name())

		return err
	}

	return nil
}
