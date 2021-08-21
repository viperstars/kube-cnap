package consts

import (
    "github.com/viperstars/kube-cnap/conf"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
)

const (
    DEFAULTPOLICY                    = v1.PullAlways
    DEFAULTSERVICETYPE               = v1.ServiceTypeClusterIP
    DEFAULTROLLINGUPDATESTRATEGYTYPE = appsv1.RollingUpdateDeploymentStrategyType
    DEFAULTSERVICEAFFINITY           = v1.ServiceAffinityClientIP
    DEFAULTMESSAGEPATH               = "/dev/termination-log"
    DEFAULTMESSAGEPOLICY             = v1.TerminationMessageFallbackToLogsOnError
    // DEFAULTSECRETNAMEFORSECRET       = "baidu-image-secret"
    DEFAULTSTICKYSESSIONKEY = "traefik.ingress.kubernetes.io/affinity"
)

var DEFAULTSECRETNAME = []v1.LocalObjectReference{{conf.Config.Registry.KeyName},
    {Name: conf.Config.RegistryInfo.KeyName}} // TODO
