/*
Copyright 2020-2021 The OpenEBS Authors

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
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/pkg/errors"
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
	if hours != 0 && currentTerms < maxTerms {
		age = age + strconv.Itoa(int(hours)) + "h"
		currentTerms++
	}
	if mins != 0 && currentTerms < maxTerms {
		age = age + strconv.Itoa(int(mins)) + "m"
		currentTerms++
	}
	if secs != 0 && currentTerms < maxTerms {
		age = age + strconv.Itoa(int(secs)) + "s"
	}
	return age
}

// PrintByTemplate of the provided template and resource
func PrintByTemplate(templateName string, resourceTemplate string, resource interface{}) error {
	genericTemplate, err := template.New(templateName).Parse(resourceTemplate)
	if err != nil {
		return errors.Wrap(err, "error creating for "+templateName)
	}
	err = genericTemplate.Execute(os.Stdout, resource)
	if err != nil {
		return errors.Wrap(err, "error displaying by template for"+templateName)
	}
	return nil
}

// TablePrinter uses cli-runtime TablePrinter to create a similar UI for the ctl
func TablePrinter(columns []metav1.TableColumnDefinition, rows []metav1.TableRow, options printers.PrintOptions) {
	table := &metav1.Table{
		ColumnDefinitions: columns,
		Rows:              rows,
	}
	out := bytes.NewBuffer([]byte{})
	printer := printers.NewTablePrinter(options)
	_ = printer.PrintObj(table, out)
	fmt.Printf("%s", out.String())
}

// TemplatePrinter uses cli-runtime TemplatePrinter to print by template without extra type
func TemplatePrinter(template string, obj runtime.Object) {
	p, _ := printers.NewGoTemplatePrinter([]byte(template))
	p.AllowMissingKeys(true)
	buffer := &bytes.Buffer{}
	_ = p.PrintObj(obj, buffer)
	fmt.Println(buffer)
}

// ConvertToIBytes humanizes all the passed units to IBytes format
func ConvertToIBytes(value string) string {
	if value == "" {
		return value
	}
	iBytes, err := humanize.ParseBytes(value)
	if err != nil {
		return value
	}
	return humanize.IBytes(iBytes)
}
