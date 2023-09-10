# infy

## Overview
Abstraction on top of pulumi.

## Features
- Describe the infra (AWS, AZURE, GCP) as yaml file
- Build the infra (AWS, AZURE, GCP) from a yaml file

## How to

First you need to install [pulumi](https://www.pulumi.com/docs/install/).
Then, we need the state management:
```
pulumi login
```
Then you can create the Yaml file by executing
```
./infy generate -c service.yaml
```

A file `service.yaml` will be create that describes the infrastructure, 
```
cat service.yaml

name: Myservice
team: MyTeam
cloud: gcp
account: project
region: us-west-2
env: development
components:
  Storage:
    name: someRandomStorageSystem
```

and then you can apply the file in the cloud by

```
./infy build -c service.yaml
```

## Limitations
- For the moment, storage and VM are supported.


## Acknowledgements

This couldn't work without

 - [Pulumi](https://github.com/pulumi/pulumi)
 - [Cobra](https://github.com/spf13/cobra/)
 - [Viper](https://github.com/spf13/viper)
 - [Goreleaser](https://github.com/Goreleaser)