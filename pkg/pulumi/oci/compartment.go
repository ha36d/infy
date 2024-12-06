package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var compartment *identity.Compartment

func (Holder) Compartment(metadata *model.Metadata, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	var err error
	compartment, err = identity.NewCompartment(ctx, metadata.Meta["name"], &identity.CompartmentArgs{
		CompartmentId: pulumi.String(metadata.Account),
		Description:   pulumi.String(metadata.Meta["name"]),
		Name:          pulumi.String(metadata.Meta["name"]),
	})
	if err != nil {
		return err
	}

	// Wait until the compartment is created
	if err := ctx.RegisterResourceOutputs(compartment, pulumi.Map{}); err != nil {
		return err
	}
	ctx.Export("compartmentId", compartment.ID())
	tracker.AddResource("compartment", metadata.Meta["name"], compartment)
	return nil
}
