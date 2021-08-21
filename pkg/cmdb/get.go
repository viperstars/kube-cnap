package cmdb

import (
    jsoniter "github.com/json-iterator/go"
    "go.uber.org/zap"
    "io/ioutil"
    "net/http"
    "strings"
)

type User struct {
    Username string  `json:"username"`
    CName string  `json:"first_name"`
}

type UserResponse struct {
    Count    int                 `json:"count"`
    Next     jsoniter.RawMessage `json:"next"`
    Previous jsoniter.RawMessage `json:"previous"`
    Result   []*User             `json:"results"`
}

type DataGetter struct {
    logger *zap.Logger
}

func (d *DataGetter) get(url string) (*UserResponse, error) {
    users := new(UserResponse)
    resp, err := http.Get(url)
    if err != nil {
        d.logger.Error("get users error", zap.Error(err), zap.String("url", url))
        return users, err
    } else {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            d.logger.Error("read resp body error", zap.Error(err))
            return users, err
        } else {
            err := jsoniter.Unmarshal(body, &users)
            return users, err
        }
    }
}

func (d *DataGetter) GetAll(url string) (*UserResponse, error) {
    users, err := d.get(url)
    if err != nil {
        return users, err
    } else {
        for string(users.Next) != "null" {
            url := strings.Trim(string(users.Next), "\"")
            s, err := d.get(url)
            if err != nil {
                d.logger.Error("pull rds error", zap.Error(err), zap.String("url", url))
                return users, err
            } else {
                users.Next = s.Next
                users.Result = append(users.Result, s.Result...)
            }
        }
        return users, nil
    }
}

func NewUserGetter(logger *zap.Logger) *DataGetter {
    return &DataGetter{
        logger: logger,
    }
}