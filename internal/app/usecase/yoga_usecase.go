package usecase

import (
	"context"
	"time"

	"onthemat/internal/app/common"
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/model"
	r "onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/elasticx"
	"onthemat/pkg/ent"
)

type YogaUsecase interface {
	CreateGroup(ctx context.Context, req *request.YogaGroupCreateBody) (err error)
	GroupList(ctx context.Context, req *request.YogaGroupListQueries) (result []*ent.YogaGroup, pagination *utils.PagenationInfo, err error)
	PutGroup(ctx context.Context, req *request.YogaGroupUpdateBody, id int) (isUpdated bool, err error)
	PatchGroup(ctx context.Context, req *request.YogaGroupPatchBody, id int) error
	DeleteGroup(ctx context.Context, ids []int) (rowAffected int, err error)

	Create(ctx context.Context, req *request.YogaCreateBody) (err error)
	List(ctx context.Context, groupId int) ([]*ent.Yoga, error)
	Recomendation(ctx context.Context, name string) ([]elasticx.ElasticSearchListBody[model.ElasticYoga], error)
	Put(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (isUpdated bool, err error)
	Delete(ctx context.Context, yogaId int) error
	Patch(ctx context.Context, req *request.YogaPatchBody, yogaId int) (err error)

	CreateRaws(ctx context.Context, names []string, userId int, userType string) error
	UpdateRaws(ctx context.Context, names []string, userId int, userType string) error
	DeleteRawAll(ctx context.Context, userId int, userType string) (err error)
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
	_, err = u.yogaRepo.CreateGroup(ctx, &ent.YogaGroup{
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
	_, err = u.yogaRepo.ElasitcDelete(ctx, ids)
	if err != nil {
		return
	}
	return u.yogaRepo.DeleteGroups(ctx, ids)
}

func (u *yogaUseCase) PutGroup(ctx context.Context, req *request.YogaGroupUpdateBody, id int) (isUpdated bool, err error) {
	isExist, err := u.yogaRepo.ExistGroup(ctx, id)
	isUpdated = isExist

	if err != nil {
		return
	}

	data := &ent.YogaGroup{
		ID:          id,
		Category:    req.Category,
		CategoryEng: req.CategoryEng,
		Description: req.Description,
	}

	if !isExist {
		_, err = u.yogaRepo.CreateGroup(ctx, data)
	} else {
		err = u.yogaRepo.UpdateGroup(ctx, data)
	}

	return
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
	y, err := u.yogaRepo.Create(ctx, &ent.Yoga{
		NameKor:     req.NameKor,
		NameEng:     req.NameEng,
		Description: req.Description,
		Level:       req.Level,
		YogaGroupID: req.YogaGroupId,
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrYogaGroupDoesNotExist, nil)
			return
		}
		return
	}

	data := []*model.ElasticYoga{
		{
			Id:   y.ID,
			Name: y.NameKor,
		},
	}

	if y.NameEng != nil {
		data = append(data, &model.ElasticYoga{
			Id:   y.ID,
			Name: *y.NameEng,
		})
	}

	for _, v := range data {
		err = u.yogaRepo.ElasticCreate(ctx, v)
		if err != nil {
			return
		}
	}

	return
}

func (u *yogaUseCase) Put(ctx context.Context, req *request.YogaUpdateBody, yogaId int) (isUpdated bool, err error) {
	isExist, err := u.yogaRepo.Exist(ctx, yogaId)
	isUpdated = !isExist
	if err != nil {
		return
	}

	data := &ent.Yoga{
		ID:          yogaId,
		YogaGroupID: req.YogaGroupId,
		NameKor:     req.NameKor,
		NameEng:     req.NameEng,
		Description: req.Description,
		Level:       req.Level,
	}

	if !isExist {
		_, err = u.yogaRepo.Create(ctx, data)
		u.yogaRepo.ElasticCreate(ctx, &model.ElasticYoga{
			Id:   yogaId,
			Name: req.NameKor,
		})
		if req.NameEng != nil {
			u.yogaRepo.ElasticCreate(ctx, &model.ElasticYoga{
				Id:   yogaId,
				Name: *req.NameEng,
			})
		}

	} else {
		err = u.yogaRepo.Update(ctx, data)
		u.yogaRepo.ElasitcDelete(ctx, []int{yogaId})
		u.yogaRepo.ElasticCreate(ctx, &model.ElasticYoga{
			Id:   yogaId,
			Name: req.NameKor,
		})

		if req.NameEng != nil {
			u.yogaRepo.ElasticCreate(ctx, &model.ElasticYoga{
				Id:   yogaId,
				Name: *req.NameEng,
			})
		}

	}
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

func (u *yogaUseCase) Delete(ctx context.Context, yogaId int) (err error) {
	_, err = u.yogaRepo.ElasitcDelete(ctx, []int{yogaId})
	if err != nil {
		return
	}
	return u.yogaRepo.Delete(ctx, yogaId)
}

func (u *yogaUseCase) List(ctx context.Context, groupId int) ([]*ent.Yoga, error) {
	return u.yogaRepo.List(ctx, groupId)
}

func (u *yogaUseCase) Recomendation(ctx context.Context, name string) ([]elasticx.ElasticSearchListBody[model.ElasticYoga], error) {
	return u.yogaRepo.ElasticList(ctx, name)
}

// ------------------- YogaRaw -------------------

func (u *yogaUseCase) CreateRaws(ctx context.Context, names []string, userId int, userType string) (err error) {
	var academyId *int
	var teacherId *int
	if userType == *model.AcademyType.ToString() {
		id, err := u.academyRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		academyId = &id
	} else if userType == *model.TeacherType.ToString() {
		id, err := u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		teacherId = &id
	}

	data := make([]*ent.YogaRaw, 0)
	for _, v := range names {
		data = append(data, &ent.YogaRaw{Name: v, AcademyID: academyId, TeacherID: teacherId})
	}
	return u.yogaRepo.CreateRaws(ctx, data)
}

func (u *yogaUseCase) UpdateRaws(ctx context.Context, names []string, userId int, userType string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var academyId *int
	var teacherId *int
	if userType == *model.AcademyType.ToString() {
		id, err := u.academyRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		academyId = &id
	} else if userType == *model.TeacherType.ToString() {
		id, err := u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		teacherId = &id
	}

	data := make([]*ent.YogaRaw, 0)
	for _, v := range names {
		data = append(data, &ent.YogaRaw{Name: v, AcademyID: academyId, TeacherID: teacherId})
	}
	return u.yogaRepo.DeleteAndCreateRaws(ctx, data, academyId, teacherId)
}

func (u *yogaUseCase) DeleteRawAll(ctx context.Context, userId int, userType string) (err error) {
	var academyId *int
	var teacherId *int
	if userType == *model.AcademyType.ToString() {
		id, err := u.academyRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		academyId = &id
	} else if userType == *model.TeacherType.ToString() {
		id, err := u.teacherRepo.GetOnlyIdByUserId(ctx, userId)
		if err != nil {
			return err
		}
		teacherId = &id
	}

	return u.yogaRepo.DeleteRawsByTeacherIdOrAcademyId(ctx, academyId, teacherId)
}
