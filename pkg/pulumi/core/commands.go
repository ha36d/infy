package corepulumi

import (
	"context"
	"fmt"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"

	"github.com/iancoleman/strcase"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	awspulumi "github.com/ha36d/infy/pkg/pulumi/aws"
	azurepulumi "github.com/ha36d/infy/pkg/pulumi/azure"
	gcppulumi "github.com/ha36d/infy/pkg/pulumi/gcp"
	model "github.com/ha36d/infy/pkg/pulumi/model"
	ocipulumi "github.com/ha36d/infy/pkg/pulumi/oci"
)

const AZURE = "azure"
const GCP = "gcp"
const AWS = "aws"
const OCI = "oci"

func Preview(ctx context.Context, meta map[string]string, cloud string,
	account string, region string, components []map[string]any) {
	metadata := model.Metadata{Meta: meta, Cloud: cloud, Account: account, Region: region}

	stack := createOrSelectObjectStack(ctx, &metadata, components)
	// wire up our update to stream progress to stdout

	stdoutStreamer := optpreview.ProgressStreams(os.Stdout)
	_, err := stack.Preview(ctx, stdoutStreamer)
	if err != nil {
		log.Printf("Failed to preview %s stack: %v\n\n", metadata.Meta["Name"], err)
		os.Exit(1)
	}
	log.Printf("%s stack preview succeeded!", metadata.Meta["Name"])
}

func Up(ctx context.Context, meta map[string]string, cloud string,
	account string, region string, components []map[string]any) {
	metadata := model.Metadata{Meta: meta, Cloud: cloud, Account: account, Region: region}

	stack := createOrSelectObjectStack(ctx, &metadata, components)
	// wire up our update to stream progress to stdout

	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	_, err := stack.Up(ctx, stdoutStreamer)
	if err != nil {
		log.Printf("Failed to update %s stack: %v\n\n", metadata.Meta["Name"], err)
		os.Exit(1)
	}
	log.Printf("%s stack update succeeded!", metadata.Meta["Name"])
}

// this function gets our object stack ready for update/destroy
func createOrSelectObjectStack(ctx context.Context, metadata *model.Metadata, components []map[string]any) auto.Stack {
	deployFunc := func(ctx *pulumi.Context) error {
		if metadata.Cloud == AZURE {
			log.Printf("preparing resourcegroup function\n")
			azurepulumi.Holder{}.Resourcegroup(metadata, ctx)
		} else if metadata.Cloud == OCI {
			log.Printf("preparing compartment function\n")
			ocipulumi.Holder{}.Compartment(metadata, ctx)
		}
		for _, component := range components {
			for key, value := range component {
				log.Printf("preparing %s function\n", key)
				values := value.(map[string]any)
				var err error

				switch metadata.Cloud {
				case GCP:
					f := reflect.ValueOf(gcppulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context) error)
					err = f(metadata, values, ctx)
				case AWS:
					f := reflect.ValueOf(awspulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context) error)
					err = f(metadata, values, ctx)
				case AZURE:
					f := reflect.ValueOf(azurepulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context) error)
					err = f(metadata, values, ctx)
				case OCI:
					f := reflect.ValueOf(ocipulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context) error)
					err = f(metadata, values, ctx)
				}
				log.Println(err)
			}
		}
		return nil
	}
	return createOrSelectStack(ctx, metadata, deployFunc)
}

// this function gets our stack ready for update/destroy by prepping the workspace, init/selecting the stack
// and doing a refresh to make sure state and cloud resources are in sync
func createOrSelectStack(ctx context.Context, metadata *model.Metadata, deployFunc pulumi.RunFunc) auto.Stack {
	stackName := fmt.Sprintf("%s-%s-%s", metadata.Meta["Team"], metadata.Meta["Name"], metadata.Meta["Env"])
	// create or select a stack with an inline Pulumi program

	s, err := auto.UpsertStackInlineSource(ctx, stackName, metadata.Meta["Team"], deployFunc)
	if err != nil {
		log.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Created/Selected stack %q\n", stackName)

	w := s.Workspace()

	switch metadata.Cloud {
	case AWS:
		log.Println("Installing the AWS plugin")
		// for inline source programs, we must manage plugins ourselves
		err = w.InstallPlugin(ctx, "aws", "6.52.0")
		if err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}

		log.Println("Successfully installed AWS plugin")
		// set stack configuration specifying the AWS region to deploy
		if err = s.SetConfig(ctx, "aws:region", auto.ConfigValue{Value: metadata.Region}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
	case GCP:
		log.Println("Installing the GCP plugin")
		if err = w.InstallPlugin(ctx, "gcp", "7.35.0"); err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the GCP account to deploy
		if err = s.SetConfig(ctx, "gcp:project", auto.ConfigValue{Value: metadata.Account}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the GCP region to deploy
		if err = s.SetConfig(ctx, "gcp:region", auto.ConfigValue{Value: metadata.Region}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
	case AZURE:
		log.Println("Installing the Azure plugin")
		if err = w.InstallPlugin(ctx, "azure", "5.89.0"); err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the Azure region to deploy
		if err = s.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: metadata.Region}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
	case OCI:
		log.Println("Installing the OCI plugin")
		if err = w.InstallPlugin(ctx, "oci", "2.10.0"); err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}
	}

	log.Println("Successfully set config")
	log.Println("Starting refresh")

	_, err = s.Refresh(ctx)
	if err != nil {
		log.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}
	log.Println("Refresh succeeded!")

	return s
}
