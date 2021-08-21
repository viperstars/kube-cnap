package server

import (
    "fmt"
    "gopkg.in/cas.v1"
    "net/http"
    "net/url"
)

type CASFilter struct {
    CallbackUrl string
}

func NewCASFilter(callback string) *CASFilter {
    return &CASFilter{
        CallbackUrl: callback,
    }
}

func NewCASClient(casUrl string) *cas.Client {
    parsed, _ := url.Parse(casUrl)
    client := cas.NewClient(&cas.Options{
        URL: parsed,
    })
    return client
}

func (c *CASFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if !cas.IsAuthenticated(r) {
        cas.RedirectToLogin(w, r)
        return
    }
    if r.URL.Path == "/logout" {
        cas.RedirectToLogout(w, r)
        return
    } else if r.URL.Path == "/login" {
        ticket, ok := r.URL.Query()["ticket"]
        if ok && len(ticket) == 1 {
            username := cas.Username(r)
            param := url.Values{}
            param.Add("username", username)
            param.Add("token", ticket[0])
            location := fmt.Sprintf(`%s?%s`, c.CallbackUrl, param.Encode())
            http.Redirect(w, r, location, 302)
        } else {
            w.WriteHeader(400)
        }
    } else {
    }
}
