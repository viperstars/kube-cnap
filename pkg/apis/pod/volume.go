package pod

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    v1 "k8s.io/api/core/v1"
)

type Volumes struct {
    CFSVolumes   []*CFSVolume   `json:"cfsVolumes"`
    BOSVolumes   []*BOSVolume   `json:"bosVolumes"`
    EmptyVolumes []*EmptyVolume `json:"emptyVolumes"`
}

func (v *Volumes) ToK8sVolumes() []v1.Volume {
    volumes := make([]v1.Volume, 0)
    if len(v.CFSVolumes) > 0 {
        for _, v := range v.CFSVolumes {
            volumes = append(volumes, v.ToK8sVolume())
        }
    }
    if len(v.BOSVolumes) > 0 {
        for _, v := range v.BOSVolumes {
            volumes = append(volumes, v.ToK8sVolume())
        }
    }
    if len(v.EmptyVolumes) > 0 {
        for _, v := range v.EmptyVolumes {
            volumes = append(volumes, v.ToK8sVolume())
        }
    }
    return volumes
}

type EmptyVolume struct {
    ID   int               `json:"id" xorm:"pk autoincr 'id'"`
    Meta basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Name string            `json:"name" xorm:"varchar(512) 'name'"`
    // Size int               `json:"size" xorm:"varchar(512) 'size'"`
    Deleted  bool          `json:"deleted" xorm:"bool"`
}

func (ev *EmptyVolume) TableName() string {
    return "empty_volume"
}

func (ev *EmptyVolume) ToK8sVolume() v1.Volume {
    // sizeString := fmt.Sprintf("%dGi", ev.Size)
    // size := resource.MustParse(sizeString)
    // fmt.Println("size",  size.String())
    volume := v1.Volume{
        Name: ev.Name,
        VolumeSource: v1.VolumeSource{
            EmptyDir: &v1.EmptyDirVolumeSource{
                Medium:    "",
                // SizeLimit: &size,
            },
        },
    }
    return volume
}

type CFSVolume struct {
    ID       int               `json:"id" xorm:"pk autoincr 'id'"`
    Meta     basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Name     string            `json:"name" xorm:"varchar(512)"`
    // Size     int               `json:"size"`
    Server   string            `json:"server" xorm:"varchar(512)"`
    ReadOnly bool              `json:"readOnly" xorm:"bool"`
    Path     string            `json:"path" xorm:"varchar(512)"`
    Deleted  bool              `json:"deleted" xorm:"bool"`
}

func (c *CFSVolume) TableName() string {
    return "cfs_volume"
}

func (c *CFSVolume) ToK8sVolume() v1.Volume {
    volume := v1.Volume{
        Name: c.Name,
        VolumeSource: v1.VolumeSource{
            NFS: &v1.NFSVolumeSource{
                Server:   c.Server,
                Path:     c.Path,
                ReadOnly: c.ReadOnly,
            },
        },
    }
    return volume
}

type BOSVolume struct {
    ID        int               `json:"id" xorm:"pk autoincr 'id'"`
    Meta      basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Name      string            `json:"name" xorm:"varchar(512) 'name'"`
    ClaimName string            `json:"claimName" xorm:"varchar(512) 'claim_name'"`
    // Size      string            `json:"size" xorm:"varchar(512) 'size'"`
    Deleted   bool              `json:"deleted" xorm:"bool"`
}

func (b *BOSVolume) TableName() string {
    return "bos_volume"
}

func (b *BOSVolume) ToK8sVolume() v1.Volume {
    volume := v1.Volume{
        Name: b.Name,
        VolumeSource: v1.VolumeSource{
            PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
                ClaimName: b.ClaimName,
                ReadOnly:  false,
            },
        },
    }
    return volume
}