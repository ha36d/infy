package azurepulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var subnet *network.Subnet

func (Holder) Network(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {
	// Virtual Network
	rg, err := core.LookupResourceGroup(ctx, &core.LookupResourceGroupArgs{
		Name: metadata.Name,
	}, nil)
	if err != nil {
		return err
	}
	net, err := network.NewVirtualNetwork(ctx, "net", &network.VirtualNetworkArgs{
		Name: pulumi.String(metadata.Name + "-network"),
		AddressSpaces: pulumi.StringArray{
			pulumi.String("10.0.0.0/16"),
		},
		Location:          pulumi.String(rg.Location),
		ResourceGroupName: pulumi.String(rg.Name),
	})
	if err != nil {
		return err
	}
	subnet, err := network.NewSubnet(ctx, "subnet", &network.SubnetArgs{
		Name:               pulumi.String(metadata.Name + "-subnet"),
		ResourceGroupName:  pulumi.String(rg.Name),
		VirtualNetworkName: net.Name,
		AddressPrefixes: pulumi.StringArray{
			pulumi.String("10.0.2.0/24"),
		},
	})
	if err != nil {
		return err
	}

	// Export VPC and subnet information
	ctx.Export("vnId", net.ID())
	ctx.Export("subnetId", subnet.ID())

	return nil
}
