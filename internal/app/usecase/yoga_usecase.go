package usecase

import (
	"context"

	"onthemat/internal/app/common"
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/model"
	r "onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
)

type YogaUsecase interface {
	CreateGroup(ctx context.Context, req *request.YogaGroupCreateBody) (err error)
	GroupList(ctx context.Context, req *request.YogaGroupListQueries) (result []*ent.YogaGroup, pagination *utils.PagenationInfo, err error)
	UpdateGroup(ctx context.Context, req *request.YogaGroupUpdateBody, yogaId int) error
	PatchGroup(ctx context.Context, req *request.YogaGroupPatchBody, id int) error
	DeleteGroup(ctx context.Context, ids []int) (rowAffected int, err error)

	Create(ctx context.Context, req *request.YogaCreateBody) (err error)
	List(ctx context.Context, groupId int) ([]*ent.Yoga, error)
	Update(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (err error)
	Delete(ctx context.Context, yogaId int) error
	Patch(ctx context.Context, req *request.YogaPatchBody, yogaId int) (err error)

	CreateRaws(ctx context.Context, names []string, userId int, userType model.UserType) error
	UpdateRaws(ctx context.Context, names []string, userId int, userType model.UserType) error
	DeleteRawAll(ctx context.Context, userId int, userType model.UserType) (err error)
}

type yogaUseCase struct {
	yogaRepo    r.YogaRepository
	academyRepo r.AcademyRepository
	teacherRepo r.TeacherRepository
}

func NewYogaUsecase(yogaRepo r.YogaRepository, academyRepo r.AcademyRepository, teacherRepo r.TeacherRepository) YogaUsecase {
	return &yogaUseCase{
		yogaRepo:    yogaRepo,
		academyRepo: academyRepo,
		teacherRepo: teacherRepo,
	}
}

// ------------------- Group -------------------

func (u *yogaUseCase) CreateGroup(ctx context.Context, req *request.YogaGroupCreateBody) (err error) {
	err = u.yogaRepo.CreateGroup(ctx, &ent.YogaGroup{
		Category:    req.Category,
		CategoryEng: req.CategoryEng,
		Description: req.Description,
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupAlreadyExist, nil)
			return
		}
		return
	}
	return
}

func (u *yogaUseCase) DeleteGroup(ctx context.Context, ids []int) (rowAffected int, err error) {
	return u.yogaRepo.DeleteGroups(ctx, ids)
}

func (u *yogaUseCase) UpdateGroup(ctx context.Context, req *request.YogaGroupUpdateBody, yogaId int) error {
	return u.yogaRepo.UpdateGroup(ctx, &ent.YogaGroup{
		ID:          yogaId,
		Category:    req.Category,
		CategoryEng: req.CategoryEng,
		Description: req.Description,
	})
}

func (u *yogaUseCase) PatchGroup(ctx context.Context, req *request.YogaGroupPatchBody, id int) error {
	return u.yogaRepo.PatchGroup(ctx, req, id)
}

func (u *yogaUseCase) GroupList(ctx context.Context, req *request.YogaGroupListQueries) (result []*ent.YogaGroup, pagination *utils.PagenationInfo, err error) {
	pgModule := utils.NewPagination(req.PageNo, req.PageSize)

	total, err := u.yogaRepo.GroupTotal(ctx, req.Category)
	if err != nil {
		return
	}
	pgModule.SetTotal(total)

	orderType := common.DESC
	if req.OrderType != nil && *req.OrderType == common.ASC {
		orderType = common.ASC
	}

	result, err = u.yogaRepo.GroupList(ctx, pgModule, req.Category, orderType)
	if err != nil {
		return
	}

	pagination = pgModule.GetInfo(len(result))
	return
}

// ------------------- Yoga -------------------

func (u *yogaUseCase) Create(ctx context.Context, req *request.YogaCreateBody) (err error) {
	err = u.yogaRepo.Create(ctx, &ent.Yoga{
		NameKor: req.NameKor,

		Description: req.Description,
		Level:       req.Level,
		Edges: ent.YogaEdges{
			YogaGroup: &ent.YogaGroup{
				ID: req.YogaGroupId,
			},
		},
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupDoesNotExist, nil)
			return
		}
		return
	}
	return
}

func (u *yogaUseCase) Update(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (err error) {
	err = u.yogaRepo.Update(ctx, &ent.Yoga{
		ID:          yogaId,
		YogaGroupID: req.YogaGroupId,
		NameKor:     req.NameKor,
		NameEng:     req.NameEng,
		Description: req.Description,
		Level:       req.Level,
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupDoesNotExist, nil)
			return
		}
		return
	}
	return
}

func (u *yogaUseCase) Patch(ctx context.Context, req *request.YogaPatchBody, yogaId int) (err error) {
	err = u.yogaRepo.Patch(ctx, req, yogaId)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupDoesNotExist, nil)
			return
		}
		return
	}
	return
}

func (u *yogaUseCase) Delete(ctx context.Context, yogaId int) error {
	return u.yogaRepo.Delete(ctx, yogaId)
}

func (u *yogaUseCase) List(ctx context.Context, groupId int) ([]*ent.Yoga, error) {
	return u.yogaRepo.List(ctx, groupId)
}

// ------------------- YogaRaw -------------------

func (u *yogaUseCase) CreateRaws(ctx context.Context, names []string, userId int, userType model.UserType) (err error) {
	var academyId int
	var teacherId int
	if userType == model.AcademyType {
		academyId, err = u.academyRepo.GetOnlyIdByUserId(ctx, userId)
	} else if userType == model.TeacherType {
		teacherId, err = u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
	}

	data := make([]*ent.YogaRaw, len(names))
	for _, v := range names {
		data = append(data, &ent.YogaRaw{Name: v, AcademyID: &academyId, TeacherID: &teacherId})
	}
	return u.yogaRepo.CreateRaws(ctx, data)
}

func (u *yogaUseCase) UpdateRaws(ctx context.Context, names []string, userId int, userType model.UserType) (err error) {
	var academyId int
	var teacherId int
	if userType == model.AcademyType {
		academyId, err = u.academyRepo.GetOnlyIdByUserId(ctx, userId)
	} else if userType == model.TeacherType {
		teacherId, err = u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
	}

	data := make([]*ent.YogaRaw, len(names))
	for _, v := range names {
		data = append(data, &ent.YogaRaw{Name: v, AcademyID: &academyId, TeacherID: &teacherId})
	}
	return u.yogaRepo.DeleteAndCreateRaws(ctx, data, &academyId, &teacherId)
}

func (u *yogaUseCase) DeleteRawAll(ctx context.Context, userId int, userType model.UserType) (err error) {
	var academyId int
	var teacherId int
	if userType == model.AcademyType {
		academyId, err = u.academyRepo.GetOnlyIdByUserId(ctx, userId)
	} else if userType == model.TeacherType {
		teacherId, err = u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
	}

	return u.yogaRepo.DeleteRawsByTeacherIdOrAcademyId(ctx, &academyId, &teacherId)
}
