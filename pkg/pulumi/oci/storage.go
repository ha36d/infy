package ocipulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/objectstorage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	namespace, err := objectstorage.GetNamespace(ctx, nil)
	if err != nil {
		return err
	}

	// Create an Object Storage bucket
	bucket, err := objectstorage.NewBucket(ctx, args["name"].(string), &objectstorage.BucketArgs{
		CompartmentId: pulumi.String("ocid1.compartment.oc1..aaaaaaaam7kgdpj5nw3hxeuofrbwzz7rozvywdtszgumosddhi3r4fqot3jq"),
		Name:          pulumi.String(args["name"].(string)),
		Namespace:     pulumi.String(namespace.Namespace),
		StorageTier:   pulumi.String("Standard"),
	})
	if err != nil {
		return err
	}

	tracker.AddResource("storage", metadata.Meta["name"], bucket)
	return nil
}
