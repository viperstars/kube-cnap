package image

import (
    "fmt"
    "testing"
)

var rc *RegistryClient

func init() {
    var err error
    rc, err = NewRegistryClient("https://docker-registry.xx.com", "", "password")
    if err != nil {
        return
    }
}

func TestRegistryClient_Tags(t *testing.T) {
    tags, err := rc.Tags("penghao-test","saver")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(tags)
}

func TestRegistryClient_NoTags(t *testing.T) {
    tags, err := rc.Tags("penghao-test","savor")
    if err == nil {
        t.Error(err)
    }
    fmt.Println(tags)
}
