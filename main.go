// Licensed to Adam Shannon under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamdecaf/windows-startstop-text/internal/sms"
)

var (
	flagConfigFilepath = flag.String("config", "examples/config.json", "Filepath of JSON config file")

	flagShutdown = flag.Bool("shutdown", false, "Send notification about system shutdown")
	flagStartup  = flag.Bool("startup", false, "Send notification about system startup")
)

type Config struct {
	Twilio sms.Config  `json:"twilio"`
	SMS    sms.Message `json:"sms"`
}

func main() {
	flag.Parse()

	log.Printf("DEBUG: starting windows-startstop-text version %s", Version)

	conf, err := readConfig(*flagConfigFilepath)
	if err != nil {
		log.Printf("ERROR: reading %s failed: %v", *flagConfigFilepath, err)
		os.Exit(1)
	}
	fmt.Printf("config: %#v\n", conf)

	// TODO(adam): Can we handle Windows event 1074 (and whatever ID for startup/boot)
	//
	// Maybe with https://pkg.go.dev/golang.org/x/sys@v0.13.0/windows#Handle
	// Examples: https://gist.github.com/glennswest/3d4bf8448193d6868baf0665a6fb1c5a

	when := time.Now()
	hostname, _ := os.Hostname()

	switch {
	case *flagShutdown:
		conf.SMS.Body = fmt.Sprintf("System %s is SHUTTING down at %v", hostname, when)
	case *flagStartup:
		conf.SMS.Body = fmt.Sprintf("System %s is STARTING up at %v", hostname, when)
	default:
		log.Print("ERROR: no action specified")
		os.Exit(1)
	}

	err = sms.Send(conf.Twilio, conf.SMS)
	if err != nil {
		log.Printf("ERROR sending sms failed: %v", err)
		os.Exit(1)
	}
}

func readConfig(where string) (Config, error) {
	var conf Config

	fd, err := os.Open(where)
	if err != nil {
		return conf, err
	}
	defer fd.Close()

	err = json.NewDecoder(fd).Decode(&conf)
	return conf, err
}
