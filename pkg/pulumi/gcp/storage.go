package gcppulumi

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/ha36d/infy/pkg/pulumi/utils"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (h Holder) Storage(metadata *model.Metadata, values map[string]any, ctx *pulumi.Context, tracker *model.ResourceTracker) error {
	// Create storage using network from tracker if needed
	storage, err := storage.NewBucket(ctx, values["name"].(string), &storage.BucketArgs{
		Name:                     pulumi.String(values["name"].(string)),
		Location:                 pulumi.String(strings.ToUpper(metadata.Region[:2])),
		ForceDestroy:             pulumi.Bool(true),
		UniformBucketLevelAccess: pulumi.Bool(true),
		Labels:                   utils.Labels(metadata),
	})
	if err != nil {
		return err
	}

	// Add the storage to the tracker
	tracker.AddResource("storage", metadata.Meta["name"], storage)
	return nil
}
