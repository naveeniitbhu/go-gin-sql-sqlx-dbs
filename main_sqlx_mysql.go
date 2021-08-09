// // mysql with sqlx example
// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/jmoiron/sqlx"
// )

// type App struct {
// 	R  *gin.Engine
// 	Db *sqlx.DB
// }

// type Quiz struct {
// 	ID          int64  `json:"id,omitempty" db:"id"`
// 	Name        string `json:"name,omitempty" db:"name"`
// 	Description string `json:"description,omitempty" db:"description"`
// }

// type Id struct {
// 	ID int64 `json:"id,omitempty" db:"id"`
// }

// func main() {
// 	db, err := sqlx.Open("mysql", "root:root@(localhost:3306)/skuad001")

// 	if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("connected to database")
// 	}

// 	app := App{
// 		R:  gin.Default(),
// 		Db: db,
// 	}

// 	app.R.GET("/api/quiz/:quiz_id", app.GetQuizId)
// 	app.R.POST("/api/quiz/", app.PostQuizDetails)

// 	app.R.Run(":8080")
// }

// func (app *App) GetQuizId(c *gin.Context) {
// 	id := c.Param("quiz_id")

// 	// TO DO: Check if id is not int
// 	// 400 failure
// 	db := app.Db
// 	row, err := db.Queryx(`Select * From quiz WHERE id=?`, id)
// 	if err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}
// 	var quiz Quiz

// 	for row.Next() {
// 		fmt.Println("entry for")
// 		err := row.StructScan(&quiz)
// 		if err != nil {
// 			log.Println(err)
// 			c.JSON(http.StatusBadRequest, gin.H{})
// 			panic(err)
// 		} else {
// 			fmt.Println("INFO: successful", quiz)
// 			c.JSON(200, gin.H{
// 				"id":          quiz.ID,
// 				"name":        quiz.Name,
// 				"description": quiz.Description,
// 			})
// 		}
// 	}
// 	fmt.Println(quiz)
// 	if quiz.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{})
// 	}
// }

// func (app *App) PostQuizDetails(c *gin.Context) {
// 	var quiz Quiz
// 	var id int64
// 	db := app.Db

// 	if err := c.ShouldBindJSON(&quiz); err == nil {
// 		log.Println("INFO: json binding is done")
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{"json binding error": err.Error()})
// 		panic(err)
// 	}

// 	res, err := db.Exec(`INSERT INTO quiz(name, description) VALUES(?, ?)`, quiz.Name, quiz.Description)
// 	id, err = res.LastInsertId()

// 	if err == nil {
// 		fmt.Printf("INFO: Quiz details inserted with id:%d & name:%s & description:%s\n", id, quiz.Name, quiz.Description)
// 		c.JSON(http.StatusOK, gin.H{
// 			"id":          id,
// 			"name":        quiz.Name,
// 			"description": quiz.Description,
// 		})
// 	} else {
// 		log.Println(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": "Quiz details not inserted properly"})
// 		panic(err)
// 	}
// }

// func checkErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
