package basemeta

import (
    "fmt"
    "go.uber.org/zap/zapcore"
)

type BaseMeta struct {
    App    string `json:"app" xorm:"varchar(256) 'app'"`
    Region string `json:"region" xorm:"varchar(256) 'region'"`
    Env    string `json:"env" xorm:"varchar(256) 'env'"`
}

func (b *BaseMeta) GetPrefix() string {
    return fmt.Sprintf("%s-%s-%s", b.Region, b.Env, b.App)
}

func (b *BaseMeta) GetLabels() map[string]string {
    labels := make(map[string]string)
    labels["app"] = b.App
    labels["region"] = b.Region
    labels["env"] = b.Env
    return labels
}

func (b *BaseMeta) GetQuery() []interface{} {
    return []interface{}{`app = ? AND region = ? AND env = ?`, b.App, b.Region, b.Env}
}

func (b *BaseMeta) MarshalLogObject(enc zapcore.ObjectEncoder) error {
    enc.AddString("app", b.App)
    enc.AddString("region", b.Region)
    enc.AddString("env", b.Env)
    return nil
}

type CNames struct {
    App    string `json:"app"`
    Region string `json:"region"`
    RCname string `json:"rCname"`
    Env    string `json:"env"`
    ECname string `json:"eCname"`
}
