package main

import (
    job "alfred-api/models"

    "fmt"
    "log"
    "net/http"
    "encoding/xml"
    "github.com/gorilla/mux"
)

type Response struct {
    Message string
}


func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func handleSms(w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Handling SMS!")

    res := Response{}
    res.Message = "Hello from Rohit @Twilio!"

    x, err := xml.MarshalIndent(res, "", "	")
    if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
    }

    w.Header().Set("Content-Type", "application/xml")
    w.Write(x)
}

func createJob(w http.ResponseWriter, r *http.Request){
    fmt.Println(w, "Creating Job")
    j := job.Job{}
    j.Data("test-job")
    j.PrintDetails()
}

func handleRequests() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homePage)
    router.HandleFunc("/sms", handleSms)
    router.HandleFunc("/jobs", createJob)
    log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
    handleRequests()
}
