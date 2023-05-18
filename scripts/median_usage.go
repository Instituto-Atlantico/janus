package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Instituto-Atlantico/janus/pkg/helper"
)

type Stats struct {
	Cpu float64
	Mem float64
}

func reduce(s []Stats, f func(a, b Stats) Stats, initValue Stats) Stats {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

func sumStats(a, b Stats) Stats {
	return Stats{Cpu: a.Cpu + b.Cpu, Mem: a.Mem + b.Mem}
}

func collectUsage(containerId string) Stats {
	command := fmt.Sprintf("docker stats --format json --no-stream %s", containerId)
	comm := helper.ParseCommand(command)
	cmd := exec.Command(comm[0], comm[1:]...)

	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	mapStats := make(map[string]string)
	json.Unmarshal(out, &mapStats)

	cpuMeasurement := strings.Split(mapStats["CPUPerc"], "%")[0]
	cpu, _ := strconv.ParseFloat(cpuMeasurement, 64)

	memMeasurement := strings.Split(mapStats["MemPerc"], "%")[0]
	mem, _ := strconv.ParseFloat(memMeasurement, 64)

	stats := Stats{Cpu: cpu, Mem: mem}
	return stats
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Container ID argument is empty")
	}
	containerId := os.Args[1]

	Statuses := []Stats{}

	startTime := time.Now()

	for {
		stats := collectUsage(containerId)

		Statuses = append(Statuses, stats)

		size := len(Statuses)

		//sum all values
		sum := reduce(Statuses, sumStats, Stats{})

		//divide by size
		sum.Cpu = sum.Cpu / float64(size)
		sum.Mem = sum.Mem / float64(size)

		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
		fmt.Println("Execution time:", time.Since(startTime))
		fmt.Printf("\nMedian Cpu: %v\nMedian Mem: %v", sum.Cpu, sum.Mem)

		time.Sleep(1 * time.Second)
	}
}
