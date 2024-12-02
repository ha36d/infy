package utils

import (
	"strings"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func StringMapLabels(metadata *model.Metadata) pulumi.StringMap {

	stringMapLabels := pulumi.StringMap{}
	mapLabels := pulumi.Map{}
	for key, value := range metadata.Meta {
		stringMapLabels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
		mapLabels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
	}
	return stringMapLabels
}

func MapLabels(metadata *model.Metadata) pulumi.Map {

	stringMapLabels := pulumi.StringMap{}
	mapLabels := pulumi.Map{}
	for key, value := range metadata.Meta {
		stringMapLabels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
		mapLabels[strings.ToLower(key)] = pulumi.String(strings.ToLower(value))
	}
	return mapLabels
}
