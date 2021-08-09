// sqlite3 with sqlx in memory example

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/mattn/go-sqlite3"
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
// 	db, err := sqlx.Open("sqlite3", ":memory:")

// 	if err != nil {
// 		log.Fatalf("cannot open an sqlite3 based database: %v", err)
// 		panic(err)
// 	} else {
// 		log.Println("connected to database")
// 	}

// 	defer db.Close()

// 	sqlstmt := `CREATE TABLE quiz(id integer NOT NULL PRIMARY KEY AUTOINCREMENT,name TEXT,description TEXT)`
// 	_, err = db.Exec(sqlstmt)

// 	if err != nil {
// 		fmt.Println("exec error")
// 		panic(err)
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

// // os.Remove("./data/1234/bogo.db")
// // os.MkdirAll("./data/1234", 0755)
// // os.Create("./data/1234/bogo.db")
