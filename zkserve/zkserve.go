package zkserve

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type ZkCenter struct {
	conn *zk.Conn
}

func New() *ZkCenter {
	return &ZkCenter{}
}

func zkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func zkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func (zkc *ZkCenter) Connect(hosts []string, sessionTimeout time.Duration) error {
	var err error
	zkc.conn, _, err = zk.Connect(hosts, sessionTimeout)
	if err != nil {
		return err
	}
	return nil
}

func (zkc *ZkCenter) ConnectWithWatcher(hosts []string, sessionTimeout time.Duration, watchEventCb func(event zk.Event)) error {
	if watchEventCb == nil {
		panic("watchEventCb is nil")
	}
	var err error
	option := zk.WithEventCallback(watchEventCb)
	zkc.conn, _, err = zk.Connect(hosts, sessionTimeout, option)
	if err != nil {
		return err
	}
	return nil
}

func (zkc *ZkCenter) ZkClose() {
	zkc.conn.Close()
}

func (zkc *ZkCenter) Create(nodePath string, flags int32, perm int32) error {
	exist, _ := zkc.Exists(nodePath)
	if !exist {
		// flags有4种取值：
		// 0:永久，除非手动删除
		// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
		// zk.FlagSequence  = 2:会自动在节点后面添加序号
		// 3:Ephemeral和Sequence，即，短暂且自动添加序号
		// var flags int32 = 0 // zk.FlagEphemeral | zk.FlagSequence
		var acl = zk.WorldACL(perm) // 表示该节点没有权限限制
		p, err := zkc.conn.Create(nodePath, nil, flags, acl)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("zookeeper path created:", p)
	}
	return nil
}

func (zkc *ZkCenter) Exists(nodePath string) (bool, *zk.Stat) {
	exist, state, err := zkc.conn.Exists(nodePath)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	return exist, state
}

// read zookeeper data
func (zkc *ZkCenter) Read(nodePath string) ([]byte, error) {
	data, s, err := zkc.conn.Get(nodePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if s.Version > 0 {
		// log.Println(zkStateStringFormat(s))
	}
	return data, nil
}

// write zookeeper data
func (zkc *ZkCenter) Write(nodePath string, data []byte) error {
	exist, s := zkc.Exists(nodePath)
	var err error
	if exist {
		_, err = zkc.conn.Set(nodePath, data, s.Version)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err = zkc.conn.Set(nodePath, data, 0)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (zkc *ZkCenter) Delete(nodePath string) error {
	exist, s, err := zkc.conn.Exists(nodePath)
	if err != nil {
		log.Println(err)
	}
	if exist {
		err = zkc.conn.Delete(nodePath, s.Version)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (zkc *ZkCenter) Watch(nodePath string) (<-chan zk.Event, error) {
	exist, _, ev, err := zkc.conn.ExistsW(nodePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if exist {

	}
	return ev, nil
}

func (zkc *ZkCenter) Children(nodePath string) ([]string, error) {
	s, _, err := zkc.conn.Children(nodePath)
	if err != nil {
		return nil, err
	}
	return s, nil
}
