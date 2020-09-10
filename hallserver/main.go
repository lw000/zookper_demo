package main

import (
	"demo/zookper_demo/hallserver/hall"
	"github.com/judwhite/go-svc/svc"
	"log"
	"time"
)

var (
	hallN int = 5
)

type Program struct {
	hallService []*hall.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}
	p.hallService = make([]*hall.Service, 0)
	for i := 0; i < hallN; i++ {
		p.hallService = append(p.hallService, hall.New())
	}
	return nil
}

func (p *Program) Start() error {
	var err error
	for _, svr := range p.hallService {
		err = svr.Start()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (p *Program) Stop() error {
	for _, svr := range p.hallService {
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
