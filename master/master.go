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

var (
	center *zkserve.ZkCenter
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

func (s *Service) init() error {
	err := center.Connect(global.ZookeeperHosts, time.Second*60)
	if err != nil {
		log.Println(err)
		return err
	}

	err = center.Create(consts.ZookeeperKeyHall, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	err = center.Create(consts.ZookeeperKeyGame, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	err = center.Create(consts.ZookeeperKeyLogin, 0, zk.PermRead|zk.PermWrite|zk.PermDelete)
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

	go s.pushGameConfig()
	go s.pushHallConfig()
	go s.pushLoginConfig()

	return nil
}

func (s *Service) Run() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master service exit")
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

func (s *Service) pushHallConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master pushHallConfig quit")
	}()

	ticker := time.NewTicker(time.Second)
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
				err = center.Write(consts.ZookeeperKeyHall, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			ticker.Stop()
			return
		}
	}
}

func (s *Service) pushGameConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master pushGameConfig quit")
	}()

	ticker := time.NewTicker(time.Second * 2)
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
				err = center.Write(consts.ZookeeperKeyGame, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			ticker.Stop()
			return
		}
	}
}

func (s *Service) pushLoginConfig() {
	defer func() {
		if x := recover(); x != nil {
			log.Println(x)
		}

		log.Println("master pushLoginConfig quit")
	}()

	ticker := time.NewTicker(time.Second * 3)
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
				err = center.Write(consts.ZookeeperKeyLogin, data)
				if err != nil {
					log.Println(err)
				}
			}
		case <-s.quit:
			ticker.Stop()
			return
		}
	}
}

func (s *Service) Stop() {
	close(s.quit)
}
