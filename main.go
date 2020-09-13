package main

import (
	"github.com/judwhite/go-svc/svc"
	"log"
	"net"
	"time"
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
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	return nil
}

func (p *Program) Start() error {
	log.Println("localIP:", localIP())
	return nil
}

func (p *Program) Stop() error {
	time.Sleep(time.Millisecond * 10)
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
