package app

type CloneInfo struct {
    App1  string       `json:"app1"`
    App2  string       `json:"app2"`
    Zones []*CloneZone `json:"zones"`
    IDs   []*CloneID
}

type CloneID struct {
    Old  int
    New  int
    Zone *CloneZone
}

type CloneZone struct {
    Old *Zone `json:"old"`
    New *Zone `json:"new"`
}

type Zone struct {
    Region string `json:"region"`
    Env    string `json:"env"`
}
