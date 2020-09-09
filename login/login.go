package login

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkserve"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type Service struct {
	center *zkserve.ZkCenter
	quit   chan int
}

func init() {

}

func New() *Service {
	return &Service{
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
	// log.Println("login >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("login path:", event.Path)
	// log.Println("login type:", event.Type)
	// log.Println("login state:", event.State)
	// log.Println("login <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := s.center.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("login:", string(data))

			_, err = s.center.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) Register(servername string) error {

	return nil
}

func (s *Service) Start() error {
	var err error
	err = s.init()
	if err != nil {
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
			log.Println("login service exit")
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
