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
	center           *zkserve.ZkCenter
	quit             chan int
}

func init() {

}

func New() *Service {
	return &Service{
		svrId:  atomic.AddInt32(&loginSvrId, 1),
		center: zkserve.New(),
		quit:   make(chan int, 1),
	}
}

func (s *Service) init() error {
	err := s.center.ConnectWithWatcher(global.ZookeeperHosts, time.Second*60, s.watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Service) watchEventCb(event zk.Event) {
	// log.Println("loginserver >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("loginserver path:", event.Path)
	// log.Println("loginserver type:", event.Type)
	// log.Println("loginserver state:", event.State)
	// log.Println("loginserver <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := s.center.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("loginserver:", string(data))

			_, err = s.center.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) register(servername string) error {
	err := s.center.Create(consts.ZookeeperKeyLoginRoot, 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	s.registerNodeName = fmt.Sprintf("%s/login_%s", consts.ZookeeperKeyLoginRoot /*servername*/, strconv.Itoa(int(s.svrId)))
	err = s.center.Create(s.registerNodeName, zk.FlagEphemeral, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	m := make(map[string]interface{})
	m["svrId"] = s.svrId
	m["register_time"] = time.Now().Format("2006-01-02 15:04:05")
	var data []byte
	data, err = json.Marshal(m)
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.center.Write(s.registerNodeName, data)
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

	err = s.register("")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.center.Watch(consts.ZookeeperKeyLogin)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.center.Create(consts.ZookeeperKeyLogin, 0, zk.PermRead)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Println(x)
			}
			log.Println("login server service exit")
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// go func() {
				// 	d, err := s.center.Children(consts.ZookeeperKeyLoginRoot)
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
	close(s.quit)
	s.center.Close()
}
