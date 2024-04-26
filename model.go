package memberaccess

import (
	"github.com/spurtcms/member"
	"gorm.io/gorm"
)

// Get all content access list
func (AccessModel) GetContentAccessList(limit, offset int, filter Filter, DB *gorm.DB) (contentAccessList []Tblaccesscontrol, count int64, err error) {

	query := DB.Table("tbl_access_controls").Select("tbl_access_controls.*,tbl_users.username,tbl_roles.name").Where("tbl_access_controls.is_deleted = 0").Order("tbl_access_controls.id DESC")

	//joins
	query = query.Joins("left join tbl_users on tbl_users.id = tbl_access_controls.created_by")
	query = query.Joins("left join tbl_roles on tbl_roles.id = tbl_users.role_id")

	if filter.Keyword != "" {

		query = query.Where("(LOWER(TRIM(tbl_access_controls.access_control_name)) ILIKE LOWER(TRIM(?)))", "%"+filter.Keyword+"%")
	}

	if limit != 0 {

		query = query.Offset(offset).Limit(limit).Find(&contentAccessList)

	} else {

		query = query.Find(&contentAccessList).Count(&count)
	}

	if err := query.Error; err != nil {

		return []Tblaccesscontrol{}, 0, err
	}

	return contentAccessList, count, nil
}

func (AccessModel) GetAccessGrantedMemberGroups(memberGroups *[]TblAccessControlUserGroup, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("is_deleted = 0 and access_control_id = ?", accessId).Find(&memberGroups).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetMemberGroupsByContentAccessMemId(memgrp *member.TblMemberGroup, id int, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_groups").Where("is_deleted = 0 and id = ?", id).First(&memgrp).Error; err != nil {

		return err

	}

	return nil

}

func (AccessModel) GetaccessGrantedPageCount(count *int64, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_controls").Distinct("tbl_access_control_pages.page_id").Joins("inner join tbl_access_control_user_groups on tbl_access_control_user_groups.access_control_id = tbl_access_controls.id").
		Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id = tbl_access_control_user_groups.id").
		Where("tbl_access_controls.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_controls.id = ? and tbl_access_control_pages.page_id!= 0", accessId).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetaccessGrantedEntriesCount(count *int64, accessId int, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_controls").Distinct("tbl_access_control_pages.entry_id").Joins("inner join tbl_access_control_user_groups on tbl_access_control_user_groups.access_control_id = tbl_access_controls.id").
		Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id = tbl_access_control_user_groups.id").
		Where("tbl_access_controls.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_controls.id = ? and tbl_access_control_pages.entry_id!= 0", accessId).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetContentAccessByAccessId(id int, DB *gorm.DB) (accesscontrol *Tblaccesscontrol, err error) {

	if err := DB.Table("tbl_access_controls").Where("is_deleted = 0 and id = ?", id).First(&accesscontrol).Error; err != nil {

		return &Tblaccesscontrol{}, err
	}

	return accesscontrol, nil
}

func (AccessModel) GetPagesAndPageGroupsInContentAccess(accessId int, DB *gorm.DB) (contentAccessPages []TblAccessControlPages, err error) {

	query := DB.Model(TblAccessControlPages{}).Select("tbl_access_control_pages.*,tbl_pages.parent_id").Where("tbl_access_controls.id = ?", accessId)

	/*Joins*/
	query.Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_groups.is_deleted = 0")

	query.Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id and tbl_access_controls.is_deleted = 0")

	query.Joins("inner join tbl_pages on tbl_pages.id = tbl_access_control_pages.page_id and tbl_access_control_pages.is_deleted = 0")

	query.Find(&contentAccessPages)

	if err := query.Error; err != nil {

		return []TblAccessControlPages{}, err
	}

	return contentAccessPages, nil
}
