package main

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
    "github.com/viperstars/bce-golang/bce_client"
    "github.com/viperstars/kube-cnap/conf"
    "github.com/viperstars/kube-cnap/log"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
    "github.com/viperstars/kube-cnap/pkg/client"
    "github.com/viperstars/kube-cnap/pkg/database"
    "github.com/viperstars/kube-cnap/pkg/image"
    "github.com/viperstars/kube-cnap/pkg/server"
    _http "github.com/viperstars/kube-cnap/pkg/server/http"
    "net/http"
)

func startServer() {
    c := bce_client.NewBCEClient(
        conf.Config.BceConfig.Ak,
        conf.Config.BceConfig.Sk,
        "bj",
        "cce",
    )
    logger := log.NewLogger(conf.Config.Logging.Path)
    db := database.NewDBClient(conf.Config.DB.Dsn, conf.Config.Server.Mode , logger)
    clients, err := client.NewK8sClients(logger, db)
    if err != nil {
        fmt.Println(err)
        panic("load from db error")
    }
    casClient := server.NewCASClient(conf.Config.Cas.Server)
    reg, err := image.NewRegistryClient(conf.Config.RegistryInfo.Url, conf.Config.RegistryInfo.Username,
        conf.Config.RegistryInfo.Password)
    if err != nil {
        fmt.Println(err)
        //panic("can not access registry")
    }
    httpServer := _http.NewServer(clients, logger, db, conf.Config.Server, casClient, c, reg)
    filter := server.NewCASFilter(conf.Config.Cas.CasCallback)
    addr, container2 := httpServer.Container(filter)
    http.Handle("/", container2)
    http.Handle("/api/sockjs/", _http.CreateAttachHandler("/api/sockjs"))
    err = http.ListenAndServe(addr, nil)
    if err != nil {
        panic(err)
    }
}
// docker create -v /var/lib/mysql --name mysql-data mysql
// docker run -d --volumes-from mysql-data -e MYSQL_ROOT_PASSWORD=mysql -p 3306:3306 mysql

func main() {
    syncTable()
    startServer()
}

func syncTable() {
    engine, err := xorm.NewEngine("mysql", conf.Config.DB.Dsn)
    err = engine.Sync2(new(container.EnvVar), new(container.ResourceRequirement), new(container.Configuration),
        new(pod.ContainerGroup), new(container.Base), new(container.VolumeMount),
        new(container.LivenessCheck), new(container.Port), new(container.ReadnessCheck),
        new(container.Command), new(pod.CFSVolume), new(pod.EmptyVolume), new(pod.BOSVolume), new(app.App),
        new(app.Region), new(app.Env), new(ingress.IngressRule), new(container.LivenessCheckHttp), new(app.History),
        new(pod.NodeSelector), new(pod.HostAlias))
    fmt.Println(err)
    engine.ShowSQL(true)
}