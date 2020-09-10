package master

import (
	"demo/zookper_demo/consts"
	"demo/zookper_demo/global"
	"demo/zookper_demo/zkserve"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net/http"
	"time"
)

type Service struct {
	center *zkserve.ZkCenter
	quit   chan int
	count  int32
}

func init() {

}

func New() *Service {
	return &Service{
		center: zkserve.New(),
		quit:   make(chan int, 1),
	}
}

func (s *Service) watchEventCb(event zk.Event) {
	log.Println("masterserver >>>>>>>>>>>>>>>>>>>>>>")
	log.Println("masterserver path:", event.Path)
	log.Println("masterserver type:", event.Type)
	log.Println("masterserver state:", event.State)
	log.Println("masterserver <<<<<<<<<<<<<<<<<<<<<<")

	// if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
	// 	go func() {
	// 		data, err := s.center.Read(event.Path)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 		log.Println("game:", string(data))
	//
	// 		_, err = s.center.Watch(event.Path)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}()
	// }
}

func (s *Service) init() error {
	err := s.center.ConnectWithWatcher(global.ZookeeperHosts, time.Second*60, s.watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) Start() error {
	err := s.init()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.center.Watch(consts.ZookeeperKeyHallRoot)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.center.Watch(consts.ZookeeperKeyGameRoot)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.center.Watch(consts.ZookeeperKeyLoginRoot)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.center.Create(consts.ZookeeperKeyHall, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.center.Create(consts.ZookeeperKeyGame, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.center.Create(consts.ZookeeperKeyLogin, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	go s.Run()

	go s.modifyGameConfig()
	go s.modifyHallConfig()
	go s.modifyLoginConfig()

	// for i := 0; i < 100; i++ {
	// 	go func(i int) {
	// 		// s.count += 1
	// 		atomic.AddInt32(&s.count, 1)
	// 		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>:", s.count)
	// 	}(i)
	// }

	return nil
}

func (s *Service) Run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("masters erver service exit")
	}()

	// gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()
	engine.GET("/api/sync", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "sync config success", "data": gin.H{}})
	})
	err := engine.Run(":9200")
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) modifyHallConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master server modifyHallConfig quit")
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := make(map[string]interface{})
			m["a"] = "10"
			m["b"] = 10
			m["c"] = 20.20
			m["d"] = true
			m["e"] = time.Now().Format("2006-01-02 15:04:05.000000")
			m["node"] = consts.ZookeeperKeyHall
			data, err := json.Marshal(m)
			if err == nil {
				err = s.center.Write(consts.ZookeeperKeyHall, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			return
		}
	}
}

func (s *Service) modifyGameConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master server modifyGameConfig quit")
	}()

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := make(map[string]interface{})
			m["a"] = "10"
			m["b"] = 10
			m["c"] = 20.20
			m["d"] = true
			m["e"] = time.Now().Format("2006-01-02 15:04:05.000000")
			m["node"] = consts.ZookeeperKeyGame
			data, err := json.Marshal(m)
			if err == nil {
				err = s.center.Write(consts.ZookeeperKeyGame, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:

			return
		}
	}
}

func (s *Service) modifyLoginConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master server modifyLoginConfig quit")
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m := make(map[string]interface{})
			m["a"] = "10"
			m["b"] = 10
			m["c"] = 20.20
			m["d"] = true
			m["e"] = time.Now().Format("2006-01-02 15:04:05.000000")
			m["node"] = consts.ZookeeperKeyLogin
			data, err := json.Marshal(m)
			if err == nil {
				err = s.center.Write(consts.ZookeeperKeyLogin, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			return
		}
	}
}

func (s *Service) Stop() {
	close(s.quit)
	s.center.Close()
}
