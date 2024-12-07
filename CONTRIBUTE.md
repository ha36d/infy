
## Resource Implementation Guidelines

1. **Resource Tracking**
   - Use `tracker.AddResource()` to track created resources
   - Use `tracker.GetResource()` to retrieve dependencies
   - Resource types should be consistent across providers

2. **Error Handling**
   - Return meaningful errors
   - Use proper error wrapping
   - Don't panic, handle errors gracefully

3. **Configuration**
   - All resource configurations come through `args map[string]any`
   - Validate required arguments
   - Provide sensible defaults when possible

4. **Dependencies**
   - Use the resource tracker to manage dependencies
   - Example retrieving a dependency:
   ```go
   networkResource, exists := tracker.GetResource("network", metadata.Meta["name"])
   if !exists {
       return fmt.Errorf("network not found")
   }
   ```

## Testing

1. **Unit Tests**
   - Test resource creation
   - Test error conditions
   - Test resource tracking
   - Test dependency handling

2. **Integration Tests**
   - Test resource creation in actual cloud environments
   - Use test fixtures for consistent testing
   - Clean up resources after tests

## Common Patterns

1. **Resource Creation**
```go
func (Holder) Resource(metadata model.Metadata, args map[string]any,
ctx pulumi.Context, tracker model.ResourceTracker) error {
// 1. Get dependencies if needed
dependency, exists := tracker.GetResource("dependency", metadata.Meta["name"])
if !exists {
return fmt.Errorf("dependency not found")
}
// 2. Create resource
resource, err := service.NewResource(ctx, args["name"].(string), &service.ResourceArgs{
// Configure resource
})
if err != nil {
return err
}
// 3. Track resource
tracker.AddResource("resourcetype", metadata.Meta["name"], resource)
return nil
}
```

2. **Resource Configuration**

```yaml
components:
NewResource:
name: myresource
property1: value1
property2: value2
```

## Best Practices

1. **Code Organization**
   - One resource type per file
   - Consistent naming across providers
   - Clear separation of concerns

2. **Documentation**
   - Document all exported functions
   - Include usage examples
   - Document any special requirements

3. **Testing**
   - Write tests for new resources
   - Include both success and failure cases
   - Test resource tracking

4. **Error Handling**
   - Use descriptive error messages
   - Handle all error cases
   - Clean up on errors when needed

## Pull Request Process

1. Create a feature branch
2. Implement the new resource
3. Add tests
4. Update documentation
5. Submit PR with description of changes

## Questions?

Feel free to open an issue for:
- Implementation questions
- Feature requests
- Bug reports
- Documentation improvements