package storage

type Storage interface {
	Update(metricName string, metricValue any)
	GetList() string
	GetValue(metricType, metricName string) (string, bool)
}
