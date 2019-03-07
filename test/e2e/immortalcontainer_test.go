package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	"github.com/flugel-it/immortalcontainer-operator/pkg/apis"
	operator "github.com/flugel-it/immortalcontainer-operator/pkg/apis/immortalcontainer/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestImmortalContainer(t *testing.T) {
	immortalContainerList := &operator.ImmortalContainerList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImmortalContainer",
			APIVersion: "immortalcontainer.flugel.it/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, immortalContainerList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("immortalcontainer-group", func(t *testing.T) {
		t.Run("Cluster", ImmortalContainerCluster)
	})
}

func ImmortalContainerCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "immortalcontainer-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = immortalContainerCreateTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}

func immortalContainerCreateTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create custom resource
	immortalContainer := &operator.ImmortalContainer{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImmortalContainer",
			APIVersion: "immortalcontainer.flugel.it/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example",
			Namespace: namespace,
		},
		Spec: operator.ImmortalContainerSpec{
			Image: "nginx:latest",
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), immortalContainer, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	// wait for pod
	err = WaitForPod(t, f.Client, namespace, "example-immortalpod", nil, retryInterval, timeout)
	if err != nil {
		return err
	}

	pod := &corev1.Pod{}
	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-immortalpod", Namespace: namespace}, pod)
	if err != nil {
		return err
	}

	pod1UID := &pod.UID

	err = f.Client.Delete(goctx.TODO(), pod)
	if err != nil {
		return err
	}

	// wait for pod recreation
	err = WaitForPod(t, f.Client, namespace, "example-immortalpod", pod1UID, retryInterval, timeout)
	if err != nil {
		return err
	}

	return nil
}
