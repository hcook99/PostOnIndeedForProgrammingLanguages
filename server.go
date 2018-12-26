package main

import (
	"net/http"
	"html/template"
	"log"
)

func handleStand(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "static/index.html")
}

func handleInput(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	formF := r.Form
	locationA := formF["urlEntry"]
	location := locationA[0]
	mapOfLanguages := getMap(location)
	if len(mapOfLanguages) == 0{
		http.ServeFile(w,r,"static/index.html")
	}
	t, err := template.ParseFiles("static/graph.html") 
	  if err != nil {
		log.Print("template parsing error: ", err)
	  }
    err = t.Execute(w, mapOfLanguages) 
	  if err != nil {
		log.Print("template executing error: ", err) 
	  }
   
}

func main() {
	http.HandleFunc("/", handleStand)
	http.HandleFunc("/input_box",handleInput)

	log.Fatal(http.ListenAndServe(":8080", nil))
}