package common

type mockAppsV1Client struct {
	createFunc func(string, string, interface{}) error
	updateFunc func(string, string, interface{}) error
	deleteFunc func(interface{}) error
	existsFunc func(string, string) bool
}

func (m *mockAppsV1Client) Create(namespace, name string, obj interface{}) error {
	return m.createFunc(namespace, name, obj)
}

func (m *mockAppsV1Client) Update(namespace, name string, obj interface{}) error {
	return m.updateFunc(namespace, name, obj)
}

func (m *mockAppsV1Client) Delete(obj interface{}) error {
	return m.deleteFunc(obj)
}

func (m *mockAppsV1Client) Exists(namespace, name string) bool {
	return m.existsFunc(namespace, name)
}
