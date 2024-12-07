package corepulumi

import (
	"context"
	"reflect"
	"testing"

	model "github.com/ha36d/infy/pkg/pulumi/model"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestPreview(t *testing.T) {
	type args struct {
		ctx        context.Context
		meta       map[string]string
		cloud      string
		account    string
		region     string
		components []map[string]map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Preview(tt.args.ctx, tt.args.meta, tt.args.cloud, tt.args.account, tt.args.region, tt.args.components); (err != nil) != tt.wantErr {
				t.Errorf("Preview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUp(t *testing.T) {
	type args struct {
		ctx        context.Context
		meta       map[string]string
		cloud      string
		account    string
		region     string
		components []map[string]map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Up(tt.args.ctx, tt.args.meta, tt.args.cloud, tt.args.account, tt.args.region, tt.args.components); (err != nil) != tt.wantErr {
				t.Errorf("Up() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrSelectObjectStack(t *testing.T) {
	type args struct {
		ctx        context.Context
		metadata   *model.Metadata
		components []map[string]map[string]any
	}
	tests := []struct {
		name string
		args args
		want auto.Stack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createOrSelectObjectStack(tt.args.ctx, tt.args.metadata, tt.args.components); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createOrSelectObjectStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createOrSelectStack(t *testing.T) {
	type args struct {
		ctx        context.Context
		metadata   *model.Metadata
		deployFunc pulumi.RunFunc
	}
	tests := []struct {
		name string
		args args
		want auto.Stack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createOrSelectStack(tt.args.ctx, tt.args.metadata, tt.args.deployFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createOrSelectStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
