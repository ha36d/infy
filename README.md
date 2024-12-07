# Infy

Infy is a multi-cloud infrastructure management tool built on top of Pulumi, supporting AWS, Azure, GCP, and OCI cloud providers.

## Features

- Multi-cloud support (AWS, Azure, GCP, OCI)
- Declarative YAML configuration
- Resource dependency management
- Consistent interface across cloud providers
- Preview infrastructure changes
- State management via Pulumi

## Supported Resources

Currently supported resources across providers:

| Resource Type | AWS | Azure | GCP | OCI |
|--------------|-----|-------|-----|-----|
| Storage      | S3  | Blob  | GCS | Object Storage |
| Compute      | EC2 | VM    | GCE | Compute Instance |
| Network      | VPC | VNet  | VPC | VCN |

## Installation

```bash
go install github.com/ha36d/infy@latest
```

## Quick Start

1. Create a configuration file (.infy.yaml):

```yaml
metadata:
  name: myservice
  team: myteam
  env: development
  org: myorg
cloud: aws
account: "123456789"
region: us-west-2
components:
  - Storage:
      name: mybucket
  - Compute:
      name: myserver
      type: t2.micro
      size: 20
```

2. Preview changes:
```bash
infy preview
```

3. Apply changes:
```bash
infy build
```

## Configuration

### Metadata
- `name`: Service name
- `team`: Team name
- `env`: Environment (development, staging, production)
- `org`: Organization name

### Cloud Providers
- `cloud`: Cloud provider (aws, azure, gcp, oci)
- `account`: Account ID/Project ID
- `region`: Region for deployment

### Components
Components are defined as a list of resources to be created. Each resource type has its own configuration options.

## Usage

### Preview Changes
```bash
infy preview [--timeout 30m]
```

### Apply Changes
```bash
infy build [--timeout 2h] [--force]
```

### Common Flags
- `--config`: Specify config file (default: .infy.yaml)
- `--verbose`: Enable verbose logging
- `--timeout`: Operation timeout

## Contributing

See [CONTRIBUTE.md](CONTRIBUTE.md) for technical details on how to add new resources and contribute to the project.

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.

This project uses [Pulumi](https://github.com/pulumi/pulumi) which is licensed under the Apache License 2.0.