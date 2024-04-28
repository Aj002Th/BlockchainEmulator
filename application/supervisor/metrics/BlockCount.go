package metrics

type UniversalMetricMsg struct {
}

type MM = UniversalMetricMsg

func BlockCountMetric() Metrics {
	var cnt = 0
	updater := func(m *MM) {
		cnt += 1
	}
	getResult := func() string {
		return "{}"
	}
	return NewMetrics("BlockCount", updater, getResult)
}

func CpuMetric() Metrics {
	var each = make([]float64, 0)
	var total = 0.0
	updater := func(m *MM) {
		each = append(each, 1)
	}
	getResult := func() string {
		desc := NewDescBuilder().AddElem("each 1", "", 0).AddElem("each2", "", total).GetDesc()
		return DescPrintJson(desc)
	}
	return NewMetrics("CpuUsage", updater, getResult)
}

func Tcl() Metrics {

}

// ---------------------------------------------------------

func lambda[T any](expr T) func() T {
	return func() T { return expr }
}

type Metrics struct {
	GetUpdater func(m *MM)
	GetResult  func() string
	GetName    func() string
}

func NewMetrics(name string, updater func(m *MM), getResult func() string) Metrics {
	return Metrics{GetUpdater: updater, GetResult: getResult, GetName: lambda(name)}
}
