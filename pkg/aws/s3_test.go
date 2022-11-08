package aws

import (
	"fmt"
	"os"
	"testing"

	"onthemat/internal/app/config"
)

func TestS3(t *testing.T) {
	file, _ := os.Open("./test_object/akmu.jpeg")
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}
	s3Module := NewS3(c)
	s3Module.SetConfig()
	output := s3Module.Upload("test/akmu.jpeg", file)
	fmt.Println(output)
}
