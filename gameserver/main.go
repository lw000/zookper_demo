package main

import (
	"demo/zookper_demo/gameserver/game"
	"github.com/judwhite/go-svc/svc"
	"log"
)

var (
	gameN int = 5
)

type Program struct {
	gameService []*game.Service
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}
	p.gameService = make([]*game.Service, 0)
	for i := 0; i < gameN; i++ {
		p.gameService = append(p.gameService, game.New())
	}

	return nil
}

func (p *Program) Start() error {
	var err error
	for _, svr := range p.gameService {
		err = svr.Start()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (p *Program) Stop() error {
	for _, svr := range p.gameService {
		svr.Stop()
	}

	return nil
}

func main() {
	var pro Program
	if err := svc.Run(&pro); err != nil {
		log.Fatalln(err)
	}
}
