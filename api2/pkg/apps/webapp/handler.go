package webapp

import (
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/apps/cms"
	"github.com/nrc-no/core-kafka/pkg/apps/iam"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"github.com/nrc-no/core-kafka/pkg/sessionmanager"
	"github.com/ory/hydra-client-go/client"
	"os"
)

type Server struct {
	hydraAdminClient  *client.OryHydra
	hydraPublicClient *client.OryHydra
	renderFactory     *RendererFactory
	sessionManager    sessionmanager.Store
	credentialsClient *auth.CredentialsClient
	iam               iam.Interface
	cms               cms.Interface
}

type Options struct {
	TemplateDirectory string
}

func NewHandler(
	options Options,
	hydraAdminClient *client.OryHydra,
	hydraPublicClient *client.OryHydra,
	sessionManager sessionmanager.Store,
	credentialsClient *auth.CredentialsClient,
	iamClient *iam.ClientSet,
	cmsClient *cms.ClientSet,
) (*Server, error) {

	renderFactory, err := NewRendererFactory(options.TemplateDirectory)
	if err != nil {
		return nil, err
	}

	e, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fmt.Println(e)

	h := &Server{
		hydraAdminClient:  hydraAdminClient,
		hydraPublicClient: hydraPublicClient,
		renderFactory:     renderFactory,
		sessionManager:    sessionManager,
		credentialsClient: credentialsClient,
		iam:               iamClient,
		cms:               cmsClient,
	}
	return h, nil
}
