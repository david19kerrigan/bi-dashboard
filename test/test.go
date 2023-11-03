package main

import (
	"encoding/json"
	"fmt"
)

type Row struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type Data struct {
	Rows []Row `json:"rows"`
}

func main() {
	jsonData := `{"rows":[{"id":"2","name":"","status":"","address":""},{"id":"2","name":"","status":"","address":""}]}`
	var data Data

	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Unpacked data:")
	for _, row := range data.Rows {
		fmt.Printf("ID: %s, Name: %s, Status: %s, Address: %s\n", row.ID, row.Name, row.Status, row.Address)
	}
}
