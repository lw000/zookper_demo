package zkserve

type ZkLock struct {
}

func (lk *ZkLock) Lock(nodePath, clientGuid string) bool {
	return false
}

func (lk *ZkLock) UnLock(nodePath, clientGuid string) {

}

func (lk *ZkLock) Exists(nodePath string) {

}
