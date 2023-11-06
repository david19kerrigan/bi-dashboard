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

var DataColumns = []string{"id", "name", "status", "address"}

type Patient struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type Column struct {
	ColumnName string `json:"columnName"`
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
	r.POST("/addColumn", addColumn)
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

func addColumn(c *gin.Context) {
	var column Column
	err := c.BindJSON(&column)

	if err != nil {
		fmt.Println("error: ", err)
	}

	db := setupDB()
	defer db.Close()

	query := "ALTER TABLE finni.patients ADD COLUMN " + column.ColumnName + " VARCHAR DEFAULT '';"
	fmt.Println(query)
	_, err = db.Exec(query)

	if err != nil {
		fmt.Println("error: ", err)
	}
}

func updateData(c *gin.Context) {
	var data map[string]interface{}

	db := setupDB()
	defer db.Close()

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println("error: ", err)
	}

	for i := range data["rows"].([]interface{}) {
		thisRow := data["rows"].([]interface{})[i]
		fmt.Printf("THIS ROW: %v\n", thisRow)
		query := "UPDATE finni.patients SET "
		for j := range thisRow.(map[string]interface{}) {
			query += " " + j + " = '" + thisRow.(map[string]interface{})[j].(string) + "',"
		}
		query = query[:len(query)-1]
		query += " WHERE id = " + thisRow.(map[string]interface{})["id"].(string)
		fmt.Printf("QUERY: %s\n", query)
		_, err := db.Exec(query)
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

	rows, err := db.Exec(query, conv)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("rows: ", rows)
}

func newData(c *gin.Context) {
	var data map[string]interface{}

	db := setupDB()
	defer db.Close()

	err := c.BindJSON(&data)
	if err != nil {
		fmt.Println("error: ", err)
	}

	for i := range data["rows"].([]interface{}) {
		thisRow := data["rows"].([]interface{})[i]
		fmt.Printf("THIS ROW: %v\n", thisRow)
		query := "INSERT INTO finni.patients (id, "
		for j := range thisRow.(map[string]interface{}) {
			query += j + ","
		}
		query = query[:len(query)-1]
		query += ") VALUES (DEFAULT, "
		for j := range thisRow.(map[string]interface{}) {
			query += "'" + thisRow.(map[string]interface{})[j].(string) + "'" + ","
		}
		query = query[:len(query)-1]
		query += ")"

		_, err := db.Exec(query)
		if err != nil {
			fmt.Println("error: ", err)
		}
	}
}

func getData() *[]map[string]interface{} {
	db := setupDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM finni.patients")
	if err != nil {
		fmt.Println("error: ", err)
	}

	var columns []string
	columns, err = rows.Columns()
	lenColumns := len(columns)
	allRows := make([]map[string]interface{}, 0)

	for rows.Next() {
		colassoc := make(map[string]interface{}, lenColumns)
		cols := make([]interface{}, lenColumns)

		for i := 0; i < lenColumns; i++ {
			cols[i] = new(interface{})
		}

		err = rows.Scan(cols...)
		for i, col := range columns {
			colassoc[col] = *cols[i].(*interface{})
		}
		allRows = append(allRows, colassoc)

		if err != nil {
			fmt.Println("error: ", err)
		}
	}
	return &allRows
}
