package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getMetrics(t *testing.T) {
	type args struct {
		pollCount Counter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive case #1",
			args: args{pollCount: Counter(123)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, getMetrics(tt.args.pollCount), "getMetrics(%v)", tt.args.pollCount)
		})
	}
}
