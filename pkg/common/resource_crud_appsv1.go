package common

import (
	"context"
	"errors"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const (
	DeploymentType  = "deployment"
	DaemonSetType   = "daemonset"
	StatefulSetType = "statefulset"
	ReplicaSetType  = "replicaset"
	// Add other resource types as needed
)

type ResourceType string

type ResourceCRUD interface {
	Exists(namespace, name string) (interface{}, error)
	Create(namespace, name string, obj interface{}) (interface{}, error)
	Update(namespace, name string, obj interface{}) (interface{}, error)
	Delete(obj interface{}) error
}

type AppsV1ResourceCRUD struct {
	Client appsv1client.AppsV1Interface

	StructType ResourceType
}

func (r *AppsV1ResourceCRUD) Exists(namespace, name string) (interface{}, error) {
	switch r.StructType {
	case DeploymentType:
		deployment, err := r.Client.Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return deployment, nil
	case DaemonSetType:
		daemonset, err := r.Client.DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return daemonset, nil
	case StatefulSetType:
		statefulset, err := r.Client.StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return statefulset, nil
	case ReplicaSetType:
		replicaset, err := r.Client.ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} 
		return replicaset, nil
	default:
		return nil, errors.New("unsupported resource type")
	}
}

func (r *AppsV1ResourceCRUD) Create(namespace, name string, obj interface{}) (interface{}, error) {
	switch resource := obj.(type) {
	case *appsv1.Deployment:
		glog.V(100).Infof("Creating deployment %s in namespace %s", name, namespace)
		return r.Client.Deployments(namespace).Create(context.TODO(), resource, metav1.CreateOptions{})
	case *appsv1.DaemonSet:
		glog.V(100).Infof("Creating daemonset %s in namespace %s", name, namespace)
		return r.Client.DaemonSets(namespace).Create(context.TODO(), resource, metav1.CreateOptions{})
	case *appsv1.StatefulSet:
		glog.V(100).Infof("Creating statefulset %s in namespace %s", name, namespace)
		return r.Client.StatefulSets(namespace).Create(context.TODO(), resource, metav1.CreateOptions{})
	case *appsv1.ReplicaSet:
		glog.V(100).Infof("Creating replicaset %s in namespace %s", name, namespace)
		return r.Client.ReplicaSets(namespace).Create(context.TODO(), resource, metav1.CreateOptions{})
	default:
		return nil, errors.New("unsupported resource type")
	}
}

func (r *AppsV1ResourceCRUD) Update(namespace, name string, obj interface{}) (interface{}, error) {
	switch resource := obj.(type) {
	case *appsv1.Deployment:
		glog.V(100).Infof("Updating deployment %s in namespace %s", name, namespace)
		return r.Client.Deployments(namespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	case *appsv1.DaemonSet:
		glog.V(100).Infof("Updating daemonset %s in namespace %s", name, namespace)
		return r.Client.DaemonSets(namespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	case *appsv1.StatefulSet:
		glog.V(100).Infof("Updating statefulset %s in namespace %s", name, namespace)
		return r.Client.StatefulSets(namespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	case *appsv1.ReplicaSet:
		glog.V(100).Infof("Updating replicaset %s in namespace %s", name, namespace)
		return r.Client.ReplicaSets(namespace).Update(context.TODO(), resource, metav1.UpdateOptions{})
	default:
		return nil, errors.New("unsupported resource type")
	}
}

func (r *AppsV1ResourceCRUD) Delete(obj interface{}) error {
	switch obj.(type) {
	case *appsv1.Deployment:
		name := obj.(*appsv1.Deployment).Name
		namespace := obj.(*appsv1.Deployment).Namespace
		glog.V(100).Infof("Deleting deployment %s in namespace %s", name, namespace)
		err := r.Client.Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if k8serrors.IsNotFound(err) {
			glog.V(100).Infof("Deployment %s in namespace %s not found, ignoring", name, namespace)
			return nil
		}
		return err
	case *appsv1.DaemonSet:
		name := obj.(*appsv1.DaemonSet).Name
		namespace := obj.(*appsv1.DaemonSet).Namespace
		glog.V(100).Infof("Deleting daemonset %s in namespace %s", name, namespace)
		err := r.Client.DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if k8serrors.IsNotFound(err) {
			glog.V(100).Infof("DaemonSet %s in namespace %s not found, ignoring", name, namespace)
			return nil
		}
		return err
	case *appsv1.StatefulSet:
		name := obj.(*appsv1.StatefulSet).Name
		namespace := obj.(*appsv1.StatefulSet).Namespace
		glog.V(100).Infof("Deleting statefulset %s in namespace %s", name, namespace)
		err := r.Client.StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if k8serrors.IsNotFound(err) {
			glog.V(100).Infof("StatefulSet %s in namespace %s not found, ignoring", name, namespace)
			return nil
		}
		return err
	case *appsv1.ReplicaSet:
		name := obj.(*appsv1.ReplicaSet).Name
		namespace := obj.(*appsv1.ReplicaSet).Namespace
		glog.V(100).Infof("Deleting replicaset %s in namespace %s", name, namespace)
		err := r.Client.ReplicaSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if k8serrors.IsNotFound(err) {
			glog.V(100).Infof("ReplicaSet %s in namespace %s not found, ignoring", name, namespace)
			return nil
		}
		return err
	default:
		return errors.New("unsupported resource type")
	}
}
