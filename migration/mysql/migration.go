package mysql

import (
	"time"

	"gorm.io/gorm"
)

type TblAccessControl struct {
	Id                int       `gorm:"primaryKey;auto_increment;"`
	AccessControlName string    `gorm:"type:varchar"`
	AccessControlSlug string     `gorm:"type:varchar"`
	CreatedOn         time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	CreatedBy         int       `gorm:"type:int"`
	ModifiedOn        time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy        int       `gorm:"DEFAULT:NULL"`
	IsDeleted         int       `gorm:"type:int"`
	DeletedOn         time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy         int       `gorm:"DEFAULT:NULL"`
}

type TblAccessControlPages struct {
	Id                       int       `gorm:"primaryKey;auto_increment;"`
	AccessControlUserGroupId int       `gorm:"type:int"`
	SpacesId                 int       `gorm:"type:int"`
	PageGroupId              int       `gorm:"type:int"`
	PageId                   int       `gorm:"type:int"`
	CreatedOn                time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	CreatedBy                int       `gorm:"type:int"`
	ModifiedOn               time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy               int       `gorm:"DEFAULT:NULL"`
	IsDeleted                int       `gorm:"type:int"`
	DeletedOn                time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy                int       `gorm:"DEFAULT:NULL"`
	ChannelId                int       `gorm:"type:int"`
	EntryId                  int       `gorm:"type:int"`
}

type TblAccessControlUserGroup struct {
	Id              int       `gorm:"primaryKey;auto_increment;"`
	AccessControlId int       `gorm:"type:int"`
	MemberGroupId   int       `gorm:"type:int"`
	CreatedOn       time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	CreatedBy       int       `gorm:"type:int"`
	ModifiedOn      time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	IsDeleted       int       `gorm:"type:int"`
	DeletedOn       time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
}

func MigrationTables(db *gorm.DB) {

	if err := db.AutoMigrate(
		&TblAccessControl{},
		&TblAccessControlUserGroup{},
		&TblAccessControlPages{},
	); err != nil {

		panic(err)
	}

}