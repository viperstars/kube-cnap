package pod

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "testing"
)

func TestBOSVolume_ToK8sVolume(t *testing.T) {
    bos := BOSVolume{
        ID: 0,
        Meta: basemeta.BaseMeta{
            App:    "test-app",
            Region: "bj",
            Env:    "test",
        },
        Name:      "bos-for-test",
        ClaimName: "pvc-name",
        //cSize:      "10Gi",
    }
    k8sBos := bos.ToK8sVolume()
    if k8sBos.Name != "bos-for-test" {
        t.Error("name error")
    }
    if k8sBos.PersistentVolumeClaim.ClaimName != "pvc-name" {
        t.Error("pvc name error")
    }
}

func TestCFSVolume_ToK8sVolume(t *testing.T) {
    cfs := CFSVolume{
        Name:     "cfs-test",
        //cSize:     "",
        Server:   "nfs.test.com",
        ReadOnly: false,
        Path:     "/",
    }
    k8sCfs := cfs.ToK8sVolume()
    if k8sCfs.Name != cfs.Name {
        t.Error("name error")
    }
    if k8sCfs.NFS.Server != cfs.Server {
        t.Error("server error")
    }
    if k8sCfs.NFS.Path != cfs.Path {
        t.Error("server error")
    }
}

func TestEmptyVolume_ToK8sVolume(t *testing.T) {
    ev := EmptyVolume{
        Name: "test-empty",
        //Size: 100,
    }
    k8sEv := ev.ToK8sVolume()
    if k8sEv.Name != ev.Name {
        t.Error("name error")
    }
    if k8sEv.EmptyDir.SizeLimit.String() != "100" {
        t.Error("size error")
    }
}
