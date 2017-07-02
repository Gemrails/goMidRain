package midapp

import (
	"github.com/Gemrails/midRainApp/midconst"
	"fmt"
	"os"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	_ "time"
	"log"
	"encoding/json"
	"bufio"
	"os/exec"
	"flag"
	"os/signal"
	"strconv"
	"net/http"
	"io/ioutil"
	"strings"
)

type ServiceInfo struct {
	ipstreams []string
	port string
	toport string
	protocol string
}

type DependServicesInfo struct {
	serviceInfo []ServiceInfo
	circuit int
}

type DependServices struct {
	infos map[string]DependServicesInfo
}

func logf(format string, a ...interface{}) {
	log.Printf("skydns: "+format, a...)
}

func KubehttpConnection() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.UnClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func KubehttpsConnection() *kubernetes.Clientset{
	kubeconfig := flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}else{
		return clientset
	}
}

func (ds *DependServices) InitServiceInfo(clientset *kubernetes.Clientset, namespace string, service_name string, ccmap map[string]int, c chan string) {
	labelname := "name" + "=" + service_name + "Service"
	endpoint, _ := clientset.CoreV1().Endpoints(namespace).List(metav1.ListOptions{LabelSelector:labelname})
	services, _ := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector:labelname})
	if len(endpoint.Items) != 0 {
		var dsi DependServicesInfo
		for key, item := range endpoint.Items {
			var si ServiceInfo
			addressList := item.Subsets[0].Addresses
			if len(addressList) == 0{
				addressList = item.Subsets[0].NotReadyAddresses
			}
			port := item.Subsets[0].Ports[0].Port
			toport := services.Items[key].Spec.Ports[0].Port
			si.protocol = "TCP"
			for _, ip := range addressList {
				stream := fmt.Sprintf(midconst.STREAMFORMAT, ip.IP, toport)
				si.ipstreams = append(si.ipstreams, stream)
				si.port = fmt.Sprintf("%d", port)
				si.toport = fmt.Sprintf("%d", toport)
			}
			dsi.serviceInfo = append(dsi.serviceInfo, si)
			if v, ok := ccmap[service_name]; ok {
				dsi.circuit = v
			}else{
				dsi.circuit = 1024
			}
		}
		ds.infos[service_name] = dsi
	}
	c <- "."
}

func (ds *DependServices) RunServiceInfo(clientset *kubernetes.Clientset, namespace string, c chan int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	listendpoints, err := clientset.CoreV1().Events(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	frontRV, _:= strconv.Atoi(listendpoints.ResourceVersion)

	watcher, err := clientset.CoreV1().Endpoints(namespace).Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for {
		select {
		case event := <- watcher.ResultChan():
		        if endpoint,ok:=event.Object.(*v1.Endpoints); ok{
				endpointRV,_ := strconv.Atoi(endpoint.ResourceVersion)
				if endpointRV > frontRV{
					c <- 1
				}
			}
		case  event:= <- signalChan:
			fmt.Printf("received ctrl+c(%v)\n", event)
			os.Exit(0)
		}
	}
}

func NewDependServices() *DependServices {
	var d DependServices
	return &d
}

func (ds *DependServices) createConfig(model *midconst.WorkModel) bool{
	var lis midconst.Listeners
	var clm midconst.ClusterManager
	admin := &midconst.Admin{
		Access_log_path: "/dev/null",
		Address: "tcp://0.0.0.0:65534",
	}
	rt := &midconst.RunTime{
		Symlink_root: "/srv/runtime_data/current",
		Subdirectory: "envoy",
		Override_subdirectory: "envoy_override",
	}
	for clustername, infos := range ds.infos {
		mcs := &midconst.MaxConnections{
			Max_connections: infos.circuit,
		}
		cbs := &midconst.CircuitBreakers{
			Default: mcs,
		}
		for _, info := range infos.serviceInfo{
			if info.port != "80" {

				pts := &midconst.PieceTcpRoute{
					Cluster: clustername+info.port,
				}
				trs := &midconst.TcpRoutes{
					Routes: []*midconst.PieceTcpRoute{
						pts,
					},
				}
				tc := &midconst.TcpConfig{
					Stat_prefix: clustername+info.port,
					Route_config: trs,
				}
				ptf := &midconst.PieceTcpFilters{
					Type: "read",
					Name: "tcp_proxy",
					Config: tc,
				}
				var hs midconst.Hosts
				for _, url := range info.ipstreams{
					phs := &midconst.PieceHosts{
						Url: url,
					}
					hs.Hosts = append(hs.Hosts, phs)

				}
				pcs := &midconst.PieceClusters{
					Name: clustername+info.port,
					Connect_timeout_ms: 250,
					Type: "static",
					Lb_type: "round_robin",
					Service_name: clustername+info.port,
					Circuit_breakers: cbs,
					Hosts: &hs.Hosts,
				}
				ptls := &midconst.PieceTcpListeners{
					Address: "tcp://0.0.0.0:" + info.port,
					Filters: []*midconst.PieceTcpFilters{
						ptf,
					},
				}
				clm.Clusters = append(clm.Clusters, pcs)
				lis.Listeners = append(lis.Listeners, ptls)
			}else{
				fmt.Println("create http config")
			}
		}
	}
	ac := &midconst.AllConfig{
		Listeners: &lis.Listeners,
		Admin: admin,
		Flags_path: "/etc/envoy/flags",
		Runtime: rt,
		Cluster_manager: &clm,
	}
	jsonAC, err := json.Marshal(ac);
	//fmt.Println(string(jsonAC))
	if err != nil {
		fmt.Println("Config init failed.")
		return false
	}
	if writeConfig(jsonAC, model){
		return true
	}else{
		fmt.Println("Write config err.")
		return false
	}
}

func writeConfig(jsonac []byte, model *midconst.WorkModel) bool{
	var confpath string
	if model.Model == 1{
		confpath = midconst.ENVOY_INIT_CONF_PATH
	}else{
		confpath = midconst.ENVOY_RUN_CONF_PATH
	}
	filename := fmt.Sprintf("%s/envoy_main.json", confpath)
	f, err := os.Create(filename)
	if err != nil{
		//fmt.Println("create configfile err.")
		return false
	}

	defer f.Close()
	_, err2 := f.Write(jsonac)
	if err2 != nil{
		fmt.Println("write configfile err.")
		return false
	}
	f.Sync()
	w := bufio.NewWriter(f)
	w.Flush()
	return true
}

func httpGet(c chan string) {
	resp, err := http.Get(midconst.Infourl)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	restart_epoch, _ := strconv.Atoi(strings.Split(string(body), " ")[5])
	c <- string(restart_epoch + 1)
}

func StartEv(model *midconst.WorkModel){
	if model.Model == 1{
		cmd := exec.Command(midconst.ENVOY_BIN, "-c", midconst.ENVOY_INIT_CONF_PATH+"/envoy_main.json")
		err := cmd.Start()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Start failed.")
		}
		fmt.Println("starting...")
	}else if model.Model == 2{
		c := make(chan string)
		go httpGet(c)
		restart_epoch := <- c
		cmd := exec.Command(midconst.ENVOY_BIN, "-c",
			midconst.ENVOY_RUN_CONF_PATH+"/envoy_main.json",
			"--restart-epoch",
			restart_epoch,
			"--parent-shutdown-time-s",
			"1")
		err := cmd.Start()
		if err != nil {
			fmt.Println("Restart failed.")
		}
		fmt.Println("restarting...")
	}
}

func StartModelService(clientset *kubernetes.Clientset, eas *EnvArgs, model *midconst.WorkModel)bool {
	/*
	DEPEND_SERVICE=grf2f1e2:ddcebdb0ce5454bba20fc95ddaf2f1e2,grdddd11:ddcebdb0ce5454bba20fc95ddaf2f1e2,
	*/
	c := make(chan string)
	ds := NewDependServices()
	ds.infos = make(map[string]DependServicesInfo)
	for _, dependservice := range eas.DependServices {
		go ds.InitServiceInfo(clientset, eas.Namespace, dependservice, eas.Circuit, c)
		val := <-c
		fmt.Println(val)
	}
	close(c)
	if ds.createConfig(model){
		return true
	}else{
		return false
	}
}

func RunningModelService(clientset *kubernetes.Clientset, eas *EnvArgs, model *midconst.WorkModel)bool {
	c := make(chan int)
	ds := NewDependServices()
	ds.infos = make(map[string]DependServicesInfo)
	go ds.RunServiceInfo(clientset, eas.Namespace, c)
	for val := range c{
		if 1 == val{
			if StartModelService(clientset, eas, model) {
				StartEv(model)
				fmt.Println("Restart success.")
				//return true
			}else{
				fmt.Println("Restart failed.")
				//return false
			}
		}
	}
	return true
}