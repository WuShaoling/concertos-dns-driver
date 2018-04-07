package restapi

import (
	"log"
	"net/http"
	"github.com/emicklei/go-restful"
	"sync"
	"github.com/concertos-dns/executor"
)

func (rest *RestApi) Start() {
	ws := new(restful.WebService)
	ws.Path("/domain").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/").To(rest.addDomain))
	ws.Route(ws.DELETE("/{domain}").To(rest.deleteDmain))
	ws.Route(ws.GET("/").To(rest.getAll))

	restful.DefaultContainer.Add(ws)

	log.Printf("Start rest api, listening on localhost:40001")
	log.Fatal(http.ListenAndServe(":40001", nil))
}

func (rest *RestApi) addDomain(request *restful.Request, response *restful.Response) {
	var doip = new(DomainIP)
	log.Println(request.ReadEntity(doip))
	domainip := doip.IP + " " + doip.Domain
	log.Println(domainip)

	if err := rest.executor.AddRecord(domainip); nil != err {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteHeader(http.StatusOK)
	}
}

func (rest *RestApi) deleteDmain(request *restful.Request, response *restful.Response) {
	if err := rest.executor.DeleteRecord(request.PathParameter("domain")); nil != err {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteHeader(http.StatusOK)
	}
}

func (rest *RestApi) getAll(request *restful.Request, response *restful.Response) {
	if lines, err := rest.executor.ReadLines(); nil != err {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteEntity(lines)
	}
}

var restApi *RestApi
var once sync.Once

type DomainIP struct {
	IP     string
	Domain string
}

type RestApi struct {
	executor *executor.Executor
}

func GetRestApi() *RestApi {
	once.Do(func() {
		restApi = &RestApi{
			executor: executor.GetExecutor(),
		}
	})
	return restApi
}
