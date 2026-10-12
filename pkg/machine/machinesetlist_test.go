package machine

import (
	"fmt"
	"testing"

	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultMachineSetNamespace  = "openshift-machine-api"
	defaultMachineSetWorkerRole = "machine.openshift.io/cluster-api-machine-role"
)

func TestListWorkerMachineSets(t *testing.T) {
	testCases := []struct {
		name          string
		namespace     string
		workerLabel   string
		listOptions   []metav1.ListOptions
		expectedError error
	}{
		{
			name:          "empty namespace",
			namespace:     "",
			workerLabel:   defaultMachineSetWorkerRole,
			expectedError: fmt.Errorf("failed to list MachineSets, 'namespace' parameter is empty"),
		},
		{
			name:          "empty workerLabel",
			namespace:     defaultMachineSetNamespace,
			workerLabel:   "",
			expectedError: fmt.Errorf("failed to list MachineSets, 'workerLabel' parameter is empty"),
		},
		{
			name:          "too many list options",
			namespace:     defaultMachineSetNamespace,
			workerLabel:   defaultMachineSetWorkerRole,
			listOptions:   []metav1.ListOptions{{}, {}},
			expectedError: fmt.Errorf("error: more than one ListOptions was passed"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testClient := clients.GetTestClients(clients.TestClientParams{})

			_, err := ListWorkerMachineSets(testClient, testCase.namespace, testCase.workerLabel,
				testCase.listOptions...)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
