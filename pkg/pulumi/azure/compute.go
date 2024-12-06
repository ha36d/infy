package azurepulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/compute"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v5/go/azure/network"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {
	// here we create the server
	os := strings.Split(args["image"].(string), "/")
	rg, err := core.LookupResourceGroup(ctx, &core.LookupResourceGroupArgs{
		Name: metadata.Meta["name"],
	}, nil)
	if err != nil {
		return err
	}
	nic, err := network.NewNetworkInterface(ctx, "nic", &network.NetworkInterfaceArgs{
		Name:              pulumi.String(args["name"].(string)),
		Location:          pulumi.String(rg.Location),
		ResourceGroupName: pulumi.String(rg.Name),
		IpConfigurations: network.NetworkInterfaceIpConfigurationArray{
			&network.NetworkInterfaceIpConfigurationArgs{
				Name:                       pulumi.String(args["name"].(string)),
				SubnetId:                   subnet.ID(),
				PrivateIpAddressAllocation: pulumi.String("Dynamic"),
			},
		},
	})
	if err != nil {
		return err
	}
	vm, err := compute.NewVirtualMachine(ctx, "main", &compute.VirtualMachineArgs{
		Name:              pulumi.String(args["name"].(string)),
		Location:          pulumi.String(rg.Location),
		ResourceGroupName: pulumi.String(rg.Name),
		NetworkInterfaceIds: pulumi.StringArray{
			nic.ID(),
		},
		VmSize: pulumi.String(args["type"].(string)),
		StorageImageReference: &compute.VirtualMachineStorageImageReferenceArgs{
			Publisher: pulumi.String(os[0]),
			Offer:     pulumi.String(os[1]),
			Sku:       pulumi.String(os[2]),
			Version:   pulumi.String("latest"),
		},
		StorageOsDisk: &compute.VirtualMachineStorageOsDiskArgs{
			Name:            pulumi.String("myosdisk"),
			Caching:         pulumi.String("ReadWrite"),
			CreateOption:    pulumi.String("FromImage"),
			ManagedDiskType: pulumi.String("Standard_LRS"),
			DiskSizeGb:      pulumi.Int(args["size"].(int)),
		},
		OsProfile: &compute.VirtualMachineOsProfileArgs{
			ComputerName:  pulumi.String("hostname"),
			AdminUsername: pulumi.String("testadmin"),
			AdminPassword: pulumi.String("Password1234!"),
		},
		OsProfileLinuxConfig: &compute.VirtualMachineOsProfileLinuxConfigArgs{
			DisablePasswordAuthentication: pulumi.Bool(false),
		},
		Tags: utils.Labels(metadata),
	})
	if err != nil {
		return err
	}
	tracker.AddResource("compute", metadata.Meta["name"], vm)
	return nil
}
