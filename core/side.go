// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package core

// IsValid returns whether the Side is one of client, server or both
func (s Side) IsValid() bool {
	return s == Side_Client || s == Side_Server || s == Side_Both
}

// ShouldInstall returns whether a mod with the given side should be installed
// based on the target installation side.
func (s Side) ShouldInstall(modSide Side) bool {
	switch s {
	case Side_Client:
		return modSide == Side_Client || modSide == Side_Both || modSide == ""
	case Side_Server:
		return modSide == Side_Server || modSide == Side_Both || modSide == ""
	case Side_Both:
		return true
	default:
		return false
	}
}
