package container

import v1 "k8s.io/api/core/v1"

type EnvVar struct {
    ID          int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int    `json:"containerID" xorm:"int 'container_id'"`
    Key         string `json:"key" xorm:"varchar(512)"`
    Value       string `json:"value" xorm:"varchar(512)"`
    Deleted     bool   `json:"deleted" xorm:"bool default 0"`
}

type EnvVars struct {
    EnvVars []*EnvVar `json:"envVars"`
}

func (e *EnvVars) ToK8sEnvVars() []v1.EnvVar {
    vars := make([]v1.EnvVar, 0, len(e.EnvVars))
    for _, ev := range e.EnvVars {
        envVar := v1.EnvVar{
            Name:  ev.Key,
            Value: ev.Value,
        }
        vars = append(vars, envVar)
    }
    return vars
}
