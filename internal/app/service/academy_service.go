package service

import (
	"encoding/json"
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

	var result struct {
		Data []struct {
			Bstt string `json:"b_stt"`
		} `json:"data"`
	}

	json.Unmarshal(resp.Body(), &result)

	if result.Data[0].Bstt == "" {
		return errors.New("사업자 번호를 다시 확인해주세요")
	}
	return nil
}
