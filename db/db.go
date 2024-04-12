package db

import (
	"database/sql"
	"fmt"
	"lab1/models"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "dasd"
)

var db = connect()

func connect() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	db, err := sql.Open("postgres", connStr)
	CheckErr(err)

	fmt.Println("Connected!")
	return db
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func GetTasks() []models.Task {
	log.Println("Получаем все задачи")
	rows, err := db.Query("SELECT * FROM tasks")
	CheckErr(err)
	defer rows.Close()
	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		var priorityId string
		var statusId string
		err = rows.Scan(
			&task.Id,
			&task.Name,
			&task.Desc,
			&task.Enddate,
			&priorityId,
			&statusId,
		)
		CheckErr(err)
		task.Priority = getPriority(priorityId)
		task.Status = getStatus(statusId)
		task.Categories = getCategories(task.Id)
		tasks = append(tasks, task)
	}
	return tasks
}

func getCategories(taskId int) []models.Category {
	log.Println("Получаем все категории задачи с id =", taskId)
	rows, err := db.Query(fmt.Sprintf("SELECT C.* FROM categories C INNER JOIN tasktocategory TC ON TC.categoryid = C.id WHERE TC.taskid = %d", taskId))
	CheckErr(err)
	defer rows.Close()
	cats := make([]models.Category, 0)
	for rows.Next() {
		var cat models.Category
		err = rows.Scan(&cat.Id, &cat.Name)
		CheckErr(err)
		cats = append(cats, cat)
	}
	return cats
}

func getStatus(statusId string) models.Status {
	log.Println("Получаем статус с id =", statusId)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM statuses WHERE id = %s", statusId))
	CheckErr(err)
	defer rows.Close()
	var status models.Status
	if rows.Next() {
		err = rows.Scan(&status.Id, &status.Name)
		CheckErr(err)
	}
	return status
}

func getPriority(priorityId string) models.Priority {
	log.Println("Получаем приоритет с id =", priorityId)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM priorities WHERE id = %s", priorityId))
	CheckErr(err)
	defer rows.Close()
	var priority models.Priority
	if rows.Next() {
		err = rows.Scan(&priority.Id, &priority.Name)
		CheckErr(err)
	}
	return priority
}

func SetTask(task models.Task, categoriesStr string) {
	log.Println("Записываем таску в бд\ntask:", task, "\ncats:", categoriesStr)
	statusId := SetStatus(task.Status)
	priorityId := SetPriority(task.Priority)
	_, err := db.Exec("INSERT INTO tasks (name, description, enddate, priorityid, statusid) VALUES ($1, $2, $3, $4, $5)", task.Name, task.Desc, task.Enddate, priorityId, statusId)
	CheckErr(err)
	taskId := getInsertedObjectId("tasks")
	categories := strings.Split(categoriesStr, ",")
	if len(categories) == 0 {
		log.Println("Категории не были указаны или указаны неверно, categoriesStr =", categoriesStr)
	} else {
		updateTaskToCategory(categories, taskId)
	}
	CheckErr(err)

}

func updateTaskToCategory(categories []string, taskId int) {
	log.Println("Обновлем связи категорий и задач")
	for _, cat := range categories {
		var category models.Category
		category.Name = cat
		categoryId := setCategory(category)
		_, err := db.Exec("INSERT INTO tasktocategory (taskid, categoryid) VALUES ($1, $2)", taskId, categoryId)
		CheckErr(err)
	}
}

func setCategory(category models.Category) int {
	log.Println("Записываем категорию в бд если ее там нет, category =", category)
	var emptyCategory models.Category
	categoryFromDb := GetCategoryByName(category.Name)
	if emptyCategory == categoryFromDb {
		_, err := db.Exec("INSERT INTO categories (name) VALUES ($1)", category.Name)
		CheckErr(err)
		return getInsertedObjectId("categories")
	} else {
		log.Println("Категория уже есть в бд")
		return categoryFromDb.Id
	}
}

func GetCategoryByName(categoryName string) models.Category {
	log.Println("Получаем категорию по имени -", categoryName)
	var category models.Category
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM categories WHERE name = '%s'", categoryName))
	CheckErr(err)
	if rows.Next() {
		err = rows.Scan(&category.Id, &category.Name)
		CheckErr(err)
	}
	return category
}

func SetStatus(status models.Status) int {
	log.Println("Записываем статус в бд если его там нет, status =", status)
	var emptyStatus models.Status
	statusFromDb := GetStatusByName(status.Name)
	if emptyStatus == statusFromDb {
		_, err := db.Exec("INSERT INTO statuses (name) VALUES ($1)", status.Name)
		CheckErr(err)
		return getInsertedObjectId("statuses")
	} else {
		log.Println("Статус уже есть в бд")
		return statusFromDb.Id
	}
}

func GetStatusByName(statusName string) models.Status {
	log.Println("Получаем статус по имени -", statusName)
	var status models.Status
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM statuses WHERE name = '%s'", statusName))
	CheckErr(err)
	if rows.Next() {
		err = rows.Scan(&status.Id, &status.Name)
		CheckErr(err)
	}
	return status
}

func SetPriority(priority models.Priority) int {
	log.Println("Записываем приоритет в бд если его там нет, priority =", priority)
	var emptyPriority models.Priority
	priorityFromDb := GetPriorityByName(priority.Name)
	if emptyPriority == priorityFromDb {
		_, err := db.Exec("INSERT INTO priorities (name) VALUES ($1)", priority.Name)
		CheckErr(err)
		return getInsertedObjectId("priorities")
	} else {
		log.Println("Приоритет уже есть в бд")
		return priorityFromDb.Id
	}
}

func getInsertedObjectId(tableName string) int {
	log.Println("Получаем id только что добавленого объекта")
	var id int
	rows, err := db.Query("SELECT MAX(id) FROM " + tableName)
	CheckErr(err)
	if rows.Next() {
		err = rows.Scan(&id)
		CheckErr(err)
	}
	return id
}

func GetPriorityByName(priorityName string) models.Priority {
	log.Println("Получаем приоритет по имени -", priorityName)
	var priority models.Priority
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM statuses WHERE name = '%s'", priorityName))
	CheckErr(err)
	if rows.Next() {
		err = rows.Scan(&priority.Id, &priority.Name)
		CheckErr(err)
	}
	return priority
}
