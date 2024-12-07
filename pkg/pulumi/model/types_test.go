package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResourceTracker(t *testing.T) {
	tracker := NewResourceTracker()
	assert.NotNil(t, tracker.resources)
	assert.NotNil(t, tracker.metadata)
}

func TestResourceTracker_AddResource(t *testing.T) {
	tracker := NewResourceTracker()

	tests := []struct {
		name         string
		resourceType string
		resourceName string
		resource     interface{}
		deps         []string
		wantErr      bool
	}{
		{
			name:         "add valid resource",
			resourceType: "storage",
			resourceName: "test-bucket",
			resource:     "mock-resource",
			deps:         nil,
			wantErr:      false,
		},
		{
			name:         "add resource with dependencies",
			resourceType: "compute",
			resourceName: "test-vm",
			resource:     "mock-vm",
			deps:         []string{"network"},
			wantErr:      false,
		},
		{
			name:         "empty resource type",
			resourceType: "",
			resourceName: "test",
			resource:     "mock",
			deps:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tracker.AddResource(tt.resourceType, tt.resourceName, tt.resource, tt.deps...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify resource was added
				resource, exists := tracker.GetResource(tt.resourceType, tt.resourceName)
				assert.True(t, exists)
				assert.Equal(t, tt.resource, resource)

				// Verify metadata was created
				state, err := tracker.GetResourceState(tt.resourceType, tt.resourceName)
				assert.NoError(t, err)
				assert.Equal(t, ResourceStatePending, state)

				// Verify dependencies
				if len(tt.deps) > 0 {
					deps, err := tracker.GetDependencies(tt.resourceType, tt.resourceName)
					assert.NoError(t, err)
					assert.Equal(t, tt.deps, deps.DependsOn)
				}
			}
		})
	}
}

func TestResourceTracker_UpdateResourceState(t *testing.T) {
	tracker := NewResourceTracker()
	resourceType := "storage"
	resourceName := "test-bucket"

	// Add a resource first
	err := tracker.AddResource(resourceType, resourceName, "mock-resource")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		resType  string
		resName  string
		newState ResourceState
		wantErr  bool
	}{
		{
			name:     "update to creating",
			resType:  resourceType,
			resName:  resourceName,
			newState: ResourceStateCreating,
			wantErr:  false,
		},
		{
			name:     "update non-existent resource",
			resType:  "invalid",
			resName:  "invalid",
			newState: ResourceStateCreating,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tracker.UpdateResourceState(tt.resType, tt.resName, tt.newState)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				state, err := tracker.GetResourceState(tt.resType, tt.resName)
				assert.NoError(t, err)
				assert.Equal(t, tt.newState, state)
			}
		})
	}
}

func TestResourceTracker_GetResourcesByType(t *testing.T) {
	tracker := NewResourceTracker()

	// Add some resources
	resources := map[string]interface{}{
		"bucket1": "mock-bucket-1",
		"bucket2": "mock-bucket-2",
	}

	for name, res := range resources {
		err := tracker.AddResource("storage", name, res)
		assert.NoError(t, err)
	}

	// Test getting resources
	got := tracker.GetResourcesByType("storage")
	assert.Equal(t, len(resources), len(got))
	for name, res := range resources {
		assert.Equal(t, res, got[name])
	}

	// Test non-existent type
	got = tracker.GetResourcesByType("invalid")
	assert.Nil(t, got)
}

func TestResourceTracker_DeleteResource(t *testing.T) {
	tracker := NewResourceTracker()

	// Add a resource
	err := tracker.AddResource("storage", "test-bucket", "mock-resource")
	assert.NoError(t, err)

	// Test deletion
	err = tracker.DeleteResource("storage", "test-bucket")
	assert.NoError(t, err)

	// Verify resource is gone
	_, exists := tracker.GetResource("storage", "test-bucket")
	assert.False(t, exists)

	// Test deleting non-existent resource
	err = tracker.DeleteResource("invalid", "invalid")
	assert.Error(t, err)
}

func TestResourceTracker_ListResourceTypes(t *testing.T) {
	tracker := NewResourceTracker()

	// Add resources of different types
	resourceTypes := []string{"storage", "compute", "network"}
	for _, rt := range resourceTypes {
		err := tracker.AddResource(rt, "test", "mock")
		assert.NoError(t, err)
	}

	// Test listing types
	types := tracker.ListResourceTypes()
	assert.Equal(t, len(resourceTypes), len(types))
	for _, rt := range resourceTypes {
		assert.Contains(t, types, rt)
	}
}
