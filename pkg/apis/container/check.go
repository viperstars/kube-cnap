package container

import (
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
    "strings"
)

type CheckParam struct {
    InitialDelaySeconds int32 `json:"initialDelaySeconds" xorm:"int"`
    TimeoutSeconds      int32 `json:"timeoutSeconds" xorm:"int"`
    PeriodSeconds       int32 `json:"periodSeconds" xorm:"int"`
    SuccessThreshold    int32 `json:"successThreshold" xorm:"int"`
    FailureThreshold    int32 `json:"failureThreshold" xorm:"int"`
}

type ReadnessCheck struct {
    ID                  int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID         int    `json:"containerID" xorm:"int 'container_id'"`
    InitialDelaySeconds int32  `json:"initialDelaySeconds" xorm:"int"`
    TimeoutSeconds      int32  `json:"timeoutSeconds" xorm:"int"`
    PeriodSeconds       int32  `json:"periodSeconds" xorm:"int"`
    SuccessThreshold    int32  `json:"successThreshold" xorm:"int"`
    FailureThreshold    int32  `json:"failureThreshold" xorm:"int"`
    Path                string `json:"path" xorm:"varchar(512) 'path'"`
    Port                int    `json:"port" xorm:"int"`
    Scheme              string `json:"scheme" xorm:"varchar(512)"`
    Deleted             bool   `json:"deleted" xorm:"bool default 0"`
}

func (r *ReadnessCheck) ToK8sProbe() *v1.Probe {
    var scheme v1.URIScheme
    if r.Scheme == strings.ToLower(string(v1.URISchemeHTTPS)) {
        scheme = v1.URISchemeHTTPS
    } else {
        scheme = v1.URISchemeHTTP
    }
    probe := &v1.Probe{
        Handler: v1.Handler{
            HTTPGet: &v1.HTTPGetAction{
                Path: r.Path,
                Port: intstr.IntOrString{
                    IntVal: int32(r.Port),
                },
                Scheme:      scheme,
                HTTPHeaders: nil,
            },
        },
        InitialDelaySeconds: r.InitialDelaySeconds,
        TimeoutSeconds:      r.TimeoutSeconds,
        PeriodSeconds:       r.PeriodSeconds,
        SuccessThreshold:    r.SuccessThreshold,
        FailureThreshold:    r.FailureThreshold,
    }
    return probe
}

type LivenessCheckHttp struct {
    ID                  int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID         int    `json:"containerID" xorm:"int 'container_id'"`
    InitialDelaySeconds int32  `json:"initialDelaySeconds" xorm:"int"`
    TimeoutSeconds      int32  `json:"timeoutSeconds" xorm:"int"`
    PeriodSeconds       int32  `json:"periodSeconds" xorm:"int"`
    SuccessThreshold    int32  `json:"successThreshold" xorm:"int"`
    FailureThreshold    int32  `json:"failureThreshold" xorm:"int"`
    Path                string `json:"path" xorm:"varchar(512) 'path'"`
    Port                int    `json:"port" xorm:"int"`
    Scheme              string `json:"scheme" xorm:"varchar(512)"`
    Deleted             bool   `json:"deleted" xorm:"bool default 0"`
}

func (r *LivenessCheckHttp) ToK8sProbe() *v1.Probe {
    var scheme v1.URIScheme
    if r.Scheme == strings.ToLower(string(v1.URISchemeHTTPS)) {
        scheme = v1.URISchemeHTTPS
    } else {
        scheme = v1.URISchemeHTTP
    }
    probe := &v1.Probe{
        Handler: v1.Handler{
            HTTPGet: &v1.HTTPGetAction{
                Path: r.Path,
                Port: intstr.IntOrString{
                    IntVal: int32(r.Port),
                },
                Scheme:      scheme,
                HTTPHeaders: nil,
            },
        },
        InitialDelaySeconds: r.InitialDelaySeconds,
        TimeoutSeconds:      r.TimeoutSeconds,
        PeriodSeconds:       r.PeriodSeconds,
        SuccessThreshold:    r.SuccessThreshold,
        FailureThreshold:    r.FailureThreshold,
    }
    return probe
}

type LivenessCheck struct {
    ID          int        `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int        `json:"containerID" xorm:"int 'container_id'"`
    InitialDelaySeconds int32 `json:"initialDelaySeconds" xorm:"int"`
    TimeoutSeconds      int32 `json:"timeoutSeconds" xorm:"int"`
    PeriodSeconds       int32 `json:"periodSeconds" xorm:"int"`
    SuccessThreshold    int32 `json:"successThreshold" xorm:"int"`
    FailureThreshold    int32 `json:"failureThreshold" xorm:"int"`
    Command     string     `json:"command" xorm:"varchar(512) 'command'"`
    Deleted     bool       `json:"deleted" xorm:"bool default 0"`
}

func (l *LivenessCheck) ToK8sProbe() *v1.Probe {
    if l.Command == "" {
        return nil
    }
    split := separator.Split(l.Command, -1)
    probe := &v1.Probe{
        Handler: v1.Handler{
            Exec: &v1.ExecAction{
                Command: split,
            },
        },
        InitialDelaySeconds: l.InitialDelaySeconds,
        TimeoutSeconds:      l.TimeoutSeconds,
        PeriodSeconds:       l.PeriodSeconds,
        SuccessThreshold:    1,
        FailureThreshold:    l.FailureThreshold,
    }
    return probe
}
