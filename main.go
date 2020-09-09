package main

import (
	"demo/zookper_demo/game"
	"demo/zookper_demo/global"
	"demo/zookper_demo/hall"
	"demo/zookper_demo/kfka"
	"demo/zookper_demo/login"
	"demo/zookper_demo/master"
	"fmt"
	"github.com/judwhite/go-svc/svc"
	"github.com/samuel/go-zookeeper/zk"
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
}

func zkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func zkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func watchEventCb(event zk.Event) {
	log.Println(">>>>>>>>>>>>>>>>>>>>>>")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type)
	log.Println("state:", event.State)
	log.Println("<<<<<<<<<<<<<<<<<<<<<<")

	if len(event.Path) > 0 && event.Type == zk.EventNodeDataChanged {
		go func() {
			data, err := global.ZKPServe.Read(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(data))
		}()

		go func() {
			_, err := global.ZKPServe.Watch(event.Path)
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

func zkWatch() {
	option := zk.WithEventCallback(watchEventCb)
	var hosts = []string{"192.168.0.115:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*60, option)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	var path1 = "/zk_test_go1"
	var data1 = []byte("zk_test_go1_data1")
	exist, s, _, err := conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("path[%s] exist[%t]\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", zkStateStringFormat(s))

	// try create
	var acls = zk.WorldACL(zk.PermAll)
	p, err_create := conn.Create(path1, data1, zk.FlagEphemeral, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return
	}
	fmt.Printf("created path[%s]\n", p)

	time.Sleep(time.Second * 2)

	exist, s, _, err = conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("path[%s] exist[%t] after create\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", zkStateStringFormat(s))
	// delete
	err = conn.Delete(path1, s.Version)
	if err != nil {
		fmt.Println(err)
		return
	}

	exist, s, _, err = conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("path[%s] exist[%t] after delete\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", zkStateStringFormat(s))
}

func zkTest() {
	var hosts = []string{"192.168.0.115:2181"}
	// var provider = zk.DNSHostProvider{}
	// err := provider.Init(hosts)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// servername,retryStart := provider.Next()

	conn, event, err := zk.Connect(hosts, time.Second*60)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	log.Println(event)

	// create
	var path = "/zk_test"
	exist, s, err := conn.Exists(path)
	if err != nil {
		log.Panic(err)
	}
	if exist {
		log.Printf("path:%s,exists:%t", path, exist)
	} else {
		// var flags = []int32{zk.PermCreate,zk.PermRead,zk.PermWrite}
		var flags int32 = 0
		var acls = zk.WorldACL(zk.PermAll)
		var data = []byte("hello")
		p, errCreate := conn.Create(path, data, flags, acls)
		if errCreate != nil {
			log.Panic(errCreate)
		}
		log.Println("created:", p)
	}

	// get
	v, s, err := conn.Get(path)
	if err != nil {
		log.Panic(err)
	}
	log.Println(v)
	log.Println("state:")
	log.Println(zkStateStringFormat(s))

	exist, s, err = conn.Exists(path)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("path:%s,exists:%t", path, exist)

	// update
	var new_data = []byte(`{"a":10,"b":123131231,"c":"123231232131313"}`)
	s, err = conn.Set(path, new_data, s.Version)
	if err != nil {
		log.Panic(err)
	}
	// get
	v, s, err = conn.Get(path)
	if err != nil {
		log.Panic(err)
		return
	}
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {

	} else {

	}

	err := global.ZKPServe.ConnectWithWatcher(global.ZKHosts, time.Second*60, watchEventCb)
	if err != nil {
		log.Println(err)
		return err
	}

	p.master_service = master.New()
	p.hall_service = hall.New()
	p.login_service = login.New()
	p.game_service = game.New()

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
