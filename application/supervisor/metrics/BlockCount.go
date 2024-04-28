package metrics

import "github.com/Aj002Th/BlockchainEmulator/application/comm"

type MM = comm.MM

func BlockCountMetric() Metrics {
	var cnt = 0
	updater := func(m *MM) {
		cnt += 1
	}
	getResult := func() Desc {
		return NewDescBuilder().GetDesc()
	}
	return NewMetrics("BlockCount", updater, getResult)
}

//----------------------------------[[ CPU度量 ]]-----------------------------------

func CpuMetric() Metrics {
	var each = make([]float64, 0)
	var total = 0.0
	updater := func(m *MM) {
		each = append(each, 1)
	}
	getResult := func() Desc {
		desc := NewDescBuilder().AddElem("each 1", "", 0).AddElem("each2", "", total).GetDesc()
		return (desc)
	}
	return NewMetrics("CpuUsage", updater, getResult)
}

//----------------------------------[[ TCL度量 ]]-----------------------------------

// func Tcl() Metrics {

// }

// ---------------------------------------------------------

func lambda[T any](expr T) func() T {
	return func() T { return expr }
}

type Metrics struct {
	OnBlockInfoMsg func(m *MM)
	GetResult      func() Desc
	GetName        func() string
}

func NewMetrics(name string, updater func(m *MM), getResult func() Desc) Metrics {
	return Metrics{OnBlockInfoMsg: updater, GetResult: getResult, GetName: lambda(name)}
}
