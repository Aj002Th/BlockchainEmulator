package service

import (
	"github.com/Aj002Th/BlockchainEmulator/application/comm"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

type Service interface {
	Start()
	Gather(m *comm.MM)
}

func GetServiceByName(name string) Service {
	switch name {
	case "SVG":
		return &Cpu{}
	case "Disk":
		return &Mem{}
	default:
		panic("")
	}
}

var services []Service

func StartAll() {
	for _, v := range params.Enabled_Measures {
		s := GetServiceByName(v)
		services = append(services, s)
		s.Start()
	}
}

func OnCommit() {
	var m comm.MM
	for _, v := range services {
		v.Gather(&m)
	}
	comm.Send(&m)
}
