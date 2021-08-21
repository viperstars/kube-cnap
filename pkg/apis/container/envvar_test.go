package container

import "testing"

func TestEnvVars_ToK8sEnvVars(t *testing.T) {
    envVars := EnvVars{EnvVars: []*EnvVar{{Key: "IP", Value: "127.0.0.1"}, {Key: "TEST", Value: "True"}}}
    k8sEnvVars := envVars.ToK8sEnvVars()
    if len(k8sEnvVars) != 2 {
        t.Error("length error")
    }
    if k8sEnvVars[0].Name != "IP" {
        t.Error("key error")
    }
    if k8sEnvVars[0].Value != "127.0.0.1" {
        t.Error("value error")
    }
}
