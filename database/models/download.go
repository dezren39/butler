package models

import (
	"fmt"
	"time"

	itchio "github.com/itchio/go-itchio"
	"github.com/jinzhu/gorm"
)

type Download struct {
	// An UUID
	ID string `json:"id" gorm:"primary_key"`

	Reason     string     `json:"reason"`
	Position   int64      `json:"position"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`

	CaveID string `json:"caveId"`

	GameID int64        `json:"gameId"`
	Game   *itchio.Game `json:"game"`

	UploadID int64          `json:"uploadId"`
	Upload   *itchio.Upload `json:"upload"`

	BuildID int64         `json:"buildId"`
	Build   *itchio.Build `json:"build"`

	StagingFolder string `json:"stagingFolder"`
}

func AllDownloads(db *gorm.DB) []*Download {
	var dls []*Download
	err := db.Order(`"position" ASC`).Find(&dls).Error
	if err != nil {
		panic(err)
	}
	return dls
}

func DownloadByID(db *gorm.DB, downloadID string) *Download {
	var dl Download
	req := db.Where("id = ?", downloadID).Find(&dl)
	if req.Error != nil {
		if req.RecordNotFound() {
			return nil
		}
		panic(req.Error)
	}
	return &dl
}

func (d *Download) Preload(db *gorm.DB) {
	if d == nil {
		return
	}
	PreloadDownloads(db, d)
}

func PreloadDownloads(db *gorm.DB, downloadOrDownloads interface{}) {
	MustPreloadSimple(db, downloadOrDownloads, "Game", "Upload", "Build")
}

func DownloadMinPosition(db *gorm.DB) int64 {
	return downloadExtremePosition(db, "min")
}

func DownloadMaxPosition(db *gorm.DB) int64 {
	return downloadExtremePosition(db, "max")
}

func downloadExtremePosition(db *gorm.DB, extreme string) int64 {
	var row = struct {
		Position int64
	}{}

	query := fmt.Sprintf(`SELECT coalesce(%s(position), 0) AS position FROM downloads`, extreme)
	err := db.Raw(query).Scan(&row).Error
	if err != nil {
		panic(err)
	}
	return row.Position
}
