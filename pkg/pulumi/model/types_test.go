package model

import (
	"reflect"
	"testing"
)

func TestNewResourceTracker(t *testing.T) {
	tests := []struct {
		name string
		want *ResourceTracker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResourceTracker(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceTracker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceTracker_AddResource(t *testing.T) {
	type fields struct {
		Resources map[string]map[string]interface{}
	}
	type args struct {
		resourceType string
		name         string
		resource     interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &ResourceTracker{
				Resources: tt.fields.Resources,
			}
			rt.AddResource(tt.args.resourceType, tt.args.name, tt.args.resource)
		})
	}
}

func TestResourceTracker_GetResource(t *testing.T) {
	type fields struct {
		Resources map[string]map[string]interface{}
	}
	type args struct {
		resourceType string
		name         string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &ResourceTracker{
				Resources: tt.fields.Resources,
			}
			got, got1 := rt.GetResource(tt.args.resourceType, tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceTracker.GetResource() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ResourceTracker.GetResource() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestResourceTracker_GetResourcesByType(t *testing.T) {
	type fields struct {
		Resources map[string]map[string]interface{}
	}
	type args struct {
		resourceType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &ResourceTracker{
				Resources: tt.fields.Resources,
			}
			if got := rt.GetResourcesByType(tt.args.resourceType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceTracker.GetResourcesByType() = %v, want %v", got, tt.want)
			}
		})
	}
}
