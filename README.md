使用人比较多的zookeeper go client有

gozk
下载地址：https://wiki.ubuntu.com/gozk
文档地址: https://wiki.ubuntu.com/gozk

go-zookeeper
下载地址：https://github.com/samuel/go-zookeeper
文档地址: http://godoc.org/github.com/samuel/go-zookeeper/zk


dataDir=./data/zookeeper
dataLogDir=./log/zookeeper

# 伪部署：
每个服务下面dataDir目录下面需要常见myid文件，填写服务编号

三个zookeeper服务器都安装在同一个服务器（platform）上，需保证clientPort不相同。
将zookeeper安装包分别解压在三个目录server1，server2，server3下，配置文件zoo.cfg
Server1配置文件 zoo.cfg，server1在data目录下增加文件myid内容为1。

dataDir=../datadata
LogDir=../dataLog
clientPort=5181
server.1=platform:5888:6888
server.2= platform:5889:6889
server.3= platform:5890:6890
Server2配置文件 zoo.cfg，server1在data目录下增加文件myid内容为2。

dataDir=../datadata
LogDir=../dataLog
clientPort=6181
server.1=platform:5888:6888
server.2= platform:5889:6889
server.3= platform:5890:6890
Server3配置文件 zoo.cfg，server1在data目录下增加文件myid内容为3。

dataDir=../datadata
LogDir=../dataLog
clientPort=7181
server.1=platform:5888:6888
server.2= platform:5889:6889
server.3= platform:5890:6890
 

server.id=host:port:port : 表示了不同的zookeeper服务器的自身标识，作为集群的一部分，每一台服务器应该知道其他服务器的信息。用户可以从“server.id=host:port:port” 中读取到相关信息。在服务器的data(dataDir参数所指定的目录)下创建一个文件名为myid的文件，这个文件的内容只有一行，指定的是自身的id值。比如，服务器“1”应该在myid文件中写入“1”。这个id必须在集群环境中服务器标识中是唯一的，且大小在1～255之间。这一样配置中，zoo1代表第一台服务器的IP地址。第一个端口号（port）是从follower连接到leader机器的端口，第二个端口是用来进行leader选举时所用的端口。所以，在集群配置过程中有三个非常重要的端口：clientPort：2181、port:2888、port:3888。