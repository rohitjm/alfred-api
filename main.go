package main

import (
    job "alfred-api/models"

    "fmt"
    "log"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/xml"
    "encoding/json"
    "github.com/gorilla/mux"
)

type Response struct {
    Message string
}

type NationalizeResponse struct {
	Name      string `json: "name"`
	Countries []struct {
		Id          string  `json:"country_id"`
		Probability float32 `json:"probability"`
	} `json:"country"`
}

type CountryCodeFile struct {
	CountryCodes []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country_codes"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func getNationalize(name string) NationalizeResponse {

    res, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
    if err != nil {
        panic(err.Error())
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err.Error())
    }

    var p NationalizeResponse

    err2 := json.Unmarshal([]byte(body), &p)
    if err2 != nil {
        panic(err2.Error())
    }
    fmt.Println(p.Countries[0].Id)
    fmt.Println(p.Countries[0].Probability)

    fmt.Println("Formatted Response: ", p.Countries)

    return p
}

func getCountryNameFromId(mapping CountryCodeFile, id string) string {

    for _, value := range mapping.CountryCodes {
        if(value.Code == id) {
            return value.Name
	}
    }

    return "Not Found"
}

func getReturnText(data NationalizeResponse) string {

    // read file
    rawdata, err := ioutil.ReadFile("./country_codes.json")
    if err != nil {
      fmt.Print(err)
    }

    var file CountryCodeFile

    // unmarshall it
    err = json.Unmarshal(rawdata, &file)
    if err != nil {
        fmt.Println("error:", err)
    }

    res := fmt.Sprintf("Hello %s your name has ties to the following countries:\n", data.Name)
    for _, value := range data.Countries {
        fmt.Println(value)
	probInt := int(value.Probability * 100)
	toAdd := fmt.Sprintf(" %s %% from %s.\n", strconv.Itoa(probInt), getCountryNameFromId(file, value.Id))
	res += toAdd
    }

    return res
}

func handleSms(w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Handling SMS!")

    body := r.FormValue("Body")

    // TODO: Input validation
    // 1. make sure single word
    // 2. Trim spaces

    nationalize_response := getNationalize(body)
    response_text := getReturnText(nationalize_response)

    fmt.Println(response_text)

    res := Response{}
    res.Message = response_text

    x, err := xml.MarshalIndent(res, "", "	")
    if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
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
