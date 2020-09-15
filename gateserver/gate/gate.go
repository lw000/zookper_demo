package gate

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkc"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	gateSvrId int32 = 0
)

type Service struct {
	registerNodeName string
	svrId            int32
	client           *zkc.ZkClient
	quit             chan int
}

func init() {
}

func New() *Service {
	return &Service{
		svrId:  atomic.AddInt32(&gateSvrId, 1),
		client: zkc.New(),
		quit:   make(chan int, 1),
	}
}

func (s *Service) init() error {
	return s.client.ConnectWithWatcher(global.ZookeeperHosts, time.Second*60, s.watchEventCb)
}

func (s *Service) watchEventCb(event zk.Event) {
	// log.Println("gate server >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("gate server path:", event.Path)
	// log.Println("gate server type:", event.Type)
	// log.Println("gate server state:", event.State)
	// log.Println("gate server <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := s.client.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("gate server id %d, data:%s\n", s.svrId, string(data))

			_, err = s.client.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) register(node string) error {
	err := s.client.Create(consts.GateServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}
	var data []byte
	data, err = json.Marshal(map[string]interface{}{
		"svrId":         s.svrId,
		"register_time": time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	s.registerNodeName = fmt.Sprintf("%s/service-%s", consts.GateServerRoot, node)
	err = s.client.Create(s.registerNodeName, data, zk.FlagEphemeral, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) Start() error {
	var err error
	if err = s.init(); err != nil {
		return err
	}

	if err = s.register(strconv.Itoa(int(s.svrId))); err != nil {
		return err
	}

	if _, err = s.client.Watch(consts.GateConfig); err != nil {
		return err
	}

	if err = s.client.Create(consts.GateConfig, []byte(""), 0, zk.PermRead); err != nil {
		return err
	}

	go s.run()

	return nil
}

func (s *Service) run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}
		log.Printf("gate server [%d] exit\n", s.svrId)
	}()

	for {
		select {
		case <-s.quit:
			return
		}
	}
}

func (s *Service) Stop() {
	s.client.Close()
	close(s.quit)
}
