package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compartment(metadata *model.Metadata, ctx *pulumi.Context) error {

	_, err := identity.NewCompartment(ctx, "exampleCompartment", &identity.CompartmentArgs{
		CompartmentId: pulumi.String(metadata.Account), // Root or parent compartment OCID
		Description:   pulumi.String(metadata.Name),
		Name:          pulumi.String(metadata.Name),
	})
	if err != nil {
		return err
	}

	return nil
}
