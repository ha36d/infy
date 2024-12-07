package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Resourcegroup(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	var err error
	compartment, err := identity.NewCompartment(ctx, metadata.Meta["name"], &identity.CompartmentArgs{
		CompartmentId: pulumi.String(metadata.Account),
		Description:   pulumi.String(args["description"].(string)),
		Name:          pulumi.String(args["name"].(string)),
	})
	if err != nil {
		return err
	}

	ctx.Export("compartmentId", compartment.ID())
	tracker.AddResource("compartment", metadata.Meta["name"], compartment)
	return nil
}
