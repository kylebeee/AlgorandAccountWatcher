package server

type Exporter interface {
	Save() error
	Close() error
}

type DefaultExporter struct {
}

func (de *DefaultExporter) Save() error {
	return nil
}

func (de *DefaultExporter) Close() error {
	return nil
}
