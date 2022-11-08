package service

import (
	"errors"

	"onthemat/pkg/openapi"
)

type AcademyService interface {
	VerifyBusinessMan(businessCode string) error
}

type academyService struct {
	businessAPI *openapi.BusinessMan
}

func NewAcademyService(businessAPI *openapi.BusinessMan) AcademyService {
	return &academyService{
		businessAPI: businessAPI,
	}
}

func (s *academyService) VerifyBusinessMan(businessCode string) error {
	resp := s.businessAPI.GetStatus(businessCode)
	if resp.StatusCode() != 200 {
		return errors.New("인증되지 않은 비즈니스 번호입니다.")
	}
	return nil
}
