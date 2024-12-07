package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runPreview(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runPreview(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runPreview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePreviewConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePreviewConfig(); (err != nil) != tt.wantErr {
				t.Errorf("validatePreviewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logPreviewParameters(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logPreviewParameters()
		})
	}
}
