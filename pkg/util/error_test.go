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
	"testing"

	"github.com/pkg/errors"
)

func TestCheckErr(t *testing.T) {
	type args struct {
		err       error
		handleErr func(string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Error not nil",
			args{
				err: errors.New("Some error occurred"),
				handleErr: func(s string) {
					fmt.Println("Handled")
				},
			},
		},
		{
			"Error nil",
			args{
				err: nil,
				handleErr: func(s string) {
					fmt.Println("Handled")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckErr(tt.args.err, tt.args.handleErr)
		})
	}
}

func TestCheckErrDefault(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Error not nil",
			args{
				err: errors.New("Some error occurred"),
				msg: "Error message",
			},
		},
		{
			"Error nil",
			args{
				err: nil,
				msg: "Error message",
			},
		},
		{
			"Message not empty",
			args{
				err: errors.New("Some error occurred"),
				msg: "Error message",
			},
		},
		{
			"Message empty",
			args{
				err: errors.New("Some error occurred"),
				msg: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckErrDefault(tt.args.err, tt.args.msg)
		})
	}
}
