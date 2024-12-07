package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_runBuild(t *testing.T) {
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
			if err := runBuild(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("runBuild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateBuildConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBuildConfig(); (err != nil) != tt.wantErr {
				t.Errorf("validateBuildConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logBuildParameters(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuildParameters()
		})
	}
}
