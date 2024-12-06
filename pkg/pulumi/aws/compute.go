package awspulumi

import (
	"fmt"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {

	available, err := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{
		State: pulumi.StringRef("available"),
		Filters: []aws.GetAvailabilityZonesFilter{
			{
				Name: "region-name",
				Values: []string{
					metadata.Region,
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}
	image, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
		MostRecent: pulumi.BoolRef(true),
		Filters: []ec2.GetAmiFilter{
			{
				Name: "name",
				Values: []string{
					args["image"].(string),
				},
			},
			{
				Name: "virtualization-type",
				Values: []string{
					"hvm",
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	networkResource, exists := tracker.GetResource("network", metadata.Meta["name"])
	if !exists {
		return fmt.Errorf("network not found")
	}

	instance, err := ec2.NewInstance(ctx, args["name"].(string), &ec2.InstanceArgs{
		Ami:          pulumi.String(image.Id),
		SubnetId:     networkResource.(*ec2.Subnet).ID(),
		InstanceType: pulumi.String(args["type"].(string)),
		Tags:         utils.Labels(metadata),
		EbsBlockDevices: ec2.InstanceEbsBlockDeviceArray{
			&ec2.InstanceEbsBlockDeviceArgs{
				DeviceName: pulumi.String("/dev/xvdb"),
				VolumeSize: pulumi.Int(args["size"].(int)),
				Tags:       utils.Labels(metadata),
			},
		},
		AvailabilityZone: pulumi.String(available.Names[0]),
	})
	if err != nil {
		return err
	}
	tracker.AddResource("compute", metadata.Meta["name"], instance)
	return nil
}
