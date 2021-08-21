package http

import (
    "github.com/viperstars/bce-golang/bce_client"
    "github.com/viperstars/kube-cnap/conf"
    "github.com/viperstars/kube-cnap/pkg/cmdb"
    "github.com/viperstars/kube-cnap/pkg/image"
    "go.uber.org/zap"
    "gopkg.in/cas.v1"
)

func NewServer(client K8sClient, logger *zap.Logger, service DBService, config *conf.Server, casClient *cas.Client,
    bceClient *bce_client.BCEClient, reg *image.RegistryClient) *Server {
    g := cmdb.NewUserGetter(logger)
    server := &Server{
        dbService:  service,
        k8sClients: client,
        config:     config,
        logger:     logger,
        cas:        casClient,
        cmdb:       g,
        bce:        bceClient,
        reg:        reg,
    }
    return server
}
