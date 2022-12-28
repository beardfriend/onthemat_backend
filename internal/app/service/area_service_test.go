package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AreaServiceTestSuite struct {
	suite.Suite
	areaService AreaService
}

func (ts *AreaServiceTestSuite) SetupSuite() {
	ts.areaService = NewAreaService()
}

func (ts *AreaServiceTestSuite) TestParseExcel() {
	ts.Run("성공", func() {
		sido, sigungu, err := ts.areaService.ParseExcelData("/Users/sehun/Downloads/행정.xlsx")
		ts.Equal(261, len(sigungu))
		ts.Equal(17, len(sido))
		ts.NoError(err)
		fmt.Println(sido)
		fmt.Println(sigungu)
		ts.GreaterOrEqual(len(sido), 1)
		ts.GreaterOrEqual(len(sigungu), 1)
	})
}

// 성능 측정 결과 :
// 엑셀 20000개 정도의 숫자에서는 비동기 처리해도 속도가 똑같다.
// 대략 20000 * 8 개 정도 되어야
// 1초 정도 차이가 벌어진다 이상.

func TestAreaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AreaServiceTestSuite))
}
