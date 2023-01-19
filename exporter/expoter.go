package exporter

type Input struct {
	Name string
	Type any
}

type Exporter struct {
}

func (e *Exporter) AddHandler(method string, path string, input Filed, output Filed) {

}

func test() {
	
	e := Exporter{}
	e.AddHandler("POST", "/download", Filed{
		Name: "downloadInput",
		Type: Array,
		Fields: []*Filed{
			{
				Name: "ok",
				Type: Object,
			},
		},
	}, Filed{
		Name: "Out",
	})
}
