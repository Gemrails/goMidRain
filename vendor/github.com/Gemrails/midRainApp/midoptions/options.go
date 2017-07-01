package midoptions

import (
	"github.com/Gemrails/midRainApp/midconst"
	"github.com/Gemrails/midRainApp/midapp"
	"flag"
	"os"
	"strconv"
	"strings"
)

type NoDealEnv struct {
	NamespaceNo string
	DependServicesNo string
	CircuitNo string
	DomainNo string
	IsHttpNo string
}

func (n *NoDealEnv)getnamespace() string{
	return n.NamespaceNo
}

func (n *NoDealEnv)getcircuit() map[string]int{
	ccs_map := make(map[string]int)
	circuits := n.CircuitNo
	if circuits != "" {
		ccs_list := strings.Split(circuits, ",")
		for _, values := range ccs_list {
			grlist := strings.Split(values, ":")
			ccs_map[grlist[0]], _ = strconv.Atoi(grlist[1])
		}
	}else{
		ccs_map["None"] = 1024
	}
	return ccs_map
}

func (n *NoDealEnv)getDomain() map[string]string{
	dm_map := make(map[string]string)
	dm_map["None"] = "none.com"
	return dm_map
}

func (n *NoDealEnv)getDependServices()[]string{
	var dslist []string
	dss := n.DependServicesNo
	if dss != ""{
		dss_list := strings.Split(dss, ",")
		for _, values := range dss_list{
			grlist := strings.Split(values, ":")
			dslist = append(dslist, grlist[0])
		}
	}
	return dslist
}

func (n *NoDealEnv)getIsHttp()string{
	return n.IsHttpNo
}

type EnvDeal interface{
	getnamespace() string
	getcircuit() map[string]int
	getDomain() map[string]string
	getDependServices() []string
	getIsHttp() string
}

func GetArgs(ed EnvDeal) midapp.EnvArgs{
	var e midapp.EnvArgs
	e.Namespace = ed.getnamespace()
	e.Circuit = ed.getcircuit()
	e.Domain = ed.getDomain()
	e.Ishttp = ed.getIsHttp()
	e.DependServices = ed.getDependServices()
	return e
}

func InitOptions() midapp.EnvArgs{
	nde := &NoDealEnv{
		NamespaceNo:os.Getenv(midconst.CONST_NAMESPACE),
		DependServicesNo:os.Getenv(midconst.CONST_SERVICES_NAME),
		CircuitNo:os.Getenv(midconst.CONST_CRICUIT),
		DomainNo:os.Getenv(midconst.CONST_DOMAIN),
		IsHttpNo:os.Getenv(midconst.CONST_IS_HTTP),
	}
	e := GetArgs(nde)
	return e
}

func Options(){
	r := flag.Bool("r", false, "Running model.Watching api-server's changes.")
	s := flag.Bool("s", false, "Start model.Use this arg when the first start.")
	flag.Parse()

	if *s {
		wm := &midconst.WorkModel{
			Model: 1,
		}
		eas := InitOptions()
		eas.Start(wm)
	} else if *r {
		wm := &midconst.WorkModel{
			Model: 2,
		}
		eas := InitOptions()
		eas.Run(wm)
	}
}
