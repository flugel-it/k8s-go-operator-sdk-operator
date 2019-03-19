package immortalcontainer

import (
	"context"
	"testing"

	immv1alpha1 "github.com/flugel-it/k8s-go-operator-sdk-operator/pkg/apis/immortalcontainer/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)


func TestImmortalContainerControllerPodCreate(t *testing.T) {
	var (
		name      = "example"
		image     = "nginx:latest"
		namespace = "testnamespace"
	)

	immortalContainer := &immv1alpha1.ImmortalContainer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: immv1alpha1.ImmortalContainerSpec{
			Image: image,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		immortalContainer,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(immv1alpha1.SchemeGroupVersion, immortalContainer)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileImmortalContainer object with the scheme and fake client.
	r := &ReconcileImmortalContainer{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
    
    res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	if res != (reconcile.Result{}) {
		t.Error("reconcile did not return an empty Result")
	}

	// Check the pod is created
	expectedPod := newPodForImmortalContainer(immortalContainer)
	pod := &corev1.Pod{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedPod.Name, Namespace: expectedPod.Namespace}, pod)
	if err != nil {
		t.Fatalf("get pod: (%v)", err)
	}

	// Check status is correctly updated
	updatedImmortalContainer := &immv1alpha1.ImmortalContainer{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: immortalContainer.Name, Namespace: immortalContainer.Namespace}, updatedImmortalContainer)
	if err != nil {
		t.Fatalf("get immortal container: (%v)", err)
	}
	if updatedImmortalContainer.Status.StartTimes != 1 {
		t.Errorf("incorrect immortal container startTimes: (%v)", updatedImmortalContainer.Status.StartTimes)
	}
	if updatedImmortalContainer.Status.CurrentPod != expectedPod.Name {
		t.Errorf("incorrect immortal container currentPod: (%v)", updatedImmortalContainer.Status.CurrentPod)
	}
}
