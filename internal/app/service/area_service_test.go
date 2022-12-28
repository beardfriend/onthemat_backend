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

func (ts *AreaServiceTestSuite) TestParseBubjungDongExcel() {
	ts.Run("성공", func() {
		r, r1, _ := ts.areaService.ParseBubjungDongExcelData("/Users/sehun/Downloads/국토교통부.xlsx")
		fmt.Println(len(r1))
		fmt.Println(len(r))
	})
}

func TestAreaServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AreaServiceTestSuite))
}
