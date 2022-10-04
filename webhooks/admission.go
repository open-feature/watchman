package webhooks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/open-feature/golang-sdk/pkg/openfeature"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// NOTE: RBAC not needed here.
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:webhook:path=/validate-v1-pod,mutating=false,failurePolicy=Ignore,groups="",resources=pods,verbs=create;update,versions=v1,name=mutate-demo.openfeature.dev,admissionReviewVersions=v1,sideEffects=NoneOnDryRun

type PodAdmission struct {
	Client   client.Client
	Log      logr.Logger
	decoder  *admission.Decoder
	OFClient *openfeature.Client
}

func (a *PodAdmission) getFlagDValue() (bool, error) {

	result := make(chan bool, 1)
	go func() {
		value, err := a.OFClient.BooleanValue("blocking", false,
			openfeature.EvaluationContext{},
			openfeature.EvaluationOptions{})
		if err != nil {
			a.Log.Error(err, "error getting flag value")
			result <- false
		} else {
			result <- value
		}
	}()
	select {
	case <-time.After(5 * time.Second):
		return false, errors.New("timeout getting flag value")
	case result := <-result:
		return result, nil
	}
}

func (a *PodAdmission) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	a.Log.Info("handling pod admission")
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	a.Log.Info("gRPC to flagD started")
	value, err := a.getFlagDValue()
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	a.Log.Info(fmt.Sprintf("blocking flag response %t", value))
	a.Log.Info("gRPC to flagD completed")

	if value {
		return admission.Denied("flag set to block admission")
	}

	return admission.Allowed("flag set to allow admission")
}

// InjectDecoder injects the decoder.
func (v *PodAdmission) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
