package clients

import (
	"context"
	"errors"
	"fmt"
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errSchemeAttach   = errors.New("scheme error")
	errSchemeAttacher = func(s *runtime.Scheme) error { return errSchemeAttach }
)

// countingSchemeAttacher returns a no-op SchemeAttacher that increments count each time it runs.
func countingSchemeAttacher(count *int) SchemeAttacher {
	return func(*runtime.Scheme) error {
		*count++

		return nil
	}
}

func TestNew(t *testing.T) {
	result := New("/nonexistent/kubeconfig.yaml")
	assert.Nil(t, result)
}

func TestSetScheme(t *testing.T) {
	s := runtime.NewScheme()
	err := SetScheme(s)
	require.NoError(t, err)
	assert.True(t, s.Recognizes(corev1.SchemeGroupVersion.WithKind("Pod")))
	assert.True(t, s.Recognizes(routev1.GroupVersion.WithKind("Route")))
}

func TestGetAPIClient(t *testing.T) {
	testCases := []struct {
		name          string
		nilClient     bool
		expectedError error
	}{
		{
			name:          "nil settings",
			nilClient:     true,
			expectedError: fmt.Errorf("APIClient cannot be nil"),
		},
		{
			name: "valid settings",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var settings *Settings
			if !testCase.nilClient {
				settings = GetTestClients(TestClientParams{})
			}

			result, err := settings.GetAPIClient()
			assert.Equal(t, testCase.expectedError, err)

			if testCase.expectedError == nil {
				assert.Same(t, settings, result)
			}
		})
	}
}

func TestAttachScheme(t *testing.T) {
	testCases := []struct {
		name          string
		nilClient     bool
		failing       bool
		expectedError error
		expectedCalls int
	}{
		{
			name:          "nil settings",
			nilClient:     true,
			expectedError: fmt.Errorf("cannot add scheme to nil client"),
		},
		{
			name:          "valid attacher",
			expectedCalls: 1,
		},
		{
			name:          "attacher returns error",
			failing:       true,
			expectedError: errSchemeAttach,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var settings *Settings
			if !testCase.nilClient {
				settings = GetTestClients(TestClientParams{})
			}

			calls := 0
			attacher := countingSchemeAttacher(&calls)

			if testCase.failing {
				attacher = errSchemeAttacher
			}

			err := settings.AttachScheme(attacher)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedCalls, calls)
		})
	}
}

func TestGetTestClients(t *testing.T) {
	testPod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "test-namespace"}}

	testCases := []struct {
		name          string
		attacherCount int
		objects       []runtime.Object
	}{
		{
			name:          "empty params",
			attacherCount: 0,
		},
		{
			name:          "with scheme attacher",
			attacherCount: 1,
		},
		{
			name:          "with multiple scheme attachers",
			attacherCount: 2,
		},
		{
			name:          "with scheme attacher and mock objects",
			attacherCount: 1,
			objects:       []runtime.Object{testPod},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			callCounts := make([]int, testCase.attacherCount)

			var attachers []SchemeAttacher
			for index := range callCounts {
				attachers = append(attachers, countingSchemeAttacher(&callCounts[index]))
			}

			result := GetTestClients(TestClientParams{K8sMockObjects: testCase.objects, SchemeAttachers: attachers})
			require.NotNil(t, result)
			assert.NotNil(t, result.Client)
			assert.NotNil(t, result.K8sClient)

			for index, count := range callCounts {
				assert.Equalf(t, 1, count, "scheme attacher %d should run exactly once", index)
			}

			for _, object := range testCase.objects {
				clientObject, ok := object.(runtimeclient.Object)
				require.True(t, ok)

				err := result.Get(context.TODO(), runtimeclient.ObjectKeyFromObject(clientObject), &corev1.Pod{})
				assert.NoError(t, err, "mock object should be served by the generic client")
			}
		})
	}
}
