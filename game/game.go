package game

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkserve"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	gameSvrId int32 = 0
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
		svrId:  atomic.AddInt32(&gameSvrId, 1),
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
	// log.Println("game >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("game path:", event.Path)
	// log.Println("game type:", event.Type)
	// log.Println("game state:", event.State)
	// log.Println("game <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := s.center.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("game:", string(data))

			_, err = s.center.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) register(servername string) error {
	err := s.center.Create(consts.ZookeeperKeyGameRoot, 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	s.registerNodeName = fmt.Sprintf("%s/game_%s", consts.ZookeeperKeyGameRoot /*servername*/, strconv.Itoa(int(s.svrId)))
	err = s.center.Create(s.registerNodeName, zk.FlagEphemeral, zk.PermAll)
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

	_, err = s.center.Watch(consts.ZookeeperKeyGame)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.center.Create(consts.ZookeeperKeyGame, 0, zk.PermRead)
	if err != nil {
		log.Println(err)
		return err
	}

	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Println(x)
			}
			log.Println("game service exit")
		}()

		for {
			select {
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
