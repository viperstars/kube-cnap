package http

import (
    "errors"
    "fmt"
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    errors2 "k8s.io/apimachinery/pkg/api/errors"
    "strconv"
)

const SUCCESSMESSAGE = "success"

var ERRORREQUEST = errors.New("bad request, record id is not null but body contains more than 1 record")

func (s *Server) ConvertStringToInt(str string) (int, error) {
    return strconv.Atoi(str)
}

func (s *Server) AddContainerCommand(request *restful.Request, response *restful.Response) {
    intID := 0
    command := new(container.Command)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&command)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            command.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    command.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(command, new(container.Command), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerEnvVars(request *restful.Request, response *restful.Response) {
    intID := 0
    envVar := new(container.EnvVar)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&envVar)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            envVar.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    envVar.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(envVar, new(container.EnvVar), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)

}

func (s *Server) AddContainerPorts(request *restful.Request, response *restful.Response) {
    intID := 0
    containerPort := new(container.Port)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&containerPort)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            containerPort.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    containerPort.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(containerPort, new(container.Port), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerConfiguration(request *restful.Request, response *restful.Response) {
    intID := 0
    configuration := new(container.Configuration)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&configuration)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            configuration.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        return
    }
    configuration.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(configuration, new(container.Configuration), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        err = s.k8sClients.UpdateConfigMap(configuration.Base, configuration.Path, configuration.Content)
        if err != nil && !errors2.IsNotFound(err){
            s.returnError(400, err, response)
            return
        }
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerVolumeMounts(request *restful.Request, response *restful.Response) {
    intID := 0
    volumeMount := new(container.VolumeMount)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&volumeMount)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            return
        } else {
            volumeMount.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    volumeMount.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(volumeMount, new(container.VolumeMount), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerResourceRequirement(request *restful.Request, response *restful.Response) {
    intID := 0
    resourceRequirement := new(container.ResourceRequirement)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&resourceRequirement)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            resourceRequirement.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    resourceRequirement.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(resourceRequirement, new(container.ResourceRequirement), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerLivenessCheck(request *restful.Request, response *restful.Response) {
    intID := 0
    livenessCheck := new(container.LivenessCheck)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&livenessCheck)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            livenessCheck.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    livenessCheck.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(livenessCheck, new(container.LivenessCheck), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerReadnessCheck(request *restful.Request, response *restful.Response) {
    intID := 0
    readnessCheck := new(container.ReadnessCheck)
    containerID := request.PathParameter("containerID")
    recordID := request.PathParameter("recordID")
    err := request.ReadEntity(&readnessCheck)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            readnessCheck.ID = intID
        }
    }
    intContainerID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    readnessCheck.ContainerID = intContainerID
    _, err = s.dbService.CreateOrUpdateRecord(readnessCheck, new(container.ReadnessCheck), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddContainerAttribute(request *restful.Request, response *restful.Response) {
    recordID := request.PathParameter("recordID")
    record := request.PathParameter("record")
    intID := 0
    r, err := s.GetRecordPointer(record)
    table, _ := s.GetRecord(record)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = request.ReadEntity(&r)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if recordID != "" {
        intID, err = s.ConvertStringToInt(recordID)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    _, err = s.dbService.UpdateContainerAttribute(intID, table, record)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeleteContainerAttribute(request *restful.Request, response *restful.Response) {
    fmt.Print("xxxxx")
    recordID := request.PathParameter("recordID")
    record := request.PathParameter("record")
    if recordID == "" {
        s.returnError(400, ERRORREQUEST, response)
        return
    }
    intID, err := s.ConvertStringToInt(recordID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    r, err := s.GetRecord(record)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    aff, err := s.dbService.DeleteContainerAttribute(intID, r)
    if err != nil {
        s.returnError(400, err, response)
        return
    } else if aff != 1 {
        s.returnError(400, err, response)
        return
    } else {
        s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
    }
}

func (s *Server) GetRecord(record string) (interface{}, error) {
    switch record {
    case "envVar":
        return container.EnvVar{}, nil
    case "readnessCheck":
        return container.ReadnessCheck{}, nil
    case "resourceRequirement":
        return container.ResourceRequirement{}, nil
    case "port":
        return container.Port{}, nil
    case "configuration":
        return container.Configuration{}, nil
    case "livenessCheck":
        return container.LivenessCheck{}, nil
    case "command":
        return container.Command{}, nil
    case "volumeMount":
        return container.VolumeMount{}, nil
    default:
        return nil, errors.New("error record type")
    }
}

func (s *Server) GetRecordPointer(record string) (interface{}, error) {
    switch record {
    case "envVar":
        return new(container.EnvVar), nil
    case "readnessCheck":
        return new(container.ReadnessCheck), nil
    case "resourceRequirement":
        return new(container.ResourceRequirement), nil
    case "port":
        return new(container.Port), nil
    case "configuration":
        return new(container.Configuration), nil
    case "livenessCheck":
        return new(container.LivenessCheck), nil
    case "command":
        return new(container.Command), nil
    case "volumeMount":
        return new(container.VolumeMount), nil
    default:
        return nil, errors.New("error record type")
    }
}

func (s *Server) GetContainerAttributes(request *restful.Request, response *restful.Response) {
    record := request.PathParameter("record")
    containerID := request.PathParameter("containerID")
    intID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    attributes, err := s.dbService.GetContainerAttributes(intID, record)

    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, attributes, response)
}

func (s *Server) GetContainerAttribute(request *restful.Request, response *restful.Response) {
    record := request.PathParameter("record")
    containerID := request.PathParameter("recordID")
    intID, err := s.ConvertStringToInt(containerID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    attributes, err := s.dbService.GetContainerAttribute(intID, record)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, attributes, response)
}
