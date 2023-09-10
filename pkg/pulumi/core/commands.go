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
)

const AZURE = "azure"
const GCP = "gcp"
const AWS = "aws"

func Preview(ctx context.Context, name string, team string, env string, cloud string,
	account string, region string, components []map[string]any) {
	metadata := model.Metadata{Name: name, Team: team, Env: env, Cloud: cloud, Account: account, Region: region}

	stack := createOrSelectObjectStack(ctx, &metadata, components)
	// wire up our update to stream progress to stdout

	stdoutStreamer := optpreview.ProgressStreams(os.Stdout)
	_, err := stack.Preview(ctx, stdoutStreamer)
	if err != nil {
		log.Printf("Failed to preview %s stack: %v\n\n", metadata.Name, err)
		os.Exit(1)
	}
	log.Printf("%s stack preview succeeded!", metadata.Name)
}

func Up(ctx context.Context, name string, team string, env string, cloud string,
	account string, region string, components []map[string]any) {
	metadata := model.Metadata{Name: name, Team: team, Env: env, Cloud: cloud, Account: account, Region: region}

	stack := createOrSelectObjectStack(ctx, &metadata, components)
	// wire up our update to stream progress to stdout

	stdoutStreamer := optup.ProgressStreams(os.Stdout)
	_, err := stack.Up(ctx, stdoutStreamer)
	if err != nil {
		log.Printf("Failed to update %s stack: %v\n\n", metadata.Name, err)
		os.Exit(1)
	}
	log.Printf("%s stack update succeeded!", metadata.Name)
}

// this function gets our object stack ready for update/destroy
func createOrSelectObjectStack(ctx context.Context, metadata *model.Metadata, components []map[string]any) auto.Stack {
	deployFunc := func(ctx *pulumi.Context) error {
		if metadata.Cloud == AZURE {
			log.Printf("preparing resourcegroup function\n")
			azurepulumi.Holder{}.Resourcegroup(metadata, ctx)
		}
		for _, component := range components {
			for key, value := range component {
				log.Printf("preparing %s function\n", key)
				values := value.(map[string]any)
				if metadata.Cloud == GCP {
					f := reflect.ValueOf(gcppulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context))
					f(metadata, values, ctx)
				} else if metadata.Cloud == AWS {
					f := reflect.ValueOf(awspulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context))
					f(metadata, values, ctx)
				} else if metadata.Cloud == AZURE {
					f := reflect.ValueOf(azurepulumi.Holder{}).
						MethodByName(strcase.ToCamel(key)).Interface().(func(*model.Metadata, map[string]any, *pulumi.Context))
					f(metadata, values, ctx)
				}
			}
		}
		return nil
	}
	return createOrSelectStack(ctx, metadata, deployFunc)
}

// this function gets our stack ready for update/destroy by prepping the workspace, init/selecting the stack
// and doing a refresh to make sure state and cloud resources are in sync
func createOrSelectStack(ctx context.Context, metadata *model.Metadata, deployFunc pulumi.RunFunc) auto.Stack {
	stackName := fmt.Sprintf("%s-%s-%s", metadata.Team, metadata.Name, metadata.Env)
	// create or select a stack with an inline Pulumi program

	s, err := auto.UpsertStackInlineSource(ctx, stackName, metadata.Name, deployFunc)
	if err != nil {
		log.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Created/Selected stack %q\n", stackName)

	w := s.Workspace()

	if metadata.Cloud == AWS {
		log.Println("Installing the AWS plugin")

		// for inline source programs, we must manage plugins ourselves
		err = w.InstallPlugin(ctx, "aws", "6.27.0")
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
	} else if metadata.Cloud == GCP {
		if err = w.InstallPlugin(ctx, "gcp", "v7.13.0"); err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the AWS region to deploy
		if err = s.SetConfig(ctx, "gcp:project", auto.ConfigValue{Value: metadata.Account}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the AWS region to deploy
		if err = s.SetConfig(ctx, "gcp:region", auto.ConfigValue{Value: metadata.Region}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
		log.Println("Installing the GCP plugin")
	} else if metadata.Cloud == AZURE {
		if err = w.InstallPlugin(ctx, "azure", "v5.73.0"); err != nil {
			log.Printf("Failed to install program plugins: %v\n", err)
			os.Exit(1)
		}
		// set stack configuration specifying the AWS region to deploy
		if err = s.SetConfig(ctx, "azure:location", auto.ConfigValue{Value: metadata.Region}); err != nil {
			log.Printf("Failed to set config: %v\n", err)
			os.Exit(1)
		}
		log.Println("Using Azure")
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
