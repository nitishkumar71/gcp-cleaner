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

	ref, err := gcrName.ParseReference(digestIdentifier)
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

	return true, nil
}
