package gcppulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (Holder) Storage(metadata *model.Metadata, args map[string]any, ctx *pulumi.Context) error {
	// here we create the bucket

	_, err := storage.NewBucket(ctx, args["name"].(string), &storage.BucketArgs{
		Name:                     pulumi.String(args["name"].(string)),
		Location:                 pulumi.String(strings.ToUpper(metadata.Region[:2])),
		ForceDestroy:             pulumi.Bool(true),
		UniformBucketLevelAccess: pulumi.Bool(true),
		Labels:                   utils.StringMapLabels(metadata),
	})
	if err != nil {
		return err
	}

	return nil
}
