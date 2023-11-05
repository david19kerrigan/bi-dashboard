package main

import (
	//"io"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Patient struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type Row struct {
	RowNum string `json:"rowNum"`
}

type Rows struct {
	Rows []Patient `json:"rows"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/newData", newData)
	r.POST("/updateData", updateData)
	r.POST("/deleteRow", deleteRow)
	r.GET("/getData", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getData()})
	})

	r.Run()
}

const (
	DB_USER = "david"
	DB_NAME = "finni"
)

// DB set up
func setupDB() *sql.DB {
	db, err := sql.Open("postgres", "user=david host=localhost port=5432 dbname=template1 sslmode=disable")
	if err != nil {
		fmt.Println("error: ", err)
	}
	return db
}

func updateData(c *gin.Context) {
	//body, _ := io.ReadAll(c.Request.Body)
	//println(string(body))

	var row Rows
	err := c.BindJSON(&row)

	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Printf("Patient:  %+v \n", row)

	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer: ", db)

	for _, patient := range row.Rows {
		query := "UPDATE finni.patients SET id = $1, name = $2, status = $3, address = $4 WHERE id = $5"
		fmt.Println("Patient: %v", patient)
		_, err := db.Exec(query, patient.Id, patient.Name, patient.Status, patient.Address, patient.Id)
		if err != nil {
			fmt.Println("error: ", err)
		}
	}

}

func deleteRow(c *gin.Context) {
	var row Row
	err := c.BindJSON(&row)
	if err != nil {
		fmt.Println("error: ", err)
	}

	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer: ", db)
	query := "DELETE FROM finni.patients WHERE id = $1"
	conv, _ := strconv.Atoi(row.RowNum)
	fmt.Println("Rownum: ", conv)

	rows, err := db.Exec(query, conv)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("rows: ", rows)
}
func newData(c *gin.Context) {
	//body, _ := io.ReadAll(c.Request.Body)
	//println(string(body))

	var row Rows
	err := c.BindJSON(&row)

	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Printf("Patient: %+v \n", row)

	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer: ", db)

	for _, patient := range row.Rows {
		query := "INSERT INTO finni.patients (id, name, status, address) VALUES (DEFAULT, $1, $2, $3);"
		fmt.Println("Patient: %v", patient)
		_, err := db.Exec(query, patient.Name, patient.Status, patient.Address)
		if err != nil {
			fmt.Println("error: ", err)
		}
	}

}

func getData() *[]Patient {
	db := setupDB()
	defer db.Close()

	fmt.Println("DB pointer: ", db)
	rows, err := db.Query("SELECT * FROM finni.patients")
	if err != nil {
		fmt.Println("error: ", err)
	}
	var result []Patient
	for rows.Next() {
		var p Patient
		err = rows.Scan(&p.Id, &p.Name, &p.Status, &p.Address)
		if err != nil {
			fmt.Println("error: ", err)
		}
		result = append(result, p)
		fmt.Println("results: ", result)
	}
	return &result
}
