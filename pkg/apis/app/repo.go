package app

type Repo struct {
    Name                 string `json:"name"`
    URL                  string `json:"url"`
    PublishMode          string `json:"publishMode"`
    DestinationDirectory string `json:"destinationDirectory"`
}
