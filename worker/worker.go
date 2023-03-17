package worker

import (
	"bytes"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/benmizrahi/godist/plugins"
	"github.com/benmizrahi/godist/plugins/contract"
	"github.com/benmizrahi/godist/protos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	ID          string
	MaxParallel int
	Master      string
	Host        string
	Port        int
	Plugins     map[string]contract.IPluginContract
	Http        *gin.Engine
}

func NewWorker(host string, port int, masterPath string) *Worker {

	w := &Worker{
		ID:          (uuid.New()).String(),
		MaxParallel: 10,
		Master:      "http://" + masterPath,
		Http:        gin.Default(),
		Host:        host,
		Port:        port,
		Plugins:     map[string]contract.IPluginContract{},
	}

	w.registerToMaster()
	w.loadBuildInPlugins()
	w.Http.GET("/api/v1/health", w.healthCheck)
	w.Http.POST("/api/v1/tasks", w.tasksHandler)
	go w.Http.Run(w.Host + ":" + strconv.Itoa(w.Port))
	logrus.Info("worker %s is listening at %s", w.ID, w.Host+":"+strconv.Itoa(w.Port))
	return w
}

func (w *Worker) loadBuildInPlugins() {
	for key, plugin := range plugins.MakeBuildIns() {
		w.Plugins[key] = plugin()
		logrus.Info("GoDist Worker, plugin loaded ", key)
	}
}

func (w *Worker) registerToMaster() {
	req := &protos.RegisterReq{
		Uuid: w.ID,
		Http: "http://" + w.Host + ":" + strconv.Itoa(w.Port),
	}
	body, err := proto.Marshal(req)
	_, err = http.Post(w.Master+"/api/register", "application/protobuf", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
}
