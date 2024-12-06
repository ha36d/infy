package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/core"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Network(metadata *model.Metadata, ctx *pulumi.Context, tracker *model.ResourceTracker) error {
	compartment, err := identity.NewCompartment(ctx, metadata.Meta["name"], &identity.CompartmentArgs{
		CompartmentId: pulumi.String(metadata.Account),
		Description:   pulumi.String(metadata.Meta["name"]),
		Name:          pulumi.String(metadata.Meta["name"]),
	})
	if err != nil {
		return err
	}

	vcn, err := core.NewVcn(ctx, metadata.Meta["name"], &core.VcnArgs{
		CidrBlock:     pulumi.String("10.0.0.0/16"),
		CompartmentId: compartment.CompartmentId,
		DisplayName:   pulumi.String(metadata.Meta["name"]),
	}, pulumi.DependsOn([]pulumi.Resource{compartment}))
	if err != nil {
		return err
	}

	subnet, err := core.NewSubnet(ctx, metadata.Meta["name"], &core.SubnetArgs{
		VcnId:         vcn.ID(),
		CidrBlock:     pulumi.String("10.0.1.0/24"),
		CompartmentId: compartment.CompartmentId,
		DisplayName:   pulumi.String(metadata.Meta["name"]),
	}, pulumi.DependsOn([]pulumi.Resource{vcn}))
	if err != nil {
		return err
	}

	// Export VPC and subnet information
	ctx.Export("vcnId", vcn.ID())
	ctx.Export("subnetId", subnet.ID())
	tracker.AddResource("network", metadata.Meta["name"], vcn)
	return nil
}
