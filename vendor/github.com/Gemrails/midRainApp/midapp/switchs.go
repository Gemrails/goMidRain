package midapp

import (
	"github.com/Gemrails/midRainApp/midconst"
	"fmt"
)

type EnvArgs struct {
	Namespace string
	DependServices []string
	Circuit map[string]int
	Domain map[string]string
	Ishttp string
}

func (eas *EnvArgs)Start(wm *midconst.WorkModel){
	fmt.Println("Start model.")
	clientset := KubehttpsConnection()
	if StartModelService(clientset, eas, wm){
		StartEv(wm)
	}
}

func (eas *EnvArgs)Run(wm *midconst.WorkModel){
	fmt.Println("Running model.")
	clientset := KubehttpsConnection()
	RunningModelService(clientset, eas, wm)
}