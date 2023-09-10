package awspulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	log "github.com/sirupsen/logrus"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) {
	// here we create the bucket
	var acl string
	if args["acl"] != nil {
		acl = args["acl"].(string)
	} else {
		acl = "private"
	}
	bucket, err := s3.NewBucket(ctx, args["name"].(string), &s3.BucketArgs{
		Bucket: pulumi.String(args["name"].(string)),
		Acl:    pulumi.String(acl),
		Tags: pulumi.StringMap{
			"team":    pulumi.String(strings.ToLower(metadata.Team)),
			"product": pulumi.String(strings.ToLower(metadata.Name)),
			"owner":   pulumi.String(strings.ToLower(metadata.Team)),
		},
	})
	if err != nil {
		log.Errorf("abnormal termination: %s", err)
	}

	// export the website URL
	ctx.Export("bucketArn", bucket.Arn)
}
