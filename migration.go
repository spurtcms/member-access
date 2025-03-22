package memberaccess

import (
	"time"

	"gorm.io/gorm"
)

type TblAccessControl struct {
	Id                int       `gorm:"primaryKey;auto_increment;type:serial"`
	AccessControlName string    `gorm:"type:character varying"`
	AccessControlSlug string    `gorm:"type:character varying"`
	CreatedOn         time.Time `gorm:"type:timestamp without time zone"`
	CreatedBy         int       `gorm:"type:integer"`
	ModifiedOn        time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy        int       `gorm:"DEFAULT:NULL"`
	IsDeleted         int       `gorm:"type:integer"`
	DeletedOn         time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy         int       `gorm:"DEFAULT:NULL"`
	TenantId          string    `gorm:"type:character varying"`
}

type TblAccessControlPages struct {
	Id                       int       `gorm:"primaryKey;auto_increment;type:serial"`
	AccessControlUserGroupId int       `gorm:"type:integer"`
	SpacesId                 int       `gorm:"type:integer"`
	PageGroupId              int       `gorm:"type:integer"`
	PageId                   int       `gorm:"type:integer"`
	CreatedOn                time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy                int       `gorm:"type:integer"`
	ModifiedOn               time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy               int       `gorm:"DEFAULT:NULL"`
	IsDeleted                int       `gorm:"type:integer"`
	DeletedOn                time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy                int       `gorm:"DEFAULT:NULL"`
	ChannelId                int       `gorm:"type:integer"`
	EntryId                  int       `gorm:"type:integer"`
	TenantId                 string    `gorm:"type:character varying"`
}

type TblAccessControlUserGroup struct {
	Id              int       `gorm:"primaryKey;auto_increment;type:serial"`
	AccessControlId int       `gorm:"type:integer"`
	MemberGroupId   int       `gorm:"type:integer"`
	CreatedOn       time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy       int       `gorm:"type:integer"`
	ModifiedOn      time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	IsDeleted       int       `gorm:"type:integer"`
	DeletedOn       time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
	TenantId        string    `gorm:"type:character varying"`
}

func Migration(db *gorm.DB) {

	if err := db.AutoMigrate(
		&TblAccessControl{},
		&TblAccessControlUserGroup{},
		&TblAccessControlPages{},
	); err != nil {

		panic(err)
	}

}
