package file

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/h3ll0kitt1/observability/internal/models"
)

func TestFileStorage_GetList(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	metric := []byte(`{"id":"testCounter","type":"counter","delta":1}`)
	file.Write(metric)

	list := []models.MetricsWithValue{
		{
			ID:    "testCounter",
			MType: "counter",
			Delta: int64(1),
		},
	}

	tests := []struct {
		name string
		want []models.MetricsWithValue
	}{
		{
			name: "get list",
			want: list,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			fs := NewStorage(file.Name())

			got, _ := fs.GetList(ctx)
			if !reflect.DeepEqual(got, list) {
				t.Errorf("GetList() = %v, want %v ", got, list)
			}
		})
	}
}

func TestFileStorage_UpdateList(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	tests := []struct {
		name string
		want string
	}{
		{
			name: "list updated",
			want: `{"id":"testCounter","type":"counter","delta":1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			fs := NewStorage(file.Name())

			list := []models.MetricsWithValue{
				{
					ID:    "testCounter",
					MType: "counter",
					Delta: int64(1),
				},
			}
			fs.UpdateList(ctx, list)

			got, _ := os.ReadFile(file.Name())
			got = got[:len(got)-1]

			if string(got) != tt.want {
				t.Errorf("UpdateList() = %v, want %v ", string(got), tt.want)
			}
		})
	}
}
