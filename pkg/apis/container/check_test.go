package container

import (
    "fmt"
    "reflect"
    "strings"
    "testing"
)

func TestLivenessCheck_ToK8sProbe(t *testing.T) {
    lc := LivenessCheck{
        ID:          0,
        ContainerID: 0,
        Param: CheckParam{
            InitialDelaySeconds: 1000,
            TimeoutSeconds:      20,
            PeriodSeconds:       60,
            SuccessThreshold:    5,
            FailureThreshold:    3,
        },
        Command: "ls /etc/init.d/nginx",
    }
    probe := lc.ToK8sProbe()
    fmt.Println(probe)
    if reflect.DeepEqual(probe.Exec.Command, []string{"ls", "/etc/init.d/nginx"}) {
        t.Error("convert error")
    }
}

func TestReadnessCheck_ToK8sProbe(t *testing.T) {
    lc := ReadnessCheck{
        ID:          0,
        ContainerID: 0,
        Param: CheckParam{
            InitialDelaySeconds: 1000,
            TimeoutSeconds:      20,
            PeriodSeconds:       60,
            SuccessThreshold:    5,
            FailureThreshold:    3,
        },
        Path: "http://127.0.0.1:8080/api/check?test=1",
    }
    probe := lc.ToK8sProbe()
    fmt.Println(probe)
    fmt.Println("port", probe.HTTPGet.Port.IntVal)
    if probe.HTTPGet.Port.IntVal != 8080 {
        t.Error("convert error")
    }
    fmt.Println("host",  probe.HTTPGet.Host)
    if probe.HTTPGet.Host != "cmdb.qiyi.so" {
        t.Error("convert error")
    }
    if probe.HTTPGet.Path != "/api/check" {
        t.Error("convert error")
    }
    if strings.ToLower(string(probe.HTTPGet.Scheme)) != "http" {
        t.Error("convert error")
    }
}

func TestReadnessCheck_NilToK8sProbe(t *testing.T) {
    lc := new(ReadnessCheck)
    probe := lc.ToK8sProbe()
    fmt.Println(probe)
    if probe.HTTPGet.Port.IntVal != 8080 {
        t.Error("convert error")
    }
    if probe.HTTPGet.Host != "cmdb.qiyi.so:8080" {
        t.Error("convert error")
    }
    if probe.HTTPGet.Path != "/api/check" {
        t.Error("convert error")
    }
    if strings.ToLower(string(probe.HTTPGet.Scheme)) != "http" {
        t.Error("convert error")
    }
}
