// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package version

import (
	"fmt"
	"runtime/debug"
)

const (
	// the version of the compiled binary
	version = "0.2.1" // x-release-please-version
)

// Version is a struct to hold the version information
type Version struct {
	Version string `json:"tag,omitempty"`    // the version tag of the release
	Commit  string `json:"commit,omitempty"` // the commit id of the release
	Dirty   bool   `json:"dirty,omitempty"`  // flag whether the binary was built from a dirty repository
}

// NewVersion returns a new version object
func NewVersion() *Version {
	id, dirty := CommitInfo()

	return &Version{
		Version: version,
		Commit:  id,
		Dirty:   dirty,
	}
}

// String returns the version string
func (v *Version) String() string {
	var suffix string
	if v.Dirty {
		suffix = "-dirty"
	}

	// if build info is not available, use the version only
	if v.Commit == "" {
		return fmt.Sprintf("v%s%s", v.Version, suffix)
	}

	return fmt.Sprintf("v%s-%s%s", v.Version, v.Commit[:7], suffix)
}

// CommitInfo returns the commit id and a flag if the repository is dirty
func CommitInfo() (id string, dirty bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			id = setting.Value
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}

	return id, dirty
}
