// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package desktop provides desktop integration utilities.
package desktop

import (
	"os/exec"
	"runtime"

	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/errs"
)

// OpenBrowser opens 'url' with the user's preferred browser.
func OpenBrowser(url string) error {
	var cmd string
	switch runtime.GOOS {
	case toolbox.MacOS:
		cmd = "open"
	case toolbox.LinuxOS:
		cmd = "xdg-open"
	case toolbox.WindowsOS:
		cmd = "explorer"
	default:
		return errs.New("Unsupported platform")
	}
	if err := exec.Command(cmd, url).Start(); err != nil {
		return errs.NewWithCause("Unable to open "+url, err)
	}
	return nil
}
