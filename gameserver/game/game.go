package game

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
	gameSvrId int32 = 0
)

type Service struct {
	registerNodeName string
	svrId            int32
	client           *zkc.ZkClient
	quit             chan int
	lock             *zk.Lock
}

func init() {

}

func New() *Service {
	return &Service{
		svrId:  atomic.AddInt32(&gameSvrId, 1),
		client: zkc.New(),
		quit:   make(chan int),
	}
}

func (s *Service) init() error {
	s.lock = s.client.Lock("game-data-lock")
	return s.client.Connect(global.ZookeeperHosts, time.Second*60)
}

func (s *Service) register(node string) error {
	var err error
	err = s.client.Create(consts.GameServerRoot, []byte(""), 0, zk.PermAll)
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
	s.registerNodeName = fmt.Sprintf("%s/service-%s", consts.GameServerRoot, node)
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

	if err = s.client.Create(consts.GameConfig, []byte(""), 0, zk.PermRead); err != nil {
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
		log.Printf("game server [%d] exit\n", s.svrId)
	}()

	go s.WatchConfigChanged()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// go func() {
			// 	d, err := s.client.Children(consts.GameServerRoot)
			// 	if err == nil {
			// 		log.Println(d)
			// 	}
			// }()
			// var err error
			// if err = s.lock.Lock(); err == nil {
			// 	global.SharedData["gamedData"] = "game data"
			// }
			// if err = s.lock.Unlock(); err != nil {
			// 	log.Println("unlock error")
			// }
		case <-s.quit:
			return
		}
	}
}

func (s *Service) WatchConfigChanged() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}
		log.Printf("game server watch confiog [%d] exit\n", s.svrId)
	}()

	for {
		ev, err := s.client.Watch(consts.GameConfig)
		if err != nil {
			return
		}

		select {
		case event := <-ev:
			if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
				data, err := s.client.Read(event.Path)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("game server id %d, data:%s\n", s.svrId, string(data))
			}
		case <-s.quit:
			return
		}
	}
}

func (s *Service) Stop() {
	s.client.Close()
	close(s.quit)
}
