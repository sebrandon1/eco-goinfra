package deployment

import (
	"fmt"
	"testing"

	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const testDeploymentLabelKey = "app"

var (
	defaultDeploymentName      = "test-deployment"
	defaultDeploymentNamespace = "test-namespace"
	otherDeploymentName        = "other-deployment"
	otherDeploymentNamespace   = "other-namespace"

	matchingLabelSelector = metav1.ListOptions{LabelSelector: testDeploymentLabelKey + "=match"}

	// testDeployments has one deployment in each of two namespaces, with different labels, so that tests can
	// observe both namespace scoping and label selector forwarding.
	testDeployments = []runtime.Object{
		buildDummyDeployment(defaultDeploymentName, defaultDeploymentNamespace, "match"),
		buildDummyDeployment(otherDeploymentName, otherDeploymentNamespace, "other"),
	}
)

func buildDummyDeployment(name, namespace, labelValue string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{testDeploymentLabelKey: labelValue},
		},
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		name          string
		nsname        string
		objects       []runtime.Object
		listOptions   []metav1.ListOptions
		expectedError error
		expected      []types.NamespacedName
	}{
		{
			name:          "empty namespace",
			nsname:        "",
			expectedError: fmt.Errorf("failed to list deployments, 'nsname' parameter is empty"),
		},
		{
			name:          "too many list options",
			nsname:        defaultDeploymentNamespace,
			listOptions:   []metav1.ListOptions{{}, {}},
			expectedError: fmt.Errorf("error: more than one ListOptions was passed"),
		},
		{
			name:   "no deployments",
			nsname: defaultDeploymentNamespace,
		},
		{
			name:     "only deployments in requested namespace",
			nsname:   otherDeploymentNamespace,
			objects:  testDeployments,
			expected: []types.NamespacedName{{Name: otherDeploymentName, Namespace: otherDeploymentNamespace}},
		},
		{
			name:        "label selector excludes non-matching deployment",
			nsname:      otherDeploymentNamespace,
			objects:     testDeployments,
			listOptions: []metav1.ListOptions{matchingLabelSelector},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testClient := clients.GetTestClients(clients.TestClientParams{
				K8sMockObjects: testCase.objects,
			})

			result, err := List(testClient, testCase.nsname, testCase.listOptions...)
			assert.Equal(t, testCase.expectedError, err)

			if testCase.expectedError == nil {
				assert.ElementsMatch(t, testCase.expected, deploymentKeys(result))
			}
		})
	}
}

func TestListInAllNamespaces(t *testing.T) {
	testCases := []struct {
		name          string
		objects       []runtime.Object
		listOptions   []metav1.ListOptions
		expectedError error
		expected      []types.NamespacedName
	}{
		{
			name:          "too many list options",
			listOptions:   []metav1.ListOptions{{}, {}},
			expectedError: fmt.Errorf("error: more than one ListOptions was passed"),
		},
		{
			name: "no deployments",
		},
		{
			name:    "deployments from all namespaces",
			objects: testDeployments,
			expected: []types.NamespacedName{
				{Name: defaultDeploymentName, Namespace: defaultDeploymentNamespace},
				{Name: otherDeploymentName, Namespace: otherDeploymentNamespace},
			},
		},
		{
			name:        "label selector filters deployments",
			objects:     testDeployments,
			listOptions: []metav1.ListOptions{matchingLabelSelector},
			expected:    []types.NamespacedName{{Name: defaultDeploymentName, Namespace: defaultDeploymentNamespace}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testClient := clients.GetTestClients(clients.TestClientParams{
				K8sMockObjects: testCase.objects,
			})

			result, err := ListInAllNamespaces(testClient, testCase.listOptions...)
			assert.Equal(t, testCase.expectedError, err)

			if testCase.expectedError == nil {
				assert.ElementsMatch(t, testCase.expected, deploymentKeys(result))
			}
		})
	}
}

func deploymentKeys(builders []*Builder) []types.NamespacedName {
	var keys []types.NamespacedName

	for _, builder := range builders {
		keys = append(keys, types.NamespacedName{Name: builder.Object.Name, Namespace: builder.Object.Namespace})
	}

	return keys
}
