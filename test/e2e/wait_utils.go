// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file is inspired in https://github.com/operator-framework/operator-sdk/pull/1129

package e2e

import (
	"context"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForPod checks to see if a given pod is running after a specified amount of time.
// If the deployment is not running after timeout * retries seconds, the function returns an error
func WaitForPod(t *testing.T, runtimeClient test.FrameworkClient, namespace, name string, ignoreUID *types.UID, retryInterval, timeout time.Duration) error {
	pod := &corev1.Pod{}
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = runtimeClient.Get(context.TODO(), client.ObjectKey{Name: name, Namespace: namespace}, pod)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s pod\n", name)
				return false, nil
			}
			return false, err
		}

		if ignoreUID != nil && pod.UID == *ignoreUID {
			return false, nil
		}

		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		t.Logf("Waiting for full availability of %s pod\n", name)
		return false, nil
	})

	if err != nil {
		return err
	}

	t.Logf("%s Pod available\n", name)
	return nil
}
