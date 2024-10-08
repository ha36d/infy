package azurepulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Resourcegroup(metadata *model.Metadata, ctx *pulumi.Context) error {
	// here we create the bucket
	_, err := core.NewResourceGroup(ctx, metadata.Name, &core.ResourceGroupArgs{
		Name:     pulumi.String(metadata.Name),
		Location: pulumi.String(metadata.Region),
	})
	if err != nil {
		return err
	}

	return nil
}
