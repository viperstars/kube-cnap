package container

import (
    "fmt"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    v1 "k8s.io/api/core/v1"
    "regexp"
    "strings"
)

var separator = regexp.MustCompile(`\s+`)

type Base struct {
    ID      int    `json:"id" xorm:"pk autoincr 'id'"`
    Name    string `json:"name" xorm:"varchar(256)"`
    Image   string `json:"image" xorm:"varchar(512)"`
    Comment string `json:"comment" xorm:"varchar(256)"`
    Deleted  bool  `json:"deleted" xorm:"bool"`
}

type Container struct {
    Name                string
    Image               string
    Command             *Command
    Ports               *Ports
    Configurations      *Configurations
    ResourceRequirement *ResourceRequirement
    EnvVar              *EnvVars
    LivenessCheck       *LivenessCheck
    LivenessCheckHttp   *LivenessCheckHttp
    ReadnessCheck       *ReadnessCheck
    VolumeMounts        *VolumeMounts
}

func (c *Container) GetConfigMapsAndVolumeMountsAndVolumes(base basemeta.BaseMeta) ([]*v1.ConfigMap, []v1.VolumeMount,
    []v1.Volume) {
    return c.Configurations.ToK8sObjects(base)
}

func (c *Container) ToK8sContainer(base basemeta.BaseMeta) (*v1.Container, []*v1.ConfigMap, []v1.Volume) {
    volumeMounts := make([]v1.VolumeMount, 0)
    configMaps := make([]*v1.ConfigMap, 0)
    volumes := make([]v1.Volume, 0)
    if c.VolumeMounts != nil {
        volumeMountForVolumes := c.VolumeMounts.ToK8sVolumeMounts()
        volumeMounts = append(volumeMounts, volumeMountForVolumes...)
    }
    if c.Configurations != nil {
        configMapsFromConf, volumeMountsFromConf, volumesFromConf := c.GetConfigMapsAndVolumeMountsAndVolumes(base)
        configMaps = append(configMaps, configMapsFromConf...)
        volumeMounts = append(volumeMounts, volumeMountsFromConf...)
        volumes = append(volumes, volumesFromConf...)
    }
    container := &v1.Container{
        Name:                     c.Name,
        Image:                    c.Image,
        Command:                  nil,
        Args:                     nil,
        WorkingDir:               "",
        Ports:                    nil,
        EnvFrom:                  nil,
        Env:                      nil,
        VolumeMounts:             volumeMounts,
        VolumeDevices:            nil,
        LivenessProbe:            nil,
        ReadinessProbe:           nil,
        Lifecycle:                nil,
        TerminationMessagePath:   consts.DEFAULTMESSAGEPATH,
        TerminationMessagePolicy: consts.DEFAULTMESSAGEPOLICY,
        ImagePullPolicy:          consts.DEFAULTPOLICY,
        SecurityContext:          nil,
        Stdin:                    false,
        StdinOnce:                false,
        TTY:                      false,
    }
    if c.Ports != nil {
        container.Ports = c.Ports.ToK8sPort()
    }
    if c.ResourceRequirement != nil {
        container.Resources = c.ResourceRequirement.ToK8sResourceRequirement()
    }
    if c.Command != nil {
        if c.Command.Command != ""  {
            container.Command = separator.Split(c.Command.Command, -1)
        }
        if c.Command.Args != "" {
            container.Args = separator.Split(c.Command.Args, -1)
        }
    }
    if c.LivenessCheck != nil {
        container.LivenessProbe = c.LivenessCheck.ToK8sProbe()
    }
    /*if c.LivenessCheckHttp != nil {
        container.LivenessProbe = c.LivenessCheckHttp.ToK8sProbe()
    }*/
    if c.ReadnessCheck != nil {
        container.ReadinessProbe = c.ReadnessCheck.ToK8sProbe()
    }
    envs := c.AddEnvs(base)
    if c.EnvVar != nil {
        envs = append(envs, c.EnvVar.ToK8sEnvVars()...)
    }
    container.Env = envs
    return container, configMaps, volumes
}

func (c *Container) AddEnvs(base basemeta.BaseMeta) []v1.EnvVar {
    envVars := make([]v1.EnvVar, 0)
    envVars = append(envVars, v1.EnvVar{
        Name:      "APP",
        Value:     base.App,
        ValueFrom: nil,
    })
    envVars = append(envVars, v1.EnvVar{
        Name:      "REGION",
        Value:     base.Region,
    })
    envVars = append(envVars, v1.EnvVar{
        Name:      "ENV",
        Value:     base.Env,
    })
    envVars = append(envVars, v1.EnvVar{
        Name:      "HOST_IP",
        ValueFrom: &v1.EnvVarSource{
            FieldRef: &v1.ObjectFieldSelector{
                FieldPath:  "status.hostIP",
            },
        },
    })
    envVars = append(envVars, v1.EnvVar{
        Name:      "POD_NAME",
        ValueFrom: &v1.EnvVarSource{
            FieldRef: &v1.ObjectFieldSelector{
                FieldPath:  "metadata.name",
            },
        },
    })
    e := strings.Split(base.Env, "-")
    if len(e) <= 1 {
        envVars = append(envVars, v1.EnvVar{
            Name: "PAAS_REGION",
            Value: fmt.Sprintf("%s-%s", base.Region, base.Env),
        })
    } else {
        envVars = append(envVars, v1.EnvVar{
            Name: "PAAS_REGION",
            Value: fmt.Sprintf("%s-%s", base.Region, e[0]),
        })
    }
    return envVars
}
