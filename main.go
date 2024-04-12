package main

import (
	"html/template"
	"lab1/db"
	"lab1/models"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/addtask", addtask)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/index.html")
	CheckErr(err)
	tasks := db.GetTasks()
	err = temp.Execute(w, tasks)
	CheckErr(err)
}

func addtask(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		CheckErr(err)
		var task models.Task
		task.Name = r.FormValue("name")
		task.Desc = r.FormValue("description")
		task.Enddate, err = time.Parse("2006-01-02", r.FormValue("enddate"))
		CheckErr(err)
		task.Status = models.Status{Name: r.FormValue("status")}
		task.Priority = models.Priority{Name: r.FormValue("priority")}
		categoriesStr := r.FormValue("categories")
		db.SetTask(task, categoriesStr)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
