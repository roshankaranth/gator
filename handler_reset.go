package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("Succesfully cleared users table!\n")
	return nil
}
