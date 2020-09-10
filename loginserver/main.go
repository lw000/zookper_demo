package main

import (
	"demo/zookper_demo/loginserver/login"
	"github.com/judwhite/go-svc/svc"
	"log"
	"time"
)

var (
	loginN int = 5
)

type Program struct {
	loginService []*login.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}
	p.loginService = make([]*login.Service, 0)
	for i := 0; i < loginN; i++ {
		p.loginService = append(p.loginService, login.New())
	}
	return nil
}

func (p *Program) Start() error {
	var err error

	for _, svr := range p.loginService {
		err = svr.Start()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (p *Program) Stop() error {
	for _, svr := range p.loginService {
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
