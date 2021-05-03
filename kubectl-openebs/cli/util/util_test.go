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
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Duration in Days",
			args{d: (time.Hour * 24 * 2) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"2d1m",
		},
		{
			"Duration in Hour",
			args{d: (time.Hour * 24 * 2) + (time.Hour * 2) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"2d2h",
		},
		{
			"Duration in Months",
			args{d: (time.Hour * 24 * 30) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"30d1m",
		},
		{
			"Duration in Years",
			args{d: (time.Hour * 24 * 365) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"365d1m",
		},
		{
			"Duration in Minutes",
			args{d: time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"1m59s",
		},
		{
			"Duration in Seconds",
			args{d: time.Second*59 + (time.Millisecond * 300)},
			"59s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.d); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
