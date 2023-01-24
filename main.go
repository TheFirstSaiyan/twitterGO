package main

import (
	"fmt"

	"go projects/user"
)

func main() {

	users := []user.User{{Name: "Hari", Email: "hari", Password: "password"},
		{Name: "Hari", Email: "hari", Password: "password"}}

	for _, value := range users {
		fmt.Println(value.email)
	}
}
