package app

type App struct {
    ID               int      `json:"id" xorm:"pk autoincr 'id'"`
    Name             string   `json:"name" xorm:"varchar(512)"`
    // Repo             *Repo    `json:"repo" xorm:"text"`
    // DependedRepos    []*Repo  `json:"dependedRepos" xorm:"text"`
    Administrators   []string `json:"administrators" xorm:"text"`
    Developers       []string `json:"developers" xorm:"text"`
    Operators        []string `json:"operators" xorm:"text"`
    TestingEngineers []string `json:"testingEngineers" xorm:"text"`
    ImageNameSpace   string   `json:"imageNamespace"`
    ImageRepo        string   `json:"imageRepo"`
    Description      string   `json:"description"`
    Deleted          bool     `json:"deleted"`
    Region           string   `json:"region" xorm:"-"`
    Env              string   `json:"env" xorm:"-"`
}

func (app *App) returnUsers(name string) []string {
    switch name {
    case "op":
        return app.Operators
    case "admin":
        return app.Administrators
    case "qa":
        return app.TestingEngineers
    case "dev":
        return app.Developers
    default:
        return app.Administrators
    }
}

func (app *App) ReturnUsers(roles []string) []string {
    all := make([]string, 0)
    for _, v := range roles {
        all = append(all, app.returnUsers(v)...)
    }
    return all
}

type History struct {
    ID         int    `json:"id" xorm:"pk autoincr 'id'"`
    App        string `json:"app"`
    Region     string `json:"region"`
    Env        string `json:"env"`
    Node       string `json:"node"`
    ImageName  string `json:"imageName"`
    ImageId    string `json:"imageId"`
    Tag        string `json:"tag"`
    CreateDate string `json:"createDate"`
    Deleted    bool   `json:"deleted"`
}
