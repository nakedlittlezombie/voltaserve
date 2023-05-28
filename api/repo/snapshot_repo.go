package repo

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SnapshotRepo interface {
	Find(id string) (model.Snapshot, error)
	Save(snapshot model.Snapshot) error
	MapWithFile(id string, fileId string) error
	DeleteMappingsForFile(fileId string) error
	FindAllDangling() ([]model.Snapshot, error)
	DeleteAllDangling() error
	GetLatestVersionForFile(fileId string) (int64, error)
}

func NewSnapshotRepo() SnapshotRepo {
	return newSnapshotRepo()
}

func NewSnapshot() model.Snapshot {
	return &snapshotEntity{}
}

type snapshotEntity struct {
	ID         string         `json:"id" gorm:"column:id;size:36"`
	Version    int64          `json:"version" gorm:"column:version"`
	Original   datatypes.JSON `json:"original,omitempty" gorm:"column:original"`
	Preview    datatypes.JSON `json:"preview,omitempty" gorm:"column:preview"`
	Text       datatypes.JSON `json:"text,omitempty" gorm:"column:text"`
	OCR        datatypes.JSON `json:"ocr,omitempty" gorm:"column:ocr"`
	Thumbnail  datatypes.JSON `json:"thumbnail,omitempty" gorm:"column:thumbnail"`
	CreateTime string         `json:"createTime" gorm:"column:create_time"`
	UpdateTime *string        `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (snapshotEntity) TableName() string {
	return "snapshot"
}

func (s *snapshotEntity) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (s *snapshotEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	s.UpdateTime = &timeNow
	return nil
}

func (s snapshotEntity) GetID() string {
	return s.ID
}

func (s snapshotEntity) GetVersion() int64 {
	return s.Version
}

func (s snapshotEntity) GetOriginal() *model.S3Object {
	if s.Original.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Original.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s snapshotEntity) GetPreview() *model.S3Object {
	if s.Preview.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Preview.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s snapshotEntity) GetText() *model.S3Object {
	if s.Text.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.Text.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s snapshotEntity) GetOCR() *model.S3Object {
	if s.OCR.String() == "" {
		return nil
	}
	var res = model.S3Object{}
	if err := json.Unmarshal([]byte(s.OCR.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s snapshotEntity) GetThumbnail() *model.Thumbnail {
	if s.Thumbnail.String() == "" {
		return nil
	}
	var res = model.Thumbnail{}
	if err := json.Unmarshal([]byte(s.Thumbnail.String()), &res); err != nil {
		log.Fatal(err)
		return nil
	}
	return &res
}

func (s *snapshotEntity) SetID(id string) {
	s.ID = id
}

func (s *snapshotEntity) SetVersion(version int64) {
	s.Version = version
}

func (s *snapshotEntity) SetOriginal(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Original.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *snapshotEntity) SetPreview(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Preview.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *snapshotEntity) SetText(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Text.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *snapshotEntity) SetOCR(m *model.S3Object) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.OCR.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s *snapshotEntity) SetThumbnail(m *model.Thumbnail) {
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := s.Thumbnail.UnmarshalJSON(b); err != nil {
		log.Fatal(err)
	}
}

func (s snapshotEntity) HasOriginal() bool {
	return s.Original != nil
}

func (s snapshotEntity) HasPreview() bool {
	return s.Preview != nil
}

func (s snapshotEntity) HasText() bool {
	return s.Text != nil
}

func (s snapshotEntity) HasOCR() bool {
	return s.OCR != nil
}

func (s snapshotEntity) HasThumbnail() bool {
	return s.Thumbnail != nil
}

func (s snapshotEntity) GetCreateTime() string {
	return s.CreateTime
}

func (s snapshotEntity) GetUpdateTime() *string {
	return s.UpdateTime
}

type snapshotRepo struct {
	db *gorm.DB
}

func newSnapshotRepo() *snapshotRepo {
	return &snapshotRepo{
		db: infra.GetDb(),
	}
}

func (repo *snapshotRepo) find(id string) (*snapshotEntity, error) {
	var res snapshotEntity
	if db := repo.db.Where("id = ?", id).First(&res); db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewSnapshotNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &res, nil
}

func (repo *snapshotRepo) Find(id string) (model.Snapshot, error) {
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *snapshotRepo) Save(snapshot model.Snapshot) error {
	if db := repo.db.Save(snapshot); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) MapWithFile(id string, fileId string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) VALUES (?, ?)", id, fileId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) DeleteMappingsForFile(fileId string) error {
	if db := repo.db.Exec("DELETE FROM snapshot_file WHERE file_id = ?", fileId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) findAllForFile(fileId string) ([]*snapshotEntity, error) {
	var res []*snapshotEntity
	db := repo.db.
		Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.file_id = ? ORDER BY s.version", fileId).
		Scan(&res)
	if db.Error != nil {
		return nil, db.Error
	}
	return res, nil
}

func (repo *snapshotRepo) FindAllDangling() ([]model.Snapshot, error) {
	var snapshots []*snapshotEntity
	db := repo.db.Raw("SELECT * FROM snapshot s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL").Scan(&snapshots)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Snapshot
	for _, s := range snapshots {
		res = append(res, s)
	}
	return res, nil
}

func (repo *snapshotRepo) DeleteAllDangling() error {
	if db := repo.db.Exec("DELETE FROM snapshot WHERE id IN (SELECT s.id FROM (SELECT * FROM snapshot) s LEFT JOIN snapshot_file sf ON s.id = sf.snapshot_id WHERE sf.snapshot_id IS NULL)"); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *snapshotRepo) GetLatestVersionForFile(fileId string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	if db := repo.db.
		Raw("SELECT coalesce(max(s.version), 0) + 1 result FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id WHERE map.file_id = ?", fileId).
		Scan(&res); db.Error != nil {
		return 0, db.Error
	}
	return res.Result, nil
}
