package routers

import (
	"ipManager/controllers"
	"net/http"
)

// type Mux *http.ServerMux
var Mux = http.NewServeMux()

func Init() {
	Mux.Handle("/nsxt/subnets", &controllers.SubnetsController{})
	// 参数Id=和cidr=
	Mux.HandleFunc("/nsxt/subnets/view", controllers.SubnetView)
	// 参数namespace=${cluster}-${namespace}
	Mux.HandleFunc("/nsxt/namespace/subnetview", controllers.NsSubnetView)
}
