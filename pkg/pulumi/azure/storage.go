package azurepulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {
	// here we create the bucket
	rg, err := core.LookupResourceGroup(ctx, &core.LookupResourceGroupArgs{
		Name: metadata.Name,
	}, nil)
	if err != nil {
		return err
	}
	_, err = storage.NewAccount(ctx, args["name"].(string), &storage.AccountArgs{
		Name:                   pulumi.String(args["name"].(string)),
		ResourceGroupName:      pulumi.String(rg.Name),
		Location:               pulumi.String(rg.Location),
		AccountTier:            pulumi.String("Standard"),
		AccountReplicationType: pulumi.String("GRS"),
		Tags: pulumi.StringMap{
			"team":    pulumi.String(strings.ToLower(metadata.Team)),
			"product": pulumi.String(strings.ToLower(metadata.Name)),
			"owner":   pulumi.String(strings.ToLower(metadata.Team)),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
