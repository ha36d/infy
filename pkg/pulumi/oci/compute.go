package ocipulumi

import (
	"fmt"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/core"
	"github.com/pulumi/pulumi-oci/sdk/v2/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	availabilityDomains, err := identity.GetAvailabilityDomains(ctx, &identity.GetAvailabilityDomainsArgs{
		CompartmentId: metadata.Account,
	})
	if err != nil {
		return err
	}

	compartmentResource, exists := tracker.GetResource("compartment", metadata.Meta["name"])
	if !exists {
		return fmt.Errorf("compartment not found")
	}

	networkResource, exists := tracker.GetResource("network", metadata.Meta["name"])
	if !exists {
		return fmt.Errorf("network not found")
	}

	instance, err := core.NewInstance(ctx, args["name"].(string), &core.InstanceArgs{
		AvailabilityDomain: pulumi.String(availabilityDomains.AvailabilityDomains[0].Name),
		CompartmentId:      compartmentResource.(*identity.Compartment).ID(),
		DisplayName:        pulumi.String(args["name"].(string)),
		Shape:              pulumi.String(args["type"].(string)),
		CreateVnicDetails: &core.InstanceCreateVnicDetailsArgs{
			SubnetId:    networkResource.(*core.Vcn).ID(),
			DisplayName: pulumi.String(args["name"].(string)),
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
	tracker.AddResource("compute", metadata.Meta["name"], instance)
	return nil
}
