package services

import (
	"fmt"
	"sync"

	"github.com/google/go-containerregistry/pkg/authn"
	gcrName "github.com/google/go-containerregistry/pkg/name"
	gcrGoogle "github.com/google/go-containerregistry/pkg/v1/google"
	gcrRemote "github.com/google/go-containerregistry/pkg/v1/remote"
)

var gcrAuth authn.Authenticator
var gcrOnce sync.Once

func getGcrAuthenticator() authn.Authenticator {
	gcrOnce.Do(func() {
		var err error
		gcrAuth, err = gcrGoogle.NewEnvAuthenticator()
		if err != nil {
			panic(err)
		}

		if gcrAuth == nil {
			fmt.Errorf("Please provide valid Google Cloud credentials file")
			return
		}
	})
	return gcrAuth
}

// DeleteDigestFromString deletes digest idenitifed from string value
func DeleteDigestFromString(repoName string, digestIdentifier string) (bool, error) {

	// gcrName "github.com/google/go-containerregistry/pkg/name"
	// gcrRemote "github.com/google/go-containerregistry/pkg/v1/remote/transport"
	// conf "github.com/nitishkumar71/gcp-cleaner/pkg/configuration"
	// repoName = "us.gcr.io/one-storage-287104/onestorage"
	// fmt.Println("Repo Name", repoName)
	// gcrrepo, err := gcrName.NewRepository(repoName)
	// scopes := []string{gcrrepo.Scope(gcrRemote.DeleteScope), gcrrepo.Scope(gcrRemote.PullScope)}
	// ref := gcrrepo.Digest(digest)
	ref, err := gcrName.ParseReference(digestIdentifier)
	fmt.Println("Reference Created: ", ref)
	if err != nil {
		return false, err
	}

	if getGcrAuthenticator() == nil {
		fmt.Println("GCR Auth is Null")
	}

	err = gcrRemote.Delete(ref, gcrRemote.WithAuth(getGcrAuthenticator()))
	if err != nil {
		return false, err
	}

	// descriptor, err := gcrRemote.Get(ref, gcrRemote.WithContext(conf.GetContext()))
	// image, err := gcrRemote.Image(ref, gcrRemote.WithContext(conf.GetContext()))

	// fmt.Println("Image ", image)
	// descriptor, err := gcrRemote.Get(ref, gcrRemote.WithContext(conf.GetContext()))
	// fmt.Println("Descriptor ", descriptor)
	// if err != nil {
	// 	return false, err
	// }

	// auth := gcrGoogle.NewJSONKeyAuthenticator("/home/nitishkumar/Practice/onestorage/one-storage-287104-f9fe7d372d64.json")

	// tags, err := gcrGoogle.List(gcrrepo, gcrGoogle.WithContext(conf.GetContext()))
	// if err != nil {
	// 	return false, err
	// }

	// fmt.Println("================Tags===================")
	// fmt.Println(tags)
	// for _, m := range tags.Manifests {
	// 	fmt.Println("Name ", m.Tags)
	// }

	// repo, err := gcrName.NewRepository(repoName)
	// if err != nil {
	// 	return false, err
	// }

	// digest := repo.Digest(digestIdentifier)
	// fmt.Println("Digest ", digest)

	return true, nil
}
