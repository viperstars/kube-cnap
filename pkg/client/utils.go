package client

import (
    apps "k8s.io/api/apps/v1"
    batch "k8s.io/api/batch/v1"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/equality"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/labels"
)

func ToLabelSelector(selector map[string]string) (labels.Selector, error) {
    labelSelector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{MatchLabels: selector})
    if err != nil {
        return nil, err
    }
    return labelSelector, nil
}

func FilterDeploymentPodsByOwnerReference(deployment apps.Deployment, allRS []apps.ReplicaSet,
    allPods []v1.Pod) []v1.Pod {
    var matchingPods []v1.Pod
    for _, rs := range allRS {
        if metav1.IsControlledBy(&rs, &deployment) {
            matchingPods = append(matchingPods, FilterPodsByControllerRef(&rs, allPods)...)
        }
    }

    return matchingPods
}

func FilterPodsByControllerRef(owner metav1.Object, allPods []v1.Pod) []v1.Pod {
    var matchingPods []v1.Pod
    for _, pod := range allPods {
        if metav1.IsControlledBy(&pod, owner) {
            matchingPods = append(matchingPods, pod)
        }
    }
    return matchingPods
}

func FilterPodsForJob(job batch.Job, pods []v1.Pod) []v1.Pod {
    result := make([]v1.Pod, 0)
    for _, pod := range pods {
        if pod.Namespace == job.Namespace {
            selectorMatch := true
            for key, value := range job.Spec.Selector.MatchLabels {
                if pod.Labels[key] != value {
                    selectorMatch = false
                    break
                }
            }
            if selectorMatch {
                result = append(result, pod)
            }
        }
    }

    return result
}

func GetContainerImages(podTemplate *v1.PodSpec) []string {
    var containerImages []string
    for _, container := range podTemplate.Containers {
        containerImages = append(containerImages, container.Image)
    }
    return containerImages
}

func GetInitContainerImages(podTemplate *v1.PodSpec) []string {
    var initContainerImages []string
    for _, initContainer := range podTemplate.InitContainers {
        initContainerImages = append(initContainerImages, initContainer.Image)
    }
    return initContainerImages
}

func GetContainerNames(podTemplate *v1.PodSpec) []string {
    var containerNames []string
    for _, container := range podTemplate.Containers {
        containerNames = append(containerNames, container.Name)
    }
    return containerNames
}

func GetInitContainerNames(podTemplate *v1.PodSpec) []string {
    var initContainerNames []string
    for _, initContainer := range podTemplate.InitContainers {
        initContainerNames = append(initContainerNames, initContainer.Name)
    }
    return initContainerNames
}

func GetNonduplicateContainerImages(podList []v1.Pod) []string {
    var containerImages []string
    for _, pod := range podList {
        for _, container := range pod.Spec.Containers {
            if noStringInSlice(container.Image, containerImages) {
                containerImages = append(containerImages, container.Image)
            }
        }
    }
    return containerImages
}

func GetNonduplicateInitContainerImages(podList []v1.Pod) []string {
    var initContainerImages []string
    for _, pod := range podList {
        for _, initContainer := range pod.Spec.InitContainers {
            if noStringInSlice(initContainer.Image, initContainerImages) {
                initContainerImages = append(initContainerImages, initContainer.Image)
            }
        }
    }
    return initContainerImages
}

func GetNonduplicateContainerNames(podList []v1.Pod) []string {
    var containerNames []string
    for _, pod := range podList {
        for _, container := range pod.Spec.Containers {
            if noStringInSlice(container.Name, containerNames) {
                containerNames = append(containerNames, container.Name)
            }
        }
    }
    return containerNames
}

func GetNonduplicateInitContainerNames(podList []v1.Pod) []string {
    var initContainerNames []string
    for _, pod := range podList {
        for _, initContainer := range pod.Spec.InitContainers {
            if noStringInSlice(initContainer.Name, initContainerNames) {
                initContainerNames = append(initContainerNames, initContainer.Name)
            }
        }
    }
    return initContainerNames
}

func noStringInSlice(str string, array []string) bool {
    for _, alreadystr := range array {
        if alreadystr == str {
            return false
        }
    }
    return true
}

func EqualIgnoreHash(template1, template2 v1.PodTemplateSpec) bool {
    labels1, labels2 := template1.Labels, template2.Labels
    if len(labels1) > len(labels2) {
        labels1, labels2 = labels2, labels1
    }
    for k, v := range labels2 {
        if labels1[k] != v && k != apps.DefaultDeploymentUniqueLabelKey {
            return false
        }
    }
    template1.Labels, template2.Labels = nil, nil
    return equality.Semantic.DeepEqual(template1, template2)
}
