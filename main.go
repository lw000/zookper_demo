package main

import (
	"demo/zookper_demo/tasks"
	"encoding/json"
	// "github.com/flier/curator.go"
	"github.com/heteddy/delaytask-go/delaytask"
	"github.com/judwhite/go-svc/svc"
	"log"
	"net"
)

// var (
// 	logger *log.Logger
// )

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

type Program struct {
	// client curator.CuratorFramework
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	return nil
}

func (p *Program) Start() error {
	log.Println("localIP:", localIP())

	engine := delaytask.NewEngine("1s", 10, "redis://@192.168.0.115:6379/4",
		"messageQ", "remote-task0:")
	engine.AddTaskCreator("OncePingTask", func(task string) delaytask.Runner {
		p := &tasks.OncePingTask{}
		if err := json.Unmarshal([]byte(task), p); err != nil {
		} else {
			return p
		}
		return nil
	})
	engine.AddTaskCreator("PeriodPingTask", func(task string) delaytask.Runner {
		t := &tasks.PeriodPingTask{}
		if err := json.Unmarshal([]byte(task), t); err != nil {
			return nil
		} else {
			return t
		}
	})
	engine.Start()

	// retryPolicy := curator.NewExponentialBackoffRetry(time.Second, 3, 15*time.Second)
	// connString := "192.168.0.115:2182"
	// p.client = curator.NewClient(connString, retryPolicy)
	// var err error
	// if err = p.client.Start(); err != nil {
	// 	return err
	// }

	return nil
}

func (p *Program) Stop() error {
	// _ = p.client.Close()
	return nil
}

func main() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c,os.Interrupt,os.Kill)
	// log.Println(<-c)

	// f, err := os.Create("./log/log.log")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	//
	// logger = log.New(f, "", log.Llongfile)
	//
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		logger.Println(time.Now().Format("2006-01-02 15:04:05.000000000"), "this is test write file")
	// 	}
	// }()

	var pro Program
	if err := svc.Run(&pro); err != nil {
		log.Fatalln(err)
	}
}
