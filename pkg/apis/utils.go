package apis

import (
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
