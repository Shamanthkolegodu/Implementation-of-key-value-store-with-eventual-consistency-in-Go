package main

import(
	"fmt"
	"time"
	"db/dbdesign"
)
var replication_factor = 5

func main()  {
	fmt.Println("Hello")
	fmt.Println(replication_factor)
	db, err := dbdesign.Create(":1234", "users", 1*time.Minute)
}

