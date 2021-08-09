// sqlite3 with sql in memory example

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	R  *gin.Engine
	Db *sql.DB
}

type Quiz struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type Question struct {
	ID             int64  `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Options        string `json:"options" db:"options"`
	Correct_Option int64  `json:"correct_option" db:"correct_option"`
	Quiz           int64  `json:"quiz" db:"quiz"`
	Points         int64  `json:"points" db:"points"`
}

type Id struct {
	ID int64 `json:"id,omitempty" db:"id"`
}

func main() {
	// os.Remove(".memory.db")
	// os.Create("./memory.db")

	db, err := sql.Open("sqlite3", "./db/migrate-goose-sqlite/memory-sqlite.db")

	if err != nil {
		log.Fatalf("cannot open an sqlite3 based database: %v", err)
		panic(err)
	} else {
		log.Println("connected to database")
	}

	defer db.Close()

	// createtablestmt := `CREATE TABLE IF NOT EXISTS quiz(id integer NOT NULL PRIMARY KEY AUTOINCREMENT,name TEXT,description TEXT)`
	// _, err = db.Exec(createtablestmt)
	// // Exec gives sql.result that gives access to stmt metadata: lastinsertid and number of rows affected.

	// if err != nil {
	// 	log.Fatal(err.Error())
	// 	panic(err)
	// } else {
	// 	log.Println("INFO: Quiz table is created.")
	// }

	app := App{
		R:  gin.Default(),
		Db: db,
	}

	app.R.GET("/api/quiz/:quiz_id", app.GetQuizId)
	app.R.POST("/api/quiz/", app.PostQuizDetails)

	app.R.GET("/api/question/:question_id", app.GetQuestion)
	app.R.POST("/api/question/", app.PostQuestionDetails)

	app.R.GET("/api/quiz-questions/:quiz_id", app.GetAllQuestions)

	app.R.Run(":8080")
}

func (app *App) GetQuizId(c *gin.Context) {
	id := c.Param("quiz_id")

	db := app.Db
	var quiz Quiz

	err := db.QueryRow(`Select name,description From quiz WHERE id=?`, id).Scan(&quiz.Name, &quiz.Description)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	} else {
		log.Println("INFO: successful-get request for name and description")
		c.JSON(200, gin.H{
			"id":          id,
			"name":        quiz.Name,
			"description": quiz.Description,
		})
	}
}

func (app *App) PostQuizDetails(c *gin.Context) {
	var quiz Quiz
	var id int64
	db := app.Db

	if err := c.ShouldBindJSON(&quiz); err == nil {
		log.Println("INFO: json binding is done")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"json binding error": err.Error()})
		panic(err)
	}

	res, err := db.Exec(`INSERT INTO quiz(name, description) VALUES(?, ?)`, quiz.Name, quiz.Description)
	id, err = res.LastInsertId()

	if err == nil {
		fmt.Printf("INFO: Quiz details inserted with id:%d & name:%s & description:%s\n", id, quiz.Name, quiz.Description)
		c.JSON(http.StatusOK, gin.H{
			"id":          id,
			"name":        quiz.Name,
			"description": quiz.Description,
		})
	} else {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Quiz details not inserted properly"})
		panic(err)
	}
}

func (app *App) GetQuestion(c *gin.Context) {
	id := c.Param("question_id")

	db := app.Db
	var question Question

	err := db.QueryRow(`Select name, options, correct_option, quiz, points From questions WHERE id=?`, id).Scan(&question.Name, &question.Options, &question.Correct_Option, &question.Quiz, &question.Points)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	} else {
		log.Println("INFO: successful-get request for question")
		c.JSON(200, gin.H{
			"id":             id,
			"name":           question.Name,
			"options":        question.Options,
			"correct_option": question.Correct_Option,
			"quiz":           question.Quiz,
			"points":         question.Points,
		})
	}
}

func (app *App) PostQuestionDetails(c *gin.Context) {
	var question Question
	var id int64
	db := app.Db

	if err := c.ShouldBindJSON(&question); err == nil {
		log.Println("INFO: json binding of question is done")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"json binding error": err.Error()})
		panic(err)
	}

	res, err := db.Exec(`INSERT INTO questions(name, options, correct_option, quiz, points) VALUES(?, ?, ?, ?, ?)`, question.Name, question.Options, question.Correct_Option, question.Quiz, question.Points)
	id, err = res.LastInsertId()

	if err == nil {
		fmt.Printf("INFO: Quiz details inserted with id:%d & name:%s & options:%s\n", id, question.Name, question.Options)
		c.JSON(http.StatusOK, gin.H{
			"id":             id,
			"name":           question.Name,
			"options":        question.Options,
			"correct_option": question.Correct_Option,
			"quiz":           question.Quiz,
			"points":         question.Points,
		})
	} else {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Question details not inserted properly"})
		panic(err)
	}
}

func (app *App) GetAllQuestions(c *gin.Context) {
	quiz_id, err := strconv.ParseInt(c.Param("quiz_id"), 10, 64)
	if err != nil {
		log.Fatal(err.Error())
	}
	db := app.Db
	var (
		question Question
		quiz     Quiz
	)

	err = db.QueryRow(`Select name, description From quiz WHERE id=?`, quiz_id).Scan(&quiz.Name, &quiz.Description)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	} else {
		log.Println("INFO: quiz details gathered")
	}
	rows, err := db.Query(`Select id, name, options, correct_option, points From questions WHERE quiz=?`, quiz_id)
	question.Quiz = quiz_id
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	type keyval map[string]interface{}

	var keyvalslice []keyval

	for rows.Next() {

		err := rows.Scan(&question.ID, &question.Name, &question.Options, &question.Correct_Option, &question.Points)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println("224: ", question)

			keyvalslice = append(keyvalslice, map[string]interface{}{
				"id":             question.ID,
				"name":           question.Name,
				"options":        question.Options,
				"correct_option": question.Correct_Option,
				"quiz":           question.Quiz,
				"points":         question.Points,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"name":        quiz.Name,
		"description": quiz.Description,
		"questions":   keyvalslice,
	})
}
