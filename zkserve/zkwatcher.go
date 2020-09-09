package zkserve

import (
	"demo/zookper_demo/consts"
	"github.com/samuel/go-zookeeper/zk"
	"log"
)

func WatcherNodeCreated(ev <-chan zk.Event, cb func(path string)) {
	event := <-ev
	switch event.Path {
	case consts.ZookeeperKey:
		cb(event.Path)
	}
	log.Println("*******************")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("-------------------")
}

func WatcherNodeChanged(ev <-chan zk.Event, cb func(path string)) {
	event := <-ev
	switch event.Type {
	case zk.EventNodeDataChanged:
		cb(event.Path)
	case zk.EventNodeCreated:
	case zk.EventNodeDeleted:
	}
	log.Println("*******************")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("-------------------")
}

func WatcherNodeDeleted(ev <-chan zk.Event, cb func(path string)) {
	event := <-ev
	switch event.Path {
	case consts.ZookeeperKey:
		cb(event.Path)
	}
	log.Println("*******************")
	log.Println("path:", event.Path)
	log.Println("type:", event.Type.String())
	log.Println("state:", event.State.String())
	log.Println("-------------------")
}
