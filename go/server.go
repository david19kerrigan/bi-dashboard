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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/postData", postData)
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
		fmt.Println("error", err)
	}
	return db
}

func deleteRow(c *gin.Context) {
	var row Row
	err := c.BindJSON(&row)
	if err != nil {
		fmt.Println("error", err)
	}
	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer ", db)
	query := "DELETE FROM finni.patients WHERE id = $1"
	conv, _ := strconv.Atoi(row.RowNum)
	fmt.Println(conv)
	rows, err := db.Exec(query, conv)
	if err != nil {
		fmt.Println("error ", err)
	}
	fmt.Println("rows ", rows)
}
func postData(c *gin.Context) {
	//body, _ := io.ReadAll(c.Request.Body)
	//println(string(body))

	var row Rows
	err := c.BindJSON(&row)

	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Printf("Patient %+v \n", row)

	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer ", db)

	for _, patient := range row.Rows {
		query := "INSERT INTO finni.patients (id, name, status, address) VALUES ($1, $2, $3, $4);"
		fmt.Println("%v", patient)
		_, err := db.Exec(query, patient.Id, patient.Name, patient.Status, patient.Address)
		if err != nil {
			fmt.Println("error ", err)
		}
	}

}

func getData() *[]Patient {
	db := setupDB()
	defer db.Close()
	fmt.Println("DB pointer ", db)
	rows, err := db.Query("SELECT * FROM finni.patients")
	if err != nil {
		fmt.Println("error", err)
	}
	var result []Patient
	for rows.Next() {
		var p Patient
		err = rows.Scan(&p.Id, &p.Name, &p.Status, &p.Address)
		if err != nil {
			fmt.Println("error ", err)
		}
		result = append(result, p)
		fmt.Println("results ", result)
	}
	return &result
}
