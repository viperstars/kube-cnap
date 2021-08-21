package http

import (
    "github.com/emicklei/go-restful"
)

type Response struct {
    Code    int         `json:"code,omitempty"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data"`
}

func (s *Server) returnError(status int, err error, response *restful.Response) {
    resp := &Response{
        Code:    status,
        Message: err.Error(),
    }
    response.WriteHeaderAndJson(status, resp, restful.MIME_JSON)
}

func (s *Server) returnResponseWithData(status int, data interface{}, response *restful.Response) {
    resp := &Response{
        Code: status,
        Data: data,
    }
    response.WriteHeaderAndJson(status, resp, restful.MIME_JSON)
}

func (s *Server) returnResponseWithMessage(status int, message string, response *restful.Response) {
    resp := &Response{
        Code:    status,
        Message: message,
    }
    response.WriteHeaderAndJson(status, resp, restful.MIME_JSON)
}

func (s *Server) stringInSlice(str string, slice []string) bool {
    for _, s := range slice {
        if str == s {
            return true
        }
    }
    return false
}

func (s *Server) Users(request *restful.Request, response *restful.Response) {
    users, err := s.cmdb.GetAll("http://cmdb.qiyi.so/api/user/?is_active=true&page_size=500")
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, users.Result, response)
}
