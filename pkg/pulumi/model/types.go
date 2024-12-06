package model

type ResourceTracker struct {
	// Map of resource type to a map of resource name to resource
	// e.g., "vpc" -> {"main-vpc": vpcResource, "secondary-vpc": vpcResource}
	Resources map[string]map[string]interface{}
}

// NewResourceTracker creates a new ResourceTracker instance
func NewResourceTracker() *ResourceTracker {
	return &ResourceTracker{
		Resources: make(map[string]map[string]interface{}),
	}
}

// AddResource adds a resource to the tracker
func (rt *ResourceTracker) AddResource(resourceType string, name string, resource interface{}) {
	if rt.Resources[resourceType] == nil {
		rt.Resources[resourceType] = make(map[string]interface{})
	}
	rt.Resources[resourceType][name] = resource
}

// GetResource retrieves a resource by type and name
func (rt *ResourceTracker) GetResource(resourceType string, name string) (interface{}, bool) {
	if resources, ok := rt.Resources[resourceType]; ok {
		resource, exists := resources[name]
		return resource, exists
	}
	return nil, false
}

// GetResourcesByType retrieves all resources of a given type
func (rt *ResourceTracker) GetResourcesByType(resourceType string) map[string]interface{} {
	if resources, ok := rt.Resources[resourceType]; ok {
		return resources
	}
	return nil
}
