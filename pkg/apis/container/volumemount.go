package container

import v1 "k8s.io/api/core/v1"

type VolumeMount struct {
    ID          int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int    `json:"containerID" xorm:"int 'container_id'"`
    VolumeName  string `json:"volumeName" xorm:"varchar(512)"`
    MountPath   string `json:"mountPath" xorm:"varchar(512)"`
    SubPath     string `json:"subPath" xorm:"varchar(512)"`
    ReadOnly    bool   `json:"readOnly" xorm:"bool"`
    Deleted     bool   `json:"deleted" xorm:"bool default 0"`
}

type VolumeMounts struct {
    VolumeMounts []*VolumeMount `json:"volumeMounts"`
}

func (vm *VolumeMounts) ToK8sVolumeMounts() []v1.VolumeMount {
    volumeMounts := make([]v1.VolumeMount, 0)
    for _, volumeMount := range vm.VolumeMounts {
        v := v1.VolumeMount{
            Name:             volumeMount.VolumeName,
            ReadOnly:         volumeMount.ReadOnly,
            MountPath:        volumeMount.MountPath,
            SubPath:          volumeMount.SubPath,
            MountPropagation: nil,
            SubPathExpr:      "",
        }
        volumeMounts = append(volumeMounts, v)
    }
    return volumeMounts
}
