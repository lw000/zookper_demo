package login

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkserve"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	loginSvrId int32 = 0
)

type Service struct {
	registerNodeName string
	svrId            int32
	client           *zkserve.ZkClient
	quit             chan int
	startTime        time.Time
}

func init() {

}

func New() *Service {
	return &Service{
		svrId:  atomic.AddInt32(&loginSvrId, 1),
		client: zkserve.New(),
		quit:   make(chan int, 1),
	}
}

func (s *Service) init() error {
	err := s.client.ConnectWithWatcher(global.ZookeeperHosts, time.Second*60, s.watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}

	s.startTime = time.Now()
	return nil
}

func (s *Service) watchEventCb(event zk.Event) {
	// log.Println("login server >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("login server path:", event.Path)
	// log.Println("login server type:", event.Type)
	// log.Println("login server state:", event.State)
	// log.Println("login server <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := s.client.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("login server id %d, data:%s\n", s.svrId, string(data))

			_, err = s.client.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) register(node string) error {
	err := s.client.Create(consts.LoginServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}
	m := make(map[string]interface{})
	m["svrId"] = s.svrId
	m["register_time"] = time.Now().Format("2006-01-02 15:04:05")
	m["startTime"] = s.startTime.Format("2006-01-02 15:04:05")

	var data []byte
	data, err = json.Marshal(m)
	if err != nil {
		log.Println(err)
		return err
	}
	s.registerNodeName = fmt.Sprintf("%s/%s", consts.LoginServerRoot, node)
	err = s.client.Create(s.registerNodeName, data, zk.FlagEphemeral, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) Start() error {
	var err error
	err = s.init()
	if err != nil {
		return err
	}

	err = s.register(fmt.Sprintf("%s_%s", "login", strconv.Itoa(int(s.svrId))))
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.client.Watch(consts.LoginConfig)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.LoginConfig, []byte(""), 0, zk.PermRead)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Println(x)
			}
			log.Printf("login server [%d] exit\n", s.svrId)
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// go func() {
				// 	d, err := s.client.Children(consts.LoginServerRoot)
				// 	if err == nil {
				// 		log.Println(d)
				// 	}
				// }()
			case <-s.quit:
				return
			}
		}
	}()

	return nil
}

func (s *Service) Stop() {
	s.client.Close()
	close(s.quit)
}
