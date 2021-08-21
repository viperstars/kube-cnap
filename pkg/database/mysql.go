package database

import (
    "errors"
    "fmt"
    "github.com/go-xorm/xorm"
    myApp "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
    "go.uber.org/zap"
    "reflect"
    "regexp"
    "strings"
)

type MySQLClient struct {
    engine *xorm.Engine
    logger *zap.Logger
}

func (m *MySQLClient) DeleteContainerAttribute(id int, table interface{}) (int64, error) {
    affected, err := m.engine.Table(table).Where(`id = ?`, id).Update(map[string]interface{}{"deleted": true})
    return affected, err
}

func (m *MySQLClient) UpdateContainerAttribute(id int, table interface{}, record interface{}) (int64, error) {
    // data := m.GetUpdateMap(record)
    affected, err := m.engine.Table(table).Where(`id = ?`, id).AllCols().Update(record)
    return affected, err
}

func (m *MySQLClient) GetUpdateMap(record interface{}) map[string]interface{} {
    d := make(map[string]interface{}, 0)
    getType := reflect.TypeOf(record)
    if getType.Kind() == reflect.Ptr {
        getType = getType.Elem()
    }
    getValue := reflect.ValueOf(record)
    if getValue.Kind() == reflect.Ptr {
        getValue = getValue.Elem()
    }
    word := regexp.MustCompile(`'.+'`)
    for i := 0; i < getType.NumField(); i++ {
        var realName string
        field := getType.Field(i)
        tag := field.Tag.Get("xorm")
        x := word.FindAll([]byte(tag), -1)
        if len(x) == 1 {
            realName = strings.Trim(string(x[0]), `'`)
            fmt.Println(realName)
        } else if tag != "extends"  {
            realName = m.engine.ColumnMapper.Obj2Table(field.Name)
        } else {
            continue
        }
        value := getValue.Field(i).Interface()
        d[realName] = value
    }
    return d
}

func (m *MySQLClient) CreateContainerAttribute(record interface{}) (int64, error) {
    affected, err := m.engine.Insert(record)
    return affected, err
}

func (m *MySQLClient) RemoveContainer(meta basemeta.BaseMeta, kind string, id int) error {
    _, err := m.engine.Delete(&container.Base{
        ID: id,
    })
    return err
}

func (m *MySQLClient) RemoveVolume(meta basemeta.BaseMeta, kind string, id int) error {
    var table interface{}
    switch kind {
    case "cfs":
        table = pod.CFSVolume{}
    case "bos":
        table = pod.BOSVolume{}
    case "empty":
        table = pod.EmptyVolume{}
    }
    _, err := m.engine.Table(table).Where(`id = ?`, id).Update(map[string]interface{}{"deleted": true})
    return err
}

func (m *MySQLClient) baseQuery(app, region, env string) []interface{} {
    return []interface{}{`app = ? AND region = ? AND env = ?`, app, region, env}
}

func (m *MySQLClient) BaseQuery(app, region, env string) []interface{} {
    return []interface{}{`app = ? AND region = ? AND env = ?`, app, region, env}
}

func (m *MySQLClient) GetVolumes(base basemeta.BaseMeta) (*pod.Volumes, error) {
    volumes := &pod.Volumes{}
    cfs, err := m.GetCFSVolumes(base)
    if err != nil {
        return volumes, err
    }
    volumes.CFSVolumes = cfs
    bos, err := m.GetBOSVolumes(base)
    if err != nil {
        return volumes, err
    }
    volumes.BOSVolumes = bos
    empty, err := m.GetEmptyVolumes(base)
    if err != nil {
        return volumes, err
    }
    volumes.EmptyVolumes = empty
    return volumes, nil
}

func (m *MySQLClient) GetCFSVolumes(base basemeta.BaseMeta) ([]*pod.CFSVolume, error) {
    cfsVolumes := make([]*pod.CFSVolume, 0)
    query := base.GetQuery()
    m.logger.Info(
        "getting volumes",
        zap.Any("meta", base),
    )
    err := m.engine.Where(query[0], query[1:]...).Where("deleted = ?", 0).Find(&cfsVolumes)
    if err != nil {
        m.logger.Error(
            "get cfs volumes error",
            zap.Any("meta", base),
            zap.Error(err),
        )
    }
    return cfsVolumes, err
}

func (m *MySQLClient) GetBOSVolumes(base basemeta.BaseMeta) ([]*pod.BOSVolume, error) {
    bosVolumes := make([]*pod.BOSVolume, 0)
    query := base.GetQuery()
    m.logger.Info(
        "getting volumes",
        zap.Any("meta", base),
    )
    err := m.engine.Where(query[0], query[1:]...).Where("deleted = ?", 0).Find(&bosVolumes)
    if err != nil {
        m.logger.Error(
            "get bos volumes error",
            zap.Any("meta", base),
            zap.Error(err),
        )
    }
    return bosVolumes, err
}

func (m *MySQLClient) GetEmptyVolumes(base basemeta.BaseMeta) ([]*pod.EmptyVolume, error) {
    emptyVolumes := make([]*pod.EmptyVolume, 0)
    query := base.GetQuery()
    m.logger.Info(
        "getting empty volumes",
        zap.Any("meta", base),
    )
    err := m.engine.Where(query[0], query[1:]...).Where("deleted = ?", 0).Find(&emptyVolumes)
    if err != nil {
        m.logger.Error(
            "get empty volumes error",
            zap.Any("meta", base),
            zap.Error(err),
        )
    }
    return emptyVolumes, err
}

func (m *MySQLClient) GetContainerGroup(base basemeta.BaseMeta) (*pod.ContainerGroup, error) {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "get container group error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return cg, err
    }
    initContainers, err := m.GetContainersFromSlice(base, cg.InitContainersIDs)
    if err != nil {
        m.logger.Error(
            "get init containers for container group error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return cg, err
    }
    sidecarContainers, err := m.GetContainersFromSlice(base, cg.SidecarContainerIDs)
    if err != nil {
        m.logger.Error(
            "get sidecar containers for container group error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return cg, err
    }
    mainContainer, err := m.GetContainer(cg.MainContainerID)
    if err != nil {
        m.logger.Error(
            "get main container for container group error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return cg, err
    }
    tolerations, _ := m.GetTolerations(base)
    if tolerations != nil {
        cg.Tolerations = tolerations
    }
    nodeSelector, err := m.GetNodeSelector(base)
    if nodeSelector != nil {
        cg.NodeSelector = nodeSelector
    }
    hostAliases, err := m.GetHostAlias(base)
    if hostAliases != nil {
        cg.HostAlias = hostAliases
    }
    cg.MainContainer = mainContainer
    cg.InitContainers = initContainers
    cg.SidecarContainers = sidecarContainers
    return cg, nil
}

func (m *MySQLClient) GetContainers(base basemeta.BaseMeta, ids []int) ([]*container.Container, error) {
    cs := make([]*container.Container, 0, len(ids))
    if len(ids) == 0 {
        m.logger.Info(
            "length of ids is 0",
            zap.Any("meta", base),
        )
        return cs, nil
    }
    for _, id := range ids {
        c, err := m.GetContainer(id)
        if err != nil {
            m.logger.Error(
                "get container error",
                zap.Any("meta", base),
                zap.Int("containerID", id),
                zap.Error(err),
            )
            return cs, err
        }
        cs = append(cs, c)
    }
    return cs, nil
}

func (m *MySQLClient) CreateContainerGroup(base basemeta.BaseMeta) (int64, error) {
    initIds := make([]int, 0)
    sidecarIds := make([]int, 0)
    containerGroup := &pod.ContainerGroup{
        Meta:                base,
        MainContainerID:     0,
        InitContainersIDs:   initIds,
        SidecarContainerIDs: sidecarIds,
        MainContainer:       nil,
        InitContainers:      nil,
        SidecarContainers:   nil,
    }
    return m.engine.Insert(containerGroup)
}

func (m *MySQLClient) GetContainersFromSlice(base basemeta.BaseMeta, ids []int) ([]*container.Container, error) {
    cs := make([]*container.Container, 0, len(ids))
    if len(ids) == 0 {
        return cs, nil
    }
    for _, id := range ids {
        c, err := m.GetContainer(id)
        if err != nil {
            return cs, err
        }
        cs = append(cs, c)
    }
    return cs, nil
}

func (m *MySQLClient) GetContainersFromMap(base basemeta.BaseMeta, ids map[int]struct{}) ([]*container.Container, error) {
    cs := make([]*container.Container, 0, len(ids))
    if len(ids) == 0 {
        return cs, nil
    }
    for id := range ids {
        c, err := m.GetContainer(id)
        if err != nil {
            return cs, err
        }
        cs = append(cs, c)
    }
    return cs, nil
}

func (m *MySQLClient) GetContainerCommand(containerID int) (*container.Command, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    command := new(container.Command)
    has, err := m.engine.Where(query[0], query[1:]...).Get(command)
    if err != nil {
        m.logger.Error(
            "get container command error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if has {
        return command, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerPorts(containerID int) ([]*container.Port, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    ports := make([]*container.Port, 0)
    err := m.engine.Where(query[0], query[1:]...).Find(&ports)
    if err != nil {
        m.logger.Error(
            "get container ports error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if len(ports) > 0 {
        return ports, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerConfigurations(containerID int) ([]*container.Configuration,
    error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    conf := make([]*container.Configuration, 0)
    err := m.engine.Where(query[0], query[1:]...).Find(&conf)
    if err != nil {
        m.logger.Error(
            "get container configurations error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if len(conf) > 0 {
        return conf, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerResourceRequirement(containerID int) (*container.ResourceRequirement, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    res := new(container.ResourceRequirement)
    has, err := m.engine.Where(query[0], query[1:]...).Get(res)
    fmt.Println("ssss", has, err)
    if err != nil {
        m.logger.Error(
            "get container resource requirement error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if has {
        return res, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerEnvVars(containerID int) ([]*container.EnvVar, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    vars := make([]*container.EnvVar, 0)
    err := m.engine.Where(query[0], query[1:]...).Find(&vars)
    if err != nil {
        m.logger.Error(
            "get container env vars error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if len(vars) > 0 {
        return vars, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerVolumeMounts(containerID int) ([]*container.VolumeMount, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    vms := make([]*container.VolumeMount, 0)
    err := m.engine.Where(query[0], query[1:]...).Find(&vms)
    if err != nil {
        m.logger.Error(
            "get container volume mounts error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return nil, err
    } else if len(vms) > 0 {
        return vms, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerReadnessCheck(containerID int) (*container.ReadnessCheck, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    rc := new(container.ReadnessCheck)
    has, err := m.engine.Where(query[0], query[1:]...).Get(rc)
    if err != nil {
        m.logger.Error(
            "get container readness check error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return rc, err
    } else if has {
        return rc, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerLivenessCheck(containerID int) (*container.LivenessCheck, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    lc := new(container.LivenessCheck)
    has, err := m.engine.Where(query[0], query[1:]...).Get(lc)
    if err != nil {
        m.logger.Error(
            "get container liveness check error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return lc, err
    } else if has {
        return lc, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainerLivenessCheckHttp(containerID int) (*container.LivenessCheckHttp, error) {
    query := []interface{}{"container_id = ? AND deleted = ?", containerID, 0}
    lc := new(container.LivenessCheckHttp)
    has, err := m.engine.Where(query[0], query[1:]...).Get(lc)
    if err != nil {
        m.logger.Error(
            "get container liveness check http error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return lc, err
    } else if has {
        return lc, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetContainer(containerID int) (*container.Container, error) {
    c := &container.Container{}
    query := []interface{}{"id = ?", containerID}
    base := new(container.Base)
    has, err := m.engine.Where(query[0], query[1:]...).Get(base)
    if err != nil {
        m.logger.Error(
            "get container base info error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return c, err
    } else if has {
        c.Name = base.Name
        c.Image = base.Image
    } else {
        return c, errors.New("no base info found for container")
    }
    command, err := m.GetContainerCommand(containerID)
    if command != nil {
        c.Command = command
    }
    ports, err := m.GetContainerPorts(containerID)
    if ports != nil {
        c.Ports = &container.Ports{Ports: ports}
    }
    confs, err := m.GetContainerConfigurations(containerID)
    if confs != nil {
        c.Configurations = &container.Configurations{Configurations: confs}
    }
    vms, err := m.GetContainerVolumeMounts(containerID)
    if vms != nil {
        c.VolumeMounts = &container.VolumeMounts{VolumeMounts: vms}
    }
    envVars, err := m.GetContainerEnvVars(containerID)
    if envVars != nil {
        c.EnvVar = &container.EnvVars{EnvVars: envVars}
    }
    lc, err := m.GetContainerLivenessCheck(containerID)
    if lc != nil {
        c.LivenessCheck = lc
    }
    /*lch, err := m.GetContainerLivenessCheckHttp(containerID)
    if lch != nil {
        c.LivenessCheckHttp = lch
    }*/
    rc, err := m.GetContainerReadnessCheck(containerID)
    if rc != nil && rc.Path != ""{
        c.ReadnessCheck = rc
    }
    rr, err := m.GetContainerResourceRequirement(containerID)
    if rr != nil {
        c.ResourceRequirement = rr
    }
    return c, nil
}

func (m *MySQLClient) CreateOrUpdateRecord(record interface{}, table interface{}, id int) (int64, error) {
    var err error
    var affected int64
    if id > 0 {
        // data := m.GetUpdateMap(record)
        affected, err = m.engine.ID(id).Table(table).AllCols().Update(record)
    } else {
        has, err := m.checkExist(record)
        if err != nil {
            return 0, err
        } else if !has {
            affected, err = m.engine.Insert(record)
        } else {
            return 0, errors.New("adding an exist record")
        }
    }
    if err != nil {
        m.logger.Error(
            "create or update error",
            zap.Error(err),
        )
    }
    return affected, err
}

func (m *MySQLClient) CreateOrUpdateRecordNoCheck(record interface{}, table interface{}, id int) (int64, error) {
    var err error
    var affected int64
    if id > 0 {
        // data := m.GetUpdateMap(record)
        affected, err = m.engine.ID(id).Table(table).AllCols().Update(record)
    } else {
        affected, err = m.engine.Insert(record)
    }
    if err != nil {
        m.logger.Error(
            "create or update error",
            zap.Error(err),
        )
    }
    return affected, err
}

func (m *MySQLClient) UpdateApp(app *myApp.App) error {
    _, err := m.engine.ID(app.ID).AllCols().Update(app)
    return err
}

func (m *MySQLClient) checkExist(record interface{}) (bool, error) {
    return m.engine.Where("deleted = ?", 0).Get(record)
}

func (m *MySQLClient) UpdateInitOrSidecarContainer(base basemeta.BaseMeta, kind string, ids []int) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "update init or sidecar containers error",
            zap.Error(err),
        )
        return err
    }
    switch kind {
    case "init":
        cg.InitContainersIDs = ids
    case "sidecar":
        cg.SidecarContainerIDs = ids
    }
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "update init or sidecar container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) AddMainContainer(base basemeta.BaseMeta, id int) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "update main container error",
            zap.Error(err),
        )
        return err
    }
    cg.MainContainerID = id
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "update container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) UpdateMainContainer(base basemeta.BaseMeta, id int) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "update main container error",
            zap.Error(err),
        )
        return err
    }
    cg.MainContainerID = id
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "update container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) AddInitOrSidecarContainer(base basemeta.BaseMeta, kind string, id int) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "add init or sidecar containers error",
            zap.Error(err),
        )
        return err
    }
    switch kind {
    case "init":
        cg.InitContainersIDs = append(cg.InitContainersIDs, id)
    case "sidecar":
        cg.SidecarContainerIDs = append(cg.SidecarContainerIDs, id)
    }
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "add init or sidecar container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) GetContainerGroupBase(base basemeta.BaseMeta) (*pod.ContainerGroup, error) {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil {
        m.logger.Error(
            "get container group base info error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return cg, err
    } else if !has {
        return cg, errors.New("no base info found for container")
    }
    if cg.MainContainerID != 0 {
        mainContainer, err := m.GetContainerBase(cg.MainContainerID)
        if err != nil {
            return cg, err
        }
        cg.MainBase = mainContainer
    }
    initContainerBases, err := m.GetContainerBaseFromSlice(cg.InitContainersIDs)
    if err != nil {
        return cg, err
    }
    cg.InitBases = initContainerBases
    sidecarContainerBases, err := m.GetContainerBaseFromSlice(cg.SidecarContainerIDs)
    if err != nil {
        return cg, err
    }
    cg.SidecarBases = sidecarContainerBases
    return cg, nil
}

func (m *MySQLClient) GetContainerBase(containerID int) (*container.Base, error) {
    query := []interface{}{"id = ?", containerID}
    containerBase := new(container.Base)
    has, err := m.engine.Where(query[0], query[1:]...).Get(containerBase)
    if err != nil {
        m.logger.Error(
            "get container base info error",
            zap.Int("containerID", containerID),
            zap.Error(err),
        )
        return containerBase, err
    } else if has {
        return containerBase, nil
    } else {
        return containerBase, errors.New("no base info found for container")
    }
}

func (m *MySQLClient) GetContainerBaseFromMap(ids map[int]struct{}) ([]*container.Base, error) {
    bases := make([]*container.Base, 0)
    for id := range ids {
        base, err := m.GetContainerBase(id)
        if err != nil {
            return bases, err
        }
        bases = append(bases, base)
    }
    return bases, nil
}

func (m *MySQLClient) GetContainerBaseFromSlice(ids []int) ([]*container.Base, error) {
    bases := make([]*container.Base, 0)
    for _, id := range ids {
        base, err := m.GetContainerBase(id)
        if err != nil {
            return bases, err
        }
        bases = append(bases, base)
    }
    return bases, nil
}

func (m *MySQLClient) AddContainer(base basemeta.BaseMeta, kind string, container2 *container.Base) error {
    var err error
    _, err = m.CreateOrUpdateRecordNoCheck(container2, new(*container.Base), 0)
    if err != nil {
        return err
    }
    // todo add container group
    _, err = m.CreateContainerGroup(base)
    if kind == "main" {
        err = m.AddMainContainer(base, container2.ID)
    } else {
        err = m.AddInitOrSidecarContainer(base, kind, container2.ID)
    }
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) RemoveMainContainer(base basemeta.BaseMeta) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil {
        m.logger.Error(
            "get container group base info error",
            zap.Any("meta", base),
            zap.Error(err),
        )
        return err
    } else if !has {
        return errors.New("no base info found for container")
    }
    cg.MainContainerID = 0
    aff, err := m.engine.Table(cg).Where(`id = ?`, cg.ID).Update(map[string]interface{}{"main_container_id": 0})
    if err != nil {
        return err
    } else if aff == 0 {
        return nil
    }
    return nil
}

func (m *MySQLClient) RemoveInitOrSidecarContainer(base basemeta.BaseMeta, kind string) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "add init or sidecar containers error",
            zap.Error(err),
        )
        return err
    }
    empty := make([]int, 0)
    switch kind {
    case "init":
        cg.InitContainersIDs = empty
    case "sidecar":
        cg.SidecarContainerIDs = empty
    }
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "add init or sidecar container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) RemoveInitOrSidecarContainerByID(base basemeta.BaseMeta, kind string, id int) error {
    cg := new(pod.ContainerGroup)
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(cg)
    if err != nil || !has {
        m.logger.Error(
            "add init or sidecar containers error",
            zap.Error(err),
        )
        return err
    }
    switch kind {
    case "init":
        cg.InitContainersIDs = m.RemoveElementFromSlice(cg.InitContainersIDs, id)
    case "sidecar":
        fmt.Print(cg.SidecarContainerIDs)
        cg.SidecarContainerIDs = m.RemoveElementFromSlice(cg.SidecarContainerIDs, id)
        fmt.Print(cg.SidecarContainerIDs)
    default:
        return errors.New("wrong container type")
    }
    _, err = m.engine.Id(cg.ID).AllCols().Update(cg)
    if err != nil {
        m.logger.Error(
            "add init or sidecar container error",
            zap.Error(err),
        )
    }
    return err
}

func (m *MySQLClient) GetHostAlias(base basemeta.BaseMeta) (*pod.HostAlias, error) {
    hls := &pod.HostAlias{}
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Get(hls)
    if err != nil {
        return hls, err
    } else if has {
        return hls, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetNodeSelector(base basemeta.BaseMeta) (*pod.NodeSelector, error) {
    ns := &pod.NodeSelector{}
    query := base.GetQuery()
    has, err := m.engine.Where(query[0], query[1:]...).Where("deleted = 0").Get(ns)
    if err != nil {
        return ns, err
    } else if has {
        return ns, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetTolerations(base basemeta.BaseMeta) (*pod.Tolerations, error) {
    tls := make([]*pod.Toleration, 0)
    tlss := &pod.Tolerations{}
    query := base.GetQuery()
    err := m.engine.Where(query[0], query[1:]...).Find(tls)
    if err != nil {
        return tlss, err
    } else if len(tls) > 0 {
        tlss.Tolerations = tls
        return tlss, nil
    } else {
        return nil, nil
    }
}

func (m *MySQLClient) GetIngressRules(base basemeta.BaseMeta, host string) ([]*ingress.IngressRule, error) {
    var err error
    rules := make([]*ingress.IngressRule, 0)
    query := base.GetQuery()
    if host != "" {
        err = m.engine.Where(query[0], query[1:]...).Where("deleted = ?", 0).Where("host = ?", host).Find(rules)
    } else {
        err = m.engine.Where(query[0], query[1:]...).Where("deleted = ?", 0).Find(&rules)
    }
    if err != nil {
        return rules, err
    } else {
        return rules, nil
    }
}

func (m *MySQLClient) DeleteIngressRule(id int, table interface{}) (int64, error) {
    affected, err := m.engine.Table(table).Where(`id = ?`, id).Update(map[string]interface{}{"deleted": true})
    return affected, err
}

func (m *MySQLClient) GetAllRegionsAndEnvs(withConfig bool) ([]*myApp.Region, error) {
    regions := make([]*myApp.Region, 0)
    err := m.engine.Where("deleted = ?", 0).Find(&regions)
    if err != nil {
        return regions, err
    }
    for index, r := range regions {
        envs, err := m.GetEnvs(r.ID, withConfig)
        if err != nil {
            return regions, err
        }
        regions[index].Envs = envs
    }

    return regions, err
}

func (m *MySQLClient) RemoveElementFromSlice(slice []int, element int) []int {
    for index, value := range slice {
        if value == element {
            slice = append(slice[:index], slice[index+1:]...)
            return slice
        }
    }
    return slice
}

func (m *MySQLClient) GetEnvs(id int, withConfig bool) ([]*myApp.Env, error) {
    envs := make([]*myApp.Env, 0)
    err := m.engine.Where("region_id = ?", id).Where("deleted = ?", 0).Desc("id").Find(&envs)
    if !withConfig {
        for i := 0; i < len(envs); i++ {
            envs[i].Config = ""
            envs[i].Namespace = ""
        }
    }
    return envs, err
}

func (m *MySQLClient) GetAllApps() ([]*myApp.App, error) {
    apps := make([]*myApp.App, 0)
    err := m.engine.Where("deleted = ?", 0).Find(&apps)
    return apps, err
}

func (m *MySQLClient) GetApp(id int) (*myApp.App, error) {
    _app := new(myApp.App)
    _, err := m.engine.ID(id).Get(_app)
    return _app, err
}

func (m *MySQLClient) GetAppByName(name string) (*myApp.App, error) {
    _app := new(myApp.App)
    has, err := m.engine.Where(`name = ?`, name).Get(_app)
    if err != nil {
        return nil, err
    } else if !has {
        return nil, errors.New("app not found")
    } else {
        return _app, err
    }
}

func (m *MySQLClient) GetContainerAttributes(containerID int, table string) (interface{}, error) {
    var data interface{}
    var err error
    switch table {
    case "readnessCheck":
        data = new(container.ReadnessCheck)
        _, err = m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Get(data)
    case "resourceRequirement":
        data = new(container.ResourceRequirement)
        _, err = m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Get(data)
    case "livenessCheck":
        data = new(container.LivenessCheck)
        _, err = m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Get(data)
    case "command":
        data = new(container.Command)
        _, err = m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Get(data)
    case "envVars":
        return m.getEnvVars(containerID)
    case "ports":
        return m.getPorts(containerID)
    case "configurations":
        return m.getConfigurations(containerID)
    case "volumeMounts":
        return m.getVolumeMounts(containerID)
    }
    return data, err
}

func (m *MySQLClient) GetContainerAttribute(recordID int, table string) (interface{}, error) {
    var data interface{}
    var err error
    switch table {
    case "readnessCheck":
        data = new(container.ReadnessCheck)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "resourceRequirement":
        data = new(container.ResourceRequirement)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "livenessCheck":
        data = new(container.LivenessCheck)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "command":
        data = new(container.Command)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "envVar":
        data = new(container.EnvVar)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "port":
        data = new(container.Port)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "configuration":
        data = new(container.Configuration)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    case "volumeMount":
        data = new(container.VolumeMount)
        _, err = m.engine.Where("id = ?", recordID).Get(data)
    }
    return data, err
}

func (m *MySQLClient) getEnvVars(containerID int) (interface{}, error) {
    data := make([]*container.EnvVar, 0)
    err := m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Find(&data)
    return data, err
}

func (m *MySQLClient) getConfigurations(containerID int) (interface{}, error) {
    data := make([]*container.Configuration, 0)
    err := m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Find(&data)
    return data, err
}

func (m *MySQLClient) getPorts(containerID int) (interface{}, error) {
    data := make([]*container.Port, 0)
    err := m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Find(&data)
    return data, err
}

func (m *MySQLClient) getVolumeMounts(containerID int) (interface{}, error) {
    data := make([]*container.VolumeMount, 0)
    err := m.engine.Where("container_id = ?", containerID).Where("deleted = ?", 0).Find(&data)
    return data, err
}

func (m *MySQLClient) GetIngressByID(recordID int) (*ingress.IngressRule, error) {
    ingressRule := new(ingress.IngressRule)
    _, err := m.engine.ID(recordID).Where("deleted = ?", 0).Get(ingressRule)
    return ingressRule, err
}

func (m *MySQLClient) GetVolumeByID(recordID int, kind string) (interface{}, error) {
    switch kind {
    case "cfs":
        data := new(pod.CFSVolume)
        _, err := m.engine.ID(recordID).Where("deleted = ?", 0).Get(data)
        return data, err
    case "bos":
        data := new(pod.BOSVolume)
        _, err := m.engine.ID(recordID).Where("deleted = ?", 0).Get(data)
        return data, err
    case "empty":
        data := new(pod.EmptyVolume)
        _, err := m.engine.ID(recordID).Where("deleted = ?", 0).Get(data)
        return data, err
    }
    return nil, errors.New("wrong volume type")
}

func(m * MySQLClient) GetBaseInfo(base basemeta.BaseMeta) (string, string, error) {
    region := new(myApp.Region)
    env := new(myApp.Env)
    has, err := m.engine.Where("name = ?", base.Region).Where("deleted != ?", 1).Get(region)
    if !has || err != nil {
        fmt.Println(err)
        return "", "", errors.New("base region not found")
    }
    has, err = m.engine.Where("name = ?", base.Env).Where("deleted != ?", 1).Get(env)
    if !has || err != nil {
        return "", "", errors.New("base env not found")
    }
    return region.CName, env.CName, nil
}

func(m *MySQLClient) CreateContainerBase(base basemeta.BaseMeta) (int64, error) {
    cg := &pod.ContainerGroup{
        Meta:                base,
        MainContainerID:     0,
        InitContainersIDs:   []int{},
        SidecarContainerIDs: []int{},
    }
    return m.engine.Insert(cg)
}

func(m *MySQLClient) getCgs(app string, region, env string) (*pod.ContainerGroup, error) {
    cg := new(pod.ContainerGroup)
    a1, err := m.GetAppByName(app)
    if err != nil {
        return cg, err
    }
    has, err := m.engine.Table(pod.ContainerGroup{}).Where(`app = ?`, a1.Name).Where("region = ?",
        region).Where("env = ?", env).Get(cg)
    if err != nil {
        return cg, err
    } else if !has {
        return cg, errors.New("cgs not found")
    } else {
        return cg, nil
    }
}

func(m *MySQLClient) cloneBase(id int) (int, error) {
    if id == 0 {
        return 0, nil
    }
    base, err := m.GetContainerBase(id)
    if err != nil {
        return 0, err
    }
    newBase := new(container.Base)
    newBase.Name = base.Name
    newBase.Comment = base.Comment
    newBase.Image = base.Image
    aff, err := m.engine.Insert(newBase)
    if err != nil {
        return 0, err
    } else if aff == 0  {
        return 0, errors.New("insert failed")
    } else {
        return newBase.ID, nil
    }
}

func(m *MySQLClient) cloneConfigurations(old, new int, app string, zone *myApp.CloneZone) error {
    configurations := make([]*container.Configuration, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Find(&configurations)
    if err != nil {
        return err
    }
    for index := range configurations {
        configurations[index].ID = 0
        configurations[index].ContainerID = new
        configurations[index].Base.App = app
        configurations[index].Base.Region = zone.New.Region
        configurations[index].Base.Env = zone.New.Env
    }
    aff, err := m.engine.Insert(configurations)
    if err != nil {
        return err
    } else if aff != int64(len(configurations)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneVms(old, new int) error {
    vms := make([]*container.VolumeMount, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Find(&vms)
    if err != nil {
        return err
    }
    for index := range vms {
        vms[index].ID = 0
        vms[index].ContainerID = new
    }
    aff, err := m.engine.Insert(vms)
    if err != nil {
        return err
    } else if aff != int64(len(vms)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneEnvvars(old, new int) error {
    vars := make([]*container.EnvVar, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Find(&vars)
    if err != nil {
        return err
    }
    for index := range vars {
        vars[index].ID = 0
        vars[index].ContainerID = new
    }
    aff, err := m.engine.Insert(vars)
    if err != nil {
        return err
    } else if aff != int64(len(vars)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) clonePorts(old, new int) error {
    ports := make([]*container.Port, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Find(&ports)
    if err != nil {
        return err
    }
    for index := range ports {
        ports[index].ID = 0
        ports[index].ContainerID = new
    }
    aff, err := m.engine.Insert(ports)
    if err != nil {
        return err
    } else if aff != int64(len(ports)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneCfs(app1, app2 string, cloneZone *myApp.CloneZone) error {
    cfs := make([]*pod.CFSVolume, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`app = ?`, app1).
        Where("region = ?", cloneZone.Old.Region).Where("env = ?", cloneZone.Old.Env).Find(&cfs)
    if err != nil {
        return err
    }
    for index := range cfs {
        cfs[index].ID = 0
        cfs[index].Meta.App = app2
        cfs[index].Meta.Region = cloneZone.New.Region
        cfs[index].Meta.Env = cloneZone.New.Env
    }
    aff, err := m.engine.Insert(cfs)
    if err != nil {
        return err
    } else if aff != int64(len(cfs)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneBos(app1, app2 string, cloneZone *myApp.CloneZone) error {
    bos := make([]*pod.BOSVolume, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`app = ?`, app1).
        Where("region = ?", cloneZone.Old.Region).Where("env = ?", cloneZone.Old.Env).Find(&bos)
    if err != nil {
        return err
    }
    for index := range bos {
        bos[index].ID = 0
        bos[index].Meta.App = app2
        bos[index].Meta.Region = cloneZone.New.Region
        bos[index].Meta.Env = cloneZone.New.Env
    }
    aff, err := m.engine.Insert(bos)
    if err != nil {
        return err
    } else if aff != int64(len(bos)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneEmpty(app1, app2 string, cloneZone *myApp.CloneZone) error {
    em := make([]*pod.EmptyVolume, 0)
    err := m.engine.Where(`deleted = ?`, 0).Where(`app = ?`, app1).
        Where("region = ?", cloneZone.Old.Region).Where("env = ?", cloneZone.Old.Env).Find(&em)
    if err != nil {
        return err
    }
    for index := range em {
        em[index].ID = 0
        em[index].Meta.App = app2
        em[index].Meta.Region = cloneZone.New.Region
        em[index].Meta.Env = cloneZone.New.Env
    }
    aff, err := m.engine.Insert(em)
    if err != nil {
        return err
    } else if aff != int64(len(em)) {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneCommand(old, new int) error {
    cmd := &container.Command{}
    find, err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Get(cmd)
    if err != nil {
        return err
    } else if find == false {
        return nil
    }
    cmd.ID = 0
    cmd.ContainerID = new
    aff, err := m.engine.Insert(cmd)
    if err != nil {
        return err
    } else if aff != 1 {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) clonelc(old, new int) error {
    livenessCheck := &container.LivenessCheck{}
    find, err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Get(livenessCheck)
    if err != nil {
        return err
    } else if find == false {
        return nil
    }
    livenessCheck.ID = 0
    livenessCheck.ContainerID = new
    aff, err := m.engine.Insert(livenessCheck)
    if err != nil {
        return err
    } else if aff != 1 {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) clonelchttp(old, new int) error {
    livenessCheckHttp := &container.LivenessCheckHttp{}
    find, err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Get(livenessCheckHttp)
    if err != nil {
        return err
    } else if find == false {
        return nil
    }
    livenessCheckHttp.ID = 0
    livenessCheckHttp.ContainerID = new
    aff, err := m.engine.Insert(livenessCheckHttp)
    if err != nil {
        return err
    } else if aff != 1 {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) clonerc(old, new int) error {
    readnessCheck := &container.ReadnessCheck{}
    find, err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Get(readnessCheck)
    if err != nil {
        return err
    } else if find == false {
        return nil
    }
    readnessCheck.ID = 0
    readnessCheck.ContainerID = new
    aff, err := m.engine.Insert(readnessCheck)
    if err != nil {
        return err
    } else if aff != 1 {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) clonerr(old, new int) error {
    resourceRequirement := &container.ResourceRequirement{}
    find, err := m.engine.Where(`deleted = ?`, 0).Where(`container_id = ?`, old).Get(resourceRequirement)
    if err != nil {
        return err
    } else if find == false {
        return nil
    }
    resourceRequirement.ID = 0
    resourceRequirement.ContainerID = new
    aff, err := m.engine.Insert(resourceRequirement)
    if err != nil {
        return err
    } else if aff != 1 {
        return errors.New("insert error")
    } else {
        return nil
    }
}

func(m *MySQLClient) cloneFuncs1(old, new int) []error {
    errs := make([]error, 0)
    funcs := []func(old, new int) error{
        m.cloneCommand,
        m.clonelc,
        m.clonelchttp,
        m.clonerc,
        m.clonerr,
        m.cloneEnvvars,
        m.clonePorts,
        m.cloneVms,
    }
    for _, f := range funcs {
        err := f(old, new)
        if err != nil {
            errs = append(errs, err)
        }
    }
    return errs
}

func(m *MySQLClient) cloneFuncs2(app1, app2 string, cloneZone *myApp.CloneZone) []error {
    errs := make([]error, 0)
    funcs := []func(app1, app2 string, cloneZone *myApp.CloneZone) error{
        m.cloneBos,
        m.cloneCfs,
        m.cloneEmpty,
    }
    for _, f := range funcs {
        err := f(app1, app2, cloneZone)
        if err != nil {
            errs = append(errs, err)
        }
    }
    return errs
}


func(m *MySQLClient) cloneCgs(info *myApp.CloneInfo) (*myApp.CloneInfo, error) {
    ids := make([]*myApp.CloneID, 0)
    for _, i := range info.Zones {
        cg1s, err := m.getCgs(info.App1, i.Old.Region, i.Old.Env)
        if err != nil {
            return info, err
        }
        idss, err := m.cloneCg(info.App2, cg1s, i)
        if err != nil {
            return info, err
        }
        ids = append(ids, idss...)
    }
    info.IDs = ids
    return info, nil
}

func(m *MySQLClient) cloneCg(app string, cg *pod.ContainerGroup, zone *myApp.CloneZone) ([]*myApp.CloneID, error) {
    ids := make([]*myApp.CloneID, 0)
    newCg := new(pod.ContainerGroup)
    yes, err := m.engine.Where("region = ?", zone.New.Region).Where("env = ?", zone.New.Env).Where("app = ?",
        app).Get(newCg)
    if err != nil {
        return ids, err
    } else if yes != true {
        return ids, errors.New("cg not found")
    }
    id, err := m.cloneBase(cg.MainContainerID)
    if err != nil {
        return ids, err
    }
    if id != 0 {
        ids = append(ids, &myApp.CloneID{cg.MainContainerID, id, zone})
    }
    sideCarIds := make([]int, 0, len(cg.SidecarContainerIDs))
    for _, base := range cg.SidecarContainerIDs {
        scId, err := m.cloneBase(base)
        if err != nil {
            return ids, err
        }
        sideCarIds = append(sideCarIds, scId)
        ids = append(ids, &myApp.CloneID{base, scId, zone})
    }
    initIDs := make([]int, 0, len(cg.InitContainersIDs))
    for _, base := range cg.InitContainersIDs {
        initId, err := m.cloneBase(base)
        if err != nil {
            return ids, err
        }
        initIDs = append(initIDs, initId)
        ids = append(ids, &myApp.CloneID{base, initId, zone})
    }
    newCg.MainContainerID = id
    newCg.SidecarContainerIDs = sideCarIds
    newCg.InitContainersIDs = initIDs
    _, err = m.engine.ID(newCg.ID).Update(newCg)
    if err != nil {
        return ids, err
    } else {
        return ids, nil
    }
}

func(m *MySQLClient) cloneAll(info *myApp.CloneInfo) []error {
    errs := make([]error, 0)
    for _, zone := range info.Zones {
        err1 := m.cloneFuncs2(info.App1, info.App2, zone)
        if err1 != nil {
            errs = append(errs, err1...)
        }
    }
    for _, value := range info.IDs {
        err2 := m.cloneFuncs1(value.Old, value.New)
        errs = append(errs, err2...)
    }
    for _, value := range info.IDs {
        err2 := m.cloneConfigurations(value.Old, value.New, info.App2, value.Zone)
        if err2 != nil {
            errs = append(errs, err2)
        }
    }

    return errs
}

func(m *MySQLClient) CloneApp(info *myApp.CloneInfo) []error {
    errs := make([]error, 0)
    info, err := m.cloneCgs(info)
    if err != nil {
        errs = append(errs, err)
        return errs
    }
    err2s := m.cloneAll(info)
    errs = append(errs, err2s...)
    return errs
}

func(m *MySQLClient) GetHistories(base basemeta.BaseMeta) ([]*myApp.History, error) {
    his := make([]*myApp.History, 0)
    query := base.GetQuery()
    err := m.engine.Where(query[0], query[1:]...).Find(&his)
    if err != nil {
        m.logger.Error("get history error", zap.Error(err))
    }
    return his, err
}

func(m *MySQLClient) getCname(base basemeta.BaseMeta) (*basemeta.CNames, error) {
    cname := new(basemeta.CNames)
    region := new(myApp.Region)
    cname.App = base.App
    has, err := m.engine.Where("name = ?", base.Region).Get(region)
    if err != nil {
        return cname, err
    }
    if !has {
        return cname, errors.New("region not found")
    }
    cname.Region = base.Region
    cname.RCname = region.CName
    env := new(myApp.Env)
    has, err = m.engine.Where("name = ?", base.Env).Where("region_id = ?", region.ID).Get(env)
    if err != nil {
        return cname, err
    }
    if !has {
        return cname, errors.New("env not found")
    }
    cname.Env = base.Env
    cname.ECname = env.CName
    return cname, nil
}

func (m *MySQLClient) getNotInitializedZones(app string) ([]*basemeta.BaseMeta, error) {
    exists := make([]*basemeta.BaseMeta, 0)
    cgs := make([]*pod.ContainerGroup, 0)
    err := m.engine.Where("app = ?", app).Find(&cgs)
    if err != nil {
        return exists, err
    }
    all, err := m.GetAllRegionsAndEnvs(false)
    if err != nil {
        return exists, err
    }
    res := make([]*basemeta.BaseMeta, 0)
    for _, r := range all {
        for _, e := range r.Envs {
            base := &basemeta.BaseMeta{
                App:    app,
                Region: r.Name,
                Env:    e.Name,
            }
            res = append(res, base)
        }
    }
    notExists := make([]*basemeta.BaseMeta, 0)
    for _, r := range res {
        in := false
        for _, e := range cgs {
            if r.Region == e.Meta.Region && r.Env == e.Meta.Env {
                in = true
            }
        }
        if !in {
            notExists = append(notExists, r)
        }
    }
    return notExists, nil
}

func (m *MySQLClient) GetAvailableZone(appName string) ([]*myApp.Region, error) {
    regions := make([]*myApp.Region, 0)
    cgs := make([]*pod.ContainerGroup, 0)
    err := m.engine.Where("app = ?", appName).OrderBy("id DESC").Find(&cgs)
    if err != nil {
        return regions, err
    }
    qs := make(map[string][]string)
    for _, cg := range cgs {
        if _, ok := qs[cg.Meta.Region]; ok {
            qs[cg.Meta.Region] = append(qs[cg.Meta.Region], cg.Meta.Env)
        } else {
            qs[cg.Meta.Region] = []string{cg.Meta.Env}
        }
    }
    for k, v := range qs {
        region, err := m.getZone(k, v)
        if err != nil {
            return regions, err
        }
        regions = append(regions, region)
    }
    return regions, nil
}

func (m *MySQLClient) getZone(region string, envs []string) (*myApp.Region, error) {
    r := new(myApp.Region)
    e := make([]*myApp.Env, 0)
    has, err := m.engine.Where("name = ? and deleted = 0", region).Get(r)
    if err != nil {
        return r, err
    }
    if !has {
        return r, errors.New("region not found")
    }
    err = m.engine.Cols().Where("deleted = 0").In("name", envs).OrderBy("id DESC").Find(&e)
    if err != nil {
        return r, err
    }
    for index := range e {
        fmt.Println(e[index].ID)
        e[index].Config = ""
    }
    r.Envs = e
    return r, nil
}

func (m *MySQLClient) GetNotInitializedZones(app string) ([]*basemeta.CNames, error) {
    cnames := make([]*basemeta.CNames, 0)
    zs, err := m.getNotInitializedZones(app)
    if err != nil {
        return cnames, err
    }
    for _, z := range zs {
        cname, err := m.getCname(*z)
        if err != nil {
            return cnames, err
        }
        cnames = append(cnames, cname)
    }
    return cnames, err
}

func (m *MySQLClient) DeleteHostAlias(id int) error {
    _, err := m.engine.ID(id).Table(pod.HostAlias{}).Update(map[string]interface{}{"deleted": true})
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) UpdateHostAlias(id int, ns *pod.HostAlias) error {
    _, err := m.engine.ID(id).Table(pod.HostAlias{}).Update(ns)
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) InsertHostAlias(base basemeta.BaseMeta, hs *pod.HostAlias) error {
    hs.Meta = base
    _, err := m.engine.Table(pod.HostAlias{}).Insert(hs)
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) DeleteNodeSelector(id int) error {
    ns := &pod.NodeSelector{}
    _, err := m.engine.ID(id).Table(ns).Update(map[string]interface{}{"deleted": true})
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) UpdateNodeSelector(id int, ns *pod.NodeSelector) error {
    _, err := m.engine.ID(id).Table(pod.NodeSelector{}).Update(ns)
    if err != nil {
        return err
    }
    return nil
}

func (m *MySQLClient) InsertNodeSelector(base basemeta.BaseMeta, ns *pod.NodeSelector) error {
    ns.Meta = base
    _, err := m.engine.Table(pod.NodeSelector{}).Insert(ns)
    if err != nil {
        return err
    }
    return nil
}