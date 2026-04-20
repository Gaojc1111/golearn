package main

import (
	"fmt"
	"proto-cluster-tutorial/smartbulb"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {

	// 创建标准Proto.Actor
	system := actor.NewActorSystem()

	// 提供服务发现，集群节点管理
	provider, _ := consul.New()

	// 负责grain的节点分配和查找
	lookup := disthash.New()

	remoteConfig := remote.Configure("localhost", 0) // 开启远程通信

	bulbKind := smartbulb.NewSmartBulbKind(func() smartbulb.SmartBulb { return &smartbulb.Bulb{} }, 0)
	clusterConfig := cluster.Configure("ProtoClusterTutorial", provider, lookup, remoteConfig,
		cluster.WithKinds(bulbKind), // 注册业务grain
	)

	// 创建集群对象 来和集群进行交互
	c := cluster.New(system, clusterConfig)

	c.StartMember()
	// 启动 gRPC 服务器（Remote）。
	// 连接 Consul 并把自己注册进去（Provider）。
	// 开始同步集群状态。
	// 此时，该节点已经准备好接收来自全球各地的 Grain 调用请求了。
	fmt.Println("Cluster member started. Press Enter to exit...")

	bulbClient := smartbulb.GetSmartBulbGrainClient(c, "kitchen")
	// if _, err := bulbClient.TurnOn(nil); err != nil {
	// 	panic(err)
	// }

	var count int
	for {
		if count > 30 {
			break
		}
		state, _ := bulbClient.GetState(nil)
		if !state.IsOn {
			bulbClient.TurnOn(nil)
			count++
			fmt.Println("Turning on")
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Scanln()
	c.Shutdown(true)
}
