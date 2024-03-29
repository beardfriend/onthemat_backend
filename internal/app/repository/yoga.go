package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"onthemat/internal/app/common"
	"onthemat/internal/app/model"
	"onthemat/internal/app/transport/request"

	"onthemat/internal/app/utils"
	"onthemat/pkg/elasticx"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/yoga"
	"onthemat/pkg/ent/yogagroup"
	"onthemat/pkg/ent/yogaraw"
	"onthemat/pkg/entx"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

type YogaRepository interface {
	CreateGroup(ctx context.Context, data *ent.YogaGroup) (*ent.YogaGroup, error)
	UpdateGroup(ctx context.Context, data *ent.YogaGroup) error
	PatchGroup(ctx context.Context, data *request.YogaGroupPatchBody, id int) error
	DeleteGroups(ctx context.Context, ids []int) (int, error)
	ExistGroup(ctx context.Context, id int) (bool, error)
	GroupTotal(ctx context.Context, category *string) (count int, err error)
	GroupList(ctx context.Context, pgModule *utils.Pagination, category *string, sorts common.Sorts) (result []*ent.YogaGroup, err error)

	Create(ctx context.Context, data *ent.Yoga) (*ent.Yoga, error)
	ElasticCreate(ctx context.Context, d *model.ElasticYoga) error
	ElasitcDelete(ctx context.Context, ids []int) (rowAfftected int, err error)
	//@@deprecated
	ElasticUpdate(ctx context.Context, data *model.ElasticYoga) (rowAfftected int, err error)
	Update(ctx context.Context, data *ent.Yoga) error
	Patch(ctx context.Context, data *request.YogaPatchBody, id int) error
	Delete(ctx context.Context, id int) error
	Exist(ctx context.Context, id int) (bool, error)
	List(ctx context.Context, groupdId int) ([]*ent.Yoga, error)
	ElasticList(ctx context.Context, name string) (resp []elasticx.ElasticSearchListBody[model.ElasticYoga], err error)

	CreateRaws(ctx context.Context, d []*ent.YogaRaw) error
	DeleteAndCreateRaws(ctx context.Context, d []*ent.YogaRaw, academyId, teacherId *int) (err error)
	DeleteRawsByTeacherIdOrAcademyId(ctx context.Context, academyId, teacherId *int) (err error)
}

type yogaRepository struct {
	db  *ent.Client
	ela *elasticsearch.Client
}

func NewYogaRepository(db *ent.Client, ela *elasticsearch.Client) YogaRepository {
	return &yogaRepository{
		db:  db,
		ela: ela,
	}
}

// -------------------  FOR Other Repositry Dependency -------------------

func (repo *yogaRepository) createRaws(ctx context.Context, db *ent.Client, d []*ent.YogaRaw, teacherId, academyId *int) error {
	bulk := make([]*ent.YogaRawCreate, len(d))

	for i, v := range d {
		clause := db.YogaRaw.Create().
			SetNillableAcademyID(academyId).
			SetNillableTeacherID(teacherId).
			SetName(v.Name)
		if v.ID != 0 {
			clause.SetID(v.ID)
		}
		bulk[i] = clause
	}

	return db.YogaRaw.CreateBulk(bulk...).Exec(ctx)
}

func (repo *yogaRepository) getRawIdsByTeacherId(ctx context.Context, db *ent.Client, teacherId int) ([]int, error) {
	return db.YogaRaw.Query().
		Where(yogaraw.TeacherIDEQ(teacherId)).
		IDs(ctx)
}

func (repo *yogaRepository) deleteRawsbyTeacherId(ctx context.Context, db *ent.Client, teacherId int) error {
	_, err := db.YogaRaw.Delete().Where(yogaraw.TeacherIDEQ(teacherId)).Exec(ctx)
	return err
}

func (repo *yogaRepository) updateRaws(ctx context.Context, db *ent.Client, value []*ent.YogaRaw, teacherId, academyId *int) (err error) {
	for _, v := range value {
		clause := db.YogaRaw.Update().
			SetNillableTeacherID(teacherId).
			SetNillableAcademyID(academyId).
			SetName(v.Name)

		mu := clause.Mutation()
		if v.AcademyID == nil {
			clause.Where(yogaraw.IDEQ(v.ID), yogaraw.TeacherIDEQ(*teacherId))
			mu.ClearAcademyID()
		}
		if v.TeacherID == nil {
			clause.Where(yogaraw.IDEQ(v.ID), yogaraw.AcademyIDEQ(*academyId))
			mu.ClearTeacherID()
		}

		err = clause.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return
}

func (repo *yogaRepository) deleteRawsByIds(ctx context.Context, db *ent.Client, ids []int, teacherId, academyId *int) (int, error) {
	clause := db.YogaRaw.Delete()

	if teacherId != nil {
		clause.Where(yogaraw.And(
			yogaraw.TeacherIDEQ(*teacherId),
			yogaraw.IDIn(ids...),
		))
	} else if academyId != nil {
		clause.Where(yogaraw.And(
			yogaraw.AcademyIDEQ(*academyId),
			yogaraw.IDIn(ids...),
		))
	}

	return clause.Exec(ctx)
}

// ------------------- Group -------------------

func (repo *yogaRepository) ExistGroup(ctx context.Context, id int) (bool, error) {
	return repo.db.YogaGroup.Query().Where(yogagroup.IDEQ(id)).Exist(ctx)
}

func (repo *yogaRepository) CreateGroup(ctx context.Context, data *ent.YogaGroup) (*ent.YogaGroup, error) {
	clause := repo.db.YogaGroup.Create().
		SetCategory(data.Category).
		SetCategoryEng(data.CategoryEng).
		SetNillableDescription(data.Description)

	if data.ID != 0 {
		clause.SetID(data.ID)
	}

	return clause.Save(ctx)
}

func (repo *yogaRepository) UpdateGroup(ctx context.Context, data *ent.YogaGroup) error {
	clause := repo.db.YogaGroup.Update().Where(yogagroup.IDEQ(data.ID))
	mu := clause.Mutation()

	if data.Description == nil {
		mu.ClearDescription()
	}

	return clause.
		SetCategory(data.Category).
		SetCategoryEng(data.CategoryEng).
		SetNillableDescription(data.Description).
		Exec(ctx)
}

func (repo *yogaRepository) PatchGroup(ctx context.Context, data *request.YogaGroupPatchBody, id int) error {
	updateableData := utils.GetUpdateableData(data, yogagroup.Columns)

	clause := repo.db.YogaGroup.
		Update().
		Where(yogagroup.IDEQ(id))

	for key, val := range updateableData {
		clause.Mutation().SetField(key, val)
	}

	return clause.Exec(ctx)
}

func (repo *yogaRepository) GroupTotal(ctx context.Context, category *string) (count int, err error) {
	clause := repo.db.YogaGroup.Query()
	clause = repo.groupConditionQuery(category, clause)

	count, err = clause.Count(ctx)
	return
}

func (repo *yogaRepository) GroupList(ctx context.Context, pgModule *utils.Pagination, category *string, sorts common.Sorts) (result []*ent.YogaGroup, err error) {
	clause := repo.db.YogaGroup.Query().
		Limit(pgModule.GetLimit()).
		Offset(pgModule.GetOffset())

	if sorts == common.ASC {
		clause = clause.Order(ent.Asc(yogagroup.FieldID))
	} else {
		clause = clause.Order(ent.Desc(yogagroup.FieldID))
	}

	clause = repo.groupConditionQuery(category, clause)

	result, err = clause.All(ctx)
	return
}

func (repo *yogaRepository) groupConditionQuery(category *string, clause *ent.YogaGroupQuery) *ent.YogaGroupQuery {
	if category != nil {
		clause.Where(yogagroup.CategoryContains(*category))
	}

	return clause
}

func (repo *yogaRepository) DeleteGroups(ctx context.Context, ids []int) (int, error) {
	return repo.db.YogaGroup.Delete().Where(yogagroup.IDIn(ids...)).Exec(ctx)
}

// ------------------- Yoga -------------------

func (repo *yogaRepository) Create(ctx context.Context, data *ent.Yoga) (*ent.Yoga, error) {
	clause := repo.db.Debug().Yoga.Create().
		SetNameKor(data.NameKor).
		SetNillableNameEng(data.NameEng).
		SetNillableLevel(data.Level).
		SetNillableDescription(data.Description).
		SetYogaGroupID(data.YogaGroupID)

	if data.ID != 0 {
		clause.SetID(data.ID)
	}

	return clause.Save(ctx)
}

func (repo *yogaRepository) ElasticCreate(ctx context.Context, d *model.ElasticYoga) error {
	uid := uuid.New().String()

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(d)

	_, err := repo.ela.Create(model.ElasticYogaIndexName, uid, &buf)
	if err != nil {
		return err
	}
	return err
}

func (repo *yogaRepository) ElasitcDelete(ctx context.Context, ids []int) (rowAfftected int, err error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"id": ids,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}
	resp, err := repo.ela.DeleteByQuery([]string{model.ElasticYogaIndexName}, &buf)
	if err != nil {
		return
	}
	var response elasticx.ElasticDeleteUpdateQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	rowAfftected = response.Deleted
	return
}

// @deprecated
func (repo *yogaRepository) ElasticUpdate(ctx context.Context, data *model.ElasticYoga) (rowAfftected int, err error) {
	var buf bytes.Buffer

	updateSource := fmt.Sprintf("ctx._source.name = '%s';", data.Name)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"name": data.Name,
						},
					},
					{
						"match": map[string]interface{}{
							"id": data.Id,
						},
					},
				},
			},
		},
		"script": map[string]interface{}{
			"source": updateSource,
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}

	resp, err := repo.ela.UpdateByQuery([]string{model.ElasticYogaIndexName}, func(r *esapi.UpdateByQueryRequest) {
		r.Body = &buf
	})
	if err != nil {
		return
	}
	var response elasticx.ElasticDeleteUpdateQueryResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	rowAfftected = response.Updated
	return
}

func (repo *yogaRepository) ElasticList(ctx context.Context, name string) (resp []elasticx.ElasticSearchListBody[model.ElasticYoga], err error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		panic(err)
	}

	res, err := esapi.SearchRequest{
		Index: []string{model.ElasticYogaIndexName},
		Body:  &buf,
	}.Do(ctx, repo.ela)
	if err != nil {
		panic(err)
	}

	var responseBody elasticx.ElasticSearchResponse[elasticx.ElasticSearchListBody[model.ElasticYoga]]
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		return
	}

	resp = responseBody.Hits.Hits
	return
}

func (repo *yogaRepository) Update(ctx context.Context, data *ent.Yoga) error {
	clause := repo.db.Yoga.Update().Where(yoga.IDEQ(data.ID))
	mu := clause.Mutation()

	if data.Level == nil {
		mu.ClearLevel()
	}
	if data.NameEng == nil {
		mu.ClearNameEng()
	}
	if data.Description == nil {
		mu.ClearDescription()
	}
	return clause.
		SetNillableDescription(data.Description).
		SetNillableLevel(data.Level).
		SetNillableNameEng(data.NameEng).
		SetNameKor(data.NameKor).
		SetYogaGroupID(data.YogaGroupID).Exec(ctx)
}

func (repo *yogaRepository) Patch(ctx context.Context, data *request.YogaPatchBody, id int) error {
	updateableData := utils.GetUpdateableData(data, yoga.Columns)

	clause := repo.db.Yoga.Update().Where(yoga.IDEQ(id))
	for key, val := range updateableData {
		clause.Mutation().SetField(key, val)
	}
	return clause.Exec(ctx)
}

func (repo *yogaRepository) Delete(ctx context.Context, id int) error {
	return repo.db.Yoga.DeleteOneID(id).Exec(ctx)
}

func (repo *yogaRepository) List(ctx context.Context, groupdId int) ([]*ent.Yoga, error) {
	// join reference
	return repo.db.YogaGroup.Query().
		Where(yogagroup.IDEQ(groupdId)).
		QueryYoga().
		Order(ent.Desc(yoga.FieldID)).
		All(ctx)
}

func (repo *yogaRepository) Exist(ctx context.Context, id int) (bool, error) {
	return repo.db.Yoga.Query().Where(yoga.IDEQ(id)).Exist(ctx)
}

// ------------------- YogaRaw -------------------

func (repo *yogaRepository) CreateRaws(ctx context.Context, d []*ent.YogaRaw) error {
	bulk := make([]*ent.YogaRawCreate, len(d))

	for i, v := range d {
		bulk[i] = repo.db.YogaRaw.Create().SetNillableAcademyID(v.AcademyID).SetNillableTeacherID(v.TeacherID).SetName(v.Name)
	}

	return repo.db.YogaRaw.CreateBulk(bulk...).Exec(ctx)
}

func (repo *yogaRepository) DeleteAndCreateRaws(ctx context.Context, d []*ent.YogaRaw, academyId, teacherId *int) (err error) {
	if academyId == nil && teacherId == nil {
		err = errors.New("academyId, teacherId 둘 중 하나는 입력해주세요")
		return
	}

	if academyId != nil && teacherId != nil {
		err = errors.New("academyId, teacherId 둘 중 하나만 입력해주세요")
		return
	}

	return entx.WithTx(ctx, repo.db, func(tx *ent.Tx) (err error) {
		clause := tx.YogaRaw.Delete()
		if academyId != nil {
			clause.Where(yogaraw.AcademyIDEQ(*academyId))
		} else if teacherId != nil {
			clause.Where(yogaraw.TeacherIDEQ(*teacherId))
		}
		_, err = clause.Exec(ctx)
		if err != nil {
			return
		}

		bulk := make([]*ent.YogaRawCreate, len(d))

		for i, v := range d {
			bulk[i] = tx.YogaRaw.Create().SetNillableAcademyID(v.AcademyID).SetNillableTeacherID(v.TeacherID).SetName(v.Name)
		}

		return tx.YogaRaw.CreateBulk(bulk...).Exec(ctx)
	})
}

func (repo *yogaRepository) DeleteRawsByTeacherIdOrAcademyId(ctx context.Context, academyId, teacherId *int) (err error) {
	if academyId == nil && teacherId == nil {
		err = errors.New("academyId, teacherId 둘 중 하나는 입력해주세요")
		return
	}

	if academyId != nil && teacherId != nil {
		err = errors.New("academyId, teacherId 둘 중 하나만 입력해주세요")
		return
	}

	clause := repo.db.YogaRaw.Delete()
	if academyId != nil {
		clause.Where(yogaraw.AcademyIDEQ(*academyId))
	} else if teacherId != nil {
		clause.Where(yogaraw.TeacherIDEQ(*teacherId))
	}
	_, err = clause.Exec(ctx)
	if err != nil {
		return
	}
	return
}
