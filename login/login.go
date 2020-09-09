package login

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkserve"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

var (
	center *zkserve.Zkp
)

type Service struct {
	quit chan int
}

func init() {
	center = zkserve.New()
}

func New() *Service {
	return &Service{
		quit: make(chan int, 1),
	}
}

func watchEventCb(event zk.Event) {
	// log.Println("login >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("login path:", event.Path)
	// log.Println("login type:", event.Type)
	// log.Println("login state:", event.State)
	// log.Println("login <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := center.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("login:", string(data))

			_, err = center.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) Start() error {
	err := center.ConnectWithWatcher(global.ZKHosts, time.Second*60, watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = center.Watch(consts.ZookeeperKeyLogin)
	if err != nil {
		log.Println(err)
		return err
	}

	err = center.Create(consts.ZookeeperKeyLogin, 0, zk.PermRead)
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
}
