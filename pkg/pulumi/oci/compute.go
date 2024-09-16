package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/core"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {

	compartments, err := identity.GetCompartments(ctx, &identity.GetCompartmentsArgs{
		CompartmentId: metadata.Account, // Replace with your parent compartment OCID (root or any other)
	})
	if err != nil {
		return err
	}

	var compartment identity.GetCompartmentsCompartment
	// Search for a specific compartment by name
	for _, compartment = range compartments.Compartments {
		if compartment.Name == args["name"] {
			break
		}
	}

	vcn, err := core.NewVcn(ctx, args["name"].(string), &core.VcnArgs{
		CidrBlock:     pulumi.String("10.0.0.0/16"),
		CompartmentId: pulumi.String(compartment.Id),
		DisplayName:   pulumi.String(args["name"].(string)),
	})
	if err != nil {
		return err
	}

	subnet, err := core.NewSubnet(ctx, args["name"].(string), &core.SubnetArgs{
		VcnId:         vcn.ID(),
		CidrBlock:     pulumi.String("10.0.1.0/24"),
		CompartmentId: pulumi.String(compartment.Id),
		DisplayName:   pulumi.String(args["name"].(string)),
	})
	if err != nil {
		return err
	}

	_, err = core.NewInstance(ctx, args["name"].(string), &core.InstanceArgs{
		AvailabilityDomain: pulumi.String(metadata.Region),
		CompartmentId:      pulumi.String(compartment.Id),
		DisplayName:        pulumi.String(args["name"].(string)),
		Shape:              pulumi.String(args["type"].(string)),
		CreateVnicDetails: &core.InstanceCreateVnicDetailsArgs{
			SubnetId: subnet.ID(),
		},
		SourceDetails: &core.InstanceSourceDetailsArgs{
			SourceType:          pulumi.String("image"),
			SourceId:            pulumi.String(args["image"].(string)),
			BootVolumeSizeInGbs: pulumi.String(args["size"].(string)),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
