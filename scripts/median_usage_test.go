package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	Statuses := []Stats{
		{
			Cpu: 1.0,
			Mem: 2.0,
		},
		{
			Cpu: 2.0,
			Mem: 2.0,
		},
		{
			Cpu: 3.0,
			Mem: 2.0,
		},
		{
			Cpu: 4.0,
			Mem: 2.0,
		},
		{
			Cpu: 5.0,
			Mem: 2.0,
		},
	}

	expectedStats := Stats{
		Cpu: 15.0,
		Mem: 10,
	}

	result := reduce(Statuses, sumStats, Stats{})

	assert.Equal(t, expectedStats, result)
}
