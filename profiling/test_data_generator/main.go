package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func NewPutGen() {
	fmt.Println("Insert record number:")
	var recNum int
	fmt.Scanf("%d", &recNum)

	//NEW PUT DATA
	newPutFile, err := os.OpenFile("../new_put_data.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening newPutFile:", err)
		return
	}
	defer newPutFile.Close()

	//Writing header
	csvLine := "id,description,value\n"
	_, err = newPutFile.WriteString(csvLine)
	if err != nil {
		fmt.Println("Error writing to newPutFile:", err)
		return
	}

	//PUT DATA FOR UPDATE
	updatePutFile, err := os.OpenFile("../update_put_data.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening updatePutFile:", err)
		return
	}
	defer newPutFile.Close()

	//Writing header
	csvLine = "id,description,value\n"
	_, err = updatePutFile.WriteString(csvLine)
	if err != nil {
		fmt.Println("Error writing to updatePutFile:", err)
		return
	}

	for i := 0; i < recNum; i++ {
		description := generateRandomString(15)
		value := rand.Float64() * math.Pow(10, float64(rand.Int()%4))
		csvLine := fmt.Sprintf("%d,\"%s\",%f\n", i, description, value)

		_, err = newPutFile.WriteString(csvLine)
		if err != nil {
			fmt.Println("Error writing to newPutFile:", err)
			return
		}

		//Repeat 2 times for update put
		description = generateRandomString(15)
		value = rand.Float64() * math.Pow(10, float64(rand.Int()%4))
		csvLine = fmt.Sprintf("%d,\"%s\",%f\n", i, description, value)

		_, err = updatePutFile.WriteString(csvLine)
		if err != nil {
			fmt.Println("Error writing to newPutFile:", err)
			return
		}

		description = generateRandomString(15)
		value = rand.Float64() * math.Pow(10, float64(rand.Int()%4))
		csvLine = fmt.Sprintf("%d,\"%s\",%f\n", i, description, value)

		_, err = updatePutFile.WriteString(csvLine)
		if err != nil {
			fmt.Println("Error writing to newPutFile:", err)
			return
		}
	}

	fmt.Println("Data generated!")

}

func main() {
	NewPutGen()
}
