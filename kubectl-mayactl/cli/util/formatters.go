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
