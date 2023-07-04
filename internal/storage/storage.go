package storage

type Storage interface {
	Update(metricName string, metricValue any)
}
