package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	managementapi "github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/managementapi"
	k8s "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	requestTimeoutSeconds = 5
	separator             = "/"
	defaultPlanSecurity   = "KEY_LESS"
	defaultPlanStatus     = "PUBLISHED"
	defaultPlanName       = "G.K.O. Default"
	origin                = "kubernetes"
	mode                  = "fully_managed"
)

type Delegate struct {
	ctx               context.Context
	managementContext *gio.ManagementContext
	apimClient        *managementapi.Client
	k8sClient         k8s.Client
	log               logr.Logger
}

func NewDelegate(ctx context.Context, client k8s.Client, log logr.Logger) *Delegate {
	return &Delegate{
		ctx, nil, nil, client, log,
	}
}

func (d *Delegate) SetManagementContext(managementContext *gio.ManagementContext) {
	if managementContext == nil {
		return
	}

	d.managementContext = managementContext

	httpClient := http.Client{Timeout: requestTimeoutSeconds * time.Second}
	d.apimClient = managementapi.NewClient(d.ctx, d.managementContext, httpClient)
}

func (d *Delegate) IsConnectedToManagementApi() bool {
	return d.apimClient != nil
}