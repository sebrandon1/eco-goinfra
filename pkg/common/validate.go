package common

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
)

// BuilderInterface is a generic interface defining Builders.
type BuilderInterface[T any] interface {
	GetDefinition() *T
	GetErrorMsg() string
	GetAPIClient() interface{}
	GetResourceType() string
}

// ValidateBuilder checks if the builder is valid.
func ValidateBuilder[T any](builder BuilderInterface[T]) (bool, error) {
	if builder == nil || reflect.ValueOf(builder).IsNil() {
		glog.V(100).Info("The builder is uninitialized or nil")

		return false, fmt.Errorf("error: received nil builder")
	}

	resourceType := builder.GetResourceType()

	definition := builder.GetDefinition()
	if definition == nil {
		glog.V(100).Infof("The %s is undefined or has a nil underlying value", resourceType)

		return false, fmt.Errorf("%s", msg.UndefinedCrdObjectErrString(resourceType))
	}

	apiClient := builder.GetAPIClient()
	if apiClient == nil || (reflect.ValueOf(apiClient).Kind() == reflect.Ptr && reflect.ValueOf(apiClient).IsNil()) {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceType)

		return false, fmt.Errorf("%s builder cannot have nil apiClient", resourceType)
	}

	if builder.GetErrorMsg() != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceType, builder.GetErrorMsg())

		return false, fmt.Errorf("%s", builder.GetErrorMsg())
	}

	return true, nil
}
