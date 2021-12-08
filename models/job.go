package models

import "fmt"

func init() {
    fmt.Println("Father package initialized")
}

// struct definition
type Job struct {
    name string
}

func (f *Job) Data(name string) {
    f.name = name
    fmt.Println("Name set to: ", f.name)
}

func (s Job) PrintDetails() {
    fmt.Println("Job Details\n---------------")
    fmt.Println("Name :", s.name)
}
