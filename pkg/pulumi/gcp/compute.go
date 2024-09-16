package gcppulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {

	available, err := compute.GetZones(ctx, &compute.GetZonesArgs{Status: pulumi.StringRef("UP")})
	if err != nil {
		return err
	}

	_, err = compute.NewInstance(ctx, strings.ToLower(args["name"].(string)), &compute.InstanceArgs{
		NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
			&compute.InstanceNetworkInterfaceArgs{
				AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
					nil,
				},
				Network: pulumi.String("default"),
			},
		},
		Name:        pulumi.String(strings.ToLower(args["name"].(string))),
		MachineType: pulumi.String(args["type"].(string)),
		Zone:        pulumi.String(available.Names[0]),
		Labels: pulumi.StringMap{
			"team":    pulumi.String(strings.ToLower(metadata.Team)),
			"product": pulumi.String(strings.ToLower(metadata.Name)),
			"owner":   pulumi.String(strings.ToLower(metadata.Team)),
		},
		BootDisk: &compute.InstanceBootDiskArgs{
			InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
				Image: pulumi.String(args["image"].(string)),
				Labels: pulumi.Map{
					"team":    pulumi.String(strings.ToLower(metadata.Team)),
					"product": pulumi.String(strings.ToLower(metadata.Name)),
					"owner":   pulumi.String(strings.ToLower(metadata.Team)),
				},
				Size: pulumi.Int(args["size"].(int)),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
