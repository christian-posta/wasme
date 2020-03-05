package istio

import (
	"github.com/pkg/errors"
	"github.com/solo-io/wasme/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	istiodDeploymentName  = "istiod"
	defaultIstioNamespace = "istio-system"
	istiodContainerName   = "discovery"
)

type VersionInspector interface {
	GetIstioVersion() (string, error)
}

type versionInspector struct {
	istioNamespace string
	kube           kubernetes.Interface
}

func (i *versionInspector) GetIstioVersion() (string, error) {
	istioNamespace := i.istioNamespace
	if istioNamespace == "" {
		istioNamespace = defaultIstioNamespace
	}
	istiodDeployment, err := i.kube.AppsV1().Deployments(istioNamespace).Get(istiodDeploymentName, metav1.GetOptions{})
	if err != nil {
		return "", nil
	}
	var istiodImage string
	for _, container := range istiodDeployment.Spec.Template.Spec.Containers {
		if container.Name == istiodContainerName {
			istiodImage = container.Image
			break
		}
	}
	if istiodImage == "" {
		return "", errors.Errorf("did not find container named %s on istiod deployment", istiodContainerName)
	}

	_, tag, err := util.SplitImageRef(istiodImage)

	return tag, err
}
