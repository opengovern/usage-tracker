package model

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"time"
)

type CspmUsage struct {
	gorm.Model

	InstallId       string    `json:"install_id" gorm:"index:install_id_hostname"`
	GatherTimestamp time.Time `json:"gather_timestamp" gorm:"index:,sort:desc"`

	Hostname             string `json:"hostname" gorm:"index:install_id_hostname"`
	NumberOfUsers        int64  `json:"number_of_users"`
	IntegrationTypeCount pgtype.JSONB
}
