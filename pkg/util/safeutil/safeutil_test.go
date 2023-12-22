package safeutil

import (
	"testing"
	"time"

	"github.com/elliotxx/safe"
	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

func TestGo(t *testing.T) {
	type args struct {
		do safe.DoFunc
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Go(tt.args.do)
			time.Sleep(100 * time.Microsecond)
		})
	}
}

func TestGoL(t *testing.T) {
	type args struct {
		do     safe.DoFunc
		logger logr.Logger
	}

	getTestingLogger := func() logr.Logger {
		logger := klog.NewKlogr()

		return logger
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "successful-recover-crash-in-safe-Go",
			args: args{
				do: func() {
					panic("ah, I'm down")
				},
				logger: getTestingLogger(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoL(tt.args.do, tt.args.logger)
			time.Sleep(time.Second * 1)
		})
	}
}
