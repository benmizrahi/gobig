package master

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/benmizrahi/gobig/internal/common"
	"github.com/benmizrahi/gobig/internal/protos"
	lzstring "github.com/daku10/go-lz-string"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// singeltone instance of master!
var lock = &sync.Mutex{}

// Singel instance
var masterInstance *Master

type Master struct {
	MasterPath string
	context    *Context
	Http       *gin.Engine
}

func NewMaster(isLocal bool, host string, port int, minWorkers int) *Master {
	if masterInstance == nil {
		lock.Lock()
		log.Info("gobig Master, Creating new master instance")
		m := &Master{
			MasterPath: host + ":" + strconv.Itoa(port),
			Http:       gin.New(),
		}

		m.Http.Use(ginlogrus.Logger(logrus.New()), gin.Recovery())

		m.Http.POST("/api/register", m.RegisterHandler)
		go m.Http.Run(m.MasterPath)
		log.Info("gobig Master, master is listening on ", m.MasterPath)

		m.context = NewContext(isLocal, minWorkers, m.MasterPath)

		m.context.InitContext()

		lock.Unlock()
		return m
	}
	return masterInstance
}

func (m *Master) RegisterHandler(c *gin.Context) {
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatalln("Failed to parse register request:", err)
	}
	req := &protos.RegisterReq{}
	if err := proto.Unmarshal(buf, req); err != nil {
		log.Fatalln("Failed to parse register request:", err)
	}

	m.context.Workers[req.Uuid] = req.Http

	data := &protos.RegisterRes{
		Ok: true,
	}

	c.ProtoBuf(http.StatusOK, data)
}

func (m *Master) Parallelize(data [][]string, option common.Options) *Mafream {
	//first row it's the columns row
	columns := data[0]
	mf := NewDataFrame(m.context, columns)

	mf.partitions = m.buildPartitons(data, 10)

	for _, p := range mf.partitions {

		p.Tasks = append(p.Tasks, &protos.Task{
			Uuid:         uuid.New().String(),
			Instactions:  []string{protos.IN_MEMORY_READ},
			CreationTime: timestamppb.Now(),
		})
	}
	return mf
}

func (m *Master) Load() *Mafream {
	mf := NewDataFrame(m.context, []string{})
	return mf
}

func (m *Master) buildPartitons(data [][]interface{}, requestedNumPartitions *int) []*protos.IPartition {
	numPartitions := m.calculatePartitons(data, requestedNumPartitions)
	partitions := make([]*protos.IPartition, numPartitions)

	// Shuffle the data
	for _, row := range data {
		dataTypes := m.recogizeTypes(row)
		partitionIndex := rand.Intn(numPartitions)
		if partitions[partitionIndex] == nil {
			partitions[partitionIndex] = &protos.IPartition{}
		}
		if partitions[partitionIndex].Data == nil {
			partitions[partitionIndex].Data = make([]*protos.Data, 0)
		}
		data, err := lzstring.Compress()
		if err != nil {
			// data = row
		}
		partitions[partitionIndex].Data = append(partitions[partitionIndex].Data, &protos.Data{
			DataTypes:     *dataTypes,
			CopressedData: data,
		})
	}

	return partitions
}
func (m *Master) calculatePartitons(data [][]interface{}, numPartitions *int) int {
	return 1
}

func (m *Master) recogizeTypes(data []interface{}) *[]protos.DataType {

	return &[]protos.DataType{}
}
