package container

type Command struct {
    ID          int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int    `json:"containerID" xorm:"int 'container_id'"`
    Command     string `json:"command" xorm:"varchar(512)"`
    Args        string `json:"args" xorm:"varchar(512)"`
    Deleted     bool   `json:"deleted" xorm:"bool default 0"`
}
