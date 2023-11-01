package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
	"net/http"
)

type Patient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	Address string `json:"address"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/getData", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ "data": getData() })
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
		err = rows.Scan(&p.ID, &p.Name, &p.Status, &p.Address)
		if err != nil {
			fmt.Println("error ", err)
		}
		result = append(result, p)
		fmt.Println("results ", result)
	}
	return &result
}
