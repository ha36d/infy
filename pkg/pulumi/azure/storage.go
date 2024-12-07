package azurepulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {
	// here we create the bucket
	rg, err := core.LookupResourceGroup(ctx, &core.LookupResourceGroupArgs{
		Name: metadata.Meta["name"],
	}, nil)
	if err != nil {
		return err
	}
	account, err := storage.NewAccount(ctx, args["name"].(string), &storage.AccountArgs{
		Name:                   pulumi.String(args["name"].(string)),
		ResourceGroupName:      pulumi.String(rg.Name),
		Location:               pulumi.String(rg.Location),
		AccountTier:            pulumi.String("Standard"),
		AccountReplicationType: pulumi.String("GRS"),
		Tags:                   utils.Labels(metadata),
	})
	if err != nil {
		return err
	}
	tracker.AddResource("storage", metadata.Meta["name"], account)
	return nil
}
