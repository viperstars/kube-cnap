package conf

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

var Config *Conf

type Conf struct {
    Server       *Server         `json:"server"`
    BceConfig    *BceConfig      `json:"bceConfig"`
    Logging      *Logging        `json:"logging"`
    DB           *DB             `json:"db"`
    Registry     *RegistryServer `json:"registry"`
    RegistryInfo *RegistryInfo   `json:"registryInfo"`
    Cas          *Cas            `json:"cas"`
}

type Server struct {
    Host string `json:"host"`
    Port int    `json:"port"`
    Mode string `json:"mode"`
}

type BceConfig struct {
    Ak string `json:"ak"`
    Sk string `json:"sk"`
}

type Logging struct {
    Path string `json:"path"`
}

type DB struct {
    Dsn string `json:"dsn"`
}

type Cas struct {
    Server      string `json:"server"`
    CasCallback string `json:"casCallback"`
}

type RegistryServer struct {
    Auth          *AuthConfig `json:"auth"`
    PullPrefix    string      `json:"pullPrefix"`
    KeyName       string      `json:"keyName"`
    ServerAddress string      `json:"serveraddress"`
}

type RegistryInfo struct {
    Auth     *AuthConfig `json:"auth"`
    Url      string `json:"url"`
    Addr     string `json:"addr"`
    KeyName       string      `json:"keyName"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthConfig struct {
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
    Auth     string `json:"auth"`
    Email    string `json:"email"`
}

func init() {
    NewConfig()
}

func NewConfig() {
    content, err := ioutil.ReadFile("conf/config.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(10)
    }
    config := &Conf{}
    err = json.Unmarshal(content, &config)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(10)
    } else {
        Config = config
    }
}
