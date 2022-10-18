package test

import "onthemat/internal/app/config"

func BeforeStart(filePath string) {
	config.Load(filePath)
}
