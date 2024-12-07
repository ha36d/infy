package utils

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Labels(metadata *model.Metadata) pulumi.StringMap {

	Labels := pulumi.StringMap{}
	mapLabels := pulumi.Map{}
	for key, value := range metadata.Meta {
		Labels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
		mapLabels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
	}
	return Labels
}
