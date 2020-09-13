package main

import (
	"demo/zookper_demo/gameserver/game"
	"encoding/json"
	"github.com/judwhite/go-svc/svc"
	"io/ioutil"
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
		// go func(svr *game.Service) {
		// 	time.AfterFunc(time.Second*time.Duration(rand.Intn(10)+10), func() {
		// 		svr.Stop()
		// 	})
		// }(svr)
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
	data, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(data, &ServiceConf)
	if err != nil {
		log.Panic(err)
	}
	var pro Program
	if err := svc.Run(&pro); err != nil {
		log.Fatalln(err)
	}
}
