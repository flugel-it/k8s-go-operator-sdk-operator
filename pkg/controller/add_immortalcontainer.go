package controller

import (
	"github.com/flugel-it/k8s-go-operator-sdk-operator/pkg/controller/immortalcontainer"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, immortalcontainer.Add)
}
