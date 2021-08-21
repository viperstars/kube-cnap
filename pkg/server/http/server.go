package http

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/emicklei/go-restful"
    "github.com/viperstars/bce-golang/bce_client"
    "github.com/viperstars/kube-cnap/conf"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
    "github.com/viperstars/kube-cnap/pkg/apis/requests"
    "github.com/viperstars/kube-cnap/pkg/apis/resourcequota"
    "github.com/viperstars/kube-cnap/pkg/apis/service"
    "github.com/viperstars/kube-cnap/pkg/client"
    "github.com/viperstars/kube-cnap/pkg/cmdb"
    "github.com/viperstars/kube-cnap/pkg/image"
    "go.uber.org/zap"
    "gopkg.in/cas.v1"
    "io/ioutil"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
    v1beta "k8s.io/api/extensions/v1beta1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type DBService interface {
    CreateOrUpdateRecord(record interface{}, table interface{}, id int) (int64, error)
    AddContainer(meta basemeta.BaseMeta, kind string, base *container.Base) error
    RemoveContainer(meta basemeta.BaseMeta, kind string, id int) error
    UpdateMainContainer(base basemeta.BaseMeta, id int) error
    UpdateInitOrSidecarContainer(base basemeta.BaseMeta, kind string, ids []int) error
    RemoveMainContainer(base basemeta.BaseMeta) error
    GetCFSVolumes(base basemeta.BaseMeta) ([]*pod.CFSVolume, error)
    GetBOSVolumes(base basemeta.BaseMeta) ([]*pod.BOSVolume, error)
    GetEmptyVolumes(base basemeta.BaseMeta) ([]*pod.EmptyVolume, error)
    GetVolumes(base basemeta.BaseMeta) (*pod.Volumes, error)
    RemoveInitOrSidecarContainer(base basemeta.BaseMeta, kind string) error
    GetContainerGroup(base basemeta.BaseMeta) (*pod.ContainerGroup, error)
    DeleteContainerAttribute(id int, table interface{}) (int64, error)
    UpdateContainerAttribute(id int, table interface{}, record interface{}) (int64, error)
    GetIngressRules(base basemeta.BaseMeta, host string) ([]*ingress.IngressRule, error)
    DeleteIngressRule(id int, table interface{}) (int64, error)
    GetAllRegionsAndEnvs(withConfig bool) ([]*app.Region, error)
    GetAllApps() ([]*app.App, error)
    GetApp(id int) (*app.App, error)
    GetAppByName(name string) (*app.App, error)
    RemoveInitOrSidecarContainerByID(base basemeta.BaseMeta, kind string, id int) error
    GetContainerAttributes(containerID int, table string) (interface{}, error)
    CreateContainerGroup(base basemeta.BaseMeta) (int64, error)
    GetContainerBase(id int) (*container.Base, error)
    GetContainerGroupBase(base basemeta.BaseMeta) (*pod.ContainerGroup, error)
    GetContainerAttribute(recordID int, table string) (interface{}, error)
    GetIngressByID(recordID int) (*ingress.IngressRule, error)
    GetVolumeByID(recordID int, kind string) (interface{}, error)
    GetBaseInfo(base basemeta.BaseMeta) (string, string, error)
    CreateContainerBase(base basemeta.BaseMeta) (int64, error)
    UpdateApp(app *app.App) error
    GetHistories(base basemeta.BaseMeta) ([]*app.History, error)
    GetNotInitializedZones(app string) ([]*basemeta.CNames, error)
    GetAvailableZone(appName string) ([]*app.Region, error)
    CloneApp(info *app.CloneInfo) []error
    DeleteNodeSelector(id int) error
    UpdateNodeSelector(id int, ns *pod.NodeSelector) error
    InsertNodeSelector(base basemeta.BaseMeta, ns *pod.NodeSelector) error
    GetNodeSelector(base basemeta.BaseMeta) (*pod.NodeSelector, error)
    DeleteHostAlias(id int) error
    UpdateHostAlias(id int, ns *pod.HostAlias) error
    InsertHostAlias(base basemeta.BaseMeta, ns *pod.HostAlias) error
    GetHostAlias(base basemeta.BaseMeta) (*pod.HostAlias, error)
}

type K8sClient interface {
    CreateDeployment(base basemeta.BaseMeta, deployment *appsv1.Deployment) error
    UpdateDeployment(base basemeta.BaseMeta, deploymentName, image string) error
    RolloutDeployment(base basemeta.BaseMeta, deployment *appsv1.Deployment) error
    ScaleDeployment(base basemeta.BaseMeta, deploymentName string, number int) error
    CreateConfigMaps(base basemeta.BaseMeta, configMaps []*v1.ConfigMap) []error
    CreateService(base basemeta.BaseMeta, service *v1.Service, labels map[string]string) error
    GetDeployments(base basemeta.BaseMeta, labels map[string]string) (*appsv1.DeploymentList, error)
    GetServices(base basemeta.BaseMeta, labels map[string]string) (*v1.ServiceList, error)
    DeleteDeployment(base basemeta.BaseMeta, deploymentName string) error
    DeleteService(base basemeta.BaseMeta, serviceName string) error
    GetDeployment(base basemeta.BaseMeta, deploymentName string) (*appsv1.Deployment, error)
    GetService(base basemeta.BaseMeta, serviceName string) (*v1.Service, error)
    GetPods(base basemeta.BaseMeta, labels map[string]string) (*v1.PodList, error)
    DeletePod(base basemeta.BaseMeta, podName string) error
    GetEventsForObject(base basemeta.BaseMeta, objectName string) (*v1.EventList, error)
    GetEventsForNode(base basemeta.BaseMeta, nodeName string, allNamespaces bool) (*v1.EventList, error)
    GetNodes(base basemeta.BaseMeta) (*v1.NodeList, error)
    UpdateLabelsToNode(base basemeta.BaseMeta, nodeName string, lbs map[string]string) error
    UpdateConfigMap(base basemeta.BaseMeta, path, content string) error
    GetClientAndConfig(base basemeta.BaseMeta) (*kubernetes.Clientset, *rest.Config, string)
    EnsureGet(base basemeta.BaseMeta) error
    GetMiniDeploymentList(base basemeta.BaseMeta) ([]*client.MiniDeployment, error)
    GetMiniDeployment(base basemeta.BaseMeta, deploymentName string) (*client.MiniDeployment, error)
    GetMiniServiceList(base basemeta.BaseMeta) ([]*client.MiniService, error)
    GetMiniService(base basemeta.BaseMeta, serviceName string) (*client.MiniService, error)
    GetMiniEventList(base basemeta.BaseMeta, objectKind string, objectName string) ([]*client.MiniEvent, error)
    GetMiniNodesList(base basemeta.BaseMeta) ([]*client.MiniNode, error)
    GetMiniNode(base basemeta.BaseMeta, nodeName string) (*client.MiniNode, error)
    GetMiniPod(base basemeta.BaseMeta, podName string) (*client.MiniPod, error)
    CreateOrUpdateIngress(base basemeta.BaseMeta, ingress *v1beta.Ingress) error
    DeleteConfigMap(base basemeta.BaseMeta, path string) error
    DeleteIngress(base basemeta.BaseMeta) error
    AddResourceQuota(base basemeta.BaseMeta, rq *v1.ResourceQuota) error
    GetResourceQuota(base basemeta.BaseMeta) (*v1.ResourceQuota, error)
    UpdateResourceQuota(base basemeta.BaseMeta, rq *v1.ResourceQuota) (*v1.ResourceQuota, error)
    CreateNamespace(base basemeta.BaseMeta) error
    CreateSecret(base basemeta.BaseMeta) error
    CreateSecret2(base basemeta.BaseMeta) error
    UpdateSecret(base basemeta.BaseMeta) error
    CreateClusterRole(base basemeta.BaseMeta) error
}


type Server struct {
    dbService  DBService
    k8sClients K8sClient
    config     *conf.Server
    logger     *zap.Logger
    cas        *cas.Client
    cmdb       *cmdb.DataGetter
    bce        *bce_client.BCEClient
    reg        *image.RegistryClient
}

func (s *Server) CreateDeployment(base basemeta.BaseMeta, createRequest *requests.CreateDeploymentRequest) error {
    createRequestJson, _ := json.Marshal(createRequest)
    containerGroup, err := s.dbService.GetContainerGroup(base)
    if err != nil {
        s.logger.Error(
            "create deployment error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("request", string(createRequestJson)),
        )
        return err
    }
    volumes, err := s.dbService.GetVolumes(base)
    if err != nil {
        s.logger.Error(
            "get volumes error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    lbs := createRequest.GetLabels(base)
    podTemplateSpec, configMaps := containerGroup.ToK8sPodTemplateSpecAndConfigMaps(createRequest.Name, lbs,
        volumes.ToK8sVolumes())
    podTemplateSpec.Spec.Containers[0].Image = createRequest.Image
    _deployment := createRequest.ToK8sDeployment(lbs, podTemplateSpec)
    errs := s.k8sClients.CreateConfigMaps(base, configMaps)
    for _, e := range errs {
        fmt.Println("errors", e)
    }
    if len(errs) != 0 {
        s.logger.Error(
            "create deployment by client error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    err = s.k8sClients.CreateDeployment(base, _deployment)
    if err != nil {
        s.logger.Error(
            "create deployment by client error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("request", string(createRequestJson)),
        )
        return err
    }
    return nil
}

func (s *Server) RolloutDeployment(base basemeta.BaseMeta, deploymentName string, image string) error {
    containerGroup, err := s.dbService.GetContainerGroup(base)
    if err != nil {
        s.logger.Error(
            "rollout deployment error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    volumes, err := s.dbService.GetVolumes(base)
    if err != nil {
        s.logger.Error(
            "get volumes error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    deployment, err := s.k8sClients.GetDeployment(base, deploymentName)
    if err != nil {
        return err
    }
    podTemplateSpec, configMaps := containerGroup.ToK8sPodTemplateSpecAndConfigMaps(deployment.Name, deployment.Labels,
        volumes.ToK8sVolumes())
    k8sDeployment := &appsv1.Deployment{
        TypeMeta: metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name:   deployment.Name,
            Labels: deployment.Labels,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: deployment.Spec.Replicas,
            Template: podTemplateSpec,
            Strategy: deployment.Spec.Strategy,
            Selector: &metav1.LabelSelector{
                MatchLabels: deployment.Spec.Selector.MatchLabels,
            },
            MinReadySeconds:         0,
            RevisionHistoryLimit:    nil,
            Paused:                  false,
            ProgressDeadlineSeconds: nil,
        },
        Status: appsv1.DeploymentStatus{},
    }
    if image != "" {
        k8sDeployment.Spec.Template.Spec.Containers[0].Image = image
    }
    errs := s.k8sClients.CreateConfigMaps(base, configMaps)
    for _, e := range errs {
        fmt.Println("errors", e)
    }
    if len(errs) != 0 {
        s.logger.Error(
            "create deployment by client error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    err = s.k8sClients.RolloutDeployment(base, k8sDeployment)
    if err != nil {
        s.logger.Error(
            "rollout deployment by client error",
            zap.Error(err),
            zap.Any("meta", base),
        )
        return err
    }
    return nil
}

func (s *Server) CreateService(base basemeta.BaseMeta, service *service.CreateServiceRequest) error {
    k8sService := service.ToK8sService(base)
    createRequestJson, _ := json.Marshal(service)
    err := s.k8sClients.CreateService(base, k8sService, nil)
    if err != nil {
        s.logger.Error(
            "create service by client error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("request", string(createRequestJson)),
        )
        return err
    }
    return nil
}

func (s *Server) ToCapacitys(base basemeta.BaseMeta, quota *v1.ResourceQuota) *resourcequota.ResourceQuota {
    rq := &resourcequota.ResourceQuota{}
    rq.Base = base
    rCpu, ok := quota.Spec.Hard["requests.cpu"]
    if ok {
        rq.RequestCPU = rCpu.Value()
    } else {
        rq.RequestCPU = 0
    }
    rMem, ok := quota.Spec.Hard["requests.memory"]
    if ok {
        rq.RequestMemory = rMem.Value() / 1024 / 1024 / 1024
    } else {
        rq.RequestMemory = 0
    }
    lCpu, ok := quota.Spec.Hard["limits.cpu"]
    if ok {
        rq.LimitCPU = lCpu.Value()
    } else {
        rq.LimitCPU = 0
    }
    lMem, ok := quota.Spec.Hard["limits.memory"]
    if ok {
        rq.LimitMemory = lMem.Value() / 1024 / 1024 / 1024
    } else {
        rq.LimitMemory = 0
    }
    region, env, err := s.dbService.GetBaseInfo(base)
    if err != nil {
        rq.Region = ""
        rq.Env = ""
    }
    rq.Region = region
    rq.Env = env
    return rq
}

func (s *Server) GetCapacities(app string) ([]*resourcequota.ResourceQuota, error) {
    rqs := make([]*resourcequota.ResourceQuota, 0)
    regions, err := s.dbService.GetAllRegionsAndEnvs(false)
    if err != nil {
        return rqs, err
    }
    for _, rg := range regions {
        for _, env := range rg.Envs {
            b := basemeta.BaseMeta{
                App: app,
                Region: rg.Name,
                Env: env.Name,
            }
            rq, err := s.k8sClients.GetResourceQuota(b)
            if err == nil {
                crq := s.ToCapacitys(b, rq)
                rqs = append(rqs, crq)
            }
        }
    }
    return rqs, nil
}

func (s *Server) GetAppCapacities(request *restful.Request, response *restful.Response) {
    appName := request.PathParameter("app")
    _, err := s.dbService.GetAppByName(appName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    rqs, err := s.GetCapacities(appName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, rqs, response)
}

func (s *Server) initAppBase(base basemeta.BaseMeta) int {
    var errNumbers = 0
    _, err := s.dbService.CreateContainerGroup(base)
    if err != nil {
        s.logger.Error("create container group error", zap.Error(err))
        errNumbers += 1
    }
    err = s.k8sClients.CreateNamespace(base)
    if err != nil {
        s.logger.Error("create container namespace error", zap.Error(err))
        errNumbers += 1
    }
    err = s.k8sClients.CreateSecret(base)
    if err != nil {
        s.logger.Error("create secret error", zap.Error(err))
        errNumbers += 1
    }
    err = s.k8sClients.CreateClusterRole(base)
    if err != nil {
        s.logger.Error("create cluster role", zap.Error(err))
        errNumbers += 1
    }
    return errNumbers
}

func (s *Server) InitApp(app string) error {
    errNumbers := 0
    regions, err := s.dbService.GetAllRegionsAndEnvs(false)
    if err != nil {
        return err
    }
    for _, region := range regions {
        for _, env := range region.Envs {
            base := basemeta.BaseMeta{
                App:    app,
                Region: region.Name,
                Env:    env.Name,
            }
            errNumbers += s.initAppBase(base)
        }
    }
    if errNumbers > 0 {
        return errors.New("init app error")
    } else {
        return nil
    }
}

func (s *Server) ImageList(request *restful.Request, response *restful.Response) {
    ns := request.PathParameter("namespace")
    repo := request.PathParameter("repo")
    param := make(map[string]string)
    if ns != "" && repo != "" {
        param["namespace"] = ns
        param["repository"] = repo
        headers := make(map[string]string)
        resp, _,  err:= s.bce.Execute("GET", "/v1/image/tags", param, headers, nil)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        data, err := s.formatResponse(b)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        for i := len(data)/2-1; i >= 0; i-- {
            opp := len(data)-1-i
            data[i], data[opp] = data[opp], data[i]
        }
        s.returnResponseWithData(200, data, response)
    } else {
        s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
        return
    }
}

func (s *Server) ImageList2(request *restful.Request, response *restful.Response) {
    ns := request.PathParameter("namespace")
    repo := request.PathParameter("repo")
    param := make(map[string]string)
    if ns != "" && repo != "" {
        param["namespace"] = ns
        param["repository"] = repo
        headers := make(map[string]string)
        resp, _, err:= s.bce.Execute("GET", "/v1/image/tags", param, headers, nil)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        b, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        data, err := s.formatResponse(b)
        if err != nil {
            s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
            return
        }
        for i := len(data)/2-1; i >= 0; i-- {  // 调整顺序显示
            opp := len(data)-1-i
            data[i], data[opp] = data[opp], data[i]
        }
        d, _ := s.reg.Tags(ns, repo)
        d = append(d, data...)
        s.returnResponseWithData(200, d, response)
    } else {
        s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
        return
    }
}

func (s *Server) formatResponse(data []byte) ([]string, error) {
    resp := new(ImageResponse)
    images := make([]string, 0)
    err := json.Unmarshal(data, &resp)
    if err != nil {
        return images, err
    }
    for _, name := range resp.Tags {
        s := fmt.Sprintf("%s/%s/%s:%s", conf.Config.Registry.PullPrefix, resp.Namespace, resp.Repository, name.Name)
        images = append(images, s)
    }
    return images, nil
}

type ImageResponse struct {
    Namespace  string `json:"namespace"`
    Repository string `json:"repository"`
    Tags       []*Tag `json:"tags"`
}

type Tag struct {
    Name string `json:"name"`
}

func (s *Server) GetUninitalizedZones(request *restful.Request, response *restful.Response) {
    appName := request.PathParameter("app")
    cnames, err := s.dbService.GetNotInitializedZones(appName)
    if err != nil {
        s.returnError(400, errors.New("namespace and repo should be non-empty"), response)
        return
    }
    s.returnResponseWithData(200, cnames, response)
}

func (s *Server) GetAvailableZones(request *restful.Request, response *restful.Response) {
    appName := request.QueryParameter("app")
    zones, err := s.dbService.GetAvailableZone(appName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, zones, response)
}