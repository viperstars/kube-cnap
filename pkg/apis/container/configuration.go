package container

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/naming"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "path"
    "path/filepath"
)

type Configuration struct {
    ID          int               `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int               `json:"containerID" xorm:"int 'container_id'"`
    Base        basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Path        string            `json:"path" xorm:"text"`
    Content     string            `json:"content" xorm:"text"`
    SubPath     bool              `json:"subPath" xorm:"bool"`
    Deleted     bool              `json:"deleted" xorm:"bool default 0"`
}

type Configurations struct {
    Configurations []*Configuration `json:"configurations"`
}

func (c *Configurations) GetSameDirnameConfigurations() map[string][]*Configuration {
    paths := make(map[string][]*Configuration)
    for _, v := range c.Configurations {
        dirname := filepath.Dir(v.Path)
        basename := filepath.Base(v.Path)
        gc := &Configuration{
            Path:    basename,
            Content: v.Content,
        }
        if _, ok := paths[dirname]; ok {
            paths[dirname] = append(paths[dirname], gc)
        } else {
            paths[dirname] = []*Configuration{gc}
        }
    }
    return paths
}

func (c *Configurations) ToK8sObjects(base basemeta.BaseMeta) ([]*v1.ConfigMap, []v1.VolumeMount, []v1.Volume) {
    configMaps := make([]*v1.ConfigMap, 0)
    volumeMounts := make([]v1.VolumeMount, 0)
    volumes := make([]v1.Volume, 0)
    paths := c.GetSameDirnameConfigurations()
    for p, v := range paths {
        labels := base.GetLabels()
        //annotations := labels
        data := make(map[string]string)
        name := naming.Namer.ConfigMapName(&base, p)
        configMap := &v1.ConfigMap{
            TypeMeta: metav1.TypeMeta{},
            ObjectMeta: metav1.ObjectMeta{
                Name:   name,
                Labels: labels,
                //  Annotations: todo annotations from labels
            },
            Data:       data,
            BinaryData: nil,
        }
        volumeName := naming.Namer.VolumeName(name, "volume")
        volume := v1.Volume{
            Name: volumeName,
            VolumeSource: v1.VolumeSource{
                ConfigMap: &v1.ConfigMapVolumeSource{
                    LocalObjectReference: v1.LocalObjectReference{
                        Name: name,
                    },
                    Items:       nil, // todo add loop to add items
                    DefaultMode: nil,
                    Optional:    nil,
                },
            },
        }
        for _, value := range v {
            configMap.Data[value.Path] = value.Content
            mountPath := path.Join(p, value.Path)
            volumeMount := v1.VolumeMount{
                Name:             volumeName,
                ReadOnly:         false,
                MountPath:        mountPath,
                SubPath:          value.Path, // todo add subpath
                MountPropagation: nil,
                SubPathExpr:      "",
            }
            volumeMounts = append(volumeMounts, volumeMount)
        }
        configMaps = append(configMaps, configMap)
        volumes = append(volumes, volume)
    }
    // d, _ := json.Marshal(configMaps)
    // fmt.Println("configMaps is: ", string(d))
    return configMaps, volumeMounts, volumes
}
