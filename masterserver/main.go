package main

import (
	"demo/zookper_demo/masterserver/master"
	"github.com/judwhite/go-svc/svc"
	"log"
	"time"
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

	return nil
}

func (p *Program) Stop() error {
	p.masterService.Stop()
	time.Sleep(time.Millisecond * 20)
	return nil
}

func main() {

	var pro Program
	if err := svc.Run(&pro); err != nil {
		log.Fatalln(err)
	}
}
