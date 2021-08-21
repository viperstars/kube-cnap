package pod

import (
    "fmt"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "testing"
)

func TestHostAlias_ToK8sHostAliases(t *testing.T) {
    hostnames := make(map[string]string)
    hostnames["192.168.1.1"] = "abc.qiyi.so, bcd.qiyi.so , test.qiyi.so,    ddd.qiyi.so"
    hs := HostAlias{
        Meta:      basemeta.BaseMeta{},
        Hostnames: hostnames,
    }
    hsK8s := hs.ToK8sHostAliases()
    for _, x := range hsK8s {
        for _, y := range x.Hostnames {
            fmt.Println(y)
        }
    }

}
