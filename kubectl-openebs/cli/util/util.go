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

const maxTerms = 2

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
// Duration return the time.Duration in no.of days,hour, mins, seconds format.
// The number of terms to be shown can be increased or decreased using maxTerms constant.
func Duration(d time.Duration) string {
	days := d / (time.Hour * 24)
	hours := d % (time.Hour * 24) / (time.Hour)
	mins := d % (time.Hour * 24) % (time.Hour) / (time.Minute)
	secs := d % (time.Hour * 24) % (time.Hour) % (time.Minute) / (time.Second)
	age := ""
	currentTerms := 0
	if days != 0 {
		age = age + strconv.Itoa(int(days)) + "d"
		currentTerms++
	}
	if hours != 0 && currentTerms < maxTerms{
		age = age + strconv.Itoa(int(hours)) + "h"
		currentTerms++
	}
	if mins != 0 && currentTerms < maxTerms{
		age = age + strconv.Itoa(int(mins)) + "m"
		currentTerms++
	}
	if secs != 0 && currentTerms < maxTerms{
		age = age + strconv.Itoa(int(secs)) + "s"
	}
	return age
}
