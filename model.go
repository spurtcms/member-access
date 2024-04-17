package memberaccess

import (
	"time"

	"github.com/spurtcms/member"
)

type tblaccesscontrol struct {
	Id                   int                     `gorm:"column:id"`
	AccessControlName    string                  `gorm:"column:access_control_name"`
	AccessControlSlug    string                  `gorm:"column:access_control_slug"`
	CreatedOn            time.Time               `gorm:"column:created_on"`
	CreatedBy            int                     `gorm:"column:created_by"`
	ModifiedOn           time.Time               `gorm:"column:modified_on"`
	ModifiedBy           int                     `gorm:"DEFAULT:NULL"`
	IsDeleted            int                     `gorm:"column:is_deleted"`
	DeletedOn            time.Time               `gorm:"DEFAULT:NULL"`
	DeletedBy            int                     `gorm:"DEFAULT:NULL"`
	Username             string                  `gorm:"column:username;<-:false"`
	Rolename             string                  `gorm:"column:name;<-:false"`
	MemberGroups         []member.TblMemberGroup `gorm:"-"`
	DateString           string                  `gorm:"-"`
	AccessGrantedModules []string                `gorm:"-"`
}

type tblaccesscontrolpages struct {
	Id                       int `gorm:"primaryKey;auto_increment"`
	AccessControlUserGroupId int
	SpacesId                 int
	PageGroupId              int
	PageId                   int
	CreatedOn                time.Time
	CreatedBy                int
	ModifiedOn               time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy               int       `gorm:"DEFAULT:NULL"`
	IsDeleted                int
	DeletedOn                time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy                int       `gorm:"DEFAULT:NULL"`
	ParentPageId             int       `gorm:"column:parent_id;<-:false"`
	ChannelId                int
	EntryId                  int
}

type tblaccesscontrolusergroup struct {
	Id              int `gorm:"primaryKey;auto_increment"`
	AccessControlId int
	MemberGroupId   int
	CreatedOn       time.Time
	CreatedBy       int
	ModifiedOn      time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy      int       `gorm:"DEFAULT:NULL"`
	IsDeleted       int
	DeletedOn       time.Time `gorm:"DEFAULT:NULL"`
	SpacesId        int       `gorm:"-:migration;<-:false"`
	PageId          int       `gorm:"-:migration;<-:false"`
	PageGroupId     int       `gorm:"-:migration;<-:false"`
	DeletedBy       int       `gorm:"DEFAULT:NULL"`
}

type Filter struct {
	Keyword string
}

type SubPage struct {
	Id       string `json:"id"`
	GroupId  string `json:"groupId"`
	ParentId string `json:"parentId"`
	SpaceId  string `json:"spaceId"`
}

type Page struct {
	Id      string `json:"id"`
	GroupId string `json:"groupId"`
	SpaceId string `json:"spaceId"`
}

type PageGroup struct {
	Id      string `json:"id"`
	SpaceId string `json:"spaceId"`
}

type Entry struct {
	Id        string `json:"id"`
	ChannelId string `json:"channelId"`
}

type MemberAccessControlRequired struct {
	Title          string
	Pages          []Page
	SubPage        []SubPage
	Group          []PageGroup
	SpacesIds      []int
	Channels       []int
	ChannelEntries []Entry
	MemberGroupIds []int
}
