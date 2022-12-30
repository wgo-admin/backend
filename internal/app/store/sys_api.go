package store

import (
  "context"
  "github.com/wgo-admin/backend/internal/app/store/model"
  "gorm.io/gorm"
)

type ISysApiStore interface {
  Create(ctx context.Context, sysApi *model.SysApiM) error
  Get(ctx context.Context, id int64) (*model.SysApiM, error)
  Update(ctx context.Context, sysApi *model.SysApiM) error
  Delete(ctx context.Context, ids []int64) error
  List(ctx context.Context, sysApi *model.SysApiM, offset, limit int) (int64, []*model.SysApiM, error)
}

type sysApis struct {
  db *gorm.DB
}

var _ ISysApiStore = (*sysApis)(nil)

func newSysApis(db *gorm.DB) *sysApis {
  return &sysApis{db: db}
}

func (s *sysApis) Create(ctx context.Context, sysApi *model.SysApiM) error {
  return s.db.Create(&sysApi).Error
}

func (s *sysApis) Get(ctx context.Context, id int64) (*model.SysApiM, error) {
  var sysApi model.SysApiM
  if err := s.db.First(&sysApi, "id = ?", id).Error; err != nil {
    return nil, err
  }
  return &sysApi, nil
}

func (s *sysApis) Update(ctx context.Context, sysApi *model.SysApiM) error {
  return s.db.Save(sysApi).Error
}

func (s *sysApis) Delete(ctx context.Context, ids []int64) error {
  if err := s.db.Where("id in (?)", ids).Delete(&model.SysApiM{}).Error; err != nil {
    return err
  }
  return nil
}

func (s *sysApis) List(ctx context.Context, sysApi *model.SysApiM, offset, limit int) (total int64, ret []*model.SysApiM, err error) {
  db := s.db.Model(&model.SysApiM{})

  if sysApi.Method != "" {
    db = db.Where("method LIKE ?", "%"+sysApi.Method+"%")
  }

  if sysApi.GroupName != "" {
    db = db.Where("group_name LIKE ?", "%"+sysApi.GroupName+"%")
  }

  if sysApi.Title != "" {
    db = db.Where("title LIKE ?", "%"+sysApi.Title+"%")
  }

  err = db.Count(&total).Error
  if err != nil {
    return 0, nil, err
  }

  db = db.Offset(offset).Limit(limit)
  err = db.Find(&ret).Error
  if err != nil {
    return 0, nil, err
  }

  return
}
