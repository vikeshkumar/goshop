package main

import (
	"testing"
)

func Test_bootstrap(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		error bool
	}{
		{"able to read configuration", "config/local.yaml", false},
		{"not able to read configuration", "config/local2.yaml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bootstrap(tt.path)
			if tt.error {
				if err == nil {
					t.Errorf("expected nil configuration and error")
				}
			}
			if !tt.error {
				if err != nil {
					t.Errorf("expected configuration and no error")
				}
			}
		})
	}
}
