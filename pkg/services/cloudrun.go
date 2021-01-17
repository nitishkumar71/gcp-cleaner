package services

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	conf "github.com/nitishkumar71/gcp-cleaner/pkg/configuration"
	"google.golang.org/api/run/v1"
)

// Google Cloud Run Delete Service
// https://cloud.google.com/run
// Enable Cloud Run Admin API - https://console.developers.google.com/apis/api/run.googleapis.com/overview
//
var runService *run.APIService
var runOnce sync.Once

func getRunService() *run.APIService {
	runOnce.Do(func() {
		var err error
		runService, err = run.NewService(conf.GetContext())
		if err != nil {
			panic(err)
		}

		if runService == nil {
			fmt.Errorf("Please provide valid Google credentials file")
			return
		}
	})
	return runService
}

type CloudRunRequest struct {
	Limit     uint   `json:"limit,omitempty"`
	Name      string `json:"name" binding:"required"`
	ProjectId string `json:"projectId" binding:"required"`
}

// DeleteCloudRunRevisionsPost is the handler for deleting google cloud revisions
func DeleteCloudRunRevisionsPost(c *gin.Context) {
	var request CloudRunRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	namespace := fmt.Sprintf("namespaces/%s", request.ProjectId)

	trafficRevisions, err := getRevisionsWithTraffic(namespace, request.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	revisions, err := deleteServiceRevisions(namespace, request.Name, int(request.Limit), trafficRevisions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"revisions": revisions, "length": len(revisions)})
}

func deleteServiceRevisions(namespace string, service string, limit int, trafficRevisions map[string]bool) ([]string, error) {
	//  revisions code
	service = fmt.Sprintf("serving.knative.dev/service=%s", service)
	namespaceService := run.NewNamespacesRevisionsService(getRunService())
	authCall := namespaceService.List(namespace)
	authCall = authCall.LabelSelector(service)
	revisions, err := authCall.Do()

	if err != nil {
		return nil, err
	}

	revisionService := run.NewProjectsLocationsRevisionsService(getRunService())
	revisionsDeleted := make([]string, 0)
	for index, revision := range revisions.Items {

		// skip limit revisions
		if index < limit {
			continue
		}

		// delete revisions not in active traffic list
		if !trafficRevisions[revision.Metadata.Name] {

			region := revision.Metadata.Labels["cloud.googleapis.com/location"]
			revisionFullPath := fmt.Sprintf("projects/%s/locations/%s/revisions/%s", revision.Metadata.Namespace, region, revision.Metadata.Name)
			deleteCall := revisionService.Delete(revisionFullPath)
			_, err = deleteCall.Do()
			if err != nil {
				return nil, err
			}

			repoName := revision.Spec.Containers[0].Image
			_, err = DeleteDigestFromString(repoName, revision.Status.ImageDigest)
			if err != nil {
				return nil, err
			}

			revisionsDeleted = append(revisionsDeleted, revision.Metadata.Name)
		}
	}
	return revisionsDeleted, nil
}

func getRevisionsWithTraffic(namespace string, service string) (map[string]bool, error) {
	// services code
	service = fmt.Sprintf("metadata.name=%s", service)
	namespaceService := run.NewNamespacesServicesService(getRunService())
	servicesListCall := namespaceService.List(namespace)
	servicesListCall = servicesListCall.FieldSelector(service)
	services, err := servicesListCall.Do()

	if err != nil {
		return nil, err
	}

	trafficRevisions := make(map[string]bool)
	for _, service := range services.Items {
		for _, target := range service.Status.Traffic {
			trafficRevisions[target.RevisionName] = true
		}
	}

	return trafficRevisions, nil
}
