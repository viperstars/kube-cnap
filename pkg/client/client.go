package client

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/viperstars/kube-cnap/conf"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/naming"
    "go.uber.org/zap"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
    "k8s.io/api/extensions/v1beta1"
    rbacv1 "k8s.io/api/rbac/v1"
    k8sError "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "net/http"
    "path/filepath"
)

var NOTFOUNDCLIENT = errors.New("region or env not found")

type DBService interface {
    GetAllRegionsAndEnvs(withConfig bool) ([]*app.Region, error)
}

type Client struct {
    client    *kubernetes.Clientset
    config    *rest.Config
    namespace string
}

func NewClient(config string, ns string) (*Client, error) {
    client := new(Client)
    restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config))
    if err != nil {
        return client, err
    }
    client.client = kubernetes.NewForConfigOrDie(restConfig)
    client.namespace = ns
    client.config = restConfig
    return client, nil
}

type Clients struct {
    ClientPool map[string]map[string]*Client
    logger     *zap.Logger
    db         DBService
}

func (c *Clients) getClient(base basemeta.BaseMeta) (*kubernetes.Clientset, string) {
    return c.ClientPool[base.Region][base.Env].client, ""
}

func (c *Clients) EnsureGet(base basemeta.BaseMeta) error {
    fmt.Println(base.Region, base.Env)
    region, ok := c.ClientPool[base.Region]
    if !ok {
        return NOTFOUNDCLIENT
    }
    _, ok = region[base.Env]
    if !ok {
        return NOTFOUNDCLIENT
    }
    return nil
}

func (c *Clients) getConfig(base basemeta.BaseMeta) (*rest.Config, string) {
    //return new(kubernetes.Clientset), ""
    return c.ClientPool[base.Region][base.Env].config, ""
}

func (c *Clients) IsNotFoundError(err error) bool {
    status, ok := err.(*k8sError.StatusError)
    if !ok {
        return false
    }
    return status.ErrStatus.Code == http.StatusNotFound
}

func (c *Clients) CreateDeployment(base basemeta.BaseMeta, deployment *appsv1.Deployment) error {
    dpJson, _ := json.Marshal(deployment)
    client, _ := c.getClient(base)
    _, err := client.AppsV1().Deployments(base.App).Create(deployment)
    c.logger.Info(
        "creating deployment",
        zap.Any("meta", base),
        zap.String("deployment", string(dpJson)),
    )
    if err != nil {
        c.logger.Error(
            "create deployment error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("deployment", string(dpJson)),
        )
    }
    return err
}

func (c *Clients) UpdateDeployment(base basemeta.BaseMeta, deploymentName, image string) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting deployment",
        zap.Any("meta", base),
        zap.String("deploymentName", deploymentName),
    )
    deployment, err := client.AppsV1().Deployments(base.App).Get(deploymentName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "get deployment error",
            zap.String("error", err.Error()),
            zap.String("deploymentName", deploymentName),
        )
        return err
    }
    deployment.Spec.Template.Spec.Containers[0].Image = image
    c.logger.Info(
        "updating deployment image",
        zap.Any("meta", base),
        zap.String("deploymentName", deploymentName),
        zap.String("image", image),
    )
    deployment, err = client.AppsV1().Deployments(base.App).Update(deployment)
    if err != nil {
        c.logger.Error("update deployment image error", zap.Any("meta", base), zap.Error(err),
            zap.String("deploymentName", deploymentName),
            zap.String("image", image),
        )
    }
    return err
}

func (c *Clients) DeleteDeployment(base basemeta.BaseMeta, deploymentName string) error {
    client, _ := c.getClient(base)
    err := client.AppsV1().Deployments(base.App).Delete(deploymentName, &metav1.DeleteOptions{})
    c.logger.Info(
        "deleting deployment",
        zap.Any("meta", base),
        zap.String("deployment", deploymentName),
    )
    if err != nil {
        c.logger.Error(
            "delete deployment error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("deployment", deploymentName),
        )
    }
    return err
}

func (c *Clients) ScaleDeployment(base basemeta.BaseMeta, deploymentName string, number int) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting deployment",
        zap.Any("meta", base),
        zap.String("deploymentName", deploymentName),
    )
    scale, err := client.AppsV1().Deployments(base.App).GetScale(deploymentName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "get deployment error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("deploymentName", deploymentName),
        )
        return err
    }
    r := int32(number)
    scale.Spec.Replicas = r
    c.logger.Info(
        "scaling deployment",
        zap.Any("meta", base),
        zap.String("deploymentName", deploymentName),
        zap.Int32("replicas", r),
    )
    _, err = client.AppsV1().Deployments(base.App).UpdateScale(deploymentName, scale)
    if err != nil {
        c.logger.Error("scale deployment error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("deploymentName", deploymentName),
            zap.Int32("replicas", r))
    }
    return err
}

func (c *Clients) CreateConfigMaps(base basemeta.BaseMeta, configMap []*v1.ConfigMap) []error {
    errorList := make([]error, 0)
    for _, configMap := range configMap {
        err := c.CreateOrUpdateConfigMaps(base, configMap)
        if err != nil {
            errorList = append(errorList, err)
        }
    }
    return errorList
}

func (c *Clients) CreateService(base basemeta.BaseMeta, service *v1.Service, labels map[string]string) error {
    client, _ := c.getClient(base)
    srvJson, _ := json.Marshal(service)
    c.logger.Info(
        "creating service",
        zap.Any("meta", base),
        zap.String("service", string(srvJson)),
    )
    _, err := client.CoreV1().Services(base.App).Create(service)
    if err != nil {
        c.logger.Error(
            "create service error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("service", string(srvJson)),
        )
    }
    return err
}

func (c *Clients) DeleteService(base basemeta.BaseMeta, serviceName string) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "deleting service",
        zap.Any("meta", base),
        zap.String("service", serviceName),
    )
    err := client.CoreV1().Services(base.App).Delete(serviceName, &metav1.DeleteOptions{})
    if err != nil {
        c.logger.Error(
            "create service error",
            zap.Any("meta", base),
            zap.Error(err), zap.String("service", serviceName),
        )
    }
    return err
}

func (c *Clients) DeletePod(base basemeta.BaseMeta, podName string) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "deleting service",
        zap.Any("meta", base),
        zap.String("pod", podName),
    )
    err := client.CoreV1().Pods(base.App).Delete(podName, &metav1.DeleteOptions{})
    if err != nil {
        c.logger.Error(
            "create service error",
            zap.Any("meta", base), zap.Error(err),
            zap.String("pod", podName),
        )
    }
    return err
}

func (c *Clients) CreateOrUpdateConfigMaps(base basemeta.BaseMeta, configMap *v1.ConfigMap) error {
    var err error
    client, _ := c.getClient(base)
    cf, err := client.CoreV1().ConfigMaps(base.App).Get(configMap.Name, metav1.GetOptions{})
    if err != nil { // is not found error
        _, err = client.CoreV1().ConfigMaps(base.App).Create(configMap)
    } else {
        cf.Data = configMap.Data
        _, err := client.CoreV1().ConfigMaps(base.App).Update(cf)
        return err
    }
    return err
}


func (c *Clients) RolloutDeployment(base basemeta.BaseMeta, deployment *appsv1.Deployment) error {
    client, _ := c.getClient(base)
    deploymentJson, _ := json.Marshal(deployment)
    deployment, err := client.AppsV1().Deployments(base.App).Update(deployment)
    if err != nil {
        c.logger.Error("update deployment image error", zap.Any("meta", base), zap.Error(err),
            zap.String("newDeployment", string(deploymentJson)),
        )
    }
    return nil
}

func (c *Clients) CreateSecret(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    data := make(map[string]string)
    secretString, err := c.CreateSecretString(conf.Config)
    if err != nil {
        return err
    }
    data[v1.DockerConfigJsonKey] = secretString
    secret := &v1.Secret{
        TypeMeta:   metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name: conf.Config.Registry.KeyName, // TODO
            Namespace: base.App,
        },
        Data:       nil,
        StringData: data,
        Type:       v1.SecretTypeDockerConfigJson,
    }
    _, err = client.CoreV1().Secrets(base.App).Create(secret)
    if err != nil {
        c.logger.Error("create secret error", zap.Any("meta", base), zap.Error(err))
    }
    err = c.CreateSecret2(base)
    if err != nil {
        c.logger.Error("create secret2 error", zap.Any("meta", base), zap.Error(err))
    }
    return nil
}

func (c *Clients) CreateSecret2(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    data := make(map[string]string)
    secretString, err := c.CreateSecretString2(conf.Config)
    if err != nil {
        return err
    }
    data[v1.DockerConfigJsonKey] = secretString
    secret := &v1.Secret{
        TypeMeta:   metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name: conf.Config.RegistryInfo.KeyName, // TODO
            Namespace: base.App,
        },
        Data:       nil,
        StringData: data,
        Type:       v1.SecretTypeDockerConfigJson,
    }
    _, err = client.CoreV1().Secrets(base.App).Create(secret)
    if err != nil {
        c.logger.Error("create secret error", zap.Any("meta", base), zap.Error(err))
    }
    return nil
}

func (c *Clients) UpdateSecret(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    data := make(map[string]string)
    secretString, err := c.CreateSecretString(conf.Config)
    if err != nil {
        return err
    }
    secret, err := client.CoreV1().Secrets(base.App).Get(conf.Config.RegistryInfo.KeyName, metav1.GetOptions{})
    if err != nil {
        return err
    }
    data[v1.DockerConfigJsonKey] = secretString
    secret.StringData = data
    _, err = client.CoreV1().Secrets(base.App).Update(secret)
    if err != nil {
        c.logger.Error("update secret error", zap.Any("meta", base), zap.Error(err))
    }
    return nil
}

func (c *Clients) CreateSecretString(conf2 *conf.Conf) (string, error) {
    auth := conf2.Registry.Auth // Todo
    data := make(map[string]interface{})
    auths := make(map[string]interface{})
    auths[conf2.Registry.ServerAddress] = auth // TODO
    data["auths"] = auths
    body, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(body[:]), nil
}

func (c *Clients) CreateSecretString2(conf2 *conf.Conf) (string, error) {
    auth := conf2.RegistryInfo.Auth
    data := make(map[string]interface{})
    auths := make(map[string]interface{})
    auths[conf2.RegistryInfo.Addr] = auth
    data["auths"] = auths
    body, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(body[:]), nil
}


func (c *Clients) CreateOrUpdateIngress(base basemeta.BaseMeta, ingress *v1beta1.Ingress) error {
    var err error
    client, _ := c.getClient(base)
    igs, err := client.ExtensionsV1beta1().Ingresses(base.App).Get(ingress.Name, metav1.GetOptions{})
    if err != nil { // is not found error
        _, err = client.ExtensionsV1beta1().Ingresses(base.App).Create(ingress)
        return err
    } else {
        copiedIgs := igs.DeepCopy()
        copiedIgs.Spec = ingress.Spec
        _, err := client.ExtensionsV1beta1().Ingresses(base.App).Update(ingress)
        return err
    }
}

func (c *Clients) CreateIngressRule(base basemeta.BaseMeta, ingress *v1beta1.Ingress) error {
    ingressJson, _ := json.Marshal(ingress)
    c.logger.Info(
        "creating ingress",
        zap.Any("meta", base),
        zap.String("ingress", string(ingressJson)),
    )
    client, _ := c.getClient(base)
    _, err := client.ExtensionsV1beta1().Ingresses(base.App).Create(ingress)
    if err != nil {
        c.logger.Error(
            "get configMap error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("ingress", string(ingressJson)),
        )
        return err
    }
    return nil
}

func (c *Clients) UpdateIngressRule(base basemeta.BaseMeta, ingress *v1beta1.Ingress) error {
    ingressJson, _ := json.Marshal(ingress)
    c.logger.Info(
        "updating ingress",
        zap.Any("meta", base),
        zap.String("ingress", string(ingressJson)),
    )
    client, _ := c.getClient(base)
    _, err := client.ExtensionsV1beta1().Ingresses(base.App).Update(ingress)
    if err != nil {
        c.logger.Error(
            "update ingress error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("ingress", string(ingressJson)),
        )
        return err
    }
    return nil
}

func (c *Clients) DeleteIngressRule(base basemeta.BaseMeta, ingressName string) error {
    client, _ := c.getClient(base)
    err := client.ExtensionsV1beta1().Ingresses(base.App).Delete(ingressName, &metav1.DeleteOptions{})
    if err != nil {
        c.logger.Error(
            "get configMap error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("ingress", ingressName),
        )
        return err
    }
    return nil
}

func (c *Clients) updateConfigMap(base basemeta.BaseMeta, configMapName, filename, content string) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting configMap",
        zap.Any("meta", base),
        zap.String("configMapName", configMapName),
    )
    configMap, err := client.CoreV1().ConfigMaps(base.App).Get(configMapName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "get configMap error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("configMapName", configMapName),
        )
        return err
    }
    configMap.Data[filename] = content
    c.logger.Info(
        "updating configMap",
        zap.Any("meta", base),
        zap.String("configMapName", configMapName),
        zap.String("key", filename),
        zap.String("content", content),
    )
    configMap, err = client.CoreV1().ConfigMaps(base.App).Update(configMap)
    if err != nil {
        c.logger.Error(
            "update configMap error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("configMapName", configMapName),
            zap.String("key", filename),
            zap.String("content", content),
        )
    }
    return err
}

func (c *Clients) deleteConfigMap(base basemeta.BaseMeta, configMapName, filename string) error {
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting configMap",
        zap.Any("meta", base),
        zap.String("configMapName", configMapName),
    )
    configMap, err := client.CoreV1().ConfigMaps(base.App).Get(configMapName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "get configMap error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("configMapName", configMapName),
        )
        return err
    } // TODO add delete all
    delete(configMap.Data, filename)
    c.logger.Info(
        "deleting configMap",
        zap.Any("meta", base),
        zap.String("configMapName", configMapName),
        zap.String("key", filename),
    )
    if len(configMap.Data) > 0 {
        configMap, err = client.CoreV1().ConfigMaps(base.App).Update(configMap)
        if err != nil {
            c.logger.Error(
                "delete configMap error",
                zap.Any("meta", base),
                zap.Error(err),
                zap.String("configMapName", configMapName),
                zap.String("key", filename),
            )
        }
        return err
    } else {
        err = client.CoreV1().ConfigMaps(base.App).Delete(configMap.Name, &metav1.DeleteOptions{})
        if err != nil {
            c.logger.Error(
                "delete configMap error",
                zap.Any("meta", base),
                zap.Error(err),
                zap.String("configMapName", configMapName),
                zap.String("key", filename),
            )
        }
        return err
    }
}

func (c *Clients) UpdateConfigMap(base basemeta.BaseMeta, path, content string) error {
    dirName := filepath.Dir(path)
    baseName := filepath.Base(path)
    configMapName := naming.Namer.ConfigMapName(&base, dirName)
    return c.updateConfigMap(base, configMapName, baseName, content)
}

func (c *Clients) DeleteConfigMap(base basemeta.BaseMeta, path string) error {
    dirName := filepath.Dir(path)
    baseName := filepath.Base(path)
    configMapName := naming.Namer.ConfigMapName(&base, dirName)
    return c.deleteConfigMap(base, configMapName, baseName)
}

func (c *Clients) GetPod(base basemeta.BaseMeta, podName string) (*v1.Pod, error) {
    var pod *v1.Pod
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting pod",
        zap.Any("meta", base),
        zap.String("podName", podName),
    )
    pod, err := client.CoreV1().Pods(base.App).Get(podName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "getting pod error",
            zap.Any("meta", base),
            zap.Error(err),
            zap.String("podName", podName),
        )
    }
    return pod, err
}

func (c *Clients) GetDeployment(base basemeta.BaseMeta, deploymentName string) (*appsv1.Deployment, error) {
    var deployment *appsv1.Deployment
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting deployment",
        zap.Any("meta", base),
        zap.String("deploymentName", deploymentName),
    )
    deployment, err := client.AppsV1().Deployments(base.App).Get(deploymentName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "getting deployment error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("deploymentName", deploymentName),
        )
    }
    return deployment, err
}

func (c *Clients) toListOptions(labels map[string]string) (metav1.ListOptions, error) {
    labelSelector, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{MatchLabels: labels})
    options := metav1.ListOptions{
        LabelSelector: labelSelector.String(),
    }
    return options, nil
}

func (c *Clients) GetDeployments(base basemeta.BaseMeta, labels map[string]string) (*appsv1.DeploymentList, error) {
    options, err := c.toListOptions(labels)
    if err != nil {
        return nil, err
    }
    client, _ := c.getClient(base)
    c.logger.Info(
        "getting deployment list",
        zap.Any("meta", base),
    )
    deploymentList, err := client.AppsV1().Deployments(base.App).List(options)
    if err != nil {
        c.logger.Error(
            "getting deployment list error",
            zap.Any("meta", base),
            zap.Error(err),
        )
    }
    return deploymentList, err
}

func (c *Clients) GetPods(base basemeta.BaseMeta, labels map[string]string) (*v1.PodList, error) {
    options, err := c.toListOptions(labels)
    if err != nil {
        return nil, err
    }
    c.logger.Info(
        "getting pod list",
        zap.Any("meta", base),
    )
    client, _ := c.getClient(base)
    podList, err := client.CoreV1().Pods(base.App).List(options)
    if err != nil {
        c.logger.Error(
            "getting pod list error",
            zap.Any("meta", base),
            zap.Error(err),
        )
    }
    return podList, err
}

func (c *Clients) GetPodsWithString(base basemeta.BaseMeta, labels string) (*v1.PodList, error) {
    client, _ := c.getClient(base)
    podList, err := client.CoreV1().Pods(base.App).List(metav1.ListOptions{LabelSelector: labels})
    return podList, err
}

func (c *Clients) GetReplicaSets(base basemeta.BaseMeta, labels string) (*appsv1.ReplicaSetList, error) {
    client, _ := c.getClient(base)
    replicaList, err := client.AppsV1().ReplicaSets(base.App).List(metav1.ListOptions{LabelSelector: labels})
    return replicaList, err
}

func (c *Clients) FilterPods(pods []v1.Pod) map[string][]v1.Pod {
    filteredPods := make(map[string][]v1.Pod)
    for _, p := range pods {
        if _, ok := filteredPods[p.Spec.Containers[0].Image]; ok {
            filteredPods[p.Spec.Containers[0].Image] = append(filteredPods[p.Spec.Containers[0].Image], p)
        } else {
            podList := make([]v1.Pod, 0)
            podList = append(podList, p)
            filteredPods[p.Spec.Containers[0].Image] = podList
        }
    }
    return filteredPods
}

func (c *Clients) GetServices(base basemeta.BaseMeta, labels map[string]string) (*v1.ServiceList, error) {
    client, _ := c.getClient(base)
    selector, _ := ToLabelSelector(labels)
    options := metav1.ListOptions{LabelSelector: selector.String()}
    return client.CoreV1().Services(base.App).List(options)
}

func (c *Clients) GetService(base basemeta.BaseMeta, serviceName string) (*v1.Service, error) {
    client, _ := c.getClient(base)
    return client.CoreV1().Services(base.App).Get(serviceName, metav1.GetOptions{})
}

func (c *Clients) GetPodListForDeployment(base basemeta.BaseMeta,
    dpList *appsv1.DeploymentList) (*SimpleDeploymentList, error) {
    simple := &SimpleDeploymentList{}
    simple.Deployments = make([]*SimpleDeployment, 0)
    for _, dep := range dpList.Items {
        lbs, _ := ToLabelSelector(dep.Spec.Selector.MatchLabels)
        pod, _ := c.GetPodsWithString(base, lbs.String())
        podInfo := GetPodInfo(*dep.Spec.Replicas, dep.Status.Replicas, pod.Items)
        replicaSets, _ := c.GetReplicaSets(base, lbs.String())
        simpleReplicaSets, _ := c.GetPodsForReplicaSets(base, replicaSets.Items)
        sp := &SimpleDeployment{
            Deployment:  &dep,
            ReplicaSets: simpleReplicaSets,
            PodInfo:     podInfo,
        }
        simple.Deployments = append(simple.Deployments, sp)
    }
    return simple, nil
}

func (c *Clients) GetPodsForReplicaSet(base basemeta.BaseMeta, replicaSet appsv1.ReplicaSet) (*SimpleReplicaSet, error) {
    var simple *SimpleReplicaSet
    lbs, _ := ToLabelSelector(replicaSet.Spec.Selector.MatchLabels)
    pod, err := c.GetPodsWithString(base, lbs.String())
    if err != nil {
        c.logger.Error(
            "get pod with string error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("label", lbs.String()),
        )
        return simple, err
    }
    pods := FilterPodsByControllerRef(&replicaSet, pod.Items)
    podInfo := GetPodInfo(*replicaSet.Spec.Replicas, replicaSet.Status.Replicas, pods)
    simple = &SimpleReplicaSet{
        ReplicaSet: &replicaSet,
        PodList:    pods,
        PodInfo:    podInfo,
    }
    return simple, nil
}

func (c *Clients) GetPodsForReplicaSets(base basemeta.BaseMeta, replicaSetList []appsv1.ReplicaSet) ([]*SimpleReplicaSet,
    error) {
    replicaSets := make([]*SimpleReplicaSet, 0)
    for _, replicaSet := range replicaSetList {
        simple, err := c.GetPodsForReplicaSet(base, replicaSet)
        if err != nil {

        }
        replicaSets = append(replicaSets, simple)
    }
    return replicaSets, nil
}

func (c *Clients) GetEventsForObject(base basemeta.BaseMeta, objectName string) (*v1.EventList, error) {
    client, _ := c.getClient(base)
    events, err := client.CoreV1().Events(base.App).List(metav1.ListOptions{
        FieldSelector: "involvedObject.name=" + objectName,
    })
    if err != nil {
        c.logger.Error(
            "get events for object error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("objectName", objectName),
        )
    }
    return events, nil
}

func (c *Clients) GetEventsForNode(base basemeta.BaseMeta, nodeName string, allNamespaces bool) (*v1.EventList,
    error) {
    scheme := runtime.NewScheme()
    groupVersion := schema.GroupVersion{Group: "", Version: "v1"}
    scheme.AddKnownTypes(groupVersion, &v1.Node{})
    client, _ := c.getClient(base)
    mc := client.CoreV1().Nodes()
    var err error
    var events *v1.EventList
    node, err := mc.Get(nodeName, metav1.GetOptions{})
    if err != nil {
        c.logger.Error(
            "get events for node error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("nodeName", nodeName),
        )
        return events, err
    }
    if allNamespaces {
        events, err = client.CoreV1().Events(v1.NamespaceAll).Search(scheme, node)
    } else {
        events, err = client.CoreV1().Events(base.App).Search(scheme, node)
    }
    if err != nil {
        c.logger.Error(
            "get events for node error",
            zap.Error(err),
            zap.Any("meta", base),
            zap.String("nodeName", nodeName),
        )
    }
    return events, err
}

func (c *Clients) UpdateLabelsToNode(base basemeta.BaseMeta, nodeName string, labels map[string]string) error {
    client, _ := c.getClient(base)
    node, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
    if err != nil {
        return err
    }
    copyedNode := node.DeepCopy()
    for k, v := range labels {
        copyedNode.Labels[k] = v
    }
    for k := range copyedNode.Labels {
        if _, ok := labels[k]; !ok {
            delete(copyedNode.Labels, k)
        }
    }
    _, err = client.CoreV1().Nodes().Update(copyedNode)
    return err
}

func (c *Clients) DeleteIngress(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    name := naming.Namer.IngressName(&base)
    return client.ExtensionsV1beta1().Ingresses(base.App).Delete(name,
        &metav1.DeleteOptions{})
}

func (c *Clients) GetNodes(base basemeta.BaseMeta) (*v1.NodeList, error) {
    client, _ := c.getClient(base)
    return client.CoreV1().Nodes().List(metav1.ListOptions{})
}

func (c *Clients) GetNode(base basemeta.BaseMeta, nodeName string) (*v1.Node, error) {
    client, _ := c.getClient(base)
    return client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
}

func (c *Clients) GetClientAndConfig(base basemeta.BaseMeta) (*kubernetes.Clientset, *rest.Config, string) {
    client, _ := c.getClient(base)
    config, _ := c.getConfig(base)
    return client, config, "ns"
}

func (c *Clients) AddResourceQuota(base basemeta.BaseMeta, rq *v1.ResourceQuota) error {
    client, _ := c.getClient(base)
    _, err := client.CoreV1().ResourceQuotas(base.App).Create(rq)
    return err
}

func (c *Clients) GetResourceQuota(base basemeta.BaseMeta) (*v1.ResourceQuota, error) {
    client, _ := c.getClient(base)
    return client.CoreV1().ResourceQuotas(base.App).Get(base.App, metav1.GetOptions{})
}

func (c *Clients) UpdateResourceQuota(base basemeta.BaseMeta, rq *v1.ResourceQuota) (*v1.ResourceQuota, error) {
    client, _ := c.getClient(base)
    return client.CoreV1().ResourceQuotas(base.App).Update(rq)
}

func (c *Clients) LoadFromDB() error {
    regions, err := c.db.GetAllRegionsAndEnvs(true)
    if err != nil {
        return err
    }
    regionClients := make(map[string]map[string]*Client)
    for _, region := range regions {
        clients := make(map[string]*Client)
        for _, env := range region.Envs {
            c, err := NewClient(env.Config, env.Namespace)
            if err != nil {
                return nil
            }
            clients[env.Name] = c
        }
        regionClients[region.Name] = clients
    }
    c.ClientPool = regionClients
    return nil
}

func (c *Clients) CreateNamespace(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    nameSpace := &v1.Namespace{
        TypeMeta:   metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name: base.App,
        },
        Spec:       v1.NamespaceSpec{},
        Status:     v1.NamespaceStatus{},
    }
    _, err :=  client.CoreV1().Namespaces().Create(nameSpace)
    return err
}

func (c *Clients) CreateClusterRole(base basemeta.BaseMeta) error {
    client, _ := c.getClient(base)
    role := &rbacv1.ClusterRoleBinding{
        TypeMeta:   metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name: base.App,
        },
        Subjects: []rbacv1.Subject{
            {Kind: "ServiceAccount", Name: "default", Namespace: base.App},
        },
        RoleRef: rbacv1.RoleRef{
            APIGroup: "rbac.authorization.k8s.io",
            Kind:     "ClusterRole",
            Name:     "cluster-admin",
        },
    }
    _, err := client.RbacV1().ClusterRoleBindings().Create(role)
    return err
}
