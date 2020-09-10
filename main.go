package main

import (
	"demo/zookper_demo/kfka"
	"demo/zookper_demo/masterserver/master"
	"github.com/judwhite/go-svc/svc"
	"log"
	"time"
)

var (
	logger *log.Logger
)

type Program struct {
	masterService *master.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	p.masterService = master.New()

	return nil
}

func (p *Program) Start() error {
	var err error

	err = p.masterService.Start()
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
	time.Sleep(time.Millisecond * 20)
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
