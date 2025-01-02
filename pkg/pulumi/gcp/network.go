package gcppulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var subnet *compute.Subnetwork

func (Holder) Network(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	network, err := compute.NewNetwork(ctx, metadata.Meta["name"]+"-network", &compute.NetworkArgs{
		AutoCreateSubnetworks: pulumi.Bool(false),
	})
	if err != nil {
		return err
	}

	// 2. Create a new subnet within the network
	subnet, err := compute.NewSubnetwork(ctx, metadata.Meta["name"]+"-subnet", &compute.SubnetworkArgs{
		IpCidrRange: pulumi.String(args["cidr"].(string)),
		Region:      pulumi.String(metadata.Region),
		Network:     network.ID(),
	})
	if err != nil {
		return err
	}

	// Convert []interface{} to []string
	tcpPorts := make([]string, 0)
	udpPorts := make([]string, 0)
	icmpPorts := make([]string, 0)
	// Separate ports by protocol (format: "tcp:80", "udp:53")
	for _, v := range args["ports"].([]interface{}) {
		portStr := v.(string)
		if len(portStr) > 4 && portStr[:4] == "tcp:" {
			tcpPorts = append(tcpPorts, portStr[4:])
		} else if len(portStr) > 4 && portStr[:4] == "udp:" {
			udpPorts = append(udpPorts, portStr[4:])
		} else if v.(string) == "icmp" {
			icmpPorts = append(icmpPorts, "icmp")
		}
	}

	// Build the firewall allows array dynamically
	var allows compute.FirewallAllowArray

	// Add TCP rules if TCP ports exist
	if len(tcpPorts) > 0 {
		allows = append(allows, &compute.FirewallAllowArgs{
			Protocol: pulumi.String("tcp"),
			Ports:    pulumi.ToStringArray(tcpPorts),
		})
	}

	// Add UDP rules if UDP ports exist
	if len(udpPorts) > 0 {
		allows = append(allows, &compute.FirewallAllowArgs{
			Protocol: pulumi.String("udp"),
			Ports:    pulumi.ToStringArray(udpPorts),
		})
	}

	// Add ICMP rule if specified in the ports array
	if len(icmpPorts) > 0 {
		allows = append(allows, &compute.FirewallAllowArgs{
			Protocol: pulumi.String("icmp"),
		})
	}

	// Only create firewall if there are rules to apply
	if len(allows) > 0 {
		firewall, err := compute.NewFirewall(ctx, metadata.Meta["name"]+"-firewall", &compute.FirewallArgs{
			Network: network.ID(),
			Allows:  allows,
			SourceRanges: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("firewallName", firewall.Name)
	}

	ctx.Export("networkName", network.Name)
	ctx.Export("subnetName", subnet.Name)
	tracker.AddResource("network", metadata.Meta["name"], network)
	return nil
}
