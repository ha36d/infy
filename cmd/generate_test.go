package cmd

import "testing"

func Test_promptGetInput(t *testing.T) {
	type args struct {
		pc promptContent
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := promptGetInput(tt.args.pc); got != tt.want {
				t.Errorf("promptGetInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_promptGetChoice(t *testing.T) {
	type args struct {
		pc    promptContent
		items []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := promptGetChoice(tt.args.pc, tt.args.items); got != tt.want {
				t.Errorf("promptGetChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateYaml(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateYaml()
		})
	}
}
