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
			"2d1m59s",
		},
		{
			"Duration in Months",
			args{d: (time.Hour * 24 * 30) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"30d1m59s",
		},
		{
			"Duration in Years",
			args{d: (time.Hour * 24 * 365) + time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"365d1m59s",
		},
		{
			"Duration in Minutes",
			args{d: time.Minute + (time.Second * 59) + (time.Millisecond * 300)},
			"1m59s",
		},
		{
			"Duration in Minutes",
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
