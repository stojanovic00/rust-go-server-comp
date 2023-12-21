package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseSpacedLine(line string) []string {
	parts := strings.Split(line, " ")

	var filteredParts []string

	for _, part := range parts {
		if part != "" {
			filteredParts = append(filteredParts, part)
		}
	}

	return filteredParts
}

func getCpuStats() (float64, float64) {

	// Open the file
	file, err := os.Open("cpu.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line, skip first 3
	skipped := 0
	count := 0
	sum := 0.0
	maxCpu := 0.0
	for scanner.Scan() {
		if skipped != 3 {
			skipped++
			continue
		}
		line := scanner.Text()

		parts := parseSpacedLine(line)

		cpuUsageStr := parts[8]
		cpuUsage, _ := strconv.ParseFloat(cpuUsageStr, 32)

		count++
		sum += cpuUsage
		if cpuUsage > maxCpu {
			maxCpu = cpuUsage
		}
	}

	avgCpu := sum / float64(count)

	return avgCpu, maxCpu
}

func getMemStats() (float64, float64) {

	// Open the file
	file, err := os.Open("mem.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 0, 0
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line, skip first 3
	skipped := 0
	count := 0
	sum := 0.0
	maxMem := 0.0
	for scanner.Scan() {
		if skipped != 3 {
			skipped++
			continue
		}
		line := scanner.Text()

		parts := parseSpacedLine(line)

		memUsageStr := parts[8]
		memUsage, _ := strconv.ParseFloat(memUsageStr, 32)

		count++
		sum += memUsage
		if memUsage > maxMem {
			maxMem = memUsage
		}
	}

	avgMem := sum / float64(count)

	return avgMem, maxMem
}

type BenchStats struct {
	//s
	TotalTestTime float64
	//ms
	PerRequestTime float64
	//KBytes/sec
	TransferRateReceived float64
	//KBytes/sec
	TransferRateSent float64
	//Start of connection till start processing
	ConnectionLatency float64
	//Start processing till connection close
	ConnectionProcessingTime float64
}

func getBenchStats() *BenchStats {
	// Open the file
	file, err := os.Open("bench.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	stats := &BenchStats{
		PerRequestTime:   -1.0,
		TransferRateSent: 0,
	}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Time taken for tests") {
			parts := strings.Split(line, ":")
			valuePart := parts[1]
			valuePart = strings.TrimSpace(valuePart)
			strValue := strings.Split(valuePart, " ")[0]
			value, _ := strconv.ParseFloat(strValue, 64)
			stats.TotalTestTime = value
		}

		if strings.Contains(line, "Time per request") {
			if stats.PerRequestTime != -1 {
				continue
			}
			parts := strings.Split(line, ":")
			valuePart := parts[1]
			valuePart = strings.TrimSpace(valuePart)
			strValue := strings.Split(valuePart, " ")[0]
			value, _ := strconv.ParseFloat(strValue, 64)
			stats.PerRequestTime = value
		}

		if strings.Contains(line, "Transfer rate") {
			parts := strings.Split(line, ":")
			valuePart := parts[1]
			valuePart = strings.TrimSpace(valuePart)
			strValue := strings.Split(valuePart, " ")[0]
			value, _ := strconv.ParseFloat(strValue, 64)
			stats.TransferRateReceived = value
		}

		if strings.Contains(line, "kb/s sent") {
			valuePart := strings.TrimSpace(line)
			strValue := strings.Split(valuePart, " ")[0]
			value, _ := strconv.ParseFloat(strValue, 64)
			stats.TransferRateSent = value
		}

		if strings.Contains(line, "Connect:") {
			parts := parseSpacedLine(line)
			valuePart := parts[2]
			valuePart = strings.TrimSpace(valuePart)
			value, _ := strconv.ParseFloat(valuePart, 64)
			stats.ConnectionLatency = value
		}

		if strings.Contains(line, "Processing:") {
			parts := parseSpacedLine(line)
			valuePart := parts[2]
			valuePart = strings.TrimSpace(valuePart)
			value, _ := strconv.ParseFloat(valuePart, 64)
			stats.ConnectionProcessingTime = value
		}
	}
	return stats
}

func main() {
	lang := os.Args[1]
	pool_size, _ := strconv.Atoi(os.Args[2])
	requests, _ := strconv.Atoi(os.Args[3])
	connections, _ := strconv.Atoi(os.Args[4])

	avgCpu, maxCpu := getCpuStats()
	avgMem, maxMem := getMemStats()
	benchStats := getBenchStats()

	csvLine := fmt.Sprintf("%s,%d,%d,%d,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,\n",
		lang,
		pool_size,
		requests,
		connections,
		avgCpu,
		maxCpu,
		avgMem,
		maxMem,
		benchStats.TotalTestTime,
		benchStats.PerRequestTime,
		benchStats.TransferRateReceived,
		benchStats.TransferRateSent,
		benchStats.ConnectionLatency,
		benchStats.ConnectionProcessingTime,
	)

	csvFile, err := os.OpenFile("profiling.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	//If file empty add header
	fileStats, err := csvFile.Stat()
	if err != nil {
		panic(err)
	}

	if fileStats.Size() == 0 {
		// File is empty, write the header
		header := "lang,pool_size,requests,connections,avg_cpu[%],max_cpu[%],avg_mem[%],max_mem[%]," +
			"total_test_time[s],per_request_mean_time[ms],transfer_rate_rcvd[kB/s],transfer_rate_sent[kB/s],connection_latency[ms],connection_processing_time[ms]\n"
		if _, err = csvFile.WriteString(header); err != nil {
			panic(err)
		}
	}

	//Finally, write aggregated results
	if _, err = csvFile.WriteString(csvLine); err != nil {
		panic(err)
	}
}
