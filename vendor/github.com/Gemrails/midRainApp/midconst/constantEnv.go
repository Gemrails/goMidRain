package midconst

const (
	CONST_SERVICES_NAME string = "DEPEND_SERVICE"
	CONST_NAMESPACE string = "TENANT_ID"
	CONST_CRICUIT string = "L7_JSON"
	CONST_DOMAIN string = "L7_DOMAIN"
	CONST_IS_HTTP string = "IS_HTTP"
	ENVOY_INIT_CONF_PATH string = "/etc"
	ENVOY_RUN_CONF_PATH string = "/opt"
	ENVOY_BIN string = "/usr/local/bin/envoy"
	READ_TIME_OUT int = 60
	//API_SERVER_URL string = "http://172.30.42.1:8080"
	//API_SERVER_URL string = "http://127.0.0.1:9000"
	API_SERVER_URL string = "http://test.goodrain.com:8181"
	STREAMFORMAT string = "tcp://%s:%d"
	DOWNCONFIG string = "http://down.goodrain.me/kube-proxy.kubeconfig"
	Infourl string = "http://127.0.0.1:65534/server_info"
)

type WorkModel struct {
	Model int
	// 1 fot start model; 2 for running model
}

type EnvDeal interface{
	getnamespace() string
	getcircuit() map[string]int
	getDomain() map[string]string
	getDependServices() []string
	getIsHttp() string
}

type PieceHttpRroutes struct {
	Timeout_ms int `json:"timeout_ms"`
	Prefix string `json:"prefix"`
	Cluster string `json:"cluster"`
}

type PieceHttpVirtualHost struct {
	Name string `json:"name"`
	Domains string `json:"domains"`
	Routes []*PieceHttpRroutes `json:"routes"`
}

type RouteConfig struct {
	Virtual_hosts []*PieceHttpVirtualHost `json:"virtual_hosts"`
}

type HttpSingleFileter struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Config map[string]string `json:"config"`
}

type HttpConfig struct {
	Codec_type string `json:"codec_type"`
	Stat_prefix string `json:"stat_prefix"`
	Route_config *RouteConfig `json:"route_config"`
	Filters []*HttpSingleFileter `json:"filters"`
}

type PieceHttpFilters struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Config *HttpConfig `json:"config"`
}

type PieceTcpRoute struct {
	Cluster string `json:"cluster"`
}

type TcpRoutes struct {
	Routes []*PieceTcpRoute `json:"routes"`
}

type TcpConfig struct {
	Stat_prefix string `json:"stat_prefix"`
	Route_config *TcpRoutes `json:"route_config"`
}

type PieceTcpFilters struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Config *TcpConfig `json:"config"`
}

type PieceHttpListeners struct {
	Address string `json:"address"`
	Filters []*PieceHttpFilters `json:"filters"`
}

type PieceTcpListeners struct {
	Address string `json:"address"`
	Filters []*PieceTcpFilters `json:"filters"`
}

type PieceHosts struct {
	Url string `json:"url"`
}

type Hosts struct {
	Hosts []*PieceHosts `json:"hosts"`
}

type MaxConnections struct {
	Max_connections int `json:"max_connections"`
}

type CircuitBreakers struct {
	Default *MaxConnections `json:"default"`
}

type PieceClusters struct {
	Name string `json:"name"`
	Connect_timeout_ms int `json:"connect_timeout_ms"`
	Type string `json:"type"`
	Lb_type string `json:"lb_type"`
	Service_name string `json:"service_name"`
	Circuit_breakers *CircuitBreakers `json:"circuit_breakers"`
	Hosts interface{} `json:"hosts"`
}

type Listeners struct {
	Listeners []interface{} `json:"listeners"`
}

type Admin struct {
	Access_log_path string `json:"access_log_path"`
	Address string `json:"address"`
}

type RunTime struct {
	Symlink_root string `json:"symlink_root"`
	Subdirectory string `json:"subdirectory"`
	Override_subdirectory string `json:"override_subdirectory"`
}

type ClusterManager struct {
	Clusters []*PieceClusters `json:"clusters"`
}

type AllConfig struct {
	Listeners interface{} `json:"listeners"`
	Admin *Admin `json:"admin"`
	Flags_path string `json:"flags_path"`
	Runtime *RunTime `json:"runtime"`
	Cluster_manager *ClusterManager `json:"cluster_manager"`
}