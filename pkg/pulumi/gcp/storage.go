package gcppulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	log "github.com/sirupsen/logrus"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) {
	// here we create the bucket

	bucket, err := storage.NewBucket(ctx, args["name"].(string), &storage.BucketArgs{
		Name:                     pulumi.String(args["name"].(string)),
		Location:                 pulumi.String(strings.ToUpper(metadata.Region[:2])),
		ForceDestroy:             pulumi.Bool(true),
		UniformBucketLevelAccess: pulumi.Bool(true),
		Labels: pulumi.StringMap{
			"team":    pulumi.String(strings.ToLower(metadata.Team)),
			"product": pulumi.String(strings.ToLower(metadata.Name)),
			"owner":   pulumi.String(strings.ToLower(metadata.Team)),
		},
	})
	if err != nil {
		log.Errorf("abnormal termination: %s", err)
	}

	// export the website URL
	ctx.Export("bucketUrl", bucket.Url)
}
