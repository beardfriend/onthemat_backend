package service

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/xuri/excelize/v2"
)

type AreaService interface {
	ParseExcelData(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error)
	ParesExcelDataV2(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error)
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

type Data struct {
	SidoName    string
	SidoCode    string
	SigunguName string
	SigunguCode string
}

func (*areaService) ParesExcelDataV2(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.NumCPU())
	f, err := excelize.OpenFile(fileUrl)
	if err != nil {
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return
	}

	keys := make(map[string]bool)
	for _, row := range rows {
		var lastAddedSidoCode string
		var lastAddedSigunguCode string
		if len(row) > 7 && len(row[7]) > 1 {
			continue
		}

		var tempSigunguCode string
		var tempSigunguName string
		var tempSidoCode string
		var tempSidoName string
		for j, colName := range row {
			if j == 0 && colName != "" {
				tempSidoCode = colName[:2]
				tempSigunguCode = colName[:5]
			}

			if j == 1 && tempSidoCode != "" && colName != "" {
				tempSidoName = colName
			}

			if j == 2 && tempSigunguCode != "" && colName != "" {
				tempSigunguName = colName
			}

			if tempSidoCode != "" && tempSidoName != "" && lastAddedSidoCode != tempSidoCode {
				lastAddedSidoCode = tempSidoCode

				if _, ok := keys[tempSidoCode]; !ok {
					keys[tempSidoCode] = true
					SidoResult = append(SidoResult, Sido{
						SidoCode: tempSidoCode,
						SidoName: tempSidoName,
					})
				}

			}

			if tempSigunguCode != "" && tempSigunguName != "" && lastAddedSigunguCode != tempSigunguCode {
				lastAddedSigunguCode = tempSigunguCode

				if _, ok := keys[tempSigunguCode]; !ok {
					keys[tempSigunguCode] = true
					SigunguResult = append(SigunguResult, Sigungu{
						SigunguCode: tempSigunguCode,
						SigunguName: tempSigunguName,
					})

				}

			}
		}
	}
	return
}

func (*areaService) ParesExcelDataV2Async(fileUrl string) (SidoResult []Sido, SigunguResult []Sigungu, err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(runtime.NumCPU())
	f, err := excelize.OpenFile(fileUrl)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return
	}

	SheetIndexHash := make(map[int]struct {
		StratNum int
		EndNum   int
	}, 0)

	chuk := len(rows) / 9
	for i := 0; i <= 8; i++ {

		startNum := (i * chuk)
		endNum := ((i + 1) * chuk) - 1
		if i == 0 {
			startNum = 1
		}
		if i == 8 {
			endNum = len(rows)
		}
		SheetIndexHash[i] = struct {
			StratNum int
			EndNum   int
		}{
			StratNum: startNum,
			EndNum:   endNum,
		}
	}
	var wg sync.WaitGroup
	mutex := new(sync.Mutex)
	keys := make(map[string]bool)
	for _, val := range SheetIndexHash {
		wg.Add(1)
		go func(val struct {
			StratNum int
			EndNum   int
		},
		) {
			defer wg.Done()
			var lastAddedSidoCode string
			var lastAddedSigunguCode string

			for i := val.StratNum; i < val.EndNum; i++ {

				if len(rows[i]) > 7 && len(rows[i][7]) > 1 {
					continue
				}

				var tempSigunguCode string
				var tempSigunguName string
				var tempSidoCode string
				var tempSidoName string
				for j, colName := range rows[i] {

					if j == 0 && colName != "" {
						tempSidoCode = colName[:2]
						tempSigunguCode = colName[:5]
					}

					if j == 1 && tempSidoCode != "" && colName != "" {
						tempSidoName = colName
					}

					if j == 2 && tempSigunguCode != "" && colName != "" {
						tempSigunguName = colName
					}

					if tempSidoCode != "" && tempSidoName != "" && lastAddedSidoCode != tempSidoCode {
						lastAddedSidoCode = tempSidoCode
						mutex.Lock()
						if _, ok := keys[tempSidoCode]; !ok {

							keys[tempSidoCode] = true
							SidoResult = append(SidoResult, Sido{
								SidoCode: tempSidoCode,
								SidoName: tempSidoName,
							})

						}
						mutex.Unlock()

					}

					if tempSigunguCode != "" && tempSigunguName != "" && lastAddedSigunguCode != tempSigunguCode {
						lastAddedSigunguCode = tempSigunguCode
						mutex.Lock()
						if _, ok := keys[tempSigunguCode]; !ok {

							keys[tempSigunguCode] = true
							SigunguResult = append(SigunguResult, Sigungu{
								SigunguCode: tempSigunguCode,
								SigunguName: tempSigunguName,
							})

						}
						mutex.Unlock()

					}
				}
			}
		}(val)
	}

	wg.Wait()

	return
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
