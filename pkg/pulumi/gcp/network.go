package gcppulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var subnet *compute.Subnetwork

func (Holder) Requirements(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	network, err := compute.NewNetwork(ctx, metadata.Meta["name"]+"-network", &compute.NetworkArgs{
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return err
	}

	// 2. Create a new subnet within the network
	subnet, err := compute.NewSubnetwork(ctx, metadata.Meta["name"]+"-subnet", &compute.SubnetworkArgs{
		IpCidrRange: pulumi.String(""),
		Region:      pulumi.String("us-central1"),
		Network:     network.ID(),
	})
	if err != nil {
		return err
	}

	firewall, err := compute.NewFirewall(ctx, metadata.Meta["name"]+"-firewall", &compute.FirewallArgs{
		Network: network.ID(),
		Allows: compute.FirewallAllowArray{
			&compute.FirewallAllowArgs{
				Protocol: pulumi.String("tcp"),
				Ports: pulumi.StringArray{
					pulumi.String("22"),
					pulumi.String("80"),
				},
			},
			&compute.FirewallAllowArgs{
				Protocol: pulumi.String("icmp"),
			},
		},
		SourceRanges: pulumi.StringArray{
			pulumi.String("0.0.0.0/0"),
		},
	})

	if err != nil {
		return err
	}

	ctx.Export("networkName", network.Name)
	ctx.Export("subnetName", subnet.Name)
	ctx.Export("firewallName", firewall.Name)
	tracker.AddResource("network", metadata.Meta["name"], network)

	return nil

}
