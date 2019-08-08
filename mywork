//项目环境变量配置
for i in `seq 1 6`; do mkdir -p /tmp/$i/objects; done

export RABBITMQ_SERVER=amqp://admin:123456@10.30.12.92:5672
echo $RABBITMQ_SERVER

ifconfig ens5:1 192.168.0.1/16
ifconfig ens5:2 192.168.0.2/16
ifconfig ens5:3 192.168.0.3/16
ifconfig ens5:4 192.168.0.4/16
ifconfig ens5:5 192.168.0.5/16
ifconfig ens5:6 192.168.0.6/16
ifconfig ens5:7 192.168.1.1/16
ifconfig ens5:8 192.168.1.2/16


python3 rabbitmqadmin declare exchange name=apiserver type=fanout
python3 rabbitmqadmin declare exchange name=dataserver type=fanout


pkill dataserver
pkill interfaceserver
LISTEN_ADDRESS=192.168.0.1:12345 STORAGE_ROOT=/tmp/1 go run dataserver.go &
LISTEN_ADDRESS=192.168.0.2:12345 STORAGE_ROOT=/tmp/2 go run dataserver.go &
LISTEN_ADDRESS=192.168.0.3:12345 STORAGE_ROOT=/tmp/3 go run dataserver.go &
LISTEN_ADDRESS=192.168.0.4:12345 STORAGE_ROOT=/tmp/4 go run dataserver.go &
LISTEN_ADDRESS=192.168.0.5:12345 STORAGE_ROOT=/tmp/5 go run dataserver.go &
LISTEN_ADDRESS=192.168.0.6:12345 STORAGE_ROOT=/tmp/6 go run dataserver.go &
LISTEN_ADDRESS=192.168.1.1:12345 go run interfaceserver.go &
LISTEN_ADDRESS=192.168.1.2:12345 go run interfaceserver.go &


curl -v 192.168.1.1:12345/objects/test2 -XPUT -d"this is object test2"

curl  192.168.1.1:12345/objects/test2 -XPUT -d"this is object test2"
curl 192.168.1.1:12345/locate/test2
curl 192.168.1.1:12345/objects/test2


pkill dataserver

ll /tmp/1/objects;ll /tmp/2/objects;ll /tmp/3/objects;ll /tmp/4/objects;ll /tmp/5/objects;ll /tmp/6/objects;


***********************

elastic search
https://fuxiaopang.gitbooks.io/learnelasticsearch/
http://www.dongming8.cn/?page_id=156
http://www.qwolf.com/?p=1387#%E6%90%9C%E7%B4%A2


安装
apt-get install rabbitmq-server
启动图形管理界面,web访问 http://10.30.12.92:15672
rabbitmq-plugins enable rabbitmq_management
查看服务启动状态
systemctl status rabbitmq-server

查看交换机列表
rabbitmqctl list_exchanges
添加交换机
python3 rabbitmqadmin declare exchange name=apiserver type=fanout

查看用户列表
rabbitmqctl list_users
添加用户
rabbitmqctl add_user  admin 123456
设置用户为管理员
rabbitmqctl set_user_tags admin administrator
设置用户权限
rabbitmqctl set_permissions -p / test ".*" ".*" ".*"

查看所有队列
rabbitmqctl list_queues

连接服务api
amqp://admin:123456@10.30.12.92:5672

amqp包路径
go get github.com/streadway/amqp
amqp学习
https://blog.csdn.net/weixin_37641832/article/details/83270778
rabbitmq教程
https://studygolang.com/articles/17119



