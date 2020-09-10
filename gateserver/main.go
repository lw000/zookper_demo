package main

import (
	"demo/zookper_demo/gateserver/gate"
	"github.com/judwhite/go-svc/svc"
	"log"
	"time"
)

var (
	gateN int = 5
)

type Program struct {
	gateService []*gate.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}
	p.gateService = make([]*gate.Service, 0)
	for i := 0; i < gateN; i++ {
		p.gateService = append(p.gateService, gate.New())
	}
	return nil
}

func (p *Program) Start() error {
	var err error
	for _, svr := range p.gateService {
		err = svr.Start()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (p *Program) Stop() error {
	for _, svr := range p.gateService {
		svr.Stop()
	}
	time.Sleep(time.Millisecond * 20)
	return nil
}

func main() {
	var pro Program
	if err := svc.Run(&pro); err != nil {
		log.Fatalln(err)
	}
}
