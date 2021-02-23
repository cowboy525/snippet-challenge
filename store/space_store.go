package sqlstore

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/topoface/snippet-challenge/model"
	"github.com/topoface/snippet-challenge/store"
	"github.com/topoface/snippet-challenge/store/sqlstore/pagination"
	"github.com/topoface/snippet-challenge/utils"
)

// SQLSpaceStore structure
type SQLSpaceStore struct {
	SQLStore
}

func newSQLSpaceStore(sqlStore SQLStore) store.SpaceStore {
	s := &SQLSpaceStore{
		sqlStore,
	}

	return s
}

// Create creates a new space
func (ss SQLSpaceStore) Create(space *model.Space) (*model.Space, *model.AppError) {
	space.PreSave()

	if err := space.IsValid(); err != nil {
		return nil, err
	}

	var count int
	err := ss.GetMaster().Model(&model.Space{}).Where("display_name = ?", space.DisplayName).Count(&count).Error
	if err != nil || count != 0 {
		errors := []map[string]interface{}{{"displayName": []*model.Error{model.NewError("store.sql_space.create_or_update.display_name_exists", nil)}}}
		return nil, model.ValidationErrorWithDetails("SpaceStore.Create", errors)
	}

	err = ss.GetMaster().Save(space).Error
	if err != nil {
		return nil, model.NewAppError("SpaceStore.Create", "store.sql_space.create", nil, err.Error(), http.StatusInternalServerError)
	}
	return space, nil
}

// Update updates an existing space
func (ss SQLSpaceStore) Update(space *model.Space) (*model.Space, *model.AppError) {
	space.PreUpdate()

	if err := space.IsValid(); err != nil {
		return nil, err
	}

	oldSpace, err := ss.Get(space.ID)
	if err != nil {
		return nil, err
	}

	if oldSpace.DisplayName != space.DisplayName {
		var count int
		err := ss.GetMaster().Model(&model.Space{}).Where("display_name = ?", space.DisplayName).Count(&count).Error
		if err != nil || count != 0 {
			errors := []map[string]interface{}{{"displayName": []*model.Error{model.NewError("store.sql_space.create_or_update.display_name_exists", nil)}}}
			return nil, model.ValidationErrorWithDetails("SpaceStore.Update", errors)
		}
	}

	err1 := ss.GetMaster().Model(space).UpdateColumns(utils.StructToMapForGorm(space)).Error
	if err1 != nil {
		return nil, model.NewAppError("SpaceStore.Update", "store.sql_space.update", nil, "id="+strconv.FormatUint(space.ID, 10)+", "+err1.Error(), http.StatusInternalServerError)
	}

	return space, nil
}

// Get return space for id
func (ss SQLSpaceStore) Get(id uint64) (*model.Space, *model.AppError) {
	space := &model.Space{}
	err := ss.GetReplica().Where(&model.Space{ID: id}).First(space).Error
	if err != nil {
		return nil, model.NewAppError("SpaceStore.Get", "store.sql_space.get", nil, "id="+strconv.FormatUint(id, 10)+", "+err.Error(), http.StatusNotFound)
	}
	return space, nil
}

// GetSpaces return spaces for given options
func (ss SQLSpaceStore) GetSpaces(options model.GetSpacesOptions) (*pagination.Paginator, *model.AppError) {
	db := ss.GetReplica().Table(model.TBL_SPACES).Select("id, display_name, avatar, avatar_small").Where("display_name LIKE ?", options.DisplayName+"%")

	var spaces []model.SpaceResponse
	paginator, err := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    options.Page,
		Limit:   options.Limit,
		URL:     options.URL,
		OrderBy: []string{"id asc"},
	}, &spaces)

	if err != nil {
		return nil, model.NewAppError("SpaceStore.GetSpaces", "store.sql_space.get_spaces", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return paginator, nil
}

// GetSpaceUsers return space users
func (ss SQLSpaceStore) GetSpaceUsers(options model.GetSpaceUsersOptions) (*pagination.Paginator, *model.AppError) {
	db := ss.GetReplica().Table(model.TBL_USERS + " u").Select("u.*")
	db = db.Where("u.space_id = ?", options.SpaceID)
	if len(options.Account) != 0 {
		db = db.Where("u.account LIKE ?", options.Account+"%")
	}
	if len(options.DisplayName) != 0 {
		db = db.Where("u.display_name LIKE ?", options.DisplayName+"%")
	}
	if len(options.Keyword) > 0 {
		db = db.Where("(u.email LIKE ? OR u.account LIKE ? OR u.display_name LIKE ?)", "%"+options.Keyword+"%", "%"+options.Keyword+"%", "%"+options.Keyword+"%")
	}
	if options.ProjectIDExact != 0 || options.ProjectIDNotExact != 0 {
		projectID := options.ProjectIDExact
		if projectID == 0 {
			projectID = options.ProjectIDNotExact
		}
		db = db.Joins("LEFT JOIN "+model.TBL_PROJECT_USERS+" pu ON u.id = pu.user_id AND pu.project_id = ?", projectID)
	}
	if options.ProjectIDExact != 0 {
		db = db.Where("(u.space_role_id IN (?) OR pu.project_id = ?)", []uint64{model.ROLE_SPACE_OWNER.ID, model.ROLE_SPACE_ADMIN.ID}, options.ProjectIDExact)
	} else if options.ProjectIDNotExact != 0 {
		db = db.Where("u.space_role_id NOT IN (?)", []uint64{model.ROLE_SPACE_OWNER.ID, model.ROLE_SPACE_ADMIN.ID})
		db = db.Where("pu.project_id IS NULL")
	}

	var users []model.User
	paginator, err := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    options.Page,
		Limit:   options.Limit,
		URL:     options.URL,
		OrderBy: []string{"id asc"},
	}, &users)

	if err != nil {
		return nil, model.NewAppError("SpaceStore.GetSpaceUsers", "store.sql_space.get_space_users", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return paginator, nil
}

// PermanentDelete deletes space by id
func (ss SQLSpaceStore) PermanentDelete(id uint64) *model.AppError {
	ss.Project().OnDeleteSpace(id)
	ss.User().OnDeleteSpace(id)
	ss.SlackTeam().OnDeleteSpace(id)

	err := ss.GetMaster().Unscoped().Delete(&model.Space{ID: id}).Error
	if err != nil {
		return model.NewAppError("SpaceStore.PermanentDelete", "store.sql_space.permanent_delete", nil, "id="+strconv.FormatUint(id, 10)+", err="+err.Error(), http.StatusInternalServerError)
	}
	return nil
}

// SearchBySubject return subject search result
func (ss SQLSpaceStore) SearchBySubject(options model.SubjectSearchOptions) (*pagination.Paginator, *model.AppError) {
	if options.SpaceID == 0 {
		return nil, model.NewAppError("SpaceStore.SearchBySubject", "store.sql_space.search_by_subject", nil, "space id required", http.StatusInternalServerError)
	}

	options.Subject = strings.ReplaceAll(options.Subject, "'", "\\'")

	taskQuery := fmt.Sprintf(`Select 'Task' as `+"`table`"+`, t.id as pk, t.subject, t.project_id, t.updated_at from `+model.TBL_TASKS+` t
		where t.subject like '%%%v%%' and t.deleted_at is null`, options.Subject)
	noteQuery := fmt.Sprintf(`Select 'Note' as `+"`table`"+`, n.id as pk, n.subject, n.project_id, n.updated_at from `+model.TBL_NOTES+` n
		where n.subject like '%%%v%%' and n.deleted_at is null`, options.Subject)
	// mediaQuery := fmt.Sprintf(`Select 'Media' as `+"`table`"+`, m.id as pk, m.name as subject, m.project_id, m.updated_at from `+model.TBL_MEDIAS+` m
	// 	where m.name like '%%%v%%'`, options.Subject)
	projectQuery := fmt.Sprintf(`Select 'Project' as`+"`table`"+`, pj.id as pk, pj.display_name as subject, pj.id as project_id, '3000-00-00' as updated_at from `+model.TBL_PROJECTS+` pj
		where pj.display_name like '%%%v%%'`, options.Subject)

	if options.JoinAllProjects {
		taskQuery = taskQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECTS+" where id = t.project_id and space_id = %v)", options.SpaceID)
		noteQuery = noteQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECTS+" where id = n.project_id and space_id = %v)", options.SpaceID)
		// mediaQuery = mediaQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECTS+" where id = m.project_id and space_id = %v)", options.SpaceID)
		projectQuery = projectQuery + fmt.Sprintf(" and pj.space_id = %v", options.SpaceID)
	} else {
		taskQuery = taskQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECT_USERS+" where user_id = %v and project_id = t.project_id)", options.UserID)
		noteQuery = noteQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECT_USERS+" where user_id = %v and project_id = n.project_id)", options.UserID)
		// mediaQuery = mediaQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECT_USERS+" where user_id = %v and project_id = m.project_id)", options.UserID)
		projectQuery = projectQuery + fmt.Sprintf(" and EXISTS(select 1 from "+model.TBL_PROJECT_USERS+" where user_id = %v and project_id = pj.id)", options.UserID)
	}

	if len(options.ProjectIDs) > 0 {
		taskQuery = taskQuery + fmt.Sprintf("and y.project_id in (%v)", utils.JoinUint64sToString(options.ProjectIDs, ","))
		noteQuery = noteQuery + fmt.Sprintf("and n.project_id in (%v)", utils.JoinUint64sToString(options.ProjectIDs, ","))
		// mediaQuery = mediaQuery + fmt.Sprintf("and m.project_id in (%v)", utils.JoinUint64sToString(options.ProjectIDs, ","))
		projectQuery = projectQuery + fmt.Sprintf("and pj.id in (%v)", utils.JoinUint64sToString(options.ProjectIDs, ","))
	}
	// mediaQuery: deprecated until File tab implementation

	db := ss.GetReplica().Table(fmt.Sprintf(`
		(
			%v
			union all
			%v
			union all
			%v
		) as res`, taskQuery, noteQuery, projectQuery)).Select("res.*")

	var searchResult []model.SubjectSearchResult
	paginator, err := pagination.Paging(&pagination.Param{
		DB:      db,
		Page:    options.Page,
		Limit:   options.Limit,
		URL:     options.URL,
		OrderBy: []string{"res.updated_at desc"},
	}, &searchResult)

	if err != nil {
		return nil, model.NewAppError("SpaceStore.SearchBySubject", "store.sql_space.search_by_subject", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return paginator, nil
}
