package ocipulumi

import (
	"testing"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestHolder_Network(t *testing.T) {
	type args struct {
		metadata *model.Metadata
		ctx      *pulumi.Context
		tracker  *model.ResourceTracker
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
			if err := h.Network(tt.args.metadata, tt.args.ctx, tt.args.tracker); (err != nil) != tt.wantErr {
				t.Errorf("Holder.Network() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
