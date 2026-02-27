package main

import (
	"fmt"

	"github.com/roshankaranth/gator/internal/config"
)

func main() {
	data, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}

	data.SetUser("roshan")

	data, err = config.Read()

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Printf("DB_URL : %s\nCurrent username : %s\n", data.Db_url, data.Current_user_name)

}
