package usecase

import (
	"context"
	"fmt"

	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/pkg/ent"
)

type AreaUsecase interface {
	CreateSiDo(ctx context.Context, fileUrl string) (err error)
}

type areaUsecase struct {
	areaRepo    repository.AreaRepository
	areaService service.AreaService
}

func NewAreaUsecase(areaRepo repository.AreaRepository, areaService service.AreaService) AreaUsecase {
	return &areaUsecase{
		areaRepo:    areaRepo,
		areaService: areaService,
	}
}

func (a *areaUsecase) CreateSiDo(ctx context.Context, fileUrl string) (err error) {
	Sido, SiGungu, err := a.areaService.ParseExcelData(fileUrl)
	version := 1
	if err != nil {
		return
	}
	for _, v := range Sido {
		var sigungu []*ent.AreaSiGungu

		for id, j := range SiGungu {
			if j.SigunguCode[:2] == v.SidoCode {

				data := &ent.AreaSiGungu{
					ID:      id + 1,
					Name:    j.SigunguName,
					AdmCode: j.SigunguCode,
				}

				lastNumber := fmt.Sprintf("%c", j.SigunguCode[4])
				if lastNumber != "0" {
					pCode := fmt.Sprintf("%s0", j.SigunguCode[:4])
					data.ParentCode = &pCode
				}
				sigungu = append(sigungu, data)
			}
		}

		err = a.areaRepo.Create(ctx, &ent.AreaSiDo{
			Name:    v.SidoName,
			AdmCode: v.SidoCode,
			Version: int32(version),
		}, sigungu)
		if err != nil {
			return
		}
	}
	return
}
