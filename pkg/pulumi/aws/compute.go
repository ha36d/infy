package awspulumi

import (
	"fmt"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {

	stackRef, err := pulumi.NewStackReference(ctx, fmt.Sprintf("%s-%s-%s", metadata.Meta["Team"], metadata.Meta["Name"], metadata.Meta["Env"]), nil)
	if err != nil {
		return err
	}

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
	_, err = ec2.NewInstance(ctx, args["name"].(string), &ec2.InstanceArgs{
		Ami: pulumi.String(image.Id),
		SubnetId: stackRef.GetOutput(pulumi.String("privateSubnetId")).ApplyT(func(id interface{}) *string {
			strId := id.(string)
			return &strId
		}).(pulumi.StringPtrOutput),
		InstanceType: pulumi.String(args["type"].(string)),
		Tags:         utils.StringMapLabels(metadata),
		EbsBlockDevices: ec2.InstanceEbsBlockDeviceArray{
			&ec2.InstanceEbsBlockDeviceArgs{
				DeviceName: pulumi.String("/dev/xvdb"),
				VolumeSize: pulumi.Int(args["size"].(int)),
				Tags:       utils.StringMapLabels(metadata),
			},
		},
		AvailabilityZone: pulumi.String(available.Names[0]),
	})
	if err != nil {
		return err
	}

	return nil
}
