package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/tealeg/xlsx"
)

type Person struct {
	Id    int
	Name  string
	Email string
	City  string
	State string
}

func main() {

	db, err := sql.Open("postgres", "user=postgres password=123456 dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	})

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("document")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = c.SaveUploadedFile(file, filepath.Join("./", file.Filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}
		// read the excel file
		xl, err := xlsx.OpenFile(file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		//iterate over the rows in the sheet

		var persons []Person

		for _, row := range xl.Sheets[0].Rows {

			id, err := row.Cells[0].Int()
			if err != nil {
				// handle the error, e.g. by logging it and skipping this row
				log.Printf("Error converting cell to integer: %s", err)
				continue
			}
			// create a person struct from the row data
			person := Person{
				Id:    id,
				Name:  row.Cells[1].String(),
				Email: row.Cells[2].String(),
				City:  row.Cells[3].String(),
				State: row.Cells[4].String(),
			}

			var abc Person
			err = db.QueryRow("SELECT * FROM persons WHERE Name=$1", person.Name).Scan(&abc.Id, &abc.Name, &abc.Email, &abc.City, &abc.State)

			if err != nil {

				_, insertErr := db.Exec("INSERT INTO persons(id,name,email,city,state)VALUES($1,$2,$3,$4,$5)", person.Id, person.Name, person.Email, person.City, person.State)
				if insertErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
					return
				}
			}

			persons = append(persons, person)

			fmt.Println(persons)
		}

		err = os.Remove(file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: Store persons in database
		c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully"})
	})

	router.DELETE("/delete", func(c *gin.Context) {
		// read the Excel file from the request body
		file, err := c.FormFile("outputfile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = c.SaveUploadedFile(file, filepath.Join("./", file.Filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}
		// read the excel file
		xl, err := xlsx.OpenFile(file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// iterate over the rows in the sheet
		for _, row := range xl.Sheets[0].Rows {
			id, err := row.Cells[0].Int()
			if err != nil {
				// handle the error, e.g. by logging it and skipping this row
				log.Printf("Error converting cell to integer: %s", err)
				continue
			}

			// create a student struct from the row data
			person := Person{
				Id:    id,
				Name:  row.Cells[1].String(),
				Email: row.Cells[2].String(),
				City:  row.Cells[3].String(),
				State: row.Cells[4].String(),
			}

			// check if there are any matching students in the database
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM persons WHERE name=$1", person.Name).Scan(&count)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// delete the matching students from the database
			if count > 0 {
				_, err = db.Exec("DELETE FROM persons WHERE Name=$1", person.Name)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
	})

	router.Run(":8080")

}
