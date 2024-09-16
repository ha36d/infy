package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/objectstorage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {

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

	namespaceName := "<YOUR_OBJECT_STORAGE_NAMESPACE>"

	// Create an Object Storage bucket
	_, err = objectstorage.NewBucket(ctx, args["name"].(string), &objectstorage.BucketArgs{
		CompartmentId: pulumi.String(compartment.Id),
		Name:          pulumi.String(args["name"].(string)),
		Namespace:     pulumi.String(namespaceName),
		StorageTier:   pulumi.String("Standard"),
	})
	if err != nil {
		return err
	}

	return nil
}
