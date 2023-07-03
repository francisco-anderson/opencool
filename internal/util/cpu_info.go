package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetCpuTemp(cpuTempFile *string) int {
	dat, err := os.ReadFile(*cpuTempFile)
	if err != nil {
		log.Fatalf("Error on open file: %v. %v", cpuTempFile, err)
	}
	cpuTemp, err := strconv.Atoi(strings.ReplaceAll(string(dat), "\n", ""))
	if err != nil {
		log.Fatalf("Error on convert string to int. %v", err)
	}
	cpuTemp = cpuTemp / 1000
	return cpuTemp
}

var prevIdleTime, prevTotalTime uint64

func GetCpuUsage() float64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()[5:] // get rid of cpu plus 2 spaces
	file.Close()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	split := strings.Fields(firstLine)
	idleTime, _ := strconv.ParseUint(split[3], 10, 64)
	totalTime := uint64(0)
	for _, s := range split {
		u, _ := strconv.ParseUint(s, 10, 64)
		totalTime += u
	}
	deltaIdleTime := idleTime - prevIdleTime
	deltaTotalTime := totalTime - prevTotalTime
	cpuUsage := (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0
	prevIdleTime = idleTime
	prevTotalTime = totalTime
	return cpuUsage
}
