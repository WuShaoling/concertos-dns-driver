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

	ws.Route(ws.PUT("/add").To(rest.addDomain))
	ws.Route(ws.DELETE("/delete/{domain}").To(rest.deleteDmain))

	restful.DefaultContainer.Add(ws)

	log.Printf("Start rest api, listening on localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func (rest *RestApi) addDomain(request *restful.Request, response *restful.Response) {
	var doip = new(DomainIP)
	request.ReadEntity(doip)
	domainip := doip.IP + " " + doip.Domain
	log.Println(domainip)

	if err := rest.executor.AddRecord(domainip); nil != err {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteHeader(http.StatusOK)
	}
}

func (rest *RestApi) deleteDmain(request *restful.Request, response *restful.Response) {
	if err := rest.executor.DeleteRecord(request.PathParameter("domain-ip")); nil != err {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteHeader(http.StatusOK)
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
