package client

import (
    "encoding/json"
    "fmt"
    "github.com/viperstars/kube-cnap/log"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "testing"
)

var client *Clients
var base = basemeta.BaseMeta{
    App:    "k8s-app",
    Region: "bj",
    Env:    "test",
}

var kubeConfig = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUR2RENDQXFTZ0F3SUJBZ0lVV0pjYnFBWFFTcmQzLytCR1pZbS9tVG02T2dBd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pERUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFXcHBibWN4RURBT0JnTlZCQWNUQjBKbAphV3BwYm1jeEREQUtCZ05WQkFvVEEyczRjekVPTUF3R0ExVUVDeE1GYVhGcGVXa3hFekFSQmdOVkJBTVRDbXQxClltVnlibVYwWlhNd0hoY05NVGt3TXpFNE1ESTFOREF3V2hjTk1qUXdNekUyTURJMU5EQXdXakJrTVFzd0NRWUQKVlFRR0V3SkRUakVRTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVNTUFvRwpBMVVFQ2hNRGF6aHpNUTR3REFZRFZRUUxFd1ZwY1dsNWFURVRNQkVHQTFVRUF4TUthM1ZpWlhKdVpYUmxjekNDCkFTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBUHZIaWdBTnZrUXViKysrM1lhejVkL0gKQ2tLZ296amFXalpxM3hyZmRzR0NpVmRUeGRpMmROZ1dkWFkyTTdvbkdaZHpGcmdlZWcwaHgvR3JoaVJRcldLcgpYaVJsN3BvWHpqQVRTZG1pSnFSQU91bEh3NTR6TGU4UVUwRVYrNGhZU09uaTEzVEdDVE4zcy95Q0VsTCtkVG53ClBNeWtBTUl3Zm9QbGJzdUxXU21nU21nNUxtcnlORmFHdHlaOGswZDV2cUVsanVrbTM2Q3ZSRWU2QlBXeW9Eb04KSytFcHl6LzNVSG1BR3dpVytGQ0tEUlBUS0h2UmJGWXUxQ1dBNG9lN0d0RVVCU3htMERDTDdneUwyV0lLampDWApSdFpNT21kc1BCOFpMQ2JPUlhubjhBeDBPMEp2RlRTam1OQTBnUFM0aTdJeUFIS3cvWXlvWjNqVElTVDY0VWNDCkF3RUFBYU5tTUdRd0RnWURWUjBQQVFIL0JBUURBZ0VHTUJJR0ExVWRFd0VCL3dRSU1BWUJBZjhDQVFJd0hRWUQKVlIwT0JCWUVGSlozSHV3dHJibS9CVUZKanovdDlWS1IzWGxCTUI4R0ExVWRJd1FZTUJhQUZKWjNIdXd0cmJtLwpCVUZKanovdDlWS1IzWGxCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBUWJVd3JJYmU0T2dtcTZQeldNclI5CklaT1N1Ykl3T2lkc0FLRkxpOWJtZ3dGYktubFk2M0pNSEtjeEtDMllQU3Boem9tZWFRS1Yvc0ZmdVlwOU1KS1AKc21aSW43R2VGaDQ2K0QrRU56ZHZHRXFjbEFkVC9BZWhUa0k3dG1OS3ZFM0Nad3BiLzZqcGlvMEIwTGIyWlRUegpRQWtpVlNvamNjZWtZT2dSVkY4UmYrRUJqUWpRd3plOUI3ZW40SFhGVVdUTFlORE90d01CVzVHQmJvOW5tQkNrCnFnRGpyWkhubXpuTy82VnU3YVlidWhzL3lXRHdCaVJkeEt5UjRiWi8yZFRUdUhqS25vc05BakhZdFRvZ1htT0IKbVlNUVBSZlZEc1gxQTFCUWRvSU0vMWZhTk9wWFdOa21FL1d3ZE9VazJmciswTERvQUtwVStQK3N3QmFEb0NNeQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://10.193.6.34:8443
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: admin
  name: kubernetes
current-context: kubernetes
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUQyekNDQXNPZ0F3SUJBZ0lVTFJseTg4elE4Wm4yR25pbXladmpDWE1ZUWFzd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pERUxNQWtHQTFVRUJoTUNRMDR4RURBT0JnTlZCQWdUQjBKbGFXcHBibWN4RURBT0JnTlZCQWNUQjBKbAphV3BwYm1jeEREQUtCZ05WQkFvVEEyczRjekVPTUF3R0ExVUVDeE1GYVhGcGVXa3hFekFSQmdOVkJBTVRDbXQxClltVnlibVYwWlhNd0hoY05NVGt3TXpFNE1ESTFPREF3V2hjTk1qa3dNekUxTURJMU9EQXdXakJxTVFzd0NRWUQKVlFRR0V3SkRUakVRTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVYTUJVRwpBMVVFQ2hNT2MzbHpkR1Z0T20xaGMzUmxjbk14RGpBTUJnTlZCQXNUQldseGFYbHBNUTR3REFZRFZRUURFd1ZoClpHMXBiakNDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFKYzV0WkxJTHFXYVhhc3AKV3NrWnpXcjg3SVEzQlhQSnEreWoraXdQQzk2NEhkKzh0WUtZSXRBWFYyVndpbzdtV094TXRubis4bzNtcUg1QQp6MEpvTXFqWFBFRzVjblZJZFlLTGVWQzlzVnpIOWkxWVZ1dWlFM1JuRFBiSzltYU9jbEl2QVNQdy9FNDdrcW1MCitkVS94cC9FazJKUEJSekdyVWE2WVcvcmRvc0RId09IWWdlcmV1WGhDK0ExZzVOaFNpR2lsUldSQ0FvVzN2dEIKN1RRTHhLOTFFNW5MZ012NXNxNjlxNHIwQS9td3BFK1dMdnZCUDUrUnhaVTM3TEx1aTM2cFloSnp3cUtuNmZTZQpPTHN1Sy9QYlZNUEdmQWNYS0hiK2JKb3BGbHFQdnQ3U29JMVpyQ3VIUHdOK3plWG5obHpPTVZsd1Z6SzREUkwyCndiWi95M3NDQXdFQUFhTi9NSDB3RGdZRFZSMFBBUUgvQkFRREFnV2dNQjBHQTFVZEpRUVdNQlFHQ0NzR0FRVUYKQndNQkJnZ3JCZ0VGQlFjREFqQU1CZ05WSFJNQkFmOEVBakFBTUIwR0ExVWREZ1FXQkJSZ21WQzRoRC9QS1hPegplaHJaeHQzYVplQ056REFmQmdOVkhTTUVHREFXZ0JTV2R4N3NMYTI1dndWQlNZOC83ZlZTa2QxNVFUQU5CZ2txCmhraUc5dzBCQVFzRkFBT0NBUUVBSzdvWHcrYkgrVXdpSmhaQ0RsSWxZK2lqMjd1WVBwRmo0Z0x3N0t4MDEwY0gKWlloTGRZZG1sT1dVZkhKeEpReGtQREo0OHhISmxNQ1hNRGZCZ1JlZllkaUFna0MxN29abksxT0NPVXhyRTJ6WQovUTM1U0s2Ri9mSWtkR3lpeTRNbkp4UUhCUGlQcWw1TlBOYWJrNG9QUEQ5T3h5MTZMQ1Exb1dFaEIxdWR4TTI3CjA0ek9iVzM4U2xIdkZic2R4ckF3Q3dOOEc5QU0rZStHdytKZEpoaENsOFNmZUl1bHhCbzNpTmhoSThlcElNOVUKcWZ1VEFFaXcyZGpCOHM2NFI0a25OZnMrUGFnVitwbjNMb0ZydGZ3dEdvQ09GUk9ydDZFOHJrMEgrMExoSDFrOQpZMmxLcTU5SXVDSmV1SWF1Y21UNTV5c0tZUStzS1ZHVTUxWTBURUpWeGc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcGdJQkFBS0NBUUVBbHptMWtzZ3VwWnBkcXlsYXlSbk5hdnpzaERjRmM4bXI3S1A2TEE4TDNyZ2QzN3kxCmdwZ2kwQmRYWlhDS2p1Wlk3RXkyZWY3eWplYW9ma0RQUW1neXFOYzhRYmx5ZFVoMWdvdDVVTDJ4WE1mMkxWaFcKNjZJVGRHY005c3IyWm81eVVpOEJJL0Q4VGp1U3FZdjUxVC9HbjhTVFlrOEZITWF0UnJwaGIrdDJpd01mQTRkaQpCNnQ2NWVFTDREV0RrMkZLSWFLVkZaRUlDaGJlKzBIdE5BdkVyM1VUbWN1QXkvbXlycjJyaXZRRCtiQ2tUNVl1Cis4RS9uNUhGbFRmc3N1NkxmcWxpRW5QQ29xZnA5SjQ0dXk0cjg5dFV3OFo4Qnhjb2R2NXNtaWtXV28rKzN0S2cKalZtc0s0Yy9BMzdONWVlR1hNNHhXWEJYTXJnTkV2YkJ0bi9MZXdJREFRQUJBb0lCQVFDRFduZ1J4OWxrdEtHWgowM0RzN29HVDVaOFc5S2ZDdkRDZWVvd0ppd1EzYjF0YmhKRndudTdXS3dBWnQxaFM2VmZoNEh3N21TeGIvemJwCmV5Zkx1YkFQSWUxUDlXR3E1OFpTSHczQUNSU3V6MjFRVThRa0pnS2FBQXl0clB1N2R3MXJ2ekpSWXJmMHlDQUwKTHU5UndIU3BQOWo0OGtReGk0emt1MjE1Qm1CUU80N3pPWnJoNm9ZaHRQUG5lTVVlckNKWXpHcFhITmdVbkRVdwpqbGhjMVdDTXJTS09LNFJrMFlPM2J4U0lzREhXcFVTZ1ZOTjAxblU5aklBVEZaVk4ybnpOQWFENmJBdkFueThoCmdzNmZUM2NHQWQ4UWt4Y0l1VXVjMmF4NG9SSzE0emZQWDEyV21md0hpVlgveE1jQ1dIMytZWVhiNnFBVThCWlAKT1kvRzVKOTVBb0dCQU1lYkk0eXpFZ1hWdFh4ZzBCdE9KK3BwZTJ5Sm9kemJCRVdralZnVFBIUm1OMTJ2cEhNSApIZFBveE81cjBhVm1TQ2tLWFRmVXlEMXhDMTROYXl5M0EvODk5bEJqWGxOdDRMWUlIMnNtUlI0NnEwN0Y0YjBUCmNXZ2xIMzkwa0EvR3ljcUNlQ1lUa2ZSeCtUTDV6a3pRRi91Sm5PaldvcUJNck1TWElQa281U0QxQW9HQkFNSHoKVzdjcWI5Y3BzT3FJRE5NS0pnSlUrVmRCajBRbFByNXRmcGxpMnljakVQTVF3cDNURkVuN0kwY0RuaERIeVRFQQplUDRFTWxlZmV5RzdiYmJxZEdqTyttNEtwQmkzODM4UEdlUVpMaFZTU2ZRZ2lmSWdKM3UvSm5CQ2dNSHYveDFYCk1KNkJ5WGYrVWE4Vi9GVGtKWC9MUm5Sby9OOGNDNGVFMFFRMHFyU3ZBb0dCQUliTEpIR2loOXc2Mm9rNDA3QnMKMGhYQnorQ1cvU0NwSXJScEVDNVhKeTh2eTluUGdBMVIwL25EcWlHYjNBS0hGTm5xTHROQ05Vc1FxTzJGd0VkOAovQTBFNmU2VmZDQjVCaFBIWG5nOGF0YWtKZ1ZYS2o5Ri93S21keVBhTW1NRkNrWmdYd1RQbUhQcjk2NU45ZHYzCmR3cWRmc0hhR0E2S1dPMlZaV1g5RU9aMUFvR0JBS0M4dVNRQ0RaSjZRTjcrUmZLWkZJc1dOVmIxUkhDcmxXWm8KaEdWR29tMjdDQThKc3VEdDBJREhtNkw5QW9EUnNwSGozR0pZeEFnT2FoTzRxK0xPU0ErY2lidXRJZlpDYlpDOQp5UzFiR1BBZXRKK1lYL3JFWHpTVlpKdmc0YWpZNThzL09WSUVLaDVDTFJ3MzBsbmdncHQ0c2psRDBWNXVkYmVvCmdUbEZGTHlOQW9HQkFNRXV4YnBQaWl6QWJSYit1L0VGQTRESStJUGdTZytTRUtaaXBjMEdDWmJaam92TmEzd3cKVXdFUmlGZDJmbGtMZjEyTnhhZktQZTRPK2RNRkhpY0VFcDFlK2NoV3p3L0tJam9tN0E4RUpZZFAvUjQrZHRhRwo0ZHlPVmlQTEpQZ0E4QiswQ2xITzRWS3FJS0xrcEZ3aDROeGF1eEc1cHY3a29NeTJydzVyZ01UcAotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
`

func init() {
    client = new(Clients)
    clients := make(map[string]map[string]*Client)
    client.ClientPool = clients
    config, _ := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
    client.ClientPool["bj"] = make(map[string]*Client)
    clientForTest := &Client{
        client: kubernetes.NewForConfigOrDie(config),
        config: config,
    }
    client.ClientPool["bj"]["test"] = clientForTest
    client.logger = log.NewLogger("./test.log")
}

func TestClients_GetDeployments(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    ls, _ := ToLabelSelector(label)
    c := client.ClientPool["bj"]["test"]
    fmt.Println(ls.String())
    deploymentlist, err := c.client.AppsV1().Deployments("kube-system").List(metav1.ListOptions{
        LabelSelector: ls.String(),
    })
    if err != nil {
        t.Error(err)
    }
    for _, dep := range deploymentlist.Items {
        fmt.Println(dep.Name)
        d, _ := json.MarshalIndent(dep.Spec.Template, "", "    ")
        fmt.Println(string(d))
    }
}

func TestClients_GetNodes(t *testing.T) {
    c := client.ClientPool["bj"]["test"]
    deploymentlist, err := c.client.CoreV1().Nodes().List(metav1.ListOptions{})
    if err != nil {
        t.Error(err)
    }
    for _, dep := range deploymentlist.Items {
        fmt.Println(dep.Name)
        d, _ := json.MarshalIndent(dep, "", "    ")
        fmt.Println(string(d))
    }
}

func TestClients_GetPods(t *testing.T) {
    c := client.ClientPool["bj"]["test"]
    deploymentlist, err := c.client.AppsV1().ReplicaSets("kube-system").List(metav1.ListOptions{})
    if err != nil {
        t.Error(err)
    }
    for _, dep := range deploymentlist.Items {
        fmt.Println(dep.Name)
        d, _ := json.MarshalIndent(dep, "", "    ")
        fmt.Println(string(d))
    }
}

func TestClients_GetReplicaSets(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    ls, _ := ToLabelSelector(label)
    c := client.ClientPool["bj"]["test"]
    fmt.Println(ls.String())
    deploymentlist, err := c.client.AppsV1().ReplicaSets("kube-system").List(metav1.ListOptions{
        LabelSelector: ls.String(),
    })
    if err != nil {
        t.Error(err)
    }
    for _, dep := range deploymentlist.Items {
        fmt.Println(dep.Name)
        d, _ := json.MarshalIndent(dep, "", "    ")
        fmt.Println(string(d))
    }
}

func TestClients_GetDeploymentPods(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    deps, err := client.GetDeployments(base, label)
    fmt.Println(err)
    resp, err := client.GetPodListForDeployment(base, deps)
    if err != nil {
        t.Error(err)
    }
    js, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Println(string(js))
}

func TestClients_GetEventsForObject(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    deps, _ := client.GetEventsForObject(base, "coredns-7c54bb4569-ts2lz")
    js, _ := json.MarshalIndent(deps, "", "  ")
    fmt.Println(string(js))
}

func TestClients_GetService(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    services, _ := client.GetServices(base, label)
    js, _ := json.MarshalIndent(services, "", "  ")
    fmt.Println(string(js))
}

func TestClients_GetMiniDeployment(t *testing.T) {
    label := make(map[string]string)
    label["k8s-app"] = "kube-dns"
    deps, err := client.GetDeployments(base, label)
    fmt.Println(err)
    resp, err := client.GetMiniDeployments(base, deps)
    if err != nil {
        t.Error(err)
    }
    js, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Println(string(js))
}

func TestClients_AddLabelToNode(t *testing.T) {
    labels := map[string]string{"k8s-master": "true"}
    err := client.UpdateLabelsToNode(base, "bj-qa-sre-k8s-1", labels)
    fmt.Println(err)
}
