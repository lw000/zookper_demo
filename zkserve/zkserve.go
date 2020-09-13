package zkserve

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type ZkClient struct {
	conn *zk.Conn
}

type ServiceNode struct {
	Services []string
	Root     string
}

func New() *ZkClient {
	return &ZkClient{}
}

func (z *ZkClient) Connect(hosts []string, sessionTimeout time.Duration) error {
	var err error
	z.conn, _, err = zk.Connect(hosts, sessionTimeout)
	if err != nil {
		return err
	}
	return nil
}

func (z *ZkClient) ConnectWithWatcher(hosts []string, sessionTimeout time.Duration, watchEventCb func(event zk.Event)) error {
	if watchEventCb == nil {
		panic("watchEventCb is nil")
	}
	var err error
	option := zk.WithEventCallback(watchEventCb)
	z.conn, _, err = zk.Connect(hosts, sessionTimeout, option)
	if err != nil {
		return err
	}
	return nil
}

func (z *ZkClient) Close() {
	z.conn.Close()
}

func (z *ZkClient) Create(path string, data []byte, flags int32, perm int32) error {
	exist, _ := z.Exists(path)
	if !exist {
		// flags有4种取值：
		// 0:永久，除非手动删除
		// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
		// zk.FlagSequence  = 2:会自动在节点后面添加序号
		// zk.Ephemeral | zk.Sequence = 3，即，短暂且自动添加序号
		// var flags int32 = 0 // zk.FlagEphemeral | zk.FlagSequence
		var acl = zk.WorldACL(perm) // zk.PermAll 表示该节点没有权限限制
		_, err := z.conn.Create(path, data, flags, acl)
		if err != nil && err != zk.ErrNodeExists {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (z *ZkClient) CreateProtectedEphemeralSequential(path string, data []byte, perm int32) error {
	exist, _ := z.Exists(path)
	if !exist {
		var acl = zk.WorldACL(perm) // zk.PermAll 表示该节点没有权限限制
		_, err := z.conn.CreateProtectedEphemeralSequential(path, data, acl)
		if err != nil && err != zk.ErrNodeExists {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (z *ZkClient) Exists(path string) (bool, *zk.Stat) {
	exist, state, err := z.conn.Exists(path)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	return exist, state
}

func (z *ZkClient) Sync(path string) (string, error) {
	return z.conn.Sync(path)
}

func (z *ZkClient) Read(path string) ([]byte, error) {
	data, s, err := z.conn.Get(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if s.Version > 0 {
		// log.Println(zkStateStringFormat(s))
	}
	return data, nil
}

func (z *ZkClient) Write(path string, data []byte) error {
	exist, s := z.Exists(path)
	var err error
	if exist {
		_, err = z.conn.Set(path, data, s.Version)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		_, err = z.conn.Set(path, data, 0)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (z *ZkClient) Delete(path string) error {
	exist, s, err := z.conn.Exists(path)
	if err != nil {
		log.Println(err)
	}
	if exist {
		err = z.conn.Delete(path, s.Version)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (z *ZkClient) Watch(path string) (<-chan zk.Event, error) {
	exist, _, ev, err := z.conn.ExistsW(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if exist {

	}
	return ev, nil
}

func (z *ZkClient) Children(path string) ([]string, error) {
	s, _, err := z.conn.Children(path)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (z *ZkClient) ChildrenW(path string) ([]string, <-chan zk.Event, error) {
	s, _, ev, err := z.conn.ChildrenW(path)
	if err != nil {
		return nil, nil, err
	}
	return s, ev, nil
}

func (z *ZkClient) Lock(path string) *zk.Lock {
	return zk.NewLock(z.conn, path, zk.WorldACL(zk.PermAll))
}

func (s *ServiceNode) GetNode() []string {

	return nil
}

func (z ZkClient) zkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func (z ZkClient) zkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}
