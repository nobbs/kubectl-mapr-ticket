// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package util

import (
	"encoding/base64"
	"os"
)

// DecodeBase64 decodes a base64 encoded string, returning the decoded bytes or an error
// if the string could not be decoded.
func DecodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// ReadFile reads the file at the given path and returns the contents of the file or an error
// if the file could not be read.
func ReadFile(path string) ([]byte, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return body, nil
}
