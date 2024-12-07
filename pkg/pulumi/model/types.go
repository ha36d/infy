package model

import (
	"fmt"
	"sync"
	"time"
)

type ResourceState string

const (
	ResourceStatePending  ResourceState = "pending"
	ResourceStateCreating ResourceState = "creating"
	ResourceStateCreated  ResourceState = "created"
	ResourceStateFailed   ResourceState = "failed"
	ResourceStateDeleted  ResourceState = "deleted"
)

type ResourceDependency struct {
	DependsOn  []string
	RequiredBy []string
}

type ResourceMetadata struct {
	State        ResourceState
	Dependencies ResourceDependency
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Resource represents a cloud resource with common attributes
type Resource interface {
	GetID() string
	GetName() string
	GetType() string
}

type ResourceTracker struct {
	// Map of resource type to a map of resource name to resource
	resources map[string]map[string]interface{}
	metadata  map[string]map[string]ResourceMetadata
	mu        sync.RWMutex // For thread safety
}

// NewResourceTracker creates a new ResourceTracker instance
func NewResourceTracker() *ResourceTracker {
	return &ResourceTracker{
		resources: make(map[string]map[string]interface{}),
		metadata:  make(map[string]map[string]ResourceMetadata),
	}
}

// AddResource adds a resource to the tracker with dependencies
func (rt *ResourceTracker) AddResource(resourceType string, name string, resource interface{}, dependencies ...string) error {
	if name == "" || resourceType == "" {
		return fmt.Errorf("resource name and type cannot be empty")
	}

	rt.mu.Lock()
	defer rt.mu.Unlock()

	// Initialize maps if they don't exist
	if rt.resources[resourceType] == nil {
		rt.resources[resourceType] = make(map[string]interface{})
		rt.metadata[resourceType] = make(map[string]ResourceMetadata)
	}

	// Create metadata
	meta := ResourceMetadata{
		State: ResourceStatePending,
		Dependencies: ResourceDependency{
			DependsOn: dependencies,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Update RequiredBy for all dependencies
	for _, dep := range dependencies {
		for rType, resources := range rt.metadata {
			for rName, rMeta := range resources {
				if rName == dep {
					rMeta.Dependencies.RequiredBy = append(rMeta.Dependencies.RequiredBy, name)
					rt.metadata[rType][rName] = rMeta
				}
			}
		}
	}

	rt.resources[resourceType][name] = resource
	rt.metadata[resourceType][name] = meta
	return nil
}

// UpdateResourceState updates the state of a resource
func (rt *ResourceTracker) UpdateResourceState(resourceType string, name string, state ResourceState) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if meta, exists := rt.metadata[resourceType][name]; exists {
		meta.State = state
		meta.UpdatedAt = time.Now()
		rt.metadata[resourceType][name] = meta
		return nil
	}
	return fmt.Errorf("resource %s of type %s not found", name, resourceType)
}

// GetResourceState gets the current state of a resource
func (rt *ResourceTracker) GetResourceState(resourceType string, name string) (ResourceState, error) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	if meta, exists := rt.metadata[resourceType][name]; exists {
		return meta.State, nil
	}
	return "", fmt.Errorf("resource %s of type %s not found", name, resourceType)
}

// GetDependencies gets the dependencies of a resource
func (rt *ResourceTracker) GetDependencies(resourceType string, name string) (ResourceDependency, error) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	if meta, exists := rt.metadata[resourceType][name]; exists {
		return meta.Dependencies, nil
	}
	return ResourceDependency{}, fmt.Errorf("resource %s of type %s not found", name, resourceType)
}

// GetResource retrieves a resource by type and name
func (rt *ResourceTracker) GetResource(resourceType string, name string) (interface{}, bool) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	if resources, ok := rt.resources[resourceType]; ok {
		resource, exists := resources[name]
		return resource, exists
	}
	return nil, false
}

// GetResourcesByType retrieves all resources of a given type
func (rt *ResourceTracker) GetResourcesByType(resourceType string) map[string]interface{} {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	if resources, ok := rt.resources[resourceType]; ok {
		// Return a copy to prevent map mutations
		result := make(map[string]interface{}, len(resources))
		for k, v := range resources {
			result[k] = v
		}
		return result
	}
	return nil
}

// DeleteResource removes a resource from the tracker
func (rt *ResourceTracker) DeleteResource(resourceType string, name string) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if resources, ok := rt.resources[resourceType]; ok {
		if _, exists := resources[name]; exists {
			delete(resources, name)
			return nil
		}
	}
	return fmt.Errorf("resource %s of type %s not found", name, resourceType)
}

// ListResourceTypes returns all resource types currently tracked
func (rt *ResourceTracker) ListResourceTypes() []string {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	types := make([]string, 0, len(rt.resources))
	for resourceType := range rt.resources {
		types = append(types, resourceType)
	}
	return types
}
