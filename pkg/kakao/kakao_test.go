package kakao

import (
	"fmt"
	"testing"

	"onthemat/internal/app/config"
)

func TestKakao(t *testing.T) {
	c := config.NewConfig()
	kaka := NewKakao(c)
	d := kaka.Authorize()
	fmt.Println(d)
}
