package container

import "testing"

func TestContainer_ToK8sVolumeMounts(t *testing.T) {
    vms := VolumeMounts{VolumeMounts: []*VolumeMount{
        {VolumeName: "cfs-test", MountPath: "/etc/host", SubPath: "abc"},
        {VolumeName: "empty-test", MountPath: "/tmp", SubPath: "test", ReadOnly: true},
    }}
    k8sVMs := vms.ToK8sVolumeMounts()
    if len(k8sVMs) != 2 {
        t.Error("length error")
    }
    if !k8sVMs[1].ReadOnly {
        t.Error("index error")
    }
}
