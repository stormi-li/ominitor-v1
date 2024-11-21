package ominitor

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/stormi-li/omiserd-v1"
)

type NodeManager struct {
	serverDiscover *omiserd.Discover
	webDiscover    *omiserd.Discover
	configDiscover *omiserd.Discover
	opts           *redis.Options
}

func NewManager(opts *redis.Options) *NodeManager {
	return &NodeManager{
		serverDiscover: omiserd.NewClient(opts, omiserd.Server).NewDiscover(),
		webDiscover:    omiserd.NewClient(opts, omiserd.Web).NewDiscover(),
		configDiscover: omiserd.NewClient(opts, omiserd.Config).NewDiscover(),
		opts:           opts,
	}
}

func (manager *NodeManager) GetServerNodes() map[string]map[string]map[string]string {
	return manager.serverDiscover.DiscoverAllServers()
}

func (manager *NodeManager) GetWebNodes() map[string]map[string]map[string]string {
	return manager.webDiscover.DiscoverAllServers()
}

func (manager *NodeManager) GetConfigNodes() map[string]map[string]map[string]string {
	return manager.configDiscover.DiscoverAllServers()
}

func (manager *NodeManager) Handler(w http.ResponseWriter, r *http.Request) {
	// 获取请求的路径并去掉开头的 '/'
	path := strings.TrimPrefix(r.URL.Path, "/")
	// 以 '/' 分割路径
	parts := strings.Split(path, "/")

	if parts[0] == command_GetWebNodes {
		w.Write([]byte(toJsonStr(manager.GetWebNodes())))
	}
	if parts[0] == command_GetServerNodes {
		w.Write([]byte(toJsonStr(manager.GetServerNodes())))
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
	var register *omiserd.Register
	defer func() {
		recover()
	}()
	if serverType == string(omiserd.Config) {
		register = omiserd.NewClient(nodeManager.opts, omiserd.Config).NewRegister(name, address)
	}
	if serverType == string(omiserd.Web) {
		register = omiserd.NewClient(nodeManager.opts, omiserd.Web).NewRegister(name, address)
	}
	if serverType == string(omiserd.Server) {
		register = omiserd.NewClient(nodeManager.opts, omiserd.Server).NewRegister(name, address)
	}
	register.SendMessage(omiserd.Command_update_weight, weight)
	register.Close()
}

func (nodeManager *NodeManager) getDetails(serverType, name, address string) map[string]string {
	var data map[string]string
	if serverType == string(omiserd.Config) {
		data = nodeManager.configDiscover.GetData(name, address)
	}
	if serverType == string(omiserd.Web) {
		data = nodeManager.webDiscover.GetData(name, address)
	}
	if serverType == string(omiserd.Server) {
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
