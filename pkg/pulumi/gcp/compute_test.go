package gcppulumi

import (
	"testing"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestHolder_Compute(t *testing.T) {
	type args struct {
		metadata *model.Metadata
		args     map[string]any
		ctx      *pulumi.Context
	}
	tests := []struct {
		name    string
		h       Holder
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Holder{}
			if err := h.Compute(tt.args.metadata, tt.args.args, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Holder.Compute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
