package app

type Region struct {
    ID      int    `json:"id" xorm:"pk autoincr 'id'"`
    Name    string `json:"name" xorm:"varchar(256)"`
    CName   string `json:"cname" xorm:"varchar(256)"`
    Deleted bool `json:"deleted"`
    Envs    []*Env `json:"envs" xorm:"-"`
    Env    *Env `json:"env" xorm:"-"`
}

type Env struct {
    ID        int    `json:"id" xorm:"pk autoincr 'id'"`
    RegionID  int    `json:"regionID" xorm:"'region_id'"`
    CName     string `json:"cname" xorm:"varchar(256)"`
    Name      string `json:"name"`
    Config    string `json:"config" xorm:"text"`
    Namespace string `json:"namespace"`
    Deleted   bool   `json:"deleted"`
}
