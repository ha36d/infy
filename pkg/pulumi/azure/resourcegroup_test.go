package azurepulumi

import (
	"testing"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestHolder_Resourcegroup(t *testing.T) {
	type args struct {
		metadata *model.Metadata
		ctx      *pulumi.Context
	}
	tests := []struct {
		name string
		h    Holder
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Holder{}
			h.Resourcegroup(tt.args.metadata, tt.args.ctx)
		})
	}
}
