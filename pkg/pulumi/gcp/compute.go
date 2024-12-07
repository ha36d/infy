package gcppulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	available, err := compute.GetZones(ctx, &compute.GetZonesArgs{Status: pulumi.StringRef("UP")})
	if err != nil {
		return err
	}

	instance, err := compute.NewInstance(ctx, strings.ToLower(args["name"].(string)), &compute.InstanceArgs{
		NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
			&compute.InstanceNetworkInterfaceArgs{
				AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
					nil,
				},
				Subnetwork: subnet.ID(),
			},
		},
		Name:        pulumi.String(strings.ToLower(args["name"].(string))),
		MachineType: pulumi.String(args["type"].(string)),
		Zone:        pulumi.String(available.Names[0]),
		Labels:      utils.Labels(metadata),
		BootDisk: &compute.InstanceBootDiskArgs{
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String(args["image"].(string)),
				Size:  pulumi.Int(args["size"].(int)),
			},
		},
	})
	if err != nil {
		return err
	}
	tracker.AddResource("compute", metadata.Meta["name"], instance)
	return nil
}
