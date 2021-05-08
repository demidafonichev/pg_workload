// Copyright 2021 demidafonichev. All rights reserved.
// Use of this source code is governed by Apache
// license that can be found in the LICENSE file.

// Package cli provides virtual command-line access
// in pgproxy include start,cli and stop action.
package cli

import (
	"os"
	"pgworkload/workload"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

const Logo = `
                               _    _                _ 
 _ __  __ _ __ __ __ ___  _ _ | |__| | ___  __ _  __| |
| '_ \/ _' |\ V  V // _ \| '_|| / /| |/ _ \/ _' |/ _' |
| .__/\__, | \_/\_/ \___/|_|  |_\_\|_|\___/\__,_|\__,_|
|_|   |___/
`

const (
	VERSION = "version-0.0.1"
)

// proxy server config struct
type ProxyConfig struct {
	ServerConfig struct {
		Host string
		Port string
	}
	DBConfig workload.DatabaseConfig `toml:"DB"`
}

func readConfig(file string) (pc ProxyConfig) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		glog.Errorln(err)
		os.Exit(int(syscall.ENOENT))
	}

	if _, err := toml.DecodeFile(file, &pc); err != nil {
		glog.Fatalln(err)
	}

	return
}
