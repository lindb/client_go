package main

import (
	"time"

	"github.com/lindb/client_go/client"
	"github.com/lindb/lindb/rpc/proto/field"
)

func main() {
	// todo create cluster, database in lindb before sending metrics.
	cli, err := client.NewClientFromConfigPath("lindb_client.toml")

	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		ml := buildMetricList(float64(i))
		if err := cli.WriteList(ml); err != nil {
			panic(err)
		}
	}

	if err := cli.Close(); err != nil {
		panic(err)
	}

}

func buildMetricList(value float64) *field.MetricList {
	return &field.MetricList{Database: "dal",
		Metrics: []*field.Metric{{
			Name:      "name",
			Timestamp: time.Now().Unix() * 1000,
			Tags:      map[string]string{"tagKey": "tagVal"},
			Fields: []*field.Field{{
				Name: "sum",
				Field: &field.Field_Sum{
					Sum: &field.Sum{
						Value: value,
					},
				},
			}},
		}}}
}
