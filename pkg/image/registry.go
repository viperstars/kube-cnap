package image

import (
    "fmt"
    "github.com/heroku/docker-registry-client/registry"
    "strings"
)

type RegistryClient struct {
    username string
    password string
    url      string
    trimmedUrl string
    hub      *registry.Registry
}

func(rc *RegistryClient) Tags(namespace, repo string) ([]string, error){
    var tags []string
    r := fmt.Sprintf("%s/%s", namespace, repo)
    tags, err := rc.hub.Tags(r)
    if err != nil {
        return tags, err
    }
    modTags := make([]string, 0, len(tags))
    for _, tag := range tags {
        t := fmt.Sprintf("%s/%s:%s", rc.trimmedUrl, r, tag)
        modTags = append(modTags, t)
    }
    return modTags, err
}

func(rc *RegistryClient) prefix() {
    rc.trimmedUrl = strings.TrimPrefix(rc.url, "https://")
}

func(rc *RegistryClient) initHub() error {
    hub, err := registry.NewInsecure(rc.url, rc.username, rc.password)
    rc.hub = hub
    return err
}

func NewRegistryClient(url, username, password string) (*RegistryClient, error) {
    rc := &RegistryClient{
        username:   username,
        password:   password,
        url:        url,
    }
    rc.prefix()
    err := rc.initHub()
    return rc, err
}
