package awspulumi

import (
	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var privateSubnet *ec2.Subnet

func (Holder) Network(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {
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
	// VPC
	vpc, err := ec2.NewVpc(ctx, metadata.Meta["Name"]+"-vpc", &ec2.VpcArgs{
		CidrBlock:          pulumi.String("10.0.0.0/16"), // VPC CIDR block
		EnableDnsSupport:   pulumi.Bool(true),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("my-vpc"),
		},
	})
	if err != nil {
		return err
	}

	// Internet Gateway for the VPC
	igw, err := ec2.NewInternetGateway(ctx, metadata.Meta["Name"]+"-igw", &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("my-igw"),
		},
	})
	if err != nil {
		return err
	}

	// Public Subnet
	publicSubnet, err := ec2.NewSubnet(ctx, metadata.Meta["Name"]+"-public-subnet", &ec2.SubnetArgs{
		VpcId:               vpc.ID(),
		CidrBlock:           pulumi.String("10.0.1.0/24"),
		AvailabilityZone:    pulumi.String(available.Names[0]),
		MapPublicIpOnLaunch: pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("public-subnet"),
		},
	})
	if err != nil {
		return err
	}

	// Private Subnet
	privateSubnet, err := ec2.NewSubnet(ctx, metadata.Meta["Name"]+"-private-subnet", &ec2.SubnetArgs{
		VpcId:            vpc.ID(),
		CidrBlock:        pulumi.String("10.0.2.0/24"),
		AvailabilityZone: pulumi.String(available.Names[0]),
		Tags: pulumi.StringMap{
			"Name": pulumi.String("private-subnet"),
		},
	})
	if err != nil {
		return err
	}

	// Route Table for the Public Subnet and associate it with the Internet Gateway
	publicRouteTable, err := ec2.NewRouteTable(ctx, metadata.Meta["Name"]+"-public-route-table", &ec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"), // Default route
				GatewayId: igw.ID(),
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String("public-route-table"),
		},
	})
	if err != nil {
		return err
	}

	// Associate the public subnet with the public route table
	_, err = ec2.NewRouteTableAssociation(ctx, metadata.Meta["Name"]+"-public-subnet-association", &ec2.RouteTableAssociationArgs{
		SubnetId:     publicSubnet.ID(),
		RouteTableId: publicRouteTable.ID(),
	})
	if err != nil {
		return err
	}

	securityGroup, err := ec2.NewSecurityGroup(ctx, metadata.Meta["Name"]+"-sg", &ec2.SecurityGroupArgs{
		VpcId: vpc.ID(),
		Ingress: ec2.SecurityGroupIngressArray{
			&ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(22), // SSH
				ToPort:     pulumi.Int(22),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
			&ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80), // HTTP
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})
	if err != nil {
		return err
	}

	// Export VPC and subnet information
	ctx.Export("vpcId", vpc.ID())
	ctx.Export("publicSubnetId", publicSubnet.ID())
	ctx.Export("privateSubnetId", privateSubnet.ID())
	ctx.Export("internetGatewayId", igw.ID())
	ctx.Export("securityGroupId", securityGroup.ID())

	return nil
}
