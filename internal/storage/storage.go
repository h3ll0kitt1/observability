package storage

type Storage interface {
	Update(metricName string, metricValue any)
	GetList() string
	GetValue(mtype, name string) (string, bool)
}
