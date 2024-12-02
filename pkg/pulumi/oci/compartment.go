package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/core"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var compartment *identity.Compartment

func (Holder) Compartment(metadata *model.Metadata, ctx *pulumi.Context) error {
	compartment, err := identity.NewCompartment(ctx, metadata.Meta["Name"], &identity.CompartmentArgs{
		CompartmentId: pulumi.String(metadata.Account),
		Description:   pulumi.String(metadata.Meta["Name"]),
		Name:          pulumi.String(metadata.Meta["Name"]),
	})

	if err != nil {
		return err
	}

	vcn, err = core.NewVcn(ctx, metadata.Meta["Name"], &core.VcnArgs{
		CidrBlock:     pulumi.String("10.0.0.0/16"),
		CompartmentId: compartment.CompartmentId,
		DisplayName:   pulumi.String(metadata.Meta["Name"]),
	}, pulumi.DependsOn([]pulumi.Resource{compartment}))

	if err != nil {
		return err
	}

	subnet, err = core.NewSubnet(ctx, metadata.Meta["Name"], &core.SubnetArgs{
		VcnId:         vcn.ID(),
		CidrBlock:     pulumi.String("10.0.1.0/24"),
		CompartmentId: compartment.CompartmentId,
		DisplayName:   pulumi.String(metadata.Meta["Name"]),
	}, pulumi.DependsOn([]pulumi.Resource{vcn}))

	if err != nil {
		return err
	}

	return nil
}
