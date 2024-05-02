package clusteroperator

import (
	"fmt"
	"testing"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	configV1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	clusterOperatorGVK = schema.GroupVersionKind{
		Group:   APIGroup,
		Version: APIVersion,
		Kind:    APIKind,
	}
	defaultClusterOperatorName = "test-co"
)

func TestClusterOperatorPull(t *testing.T) {
	generateClusterOperator := func(name string) *configV1.ClusterOperator {
		return &configV1.ClusterOperator{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: configV1.ClusterOperatorSpec{},
		}
	}

	testCases := []struct {
		name                string
		addToRuntimeObjects bool
		expectedError       error
		client              bool
	}{
		{
			name:                "etcd",
			addToRuntimeObjects: true,
			expectedError:       nil,
			client:              true,
		},
		{
			name:                "",
			addToRuntimeObjects: true,
			expectedError:       fmt.Errorf("clusterOperator 'clusterOperatorName' cannot be empty"),
			client:              true,
		},
		{
			name:                "cotest",
			addToRuntimeObjects: false,
			expectedError:       fmt.Errorf("clusterOperator object cotest does not exist"),
			client:              true,
		},
		{
			name:                "cotest",
			addToRuntimeObjects: true,
			expectedError:       fmt.Errorf("clusterOperator 'apiClient' cannot be empty"),
			client:              false,
		},
	}

	for _, testCase := range testCases {
		// Pre-populate the runtime objects
		var runtimeObjects []runtime.Object

		var testSettings *clients.Settings

		testClusterOperator := generateClusterOperator(testCase.name)

		if testCase.addToRuntimeObjects {
			runtimeObjects = append(runtimeObjects, testClusterOperator)
		}

		if testCase.client {
			testSettings = clients.GetTestClients(clients.TestClientParams{
				K8sMockObjects: runtimeObjects,
			})
		}

		builderResult, err := Pull(testSettings, testCase.name)
		assert.Equal(t, testCase.expectedError, err)

		if testCase.expectedError != nil {
			assert.Equal(t, testCase.expectedError, err)
		} else {
			assert.Equal(t, testClusterOperator.Name, builderResult.Object.Name)
		}
	}
}

func TestClusterOperatorExist(t *testing.T) {
	testCases := []struct {
		testClusterOperator *Builder
		expectedStatus      bool
	}{
		{
			testClusterOperator: buildValidClusterOperatorBuilder(buildClusterOperatorClientWithDummyObject()),
			expectedStatus:      true,
		},
		{
			testClusterOperator: buildInValidClusterOperatorBuilder(buildClusterOperatorClientWithDummyObject()),
			expectedStatus:      false,
		},
		{
			testClusterOperator: buildValidClusterOperatorBuilder(clients.GetTestClients(clients.TestClientParams{})),
			expectedStatus:      false,
		},
	}

	for _, testCase := range testCases {
		exist := testCase.testClusterOperator.Exists()
		assert.Equal(t, testCase.expectedStatus, exist)
	}
}

func TestClusterOperatorGet(t *testing.T) {
	testCases := []struct {
		testClusterOperator *Builder
		expectedError       error
	}{
		{
			testClusterOperator: buildValidClusterOperatorBuilder(buildClusterOperatorClientWithDummyObject()),
			expectedError:       nil,
		},
		{
			testClusterOperator: buildInValidClusterOperatorBuilder(buildClusterOperatorClientWithDummyObject()),
			expectedError:       fmt.Errorf("the clusterOperator 'name' cannot be empty"),
		},
		{
			testClusterOperator: buildValidClusterOperatorBuilder(clients.GetTestClients(clients.TestClientParams{})),
			expectedError:       fmt.Errorf("clusteroperators.config.openshift.io \"test-co\" not found"),
		},
	}

	for _, testCase := range testCases {
		clusterOperatorObj, err := testCase.testClusterOperator.Get()

		if testCase.expectedError == nil {
			assert.Equal(t, clusterOperatorObj, testCase.testClusterOperator.Definition)
		} else {
			assert.Equal(t, testCase.expectedError.Error(), err.Error())
		}
	}
}

func buildValidClusterOperatorBuilder(apiClient *clients.Settings) *Builder {
	return newBuilder(apiClient, defaultClusterOperatorName)
}

func buildInValidClusterOperatorBuilder(apiClient *clients.Settings) *Builder {
	return newBuilder(apiClient, "")
}

func buildClusterOperatorClientWithDummyObject() *clients.Settings {
	return clients.GetTestClients(clients.TestClientParams{
		K8sMockObjects: buildDummyClusterOperatorConfig(),
		GVK:            []schema.GroupVersionKind{clusterOperatorGVK},
	})
}

func buildDummyClusterOperatorConfig() []runtime.Object {
	return append([]runtime.Object{}, &configV1.ClusterOperator{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultClusterOperatorName,
		},
		Spec: configV1.ClusterOperatorSpec{},
	})
}

// newBuilder method creates new instance of builder (for the unit test propose only).
func newBuilder(apiClient *clients.Settings, name string) *Builder {
	glog.V(100).Infof("Initializing new Builder structure with the name: %s", name)

	builder := &Builder{
		apiClient: apiClient.Client,
		Definition: &configV1.ClusterOperator{
			TypeMeta: metav1.TypeMeta{
				Kind:       APIKind,
				APIVersion: fmt.Sprintf("%s/%s", APIGroup, APIVersion),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            name,
				ResourceVersion: "999",
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the clusterOperator is empty")

		builder.errorMsg = "the clusterOperator 'name' cannot be empty"

		return builder
	}

	return builder
}
