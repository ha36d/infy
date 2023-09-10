package awspulumi

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Compute(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) {
	// here we create the server
	zone := []string{
		"a",
		"b",
		"c",
	}
	atoiName, _ := strconv.Atoi(args["name"].(string))
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
		log.Println(err)
	}
	_, err = ec2.NewInstance(ctx, "web", &ec2.InstanceArgs{
		Ami:          pulumi.String(image.Id),
		InstanceType: pulumi.String(args["type"].(string)),
		Tags: pulumi.StringMap{
			"team":    pulumi.String(strings.ToLower(metadata.Team)),
			"product": pulumi.String(strings.ToLower(metadata.Name)),
			"owner":   pulumi.String(strings.ToLower(metadata.Team)),
		},
		EbsBlockDevices: ec2.InstanceEbsBlockDeviceArray{
			&ec2.InstanceEbsBlockDeviceArgs{
				DeviceName: pulumi.String("/dev/xvdb"),
				VolumeSize: pulumi.Int(args["size"].(int)),
				Tags: pulumi.StringMap{
					"team":    pulumi.String(strings.ToLower(metadata.Team)),
					"product": pulumi.String(strings.ToLower(metadata.Name)),
					"owner":   pulumi.String(strings.ToLower(metadata.Team)),
				},
			},
		},
		AvailabilityZone: pulumi.String(metadata.Region + zone[atoiName%3]),
	})
	if err != nil {
		log.Println(err)
	}
}
