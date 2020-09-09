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
	"text/tabwriter"
	"text/template"

	"github.com/ryanuber/columnize"
)

// FormatList takes a set of strings and formats them into properly
// aligned output, replacing any blank fields with a placeholder
// for awk-ability.
func FormatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

// Print binds the object with go template and executes it
func Print(format string, obj interface{}) error {
	// New Instance of tabwriter
	w := tabwriter.NewWriter(os.Stdout, MinWidth, MaxWidth, Padding, ' ', 0)
	// New Instance of template
	tmpl, err := template.New("Template").Parse(format)
	if err != nil {
		return fmt.Errorf("Error in parsing replica template, found error : %v", err)
	}
	// Parse struct with template
	err = tmpl.Execute(w, obj)
	if err != nil {
		return fmt.Errorf("Error in executing replica template, found error : %v", err)
	}
	return w.Flush()
}
