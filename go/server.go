package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Patient struct {
	id      string `json:"id"`
	name    string `json:"name"`
	status  string `json:"status"`
	address string `json:"address"`
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

	r.GET("/getData", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getData()})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
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

func postData(c *gin.Context) {
	var patient Patient
	err := c.BindJSON(&patient)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("Patient ", patient.name)

	body, _ := ioutil.ReadAll(c.Request.Body)
    fmt.Println(string(body))

	db := setupDB()
	fmt.Println("DB pointer ", db)
	query := fmt.Sprintf("INSERT INTO finni.patients (id, name, status, address) VALUES (&s, &s, &s, &s);", patient.id, patient.name, patient.status, patient.address)
	fmt.Println("Query ", query)
	//rows, err := db.Query(query)
	//if err != nil {
	//	fmt.Println("error", err)
	//}
}

func getData() *[]Patient {
	db := setupDB()
	fmt.Println("DB pointer ", db)
	rows, err := db.Query("SELECT * FROM finni.patients")
	if err != nil {
		fmt.Println("error", err)
	}
	var result []Patient
	for rows.Next() {
		var p Patient
		err = rows.Scan(&p.id, &p.name, &p.status, &p.address)
		if err != nil {
			fmt.Println("error ", err)
		}
		result = append(result, p)
		fmt.Println("results ", result)
	}
	return &result
}
