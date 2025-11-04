// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package core

import (
	"log"
	"runtime"

	"github.com/containifyci/go-self-update/pkg/updater"
	"github.com/thatgurkangurk/packwiz-installer/pkg/build"
)

func Update() {
	// Determine correct binary name for current OS/arch
	binaryName := "packwiz-installer_{{.OS}}_{{.Arch}}"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	u := updater.NewUpdater(
		binaryName,
		"thatgurkangurk",    // GitHub repo owner
		"packwiz-installer", // GitHub repo name
		build.Version,       // Current version
	)

	updated, err := u.SelfUpdate()
	if err != nil {
		log.Fatalf("failed to update: %v", err)
	}

	if updated {
		log.Println("updated to the latest version!")
	} else {
		log.Println("you're already on the latest version.")
	}
}
