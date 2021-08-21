package image

import (
    "encoding/json"
    "github.com/viperstars/bce-golang/bce_client"
    "io/ioutil"
)

type Response struct {
    Namespace  string `json:"namespace"`
    Repository string `json:"repository"`
    Tags       []*Tag `json:"tags"`
}

type Tag struct {
    Name        string `json:"name"`
    Digest      string `json:"digest"`
    Description string `json:"description"`
    CreateAt    string `json:"create_at"`
    UpdateAt    string `json:"update_at"`
}

type Getter struct {
    AccessKey string
    SecretKey string
}

func (g *Getter) initClient(region string, instance string) *bce_client.BCEClient {
    return bce_client.NewBCEClient(
        g.AccessKey,
        g.SecretKey,
        region,
        instance,
    )
}

func (g *Getter) GetData(region, instance, method, url, namespace, repo string) (*Response, error) {
    tagsResponse := new(Response)
    client := g.initClient(region, "cce")
    param := make(map[string]string)
    param["namespace"] = namespace
    param["repository"] = repo
    resp, _, err := client.Execute(method, url, param, nil, nil)
    if err != nil {
        return tagsResponse, err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return tagsResponse, err
    }
    err = json.Unmarshal(body, &tagsResponse)
    return tagsResponse, err
}
