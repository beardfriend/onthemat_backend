package service

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type AreaService interface {
	ParseExcelData(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error)
}

type areaService struct{}

func NewAreaService() AreaService {
	return &areaService{}
}

type Sido struct {
	SidoName string
	SidoCode string
}

type Sigungu struct {
	SigunguName string
	SigunguCode string
}

// 한국행정구역분류 항목표 엑셀파일 파싱
// http://kssc.kostat.go.kr/ksscNew_web/kssc/common/CommonBoardList.do?gubun=1&strCategoryNameCode=019&strBbsId=kascrr&categoryMenu=014
func (*areaService) ParseExcelData(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error) {
	f, err := excelize.OpenFile(fileUrl)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("1. 항목표")
	if err != nil {
		return
	}

	// 가장 마지막으로 배열에 등록한 행정구역 코드
	var lastAddedSidoCode string
	var lastAddedSigunguCode string

	for i, row := range rows {
		// 네 번째 줄 부터 데이터가 존재함.
		if i < 3 {
			continue
		}

		// 배열에 저장하기 전, 임시로 데이터를 담아놓음.
		var tempSidoCode string
		var tempSidoName string
		var tempSigunguCode string
		var tempSigunguName string
		for j, colCell := range row {
			// j (1, 시도 코드), (2,시도 이름), (3, 시군구 코드), (4, 시군구 이름)

			// 시도 코드 및 이름 저장
			if j == 1 && colCell != "" {
				tempSidoCode = colCell
			}
			if j == 2 && tempSidoCode != "" {
				tempSidoName = colCell
			}

			// 시군구 코드 및 이름 저장
			if j == 3 && colCell != "" {
				tempSigunguCode = colCell
			}
			if j == 4 && tempSigunguCode != "" {
				tempSigunguName = colCell
			}
		}

		if tempSidoCode != "" && lastAddedSidoCode != tempSidoCode {
			lastAddedSidoCode = tempSidoCode
			SidoResult = append(SidoResult, Sido{
				SidoCode: tempSidoCode,
				SidoName: tempSidoName,
			})
		}

		if tempSigunguCode != "" && lastAddedSigunguCode != tempSigunguCode {
			lastAddedSigunguCode = tempSigunguCode
			SigunguResult = append(SigunguResult, Sigungu{
				SigunguCode: tempSigunguCode,
				SigunguName: tempSigunguName,
			})
		}
	}
	return
}
