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
	client *zkserve.ZkClient
	quit   chan int
	count  int32
}

func init() {

}

func New() *Service {
	return &Service{
		client: zkserve.New(),
		quit:   make(chan int, 1),
	}
}

func (s *Service) watchEventCb(event zk.Event) {
	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		log.Println("master server >>>>>>>>>>>>>>>>>>>>>>")
		log.Println("master server path:", event.Path)
		log.Println("master server type:", event.Type)
		log.Println("master server state:", event.State)
		log.Println("master server <<<<<<<<<<<<<<<<<<<<<<")
	}
}

func (s *Service) init() error {
	err := s.client.ConnectWithWatcher(global.ZookeeperHosts, time.Second*60, s.watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.initSystemConfig()
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

	go s.Run()

	go s.modifyGameConfig()
	go s.modifyHallConfig()
	go s.modifyLoginConfig()
	go s.modifyGateConfig()

	// for i := 0; i < 100; i++ {
	// 	go func(i int) {
	// 		// s.count += 1
	// 		atomic.AddInt32(&s.count, 1)
	// 		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>:", s.count)
	// 	}(i)
	// }

	return nil
}

func (s *Service) initSystemConfig() error {
	var err error
	err = s.client.Create(consts.Root, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.RootConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.client.Watch(consts.MasterConfig)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.MasterConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.HallConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.GameConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.LoginConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.GateConfig, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.HallServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.GameServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.LoginServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.client.Create(consts.GateServerRoot, []byte(""), 0, zk.PermAll)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.client.ChildrenW(consts.LoginServerRoot)
	if err != nil {
		log.Println(err)
		return err
	}
	//
	// _, err = s.client.ChildrenW(consts.GateServerRoot)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	return nil
}

func (s *Service) Run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}
		log.Printf("master server exit\n")
	}()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				go func() {
					childes, err := s.client.Children(consts.GameServerRoot)
					if err == nil {
						log.Println("GameServerRoot", childes)
						for _, child := range childes {
							path := consts.GameServerRoot + "/" + child
							data, err := s.client.Read(path)
							if err == nil {
								log.Println(path, string(data))
							}
						}
					}
				}()

				go func() {
					childes, err := s.client.Children(consts.HallServerRoot)
					if err == nil {
						log.Println("HallServerRoot", childes)
					}
				}()

				go func() {
					childes, err := s.client.Children(consts.GateServerRoot)
					if err == nil {
						log.Println("GateServerRoot", childes)
					}
				}()

				go func() {
					childes, err := s.client.Children(consts.LoginServerRoot)
					if err == nil {
						log.Println("LoginServerRoot", childes)
					}
				}()

			case <-s.quit:
				return
			}
		}
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
			m["node"] = consts.HallConfig
			data, err := json.Marshal(m)
			if err == nil {
				err = s.client.Write(consts.HallConfig, data)
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
			m["node"] = consts.GameConfig
			data, err := json.Marshal(m)
			if err == nil {
				err = s.client.Write(consts.GameConfig, data)
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
			m["node"] = consts.LoginConfig
			data, err := json.Marshal(m)
			if err == nil {
				err = s.client.Write(consts.LoginConfig, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			return
		}
	}
}

func (s *Service) modifyGateConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master server modifyGateConfig quit")
	}()

	ticker := time.NewTicker(time.Second * 5)
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
			m["node"] = consts.GateConfig
			data, err := json.Marshal(m)
			if err == nil {
				err = s.client.Write(consts.GateConfig, data)
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
	s.client.Close()
	close(s.quit)
}
