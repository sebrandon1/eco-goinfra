package assisted

import (
	"testing"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	agentInstallV1Beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestNewAgentBuilder(t *testing.T) {
	testCases := []struct {
		apiClientNil  bool
		definitionNil bool
	}{
		{
			apiClientNil:  true,
			definitionNil: false,
		},
		{
			apiClientNil:  false,
			definitionNil: true,
		},
		{
			apiClientNil:  false,
			definitionNil: false,
		},
	}

	for _, testCase := range testCases {

		var (
			testApiClient  goclient.Client
			testDefinition *agentInstallV1Beta1.Agent
		)

		if testCase.apiClientNil {
			testApiClient = nil
		} else {
			testApiClient, _ = goclient.New(nil, goclient.Options{})
		}

		if testCase.definitionNil {
			testDefinition = nil
		} else {
			testDefinition = &agentInstallV1Beta1.Agent{}
		}

		testBuilder := newAgentBuilder(testApiClient, testDefinition)

		if testCase.apiClientNil || testCase.definitionNil {
			assert.Nil(t, testBuilder)
		} else {
			assert.NotNil(t, testBuilder)
		}
	}
}

func TestPullAgent(t *testing.T) {
	generateAgent := func(name, namespace string) *agentInstallV1Beta1.Agent {
		return &agentInstallV1Beta1.Agent{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		agentName           string
		agentNamespace      string
		expectedError       bool
		addToRuntimeObjects bool
		expectedErrorText   string
	}{
		{
			agentName:           "test-agent",
			agentNamespace:      "test-namespace",
			expectedError:       false,
			addToRuntimeObjects: true,
			expectedErrorText:   "",
		},
		{
			agentName:           "test-agent",
			agentNamespace:      "test-namespace",
			expectedError:       true,
			addToRuntimeObjects: false,
			expectedErrorText:   "agent object test-agent does not exist in namespace test-namespace",
		},
		{
			agentName:           "",
			agentNamespace:      "test-namespace",
			expectedError:       true,
			addToRuntimeObjects: false,
			expectedErrorText:   "agent 'name' cannot be empty",
		},
		{
			agentName:           "test-agent",
			agentNamespace:      "",
			expectedError:       true,
			addToRuntimeObjects: false,
			expectedErrorText:   "agent 'namespace' cannot be empty",
		},
	}

	for _, testCase := range testCases {
		var (
			runtimeObjects []runtime.Object
			testSettings   *clients.Settings
		)

		testAgent := generateAgent(testCase.agentName, testCase.agentNamespace)

		if testCase.addToRuntimeObjects {
			runtimeObjects = append(runtimeObjects, testAgent)
		}

		testSettings = clients.GetTestClients(clients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})

		result, err := PullAgent(testSettings, testCase.agentName, testCase.agentNamespace)

		if testCase.expectedError {
			assert.NotNil(t, err)
			assert.Equal(t, testCase.expectedErrorText, err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, testAgent.Name, result.Definition.Name)
			assert.Equal(t, testAgent.Namespace, result.Definition.Namespace)
			assert.Equal(t, testAgent.Name, result.Object.Name)
			assert.Equal(t, testAgent.Namespace, result.Object.Namespace)
			assert.Equal(t, result.Definition, result.Object)
		}
	}
}

func TestAgentWithHostName(t *testing.T) {
	generateAgent := func(name, namespace string) *agentInstallV1Beta1.Agent {
		return &agentInstallV1Beta1.Agent{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		agentName           string
		agentNamespace      string
		hostname            string
		expectedError       bool
		addToRuntimeObjects bool
		expectedErrorText   string
	}{
		{
			agentName:           "test-agent",
			agentNamespace:      "test-namespace",
			hostname:            "test-hostname",
			expectedError:       false,
			addToRuntimeObjects: true,
			expectedErrorText:   "",
		},
	}

	for _, testCase := range testCases {
		var (
			runtimeObjects []runtime.Object
			testSettings   *clients.Settings
		)

		testAgent := generateAgent(testCase.agentName, testCase.agentNamespace)

		if testCase.addToRuntimeObjects {
			runtimeObjects = append(runtimeObjects, testAgent)
		}

		testSettings = clients.GetTestClients(clients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})

		if testCase.expectedError {
			assert.NotNil(t, err)
			assert.Equal(t, testCase.expectedErrorText, err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, testAgent.Name, result.Definition.Name)
			assert.Equal(t, testAgent.Namespace, result.Definition.Namespace)
			assert.Equal(t, testAgent.Name, result.Object.Name)
			assert.Equal(t, testAgent.Namespace, result.Object.Namespace)
			assert.Equal(t, result.Definition, result.Object)
			assert.Equal(t, testCase.hostname, result.Object.Spec.Hostname)
		}
	}
}

func buildTestBuilderWithFakeObjects(runtimeObjects []runtime.Object) (*agentBuilder, *clients.Settings) {
	testSettings := clients.GetTestClients(clients.TestClientParams{
		K8sMockObjects: runtimeObjects,
	})

	testBuilder := newAgentBuilder(testSettings.Client, nil)

	return testBuilder, testSettings
}
