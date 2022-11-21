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
		ts.NoError(err)
		fmt.Println(sido)
		fmt.Println(sigungu)
		ts.GreaterOrEqual(len(sido), 1)
		ts.GreaterOrEqual(len(sigungu), 1)
	})
}

func TestAreaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AreaServiceTestSuite))
}
