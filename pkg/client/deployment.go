package client

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "go.uber.org/zap"
    appsv1 "k8s.io/api/apps/v1"
)

func (c *Clients) GetMiniDeployments(base basemeta.BaseMeta, dpList *appsv1.DeploymentList) ([]*MiniDeployment,
    error) {
    miniDeploymentList := make([]*MiniDeployment, 0)
    for _, d := range dpList.Items {
        miniDeployment, _ := c.GetMiniDeploymentDetail(base, d, false)
        miniDeploymentList = append(miniDeploymentList, miniDeployment)
    }
    return miniDeploymentList, nil
}

func (c *Clients) GetMiniDeploymentDetail(base basemeta.BaseMeta, deployment appsv1.Deployment,
    detail bool) (*MiniDeployment, error) {
    miniDeployment := ToMiniDeployment(deployment)
    if detail {
        lbs, _ := ToLabelSelector(deployment.Spec.Selector.MatchLabels)
        rsList, _ := c.GetReplicaSets(base, lbs.String())
        newest, others := FilterReplicaSet(deployment, rsList.Items)
        newestRs, _ := c.GetMiniPodsForReplicaSet(base, newest)
        otherRs, _ := c.GetMiniPodsForReplicaSets(base, others)
        miniDeployment.NewestRs = newestRs
        miniDeployment.OtherRs = otherRs
    }
    return miniDeployment, nil
}

func (c *Clients) GetPodsForReplicaSetsList(base basemeta.BaseMeta, newest appsv1.ReplicaSet,
    replicaSetList []appsv1.ReplicaSet) ([]*SimpleReplicaSet, error) {
    replicaSets := make([]*SimpleReplicaSet, 0)
    for _, replicaSet := range replicaSetList {
        simple, err := c.GetPodsForReplicaSet(base, replicaSet)
        if err != nil {

        }
        replicaSets = append(replicaSets, simple)
    }
    return replicaSets, nil
}

func (c *Clients) GetMiniPodsForReplicaSet(base basemeta.BaseMeta, replicaSet appsv1.ReplicaSet) (*MiniReplicaSet,
    error) {
    var miniReplicaSet *MiniReplicaSet
    var image string
    lbs, _ := ToLabelSelector(replicaSet.Spec.Selector.MatchLabels)
    pod, err := c.GetPodsWithString(base, lbs.String())
    if err != nil {
        c.logger.Error("get pod with string error", zap.Error(err), zap.Any("meta", base), zap.String("label",
            lbs.String()))
        return miniReplicaSet, err
    }
    pods := FilterPodsByControllerRef(&replicaSet, pod.Items)
    if len(pods) > 0 && len(pods[0].Spec.Containers) > 0 {
        image = pods[0].Spec.Containers[0].Image // todo 判断 length
    }
    miniPodList := ToMiniPodList(pods)
    podInfo := GetPodInfo(*replicaSet.Spec.Replicas, replicaSet.Status.Replicas, pods)
    miniReplicaSet = &MiniReplicaSet{
        Name:    replicaSet.Name,
        Image:   image,
        PodList: miniPodList,
        PodInfo: &podInfo,
    }
    return miniReplicaSet, nil
}

func (c *Clients) GetMiniPodsForReplicaSets(base basemeta.BaseMeta, replicaSet []appsv1.ReplicaSet) ([]*MiniReplicaSet,
    error) {
    replicaSetList := make([]*MiniReplicaSet, 0)
    for _, r := range replicaSet {
        miniReplicaSet, err := c.GetMiniPodsForReplicaSet(base, r)
        if err != nil {
            return replicaSetList, err
        }
        if miniReplicaSet.PodInfo.Current > 0 {
            replicaSetList = append(replicaSetList, miniReplicaSet)
        }
    }
    return replicaSetList, nil
}
