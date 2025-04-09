package common

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
)

type BuilderInterface interface {
	GetDefinition() interface{}
	GetErrorMsg() string
	GetAPIClient() interface{}
	GetResourceType() string
}

func ValidateBuilder(builder BuilderInterface) (bool, error) {
	if builder == nil {
		glog.V(100).Info("The builder is uninitialized")
		return false, fmt.Errorf("error: received nil builder")
	}

	resourceType := builder.GetResourceType()

	if builder.GetDefinition() == nil {
		glog.V(100).Infof("The %s is undefined", resourceType)
		return false, fmt.Errorf("%s", msg.UndefinedCrdObjectErrString(resourceType))
	}

	if builder.GetAPIClient() == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceType)
		return false, fmt.Errorf("%s builder cannot have nil apiClient", resourceType)
	}

	if builder.GetErrorMsg() != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceType, builder.GetErrorMsg())
		return false, fmt.Errorf("%s", builder.GetErrorMsg())
	}

	return true, nil
}
