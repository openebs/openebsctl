/*
Copyright 2020 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog"
)

const day = time.Minute * 60 * 24

// Fatal prints the message (if provided) and then exits. If V(2) or greater,
// klog.Fatal is invoked for extended information.
func Fatal(msg string) {
	if klog.V(2) {
		klog.FatalDepth(2, msg)
	}
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(1)
}

func Duration(d time.Duration) string {
	days := d / (time.Hour * 24)
	hours := d % (time.Hour * 24) / (time.Hour)
	mins := d % (time.Hour * 24) % (time.Hour) / (time.Minute)
	secs := d % (time.Hour * 24) % (time.Hour) % (time.Minute) / (time.Second)
	age := ""
	if days != 0 {
		age = age + strconv.Itoa(int(days)) + "d"
	}
	if hours != 0 {
		age = age + strconv.Itoa(int(hours)) + "h"
	}
	if mins != 0 {
		age = age + strconv.Itoa(int(mins)) + "m"
	}
	if secs != 0 {
		age = age + strconv.Itoa(int(secs)) + "s"
	}
	return age
}
