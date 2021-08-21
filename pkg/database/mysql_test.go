package database

import (
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
    v1 "k8s.io/api/core/v1"
    "os"
    "testing"
)

var client *MySQLClient

var base = basemeta.BaseMeta{
    "test-app", "bj", "test",
}

func init() {
    fmt.Println(os.Getwd())
    client = new(MySQLClient)
    engine, err := xorm.NewEngine("mysql", "root:mysql@tcp(127.0.0.1:3306)/test?charset=utf8")
    err = engine.Sync2(new(container.EnvVar), new(container.ResourceRequirement), new(container.Configuration),
        new(pod.ContainerGroup), new(container.Base), new(container.VolumeMount),
        new(container.LivenessCheck), new(container.Port), new(container.ReadnessCheck),
        new(container.Command), new(pod.CFSVolume), new(pod.EmptyVolume), new(pod.BOSVolume), new(app.App),
        new(app.Region), new(app.Env), new(ingress.IngressRule), new(container.LivenessCheckHttp))
    fmt.Println(err)
    engine.ShowSQL(true)
    client.engine = engine
}

func TestMySQLClient_AddAttributesForContainer_Add(t *testing.T) {
    envVar := new(container.EnvVar)
    envVar.ID = 0
    envVar.Value = "test-value"
    envVar.Key = "test-key"
    envVar.ContainerID = 100
    _, err := client.CreateOrUpdateRecord(envVar, []interface{}{envVar}, 0)
    if err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
    if envVar.ID == 0 {
        t.Error("generate id error")
    }
}

func TestMySQLClient_AddAttributesForContainer_update(t *testing.T) {
    envVar := new(container.EnvVar)
    envVar.ID = 1
    envVar.Value = "test-value-update"
    envVar.Key = "test-key"
    envVar.ContainerID = 100
    _, err := client.CreateOrUpdateRecord(envVar, []interface{}{envVar}, 0)
    if err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
    if envVar.ID != 1 {
        t.Error("generate id error")
    }
}

func TestMySQLClient_AddContainerGroup(t *testing.T) {
    affected, err := client.CreateContainerGroup(base)
    if err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
    if affected == 0 {
        t.Error("affected row number error")
    }
}

func TestMySQLClient_AddInitContainers(t *testing.T) {
    containerbase := &container.Base{
        Name:    "nginx",
        Image:   "nginx/1.3.2",
        Comment: "nginx-for-test",
    }
    aff, err := client.CreateOrUpdateRecord(containerbase, []interface{}{containerbase}, 0)
    if err != nil || aff == 0 {
        t.Error("add container base error")
    }
    err = client.UpdateInitOrSidecarContainer(base, "init", []int{containerbase.ID})
    if err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
    if containerbase.ID == 0 {
        t.Error("add container error")
    }
}

func TestMySQLClient_AddSidecarContainers(t *testing.T) {
    containerbase := &container.Base{
        Name:    "php",
        Image:   "php/1.3.2",
        Comment: "php-for-main",
    }
    aff, err := client.CreateOrUpdateRecord(containerbase, []interface{}{containerbase}, 0)
    if err != nil || aff == 0 {
        t.Error("add container base error")
    }
    err = client.AddMainContainer(base, containerbase.ID)
    if err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
    if containerbase.ID == 0 {
        t.Error("add container error")
    }
}

func TestMySQLClient_TestGet(t *testing.T) {
    containerGroup := new(pod.ContainerGroup)
    query := client.BaseQuery("test-app", "bj", "testa")
    has, err := client.engine.Where(query[0], query[1:]...).Get(containerGroup)
    if err != nil {
        fmt.Println(err)
        fmt.Println(containerGroup)
    }
    fmt.Println(containerGroup)

    fmt.Println("has", has)
    fmt.Println(containerGroup == nil)
}

func TestMySQLClient_AddAttributesForContainer(t *testing.T) {
    id := 453
    command := &container.Command{
        ID:          0,
        ContainerID: id,
        Command:     `ls`,
        Args:        "-l",
    }
    port := &container.Port{
        ContainerID: id,
        PortNumber:  8080,
    }
    vms := &container.VolumeMount{
        ContainerID: id,
        VolumeName:  "bos-volume",
        MountPath:   "test",
        SubPath:     "",
        ReadOnly:    false,
    }
    confs := &container.Configuration{ContainerID: id,
        Path: "/etc/init.d/nginx", Content: "abcd"}
    rc := &container.ReadnessCheck{
        ID:          0,
        ContainerID: id,
    }
    rr := &container.ResourceRequirement{
        ContainerID:   id,
        CPULimit:      1000,
        MemoryLimit:   1000,
        CPURequest:    200,
        MemoryRequest: 200,
    }
    cb := &container.Base{
        ID:      id,
        Name:    "nginx",
        Image:   "nginx/1.3.2",
        Comment: "nginx for php",
    }
    objs := []interface{}{command, cb, rr, rc, port, confs, vms}
    aff, err := client.CreateOrUpdateRecord(objs, objs, 0)
    if aff != 7 || err != nil {
        fmt.Println(err)
        t.Error("add error")
    }
}

func TestMySQLClient_GetContainer(t *testing.T) {
    id := 453
    containerb, err := client.GetContainer(id)
    if err != nil {
        t.Error(err)
    }
    meta := basemeta.BaseMeta{
        Region: "bj",
        Env:    "test",
        App:    "test-app",
    }
    pod, cf, vm := containerb.ToK8sContainer(meta)
    js, _ := json.Marshal(pod)
    cfjs, _ := json.Marshal(cf)
    vmjs, _ := json.Marshal(vm)
    fmt.Println(string(js))
    fmt.Println(string(cfjs))
    fmt.Println(string(vmjs))
}

func TestMySQLClient_AddCFSVolume(t *testing.T) {
    cfs := &pod.CFSVolume{
        Meta: basemeta.BaseMeta{
            App:    "test-app",
            Region: "bj",
            Env:    "test",
        },
        Name:     "cfs-1",
        Server:   "cfs.baidu.com",
        ReadOnly: false,
        Path:     "/",
    }
    aff, err := client.CreateOrUpdateRecord(cfs, []interface{}{cfs}, 0)
    if aff != 1 || err != nil {
        t.Error(err)
    }
}

func TestMySQLClient_AddBosVolume(t *testing.T) {
    bos := &pod.BOSVolume{
        ID: 0,
        Meta: basemeta.BaseMeta{
            App:    "test-app",
            Region: "bj",
            Env:    "test",
        },
        Name:      "bos-volume",
        ClaimName: "bos-volume",
    }
    aff, err := client.CreateOrUpdateRecord(bos, []interface{}{bos}, 0)
    if aff != 1 || err != nil {
        t.Error(err)
    }
}

func TestMySQLClient_InsertContainerGroup(t *testing.T) {
    id := 453
    cg := pod.ContainerGroup{
        ID: 0,
        Meta: basemeta.BaseMeta{
            App:    "test-app",
            Region: "bj",
            Env:    "test",
        },
        MainContainerID:     id,
        InitContainersIDs:   make([]int, 0),
        SidecarContainerIDs: make([]int, 0),
        MainContainer:       nil,
        InitContainers:      nil,
        SidecarContainers:   nil,
    }
    aff, err := client.CreateOrUpdateRecord(cg, []interface{}{cg}, 0)
    if aff != 1 || err != nil {
        t.Error(err)
    }
    //client.GetVolumes("test-app")
}

func TestMySQLClient_GetContainerGroupAndVolumes(t *testing.T) {
    cg, err := client.GetContainerGroup(base)
    if err != nil {
        t.Error(err)
    }
    volumes, err := client.GetVolumes(base)
    v1Volumes := volumes.ToK8sVolumes()
    lbs := make(map[string]string)
    depName := "test-deployment"
    podspec, conf := cg.ToK8sPodTemplateSpecAndConfigMaps(depName, lbs, v1Volumes)
    p, _ := json.Marshal(podspec)
    cf, _ := json.Marshal(conf)
    fmt.Println(string(p))
    fmt.Println(string(cf))
    //client.GetVolumes("test-app")
}

func TestMySQLClient_GetContainerGroupAndEmptyVolumes(t *testing.T) {
    cg, err := client.GetContainerGroup(base)
    if err != nil {
        t.Error(err)
    }
    //volumes, err := client.GetVolumes("test-app", "bj", "test")
    //v1Volumes := volumes.ToK8sVolumes()
    depName := "test-deployment"

    lbs := make(map[string]string)
    podspec, conf := cg.ToK8sPodTemplateSpecAndConfigMaps(depName, lbs, []v1.Volume{})
    p, _ := json.Marshal(podspec)
    cf, _ := json.Marshal(conf)
    fmt.Println(string(p))
    fmt.Println(string(cf))
    //client.GetVolumes("test-app")
}

func TestCloneBase(t *testing.T) {
    base, err := client.cloneBase(4)
    if err != nil {
        fmt.Println("error")
    }
    fmt.Println(base)
}

func TestCloneConfig(t *testing.T) {
    err := client.cloneConfigurations(4, 20, "clone-test", nil)
    if err != nil {
        fmt.Println(err)
    }
}

func TestCloneCommand(t *testing.T) {
    err := client.clonelchttp(4, 20)
    if err != nil {
        fmt.Println(err)
    }
}

func TestCloneCommandNotExist(t *testing.T) {
    err := client.clonelchttp(5, 20)
    if err != nil {
        fmt.Println(err)
    }
}

func TestCloneCfs(t *testing.T) {
    err := client.cloneCfs("penghao-test", "clone-app",nil)
    if err != nil {
        fmt.Println(err)
    }
}

func TestCloneCg(t *testing.T) {
    meta := basemeta.BaseMeta{
        App:    "penghao-test",
        Region: "bj",
        Env:    "test",
    }
    cg, err := client.GetContainerGroupBase(meta)
    if err != nil {
        t.Error("get cg error")
    }
    ids, err := client.cloneCg("clone-app", cg, nil)
    fmt.Println(ids)
}

func TestMySQLClient_CloneApp(t *testing.T) {
    info := &CloneInfo{
        App1: "penghao-test",
        App2: "penghao-test",
        Zone: []*CloneZone{
            &CloneZone{
                Old: &Zone{"bj", "test"},
                New: &Zone{"bj", "cce"},
            },
        },
    }
    errs := client.CloneApp(info)
    for _, err := range errs {
        fmt.Println(err)
    }
}


func TestMySQLClient_CloneDifferentApp(t *testing.T) {
    info := &app.CloneInfo{
        App1: "penghao-test",
        App2: "penghao-test",
        Zones: []*app.CloneZone{
            &app.CloneZone{
                Old: &app.Zone{"bj", "cce"},
                New: &app.Zone{"bj", "test"},
            },
        },
    }
    errs := client.CloneApp(info)
    for _, err := range errs {
        fmt.Println(err)
    }
}

func TestMySQLClient_GetNotInitializedZones(t *testing.T) {
    a, err := client.GetNotInitializedZones("penghao-test")
    if err != nil {
        t.Error(err)
    }
    for _, v := range a {
        fmt.Println(v.App, v.Region, v.Env)
    }
}

func TestMySQLClient_GetCname(t *testing.T) {
    a, err := client.getCname(basemeta.BaseMeta{"penghao-test", "bj", "test"})
    if err != nil {
        t.Error(err)
    }
    fmt.Println(a)
}