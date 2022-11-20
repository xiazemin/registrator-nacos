package main

import (
	_ "github.com/xiazemin/registrator-nacos/nacos"
	_ "github.com/xiazemin/registrator/consul"
	_ "github.com/xiazemin/registrator/consulkv"
	_ "github.com/xiazemin/registrator/etcd"
	_ "github.com/xiazemin/registrator/skydns2"
	_ "github.com/xiazemin/registrator/zookeeper"
)
