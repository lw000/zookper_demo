package hall

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

func (s *Service) watchEventCb(event zk.Event) {
	// log.Println("hall >>>>>>>>>>>>>>>>>>>>>>")
	// log.Println("hall path:", event.Path)
	// log.Println("hall type:", event.Type)
	// log.Println("hall state:", event.State)
	// log.Println("hall <<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := center.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("hall:", string(data))

			_, err = center.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func (s *Service) Start() error {

	err := center.ConnectWithWatcher(global.ZKHosts, time.Second*60, s.watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = center.Watch(consts.ZookeeperKeyHall)
	if err != nil {
		log.Println(err)
		return err
	}

	err = center.Create(consts.ZookeeperKeyHall, 0, zk.PermRead)
	if err != nil {
		log.Println(err)
		return err
	}

	go s.Run()

	return nil
}

func (s *Service) Run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("hall service exit")
	}()

	for {
		select {
		case <-s.quit:
			return
		}
	}
}

func (s *Service) Stop() {
	close(s.quit)
}
