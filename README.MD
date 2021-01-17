# GCP Cleaner

## Read me is in progress

A Rest API for cleaning the unused GCP Cloud Resources.

## Docker Repository
Install Latest Docker Images using
``` 
docker pull ghcr.io/nitishkumar71/gcp-cleaner:latest 
```
To Download other versions, follow the <a href="https://github.com/users/nitishkumar71/packages/container/package/gcp-cleaner" target="_blank">Link</a>

### Supported Services
- [Cloud Run](#cloud-run)



## Cloud Run
API will allow to clean up the older cloud revisions and associated container images with them
- POST - /cloudrun
#### Request Body
```
{
	"name":"cloud service name",
	"limit": 62, # no of images to retain
	"projectId": "project ID"
	"retainRevisionImage": true/false # Optional
}
```