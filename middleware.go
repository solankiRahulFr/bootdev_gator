package main

import (
	"fmt"
)

func middlewareLoggedIn(handler func(*state, command) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CurrentUser == "" {
			return fmt.Errorf("you must be logged in to run this command")
		}
		return handler(s, cmd)
	}
}	