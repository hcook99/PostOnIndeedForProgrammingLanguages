package main

import (
	"net/http"
	"html/template"
	"log"
	"os"
	"fmt"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
	  return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
  }

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
	addr,_ := determineListenAddress()
	http.HandleFunc("/", handleStand)
	http.HandleFunc("/input_box",handleInput)

	log.Fatal(http.ListenAndServe(addr, nil))
}