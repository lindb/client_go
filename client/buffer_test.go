package client

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/rpc/proto/field"
)

func TestBuffer_In(t *testing.T) {
	buffer := NewBuffer("db", 1)

	metric := &field.Metric{
		Name: "name",
	}

	err := buffer.In(metric)
	if err != nil {
		t.Fatal(err)
	}

	err = buffer.In(&field.Metric{
		Name: "name",
	})
	if err != errBufferFull {
		t.Fatal(err)
	}

	batch := buffer.Out()

	assert.Equal(t, 1, len(batch))
	assert.Equal(t, metric, batch[0])
}

func TestBuffer_Out(t *testing.T) {
	size := 10
	buffer := NewBuffer("db", size)

	var err error

	for round := 0; round < 3; round++ {
		for i := 0; i < size; i++ {
			metric := &field.Metric{
				Name: "name",
				Tags: map[string]string{"key": strconv.Itoa(i)},
			}
			err = buffer.In(metric)
			if err != nil {
				t.Fatal(err)
			}
		}

		batch := buffer.Out()

		assert.Equal(t, size, len(batch))

		for i := 0; i < size; i++ {
			assert.Equal(t, strconv.Itoa(i), batch[i].Tags["key"])
		}
	}
}
