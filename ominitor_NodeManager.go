package ominitor

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
	omiconst "github.com/stormi-li/omiserd-v1/omiserd_const"
	discover "github.com/stormi-li/omiserd-v1/omiserd_discover"
	register "github.com/stormi-li/omiserd-v1/omiserd_register"
)

type NodeManager struct {
	serverDiscover *discover.Discover
	webDiscover    *discover.Discover
	configDiscover *discover.Discover
	opts           *redis.Options
}

func NewManager(opts *redis.Options) *NodeManager {
	return &NodeManager{
		serverDiscover: omiserd.NewClient(opts, omiconst.Server).NewDiscover(),
		webDiscover:    omiserd.NewClient(opts, omiconst.Web).NewDiscover(),
		configDiscover: omiserd.NewClient(opts, omiconst.Config).NewDiscover(),
		opts:           opts,
	}
}

func (manager *NodeManager) GetNodes(discover *discover.Discover) map[string]map[string]map[string]string {
	addMap := discover.GetAll()
	res := map[string]map[string]map[string]string{}
	for name, addrs := range addMap {
		if res[name] == nil {
			res[name] = map[string]map[string]string{}
		}
		for _, addr := range addrs {
			data := discover.GetData(name, addr)
			res[name][addr] = data
		}
	}
	return res
}

func (manager *NodeManager) GetServerNodes() map[string]map[string]map[string]string {
	return manager.GetNodes(manager.serverDiscover)
}

func (manager *NodeManager) GetWebNodes() map[string]map[string]map[string]string {
	return manager.GetNodes(manager.webDiscover)
}

func (manager *NodeManager) GetConfigNodes() map[string]map[string]map[string]string {
	return manager.GetNodes(manager.configDiscover)
}

func (manager *NodeManager) Handler(w http.ResponseWriter, r *http.Request) {
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	if parts[0] == command_GetServerNodes {
		var data = manager.GetServerNodes()
		w.Write([]byte(toJsonStr(data)))
	}
	if parts[0] == command_GetWebNodes {
		var data = manager.GetWebNodes()
		w.Write([]byte(toJsonStr(data)))
	}
	if parts[0] == command_GetConfigNodes {
		w.Write([]byte(toJsonStr(manager.GetConfigNodes())))
	}
	if parts[0] == command_UpdateWeight {
		serverType := r.URL.Query().Get("type")
		name := r.URL.Query().Get("name")
		address := r.URL.Query().Get("address")
		weight := r.URL.Query().Get("weight")
		manager.updateWeight(serverType, name, address, weight)
	}
	if parts[0] == command_GetDetails {
		serverType := r.URL.Query().Get("type")
		name := r.URL.Query().Get("name")
		address := r.URL.Query().Get("address")
		manager.getDetails(serverType, name, address)
		w.Write([]byte(mapToJsonStr(manager.getDetails(serverType, name, address))))
	}
}

func (nodeManager *NodeManager) updateWeight(serverType, name, address, weight string) {
	var register *register.Register
	defer func() {
		recover()
	}()
	if serverType == string(omiconst.Config) {
		register = omiserd.NewClient(nodeManager.opts, omiconst.Config).NewRegister(name, address)
	}
	if serverType == string(omiconst.Web) {
		register = omiserd.NewClient(nodeManager.opts, omiconst.Web).NewRegister(name, address)
	}
	if serverType == string(omiconst.Server) {
		register = omiserd.NewClient(nodeManager.opts, omiconst.Server).NewRegister(name, address)
	}
	register.SendMessage(omiconst.Command_update_weight, weight)
	register.RedisClient.Close()
}

func (nodeManager *NodeManager) getDetails(serverType, name, address string) map[string]string {
	var data map[string]string
	if serverType == string(omiconst.Config) {
		data = nodeManager.configDiscover.GetData(name, address)
	}
	if serverType == string(omiconst.Web) {
		data = nodeManager.webDiscover.GetData(name, address)
	}
	if serverType == string(omiconst.Server) {
		data = nodeManager.serverDiscover.GetData(name, address)
	}
	return data
}

func toJsonStr(nodes map[string]map[string]map[string]string) string {
	res := [][]string{}
	for name, addresses := range nodes {
		for address, details := range addresses {
			weight := details["weight"]
			res = append(res, []string{name, address, weight})
		}
	}
	return sliceToJsonStr(res)
}

func sliceToJsonStr(data [][]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}

func mapToJsonStr(data map[string]string) string {
	jsonStr, _ := json.MarshalIndent(data, " ", "  ")
	return string(jsonStr)
}
