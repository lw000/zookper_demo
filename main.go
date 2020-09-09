package main

import (
	"demo/zookper_demo/game"
	"demo/zookper_demo/gate"
	"demo/zookper_demo/hall"
	"demo/zookper_demo/kfka"
	"demo/zookper_demo/login"
	"demo/zookper_demo/master"
	"demo/zookper_demo/zkserve"
	"github.com/judwhite/go-svc/svc"
	"log"
	"os"
	"time"
)

var (
	logger *log.Logger
)

type Program struct {
	center        *zkserve.ZkCenter
	masterService *master.Service
	hallService   *hall.Service
	loginService  *login.Service
	gameService   *game.Service
	gateService   *gate.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}
	p.center = zkserve.New()

	p.masterService = master.New()
	p.hallService = hall.New()
	p.loginService = login.New()
	p.gameService = game.New()
	p.gateService = gate.New()

	return nil
}

func (p *Program) Start() error {
	var err error

	err = p.masterService.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.hallService.Start()
	if err != nil {
		log.Println(err)
		return err
	}
	// err = p.hallService.Register("a")
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	err = p.loginService.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.gameService.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.gateService.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = kfka.Conn()
	if err != nil {
		log.Println(err)
		return err
	}
	defer kfka.Close()

	return nil
}

func (p *Program) Stop() error {
	p.masterService.Stop()
	p.hallService.Stop()
	p.loginService.Stop()
	p.gameService.Stop()
	p.gateService.Stop()
	time.Sleep(time.Millisecond * 20)
	return nil
}

func main() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c,os.Interrupt,os.Kill)
	// log.Println(<-c)

	f, err := os.Create("./log/log.log")
	if err != nil {
		log.Fatalln(err)
	}

	logger = log.New(f, "", log.Llongfile)

	go func() {
		for i := 0; i < 100; i++ {
			logger.Println(time.Now().Format("2006-01-02 15:04:05.000000000"), "this is test write file")
		}
	}()

	var pro Program
	err = svc.Run(&pro)
	if err != nil {
		log.Panic(err)
	}
}
