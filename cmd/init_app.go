package main

import "github.com/viperstars/kube-cnap/pkg/apis/basemeta"

type dbService interface {
    CreateContainerGroup(base basemeta.BaseMeta) (int64, error)
}
type k8sClient interface {
    CreateNamespace(base basemeta.BaseMeta) error
    CreateSecret(base basemeta.BaseMeta) error
    CreateClusterRole(base basemeta.BaseMeta) error
}

func initApp(base basemeta.BaseMeta) {

}