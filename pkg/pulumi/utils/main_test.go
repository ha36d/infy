package utils

import (
	"reflect"
	"testing"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestLabels(t *testing.T) {
	type args struct {
		metadata *model.Metadata
	}
	tests := []struct {
		name string
		args args
		want pulumi.StringMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Labels(tt.args.metadata); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Labels() = %v, want %v", got, tt.want)
			}
		})
	}
}
