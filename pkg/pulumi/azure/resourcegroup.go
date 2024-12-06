package azurepulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Resourcegroup(metadata *model.Metadata, ctx *pulumi.Context, tracker *model.ResourceTracker) error {
	// here we create the bucket
	rg, err := core.NewResourceGroup(ctx, metadata.Meta["name"], &core.ResourceGroupArgs{
		Name:     pulumi.String(metadata.Meta["name"]),
		Location: pulumi.String(metadata.Region),
	})
	if err != nil {
		return err
	}
	tracker.AddResource("resourcegroup", metadata.Meta["name"], rg)
	return nil
}
