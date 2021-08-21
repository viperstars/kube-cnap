package http

import (
    "fmt"
    "github.com/emicklei/go-restful"
    "github.com/emicklei/go-restful-swagger12"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    container2 "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
    "github.com/viperstars/kube-cnap/pkg/apis/requests"
    "github.com/viperstars/kube-cnap/pkg/apis/resourcequota"
    "github.com/viperstars/kube-cnap/pkg/apis/service"
    "net/http"
)

const mineType = "application/json"

func (s *Server) Register(container *restful.Container) {
    ws := new(restful.WebService)
    ws.Path("/api/v1").Consumes(mineType).Produces(mineType)
    ws.Route(ws.GET("/{app}/{region}/{env}/deployments").To(s.GetDeploymentsHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/deployments").To(s.CreateDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).Filter(s.validateCapacity).
        //Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(requests.CreateDeploymentRequest{}).Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/deployment/{deploymentName}").To(s.GetDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/deployment/{deploymentName}").To(s.UpdateDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Reads(requests.UpdateDeploymentRequest{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/deployment/{deploymentName}/scale").To(s.ScaleDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Reads(requests.ScaleDeploymentRequest{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/deployment/{deploymentName}/rollout").To(s.RolloutDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Reads(struct{}{}).Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/services").To(s.GetServicesHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/services").To(s.CreateServiceHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(service.CreateServiceRequest{}).Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/service/{serviceName}").To(s.GetServiceHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("serviceName", "service name")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/{objectKind}/{objectName}/events").To(s.GetEventsHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("objectKind", "object kind")).
        Param(ws.PathParameter("objectName", "object name")).
        Writes(Response{}))

    ws.Route(ws.DELETE("/{app}/{region}/{env}/deployment/{deploymentName}").To(s.DeleteDeploymentHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/service/{serviceName}").To(s.DeleteServiceHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("serviceName", "service name")).
        Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/pod/{podName}").To(s.DeletePodHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("podName", "pod name")).
        Writes(Response{}))

    ws.Route(ws.POST("/{app}/{region}/{env}/containerGroup/containerBase/{kind}").To(s.AddContainerBase).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("kind", "kind of container")).
        Reads(container2.Base{}).Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/containerGroup/containerBase/{kind}/{id}").To(s.
        RemoveContainerBase).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("kind", "kind of container")).
        Param(ws.PathParameter("id", "id")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/containerGroup/containerBase/order/{kind}").To(s.
        UpdateInitOrSidecarContainers).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("kind", "kind of container")).
        Reads(IDs{}).Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/nodes").To(s.GetNodesHandler).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/node/{nodeName}").To(s.UpdateNodeLabels).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("nodeName", "node name")).
        Reads(make(map[string]string)).Writes(Response{}))
    // ws.Route(ws.POST("/containerBase").To(s.AddContainerBase).
        //Reads(container2.Base{}).Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/containerGroup").To(s.GetContainerGroup).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/ingresses").To(s.GetIngress).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/ingresses").To(s.CreateOrUpdateIngress).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(ingress.IngressRule{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/ingress/{id}").To(s.CreateOrUpdateIngress).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "ingress id")).
        Reads(ingress.IngressRule{}).Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/ingress/{id}").To(s.DeleteIngress).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "ingress id")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/ingress/{id}").To(s.GetIngressByID).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "ingress id")).
        Writes(Response{}))
    ws.Route(ws.POST("/containerBase/{id}").To(s.UpdateContainerBase).
        Param(ws.PathParameter("id", "container id")).
        Reads(container2.Base{}).Writes(Response{}))
    ws.Route(ws.GET("/containerBase/{id}").To(s.GetContainerBase).
        Param(ws.PathParameter("id", "container id")).
        Reads(container2.Base{}).Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/deployments").To(s.WrappedGetDeployments).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/deployment/{deploymentName}").To(s.WrappedGetDeployment).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("deploymentName", "deployment name")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/services").To(s.WrappedGetServices).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/service/{serviceName}").To(s.WrappedGetService).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("serviceName", "service name")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/{objectKind}/{objectName}/events").To(s.WrappedGetMiniEventList).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("objectName", "object name")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/nodes").To(s.WrappedGetMiniNodeList).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/node/{nodeName}").To(s.WrappedGetMiniNode).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/{shell}/{podName}/{containerName}").To(s.handleExecShell).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("podName", "pod name")).
        Param(ws.PathParameter("shell", "shell type")).
        Param(ws.PathParameter("containerName", "container name")).
        Writes(Response{}))

    ws.Route(ws.POST("/container/{containerID}/command").To(s.AddContainerCommand).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.Command{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/ports").To(s.AddContainerPorts).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.Port{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/configurations").To(s.AddContainerConfiguration).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.Configuration{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/livenessCheck").To(s.AddContainerLivenessCheck).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.LivenessCheck{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/readnessCheck").To(s.AddContainerReadnessCheck).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.ReadnessCheck{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/envVars").To(s.AddContainerEnvVars).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.EnvVar{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/resourceRequirement").To(s.AddContainerResourceRequirement).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.ResourceRequirement{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/volumeMounts").To(s.AddContainerVolumeMounts).
        Param(ws.PathParameter("containerID", "container ID")).
        Reads(container2.VolumeMount{}).Writes(Response{}))

    ws.Route(ws.GET("/container/{containerID}/{record}").To(s.GetContainerAttributes).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("record", "record type")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/volumes/cfs").To(s.GetCFSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/volumes/bos").To(s.GetBOSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/volumes/empty").To(s.GetEmptyVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volumes/cfs").To(s.CreateOrUpdateCFSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(pod.CFSVolume{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volumes/bos").To(s.CreateOrUpdateBOSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(pod.BOSVolume{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volumes/empty").To(s.CreateOrUpdateEmptyVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(pod.EmptyVolume{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volume/cfs/{id}").To(s.CreateOrUpdateCFSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "id of volume")).
        Reads(pod.CFSVolume{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volume/bos/{id}").To(s.CreateOrUpdateBOSVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "id of volume")).
        Reads(pod.BOSVolume{}).Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/volume/empty/{id}").To(s.CreateOrUpdateEmptyVolumes).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "id of volume")).
        Reads(pod.EmptyVolume{}).Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/volume/{kind}/{id}").To(s.DeleteVolume).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("kind", "kind of volume")).
        Param(ws.PathParameter("id", "id of volume")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/volume/{kind}/{id}").To(s.GetVolumeByID).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("kind", "kind of volume")).
        Param(ws.PathParameter("id", "id of volume")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/resourceQuota").To(s.GetResourceQuota).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/wrapped/resourceQuota").To(s.WrappedGetResourceQuota).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/resourceQuota").To(s.CreateResourceQuota).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(resourcequota.ResourceQuota{}).Writes(Response{}))

    ws.Route(ws.POST("/{app}/{region}/{env}/resourceQuota/update").To(s.UpdateResourceQuota).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(resourcequota.ResourceQuota{}).Writes(Response{}))

    ws.Route(ws.GET("/{app}/resourceQuota").To(s.GetAppCapacities).
        //Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/dashboard").To(s.dashboard).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))

    ws.Route(ws.POST("/container/{containerID}/command/{recordID}").To(s.AddContainerCommand).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.Command{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/port/{recordID}").To(s.AddContainerPorts).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.Port{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/configuration/{recordID}").
        To(s.AddContainerConfiguration).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.Configuration{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/livenessCheck/{recordID}").To(s.AddContainerLivenessCheck).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.LivenessCheck{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/readnessCheck/{recordID}").To(s.AddContainerReadnessCheck).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.ReadnessCheck{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/envVar/{recordID}").To(s.AddContainerEnvVars).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.EnvVar{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/volumeMount/{recordID}").To(s.AddContainerVolumeMounts).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.EnvVar{}).Writes(Response{}))
    ws.Route(ws.POST("/container/{containerID}/resourceRequirement/{recordID}").
        To(s.AddContainerResourceRequirement).
        Param(ws.PathParameter("containerID", "container ID")).
        Param(ws.PathParameter("recordID", "record ID")).
        Reads(container2.ResourceRequirement{}).Writes(Response{}))

    ws.Route(ws.GET("/apps").To(s.GetAllApps).
        Writes(Response{}))
    ws.Route(ws.POST("/apps").To(s.AddOrUpdateApp).
        Reads(app.App{}).Writes(Response{}))
    ws.Route(ws.GET("/app/{appID}").To(s.GetApp).
        Param(ws.PathParameter("appID", "app ID")).
        Writes(Response{}))
    ws.Route(ws.GET("/app/{name}").To(s.GetAppByName).
        Param(ws.PathParameter("name", "app name")).
        Writes(Response{}))
    ws.Route(ws.POST("/app/{appID}").To(s.AddOrUpdateApp).
        Param(ws.PathParameter("appID", "app ID")).
        Reads(app.App{}).Writes(Response{}))
    ws.Route(ws.DELETE("/app/{appID}").To(s.AddOrUpdateApp). // todo delete app
        Param(ws.PathParameter("appID", "app ID")).
        Reads(app.App{}).Writes(Response{}))

    ws.Route(ws.GET("/{app}/uninit").To(s.GetUninitalizedZones). // todo delete app
        Param(ws.PathParameter("app", "app name")).
        Writes(Response{}))

    ws.Route(ws.POST("/{app}/{region}/{env}/init").To(s.initZone).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(struct {}{}).Writes(Response{}))

    ws.Route(ws.GET("/users").To(s.Users).
        Writes(Response{}))

    ws.Route(ws.GET("/regions/envs").To(s.GetAvailableZones).
        Param(ws.QueryParameter("app", "app name")).
        Writes(Response{}))
    ws.Route(ws.GET("/regions/envs/all").To(s.GetAllRegions).
        Param(ws.QueryParameter("app", "app name")).
        Writes(Response{}))
    ws.Route(ws.POST("/regions").To(s.AddOrUpdateRegion).
        Writes(Response{}))
    ws.Route(ws.POST("/region/{regionID}").To(s.AddOrUpdateRegion).
        Param(ws.PathParameter("regionID", "region ID")).
        Writes(Response{}))
    ws.Route(ws.POST("/region/{regionID}").To(s.AddOrUpdateRegion).
        Param(ws.PathParameter("regionID", "region ID")).
        Param(ws.PathParameter("regionID", "region ID")).
        Writes(Response{}))
    ws.Route(ws.POST("/region/{regionID}/envs").To(s.AddOrUpdateEnvForRegion).
        Reads(app.Env{}).Writes(Response{}))
    ws.Route(ws.POST("/region/{regionID}/env/{envID}").To(s.AddOrUpdateEnvForRegion).
        Reads(app.Env{}).Writes(Response{}))
    ws.Route(ws.DELETE("/region/{regionID}/env/{envID}").To(s.RemoveEnvForRegion).
        Reads(app.Env{}).Writes(Response{}))
   /* ws.Route(ws.GET("/app/{appID}").To(s.GetApp).
        Param(ws.PathParameter("appID", "app ID")).
        Writes(Response{}))*/
    ws.Route(ws.GET("/app/{name}").To(s.GetAppByName).
        Param(ws.PathParameter("name", "app name")).
        Writes(Response{}))
    ws.Route(ws.DELETE("/container/{containerID}/configuration/{recordID}").To(s.DeleteConfiguration).
        Param(ws.PathParameter("recordID", "record ID")).
        Param(ws.PathParameter("containerID", "container ID")).
        Writes(Response{}))
    ws.Route(ws.DELETE("/container/{containerID}/{record}/{recordID}").To(s.DeleteContainerAttribute).
        Param(ws.PathParameter("record", "record type")).
        Param(ws.PathParameter("recordID", "record ID")).
        Param(ws.PathParameter("containerID", "container ID")).
        Writes(Response{}))
    ws.Route(ws.GET("/container/{containerID}/{record}/{recordID}").To(s.GetContainerAttribute).
        Param(ws.PathParameter("record", "record type")).
        Param(ws.PathParameter("recordID", "record ID")).
        Param(ws.PathParameter("containerID", "container ID")).
        Writes(Response{}))
    ws.Route(ws.GET("/images/{namespace}/{repo}").To(s.ImageList).
        Param(ws.PathParameter("namespace", "namespace")).
        Param(ws.PathParameter("repo", "repo ID")).
        Writes(Response{}))
    ws.Route(ws.GET("/images2/{namespace}/{repo}").To(s.ImageList2).
        Param(ws.PathParameter("namespace", "namespace")).
        Param(ws.PathParameter("repo", "repo ID")).
        Writes(Response{}))
    ws.Route(ws.POST("/history").To(s.AddHistory).
        Reads(app.History{}).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/histories").To(s.GetHistories).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/secret").To(s.CreateSecret2).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(struct {}{}).
        Writes(Response{}))
    ws.Route(ws.GET("/{app}/{region}/{env}/nodeSelector").To(s.getNodeSelector).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/nodeSelector").To(s.insertNodeSelector).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(pod.NodeSelector{}).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/nodeSelector/{id}").To(s.updateNodeSelector).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "record")).
        Reads(pod.NodeSelector{}).
        Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/nodeSelector/{id}").To(s.deleteNodeSelector).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "record")).
        Writes(Response{}))

    ws.Route(ws.GET("/{app}/{region}/{env}/hostAlias").To(s.getHostAlias).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/hostAlias").To(s.insertHostAlias).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Reads(pod.HostAlias{}).
        Writes(Response{}))
    ws.Route(ws.POST("/{app}/{region}/{env}/hostAlias/{id}").To(s.updateHostAlias).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "record")).
        Reads(pod.HostAlias{}).
        Writes(Response{}))
    ws.Route(ws.DELETE("/{app}/{region}/{env}/hostAlias/{id}").To(s.deleteHostAlias).
        Filter(s.getMeta).Filter(s.validateMeta).
        Param(ws.PathParameter("app", "app name")).
        Param(ws.PathParameter("region", "region")).
        Param(ws.PathParameter("env", "env")).
        Param(ws.PathParameter("id", "record")).
        Writes(Response{}))
    ws.Route(ws.POST("/clone").To(s.cloneApp).
        Reads(app.CloneInfo{}).
        Writes(Response{}))
    container.Add(ws)
}

func (s *Server) RegisterSwagger(container *restful.Container) {
    if s.config.Mode == "dev" {
        addr := fmt.Sprintf("http://%s:%d", s.config.Host, s.config.Port)
        config := swagger.Config{
            WebServices:     container.RegisteredWebServices(),
            WebServicesUrl:  addr,
            ApiPath:         "/docs.json",
            SwaggerPath:     "/docs/",
            SwaggerFilePath: "swagger-ui/dist"}
        swagger.RegisterSwaggerService(config, container)
    }
}

func (s *Server) Container(handler http.Handler) (string, *restful.Container) {
    addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
    container := restful.NewContainer()
    container.HandleWithFilter("/", s.cas.Handle(handler))
    cors := restful.CrossOriginResourceSharing{
        ExposeHeaders:  []string{"Content-Type", "Accept"},
        AllowedHeaders: []string{"Content-Type", "Accept", "Token", "User"},
        AllowedMethods: []string{"GET", "POST", "PUT", "OPTION", "DELETE"},
        CookiesAllowed: false,
        Container:      container}
    // container.Filter(container.OPTIONSFilter)
    container.Filter(cors.Filter)
    s.Register(container)
    s.RegisterSwagger(container)

    return addr, container
}
