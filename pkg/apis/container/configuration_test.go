package container

import (
    "fmt"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "testing"
)

var m = &basemeta.BaseMeta{
    Region: "bj",
    Env:    "test",
    App:    "test-app",
}

var configurations = []*Configuration{
    {
        Path: "/etc/test/mysql.conf",
        Content: `[client]
port = 3306
socket = /tmp/mysql.sock 

[mysqld]

server-id = 1

user = mysql`,
    },
    {
        Path: "/etc/test/php.conf",
        Content: `[client]
port = 3306
socket = /tmp/mysql.sock
server-id = 1
user = mysql`,
    },
    {
        Path: "/etc/abc/my.conf",
        Content: `[client]
port = 3306
socket = /tmp/mysql.sock
server-id = 1

user = mysql`,
    },
}

func TestConfigurations_GetSameParentConfigurations(t *testing.T) {
    cfgs := &Configurations{
        Configurations: configurations,
    }
    c := cfgs.GetSameDirnameConfigurations()
    if len(c) != 2 {
        t.Error("get same dir error")
    }
    if value, ok := c["/etc/test"]; ok {
        if len(value) != 2 {
            t.Error("length error")
        }
    } else {
        t.Error("dir error")
    }
    if value, ok := c["/etc/abc"]; ok {
        if len(value) != 1 {
            t.Error("length error")
        }
    } else {
        t.Error("dir error")
    }
}

func TestConfigurations_ToK8sObjects(t *testing.T) {
    cfgs := &Configurations{
        Configurations: configurations,
    }
    configMaps, volumeMounts, volumes := cfgs.ToK8sObjects(*m)
    if len(configMaps) != 2 {
        t.Error("configMaps error")
    } else {
        for _, configMap := range configMaps {
            fmt.Println(configMap.Name)
            fmt.Println(configMap.Data)
        }
    }
    if len(volumes) != 2 {
        t.Error("volumes error")
    } else {
        for _, volume := range volumes {
            fmt.Println(volume.Name)
            fmt.Println(volume.ConfigMap.Name)
        }
    }
    if len(volumeMounts) != 2 {
        t.Error("volumeMounts error")
    } else {
        for _, volumeMount := range volumeMounts {
            fmt.Println(volumeMount.Name)
            fmt.Println(volumeMount.MountPath)
        }
    }
}
