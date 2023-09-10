package gcppulumi

import (
	"log"
	"strconv"
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) {
	// here we create the server
	zone := []string{
		"a",
		"b",
		"c",
	}
	atoiName, _ := strconv.Atoi(args["name"].(string))

	_, err := compute.NewInstance(ctx, strings.ToLower(args["name"].(string)), &compute.InstanceArgs{
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
		Zone:        pulumi.String(metadata.Region + "-" + zone[atoiName%3]),
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
		log.Println(err)
	}
}
