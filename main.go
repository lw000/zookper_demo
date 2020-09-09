package main

import (
	"demo/zookper_demo/game"
	"demo/zookper_demo/gate"
	"demo/zookper_demo/hall"
	"demo/zookper_demo/kfka"
	"demo/zookper_demo/login"
	"demo/zookper_demo/master"
	"github.com/judwhite/go-svc/svc"
	"log"
	"os"
	"time"
)

var (
	logger *log.Logger
)

type Program struct {
	master_service *master.Service
	hall_service   *hall.Service
	login_service  *login.Service
	game_service   *game.Service
	gate_service   *gate.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	p.master_service = master.New()
	p.hall_service = hall.New()
	p.login_service = login.New()
	p.game_service = game.New()
	p.gate_service = gate.New()
	return nil
}

func (p *Program) Start() error {
	var err error

	err = p.master_service.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.hall_service.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.login_service.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.game_service.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	err = p.gate_service.Start()
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
	p.master_service.Stop()
	p.hall_service.Stop()
	p.login_service.Stop()
	p.game_service.Stop()
	p.gate_service.Stop()
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
