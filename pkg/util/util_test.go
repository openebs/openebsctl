/*
Copyright 2020-2022 The OpenEBS Authors

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

func TestConvertToIBytes(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test with GB values",
			args{value: "1.65GB"},
			"1.5GiB",
		},
		{
			"Test with MB values",
			args{value: "1.65MB"},
			"1.6MiB",
		},
		{
			"Test with KB values",
			args{value: "1.65KB"},
			"1.6KiB",
		},
		{
			"Test with K values",
			args{value: "1.65K"},
			"1.6KiB",
		},
		{
			"Test with M values",
			args{value: "1.65M"},
			"1.6MiB",
		},
		{
			"Test with MiB values",
			args{value: "1.65MiB"},
			"1.6MiB",
		},
		{
			"Test with Mi values",
			args{value: "1.65Mi"},
			"1.6MiB",
		},
		{
			"Test with invalid",
			args{value: ""},
			"",
		},
		{
			"Test with only numeric value",
			args{value: "1766215"},
			"1.7MiB",
		},
		{
			"Test with only invalid unit",
			args{value: "1766215CiB"},
			"1766215CiB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToIBytes(tt.args.value); got != tt.want {
				t.Errorf("ConvertToIBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsedPercentage(t *testing.T) {
	type args struct {
		total string
		used  string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"Test with same units",
			args{total: "12 GiB", used: "1 GiB"},
			8.333333333333332,
		},
		{
			"Test with different units",
			args{total: "12 GiB", used: "100 MiB"},
			0.8138020833333334,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUsedPercentage(tt.args.total, tt.args.used); got != tt.want {
				t.Errorf("GetUsedPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvailableCapacity(t *testing.T) {
	type args struct {
		total string
		used  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Test with different units",
			args{
				total: "20GiB",
				used:  "655MiB",
			},
			"19.36GiB",
		},
		{
			"Test with same units with mutiple precisions",
			args{
				total: "21.66GiB",
				used:  "12.221GiB",
			},
			"9.439GiB",
		},
		{
			"Test with different units with mutiple precisions",
			args{
				total: "21.66GiB",
				used:  "12.221MiB",
			},
			"21.65GiB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAvailableCapacity(tt.args.total, tt.args.used); got != tt.want {
				t.Errorf("GetAvailableCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}
