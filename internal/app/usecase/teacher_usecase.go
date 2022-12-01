package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/pkg/ent"
)

type TeacherUsecase interface {
	Create(ctx context.Context, academy *request.TeacherCreateBody, userId int) error
	Update(ctx context.Context, d *request.TeacherUpdateBody, id, userId int) (isUpdated bool, err error)
	Patch(ctx context.Context, d *request.TeacherPatchBody, id, userId int) (isUpdated bool, err error)
	Get(ctx context.Context, id, userId int) (*ent.Teacher, error)
}

type teacherUseCase struct {
	teacherRepo repository.TeacherRepository
	userRepo    repository.UserRepository
}

func NewTeacherUsecase(
	teacherRepo repository.TeacherRepository,
	userRepo repository.UserRepository,
) TeacherUsecase {
	return &teacherUseCase{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
	}
}

func (u *teacherUseCase) Get(ctx context.Context, id, userId int) (result *ent.Teacher, err error) {
	result, err = u.teacherRepo.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {

			err = ex.NewNotFoundError(ex.ErrTeacherNotFound, nil)
			return
		}
		return
	}

	if !result.IsProfileOpen && result.UserID != userId {
		err = ex.NewForbiddenError(ex.ErrOnlyOwnUser, nil)
		return
	}

	return
}

func (u *teacherUseCase) Create(ctx context.Context, d *request.TeacherCreateBody, userId int) (err error) {
	// Already Exisit
	getUser, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return
	}

	// Check User Type
	if getUser.Type != nil {
		err = ex.NewConflictError(ex.ErrUserTypeAlreadyRegisted, nil)
		return
	}

	// Prepare Data
	info := d.TeacherInfo
	yoga := make([]*ent.Yoga, 0)
	for v := range d.YogaIds {
		yoga = append(yoga, &ent.Yoga{ID: v})
	}

	sigungu := make([]*ent.AreaSiGungu, 0)
	for v := range d.SigunguIds {
		sigungu = append(sigungu, &ent.AreaSiGungu{ID: v})
	}

	certifications := make([]*ent.TeacherCertification, 0)
	for _, v := range d.Certifications {
		certifications = append(certifications, &ent.TeacherCertification{
			AgencyName:   v.AgencyName,
			ImageUrl:     v.ImageUrl,
			ClassStartAt: v.ClassStartAt,
			ClassEndAt:   v.ClassEndAt,
			Description:  v.Description,
		})
	}

	workExperiences := make([]*ent.TeacherWorkExperience, 0)
	for _, v := range d.WorkExperiences {
		workExperiences = append(workExperiences, &ent.TeacherWorkExperience{
			AcademyName: v.AcademyName,
			WorkStartAt: v.WorkStartAt,
			WorkEndAt:   v.WorkEndAt,
			Description: v.Description,
		})
	}

	yogaRaws := make([]*ent.YogaRaw, 0)
	for _, v := range d.YogaRaws {
		yogaRaws = append(yogaRaws, &ent.YogaRaw{
			Name: v.Name,
		})
	}

	data := &ent.Teacher{
		UserID:          userId,
		ProfileImageUrl: info.ProfileImageURL,
		Age:             info.Age,
		Name:            info.Name,
		Introduce:       info.Introduce,
		IsProfileOpen:   info.IsProfileOpen,
		Edges: ent.TeacherEdges{
			Yoga:           yoga,
			YogaRaw:        yogaRaws,
			Sigungu:        sigungu,
			Certification:  certifications,
			WorkExperience: workExperiences,
		},
	}

	// Do
	err = u.teacherRepo.Create(ctx, data)

	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
			return
		}
		return
	}
	return
}

func (u *teacherUseCase) Patch(ctx context.Context, d *request.TeacherPatchBody, id, userId int) (isUpdated bool, err error) {
	isCreated, err := u.teacherRepo.Patch(ctx, d, id, userId)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
			return

		}
		return
	}
	isUpdated = !isCreated
	return
}

func (u *teacherUseCase) Update(ctx context.Context, d *request.TeacherUpdateBody, id, userId int) (isUpdated bool, err error) {
	// Prepare Data

	yoga := make([]*ent.Yoga, 0)
	for v := range d.YogaIds {
		yoga = append(yoga, &ent.Yoga{ID: v})
	}

	sigungu := make([]*ent.AreaSiGungu, 0)
	for v := range d.SigunguIds {
		sigungu = append(sigungu, &ent.AreaSiGungu{ID: v})
	}

	certifications := make([]*ent.TeacherCertification, 0)
	for _, v := range d.Certifications {
		certifications = append(certifications, &ent.TeacherCertification{
			ID:           v.Id,
			AgencyName:   v.AgencyName,
			ImageUrl:     v.ImageUrl,
			ClassStartAt: v.ClassStartAt,
			ClassEndAt:   v.ClassEndAt,
			Description:  v.Description,
		})
	}

	workExperiences := make([]*ent.TeacherWorkExperience, 0)
	for _, v := range d.WorkExperiences {
		workExperiences = append(workExperiences, &ent.TeacherWorkExperience{
			ID:          v.Id,
			AcademyName: v.AcademyName,
			WorkStartAt: v.WorkStartAt,
			WorkEndAt:   v.WorkEndAt,
			Description: v.Description,
		})
	}

	yogaRaws := make([]*ent.YogaRaw, 0)
	for _, v := range d.YogaRaws {
		yogaRaws = append(yogaRaws, &ent.YogaRaw{
			ID:   v.Id,
			Name: v.Name,
		})
	}
	info := d.TeacherInfo
	data := &ent.Teacher{
		ID:              id,
		UserID:          userId,
		ProfileImageUrl: info.ProfileImageURL,
		Age:             info.Age,
		Name:            info.Name,
		Introduce:       info.Introduce,
		IsProfileOpen:   info.IsProfileOpen,
		Edges: ent.TeacherEdges{
			Yoga:           yoga,
			YogaRaw:        yogaRaws,
			Sigungu:        sigungu,
			Certification:  certifications,
			WorkExperience: workExperiences,
		},
	}

	// Already Exisit
	isExist, err := u.teacherRepo.Exist(ctx, id)
	isUpdated = !isExist
	if err != nil {
		return
	}

	// Do
	if !isExist {
		err = u.teacherRepo.Create(ctx, data)
	} else {
		err = u.teacherRepo.Update(ctx, data)
	}

	if err != nil {
		if ent.IsConstraintError(err) {
			err = foreignKeyConstraintError(err)
			return
		}
		return
	}
	return
}
