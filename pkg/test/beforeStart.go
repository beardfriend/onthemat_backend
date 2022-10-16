package test

import "onthemat/internal/app/config"

func BeforeStart() {
	config.Load()
}
