package memberaccess

import (
	"github.com/spurtcms/member"
	"gorm.io/gorm"
)

// Get all content access list
func (AccessM AccessModel) GetContentAccessList(limit, offset int, filter Filter, DB *gorm.DB, tenantid int) (contentAccessList []Tblaccesscontrol, count int64, err error) {

	query := DB.Table("tbl_access_controls").Select("tbl_access_controls.*,tbl_users.username,tbl_roles.name").Where("tbl_access_controls.is_deleted = 0 and  tbl_access_controls.tenant_id=?", tenantid).Order("tbl_access_controls.id DESC")

	//joins
	query = query.Joins("left join tbl_users on tbl_users.id = tbl_access_controls.created_by")
	query = query.Joins("left join tbl_roles on tbl_roles.id = tbl_users.role_id")

	if AccessM.DataAccess == 1 {

		query = query.Where("tbl_access_controls.created_by = ?", AccessM.UserId)
	}

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

func (AccessModel) GetAccessGrantedMemberGroups(memberGroups *[]TblAccessControlUserGroup, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("is_deleted = 0 and access_control_id = ? and  tenant_id = ?", accessId, tenantid).Find(&memberGroups).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetMemberGroupsByContentAccessMemId(memgrp *member.TblMemberGroup, id int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_member_groups").Where("is_deleted = 0 and id = ? and  tenant_id = ?", id, tenantid).First(&memgrp).Error; err != nil {

		return err

	}

	return nil

}

func (AccessModel) GetaccessGrantedPageCount(count *int64, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_controls").Distinct("tbl_access_control_pages.page_id").Joins("inner join tbl_access_control_user_groups on tbl_access_control_user_groups.access_control_id = tbl_access_controls.id").
		Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id = tbl_access_control_user_groups.id").
		Where("tbl_access_controls.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_controls.id = ? and tbl_access_control_pages.page_id!= 0 and  tbl_access_controls.tenant_id=?)", accessId, tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetaccessGrantedEntriesCount(count *int64, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_controls").Distinct("tbl_access_control_pages.entry_id").Joins("inner join tbl_access_control_user_groups on tbl_access_control_user_groups.access_control_id = tbl_access_controls.id").
		Joins("inner join tbl_access_control_pages on tbl_access_control_pages.access_control_user_group_id = tbl_access_control_user_groups.id").
		Where("tbl_access_controls.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_controls.id = ? and tbl_access_control_pages.entry_id!= 0 and  tbl_access_controls.tenant_id=?)", accessId, tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetContentAccessByAccessId(id int, DB *gorm.DB, tenantid int) (accesscontrol *Tblaccesscontrol, err error) {

	if err := DB.Table("tbl_access_controls").Where("is_deleted = 0 and id = ? and  tenant_id = ?", id, tenantid).First(&accesscontrol).Error; err != nil {

		return &Tblaccesscontrol{}, err
	}

	return accesscontrol, nil
}

func (AccessModel) GetPagesAndPageGroupsInContentAccess(accessId int, DB *gorm.DB, tenantid int) (contentAccessPages []Tblaccesscontrolpages, err error) {

	query := DB.Model(TblAccessControlPages{}).Select("tbl_access_control_pages.*,tbl_pages.parent_id").Where("tbl_access_controls.id = ? and  tbl_access_controls.tenant_id=?)", accessId, tenantid)

	/*Joins*/
	query.Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_groups.is_deleted = 0")

	query.Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id and tbl_access_controls.is_deleted = 0")

	query.Joins("inner join tbl_pages on tbl_pages.id = tbl_access_control_pages.page_id and tbl_access_control_pages.is_deleted = 0")

	query.Find(&contentAccessPages)

	if err := query.Error; err != nil {

		return []Tblaccesscontrolpages{}, err
	}

	return contentAccessPages, nil
}

func (AccessModel) GetPageGroupsInContentAccess(accessId int, DB *gorm.DB, tenantid int) (contentAccessPages []Tblaccesscontrolpages, err error) {

	query := DB.Model(TblAccessControlPages{}).Select("tbl_access_control_pages.page_group_id,tbl_access_control_pages.id").Where("tbl_access_controls.id = ? and  tbl_access_controls.tenant_id=?)", accessId, tenantid)

	/*Joins*/
	query.Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_groups.is_deleted = 0")

	query.Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id and tbl_access_controls.is_deleted = 0")

	query.Group("page_group_id")

	query.Find(&contentAccessPages)

	if err := query.Error; err != nil {

		return []Tblaccesscontrolpages{}, err
	}

	return contentAccessPages, nil
}

func (AccessModel) GetSelectedSpaces(accessId int, DB *gorm.DB, tenantid int) (contentAccessPages []Tblaccesscontrolpages, err error) {

	query := DB.Model(TblAccessControlPages{}).Select("spaces_id").Where("tbl_access_controls.id = ? and  tbl_access_controls.tenant_id=?)", accessId, tenantid)

	/*Joins*/
	query.Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_groups.is_deleted = 0")

	query.Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id and tbl_access_controls.is_deleted = 0")

	query.Group("spaces_id")

	query.Find(&contentAccessPages)

	if err := query.Error; err != nil {

		return []Tblaccesscontrolpages{}, err
	}

	return contentAccessPages, nil
}

func (AccessModel) GetSelectedEntries(accessId int, DB *gorm.DB, tenantid int) (contentAccessPages []Tblaccesscontrolpages, err error) {

	query := DB.Model(TblAccessControlPages{}).Where("tbl_access_controls.id = ? and  tbl_access_controls.tenant_id=?)", accessId, tenantid)

	/*Joins*/
	query.Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id and tbl_access_control_user_groups.is_deleted = 0")

	query.Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id and tbl_access_controls.is_deleted = 0")

	query.Find(&contentAccessPages)

	if err := query.Error; err != nil {

		return []Tblaccesscontrolpages{}, err
	}

	return contentAccessPages, nil
}

func (AccessModel) CreateMemberGroupRestrict(access TblAccessControlUserGroup, DB *gorm.DB) (TblAccessControlUserGroup, error) {

	if err := DB.Model(TblAccessControlUserGroup{}).Create(&access).Error; err != nil {

		return TblAccessControlUserGroup{}, err
	}

	return access, nil
}

func (AccessModel) CreatePage(access *TblAccessControlPages, DB *gorm.DB) error {

	if err := DB.Model(TblAccessControlPages{}).Create(&access).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) DeletePage(pg_access *TblAccessControlPages, id []int, pgids []int, DB *gorm.DB, tenantid int) error {

	if err := DB.Model(TblAccessControlPages{}).Where("access_control_user_group_id in (?) and page_id in (?) and  tenant_id = ?", id, pgids, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": pg_access.IsDeleted, "deleted_on": pg_access.DeletedOn, "deleted_by": pg_access.DeletedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) DeleteGroup(pg_access *TblAccessControlPages, id []int, grpid []int, DB *gorm.DB, tenantid int) error {

	if err := DB.Model(TblAccessControlPages{}).Where("access_control_user_group_id in (?) and page_group_id in (?) and  tenant_id = ?", id, grpid, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": pg_access.IsDeleted, "deleted_on": pg_access.DeletedOn, "deleted_by": pg_access.DeletedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) DeleteSubPage(pg_access *TblAccessControlPages, id []int, spacesid []int, DB *gorm.DB, tenantid int) error {

	if err := DB.Model(TblAccessControlPages{}).Where("access_control_user_group_id in (?) and spaces_id in (?) and  tenant_id = ?", id, spacesid, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": pg_access.IsDeleted, "deleted_on": pg_access.DeletedOn, "deleted_by": pg_access.DeletedBy}).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetGroupsByAccessId(accessid int, DB *gorm.DB, tenantid int) (usergroups []TblAccessControlUserGroup, er error) {

	if err := DB.Model(TblAccessControlUserGroup{}).Where("access_control_id=? and  tenant_id=?", accessid, tenantid).Find(usergroups).Error; err != nil {

		return []TblAccessControlUserGroup{}, nil
	}

	return usergroups, nil
}

func (AccessModel) UpdateContentAccessId(contentAccess *TblAccessControl, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_controls").Where("is_deleted = 0 and id = ? and  tenant_id=?", contentAccess.Id, tenantid).UpdateColumns(map[string]interface{}{"access_control_name": contentAccess.AccessControlName, "access_control_slug": contentAccess.AccessControlSlug, "modified_on": contentAccess.ModifiedOn, "modified_by": contentAccess.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil

}

/*Create Access*/
func (AccessModel) NewContentAccessEntry(contentAccess *TblAccessControl, DB *gorm.DB) error {

	if err := DB.Table("tbl_access_controls").Create(&contentAccess).Error; err != nil {

		return err
	}

	return nil
}

// Delete Access Control tbl
func (AccessModel) DeleteControlAccess(accesscontrol *TblAccessControl, id int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_controls").Where("id = ? and  tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"deleted_by": accesscontrol.DeletedBy, "deleted_on": accesscontrol.DeletedOn, "is_deleted": accesscontrol.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

// get access membergroup
func (AccessModel) GetAccessGrantedMemberGroupsList(memgrps *[]int, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_member_groups").Select("tbl_member_groups.id").
		Joins("left join tbl_access_control_user_groups on tbl_access_control_user_groups.member_group_id =  tbl_member_groups.id and tbl_access_control_user_groups.is_deleted = 0 ").
		Where("tbl_member_groups.is_deleted = 0 and tbl_access_control_user_groups.access_control_id = ? and  tbl_member_groups.tenant_id=?", accessId, tenantid).Find(&memgrps).Error; err != nil {

		return err

	}

	return nil

}

func (AccessModel) GetAccessGrantedEntries(AccessEntries *[]TblAccessControlPages, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_pages").Select("distinct on (tbl_access_control_pages.entry_id),(tbl_access_control_pages.*").
		Joins("inner join tbl_access_control_user_groups on tbl_access_control_user_groups.id = tbl_access_control_pages.access_control_user_group_id").
		Joins("inner join tbl_access_controls on tbl_access_controls.id = tbl_access_control_user_groups.access_control_id").
		Where("tbl_access_controls.is_deleted = 0 and tbl_access_control_pages.is_deleted = 0 and tbl_access_controls.id = ? and tbl_access_control_pages.entry_id!= 0 and tbl_access_control_pages.tenant_id=?", accessId, tenantid).Find(&AccessEntries).Error; err != nil {

		return err
	}

	return nil
}

func (AccessModel) GetEntriesCountUnderChannel(count *int64, channelId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_channel_entries").Where("is_deleted = 0 and status = 1 and channel_id = ? and  tenant_id=?", channelId, tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

// check membergroup access
func (AccessModel) CheckPresenceOfAccessGrantedMemberGroups(count *int64, mem_id, accessId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("is_deleted = 0 and member_group_id = ? and access_control_id = ? and  tenant_id=?", mem_id, accessId, tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil
}

// create member group access
func (AccessModel) GrantAccessToMemberGroups(memberGrpAccess *TblAccessControlUserGroup, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Create(&memberGrpAccess).Error; err != nil {

		return err
	}

	return nil
}

// update membergroup access
func (AccessModel) UpdateContentAccessMemberGroup(accessmemgrp *TblAccessControlUserGroup, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("is_deleted = 0 and access_control_id = ? and member_group_id = ? and  tenant_id=?", accessmemgrp.AccessControlId, accessmemgrp.MemberGroupId, tenantid).UpdateColumns(map[string]interface{}{"modified_on": accessmemgrp.ModifiedOn, "modified_by": accessmemgrp.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

// get member groupby acessid
func (AccessModel) GetMemberGrpByAccessControlId(memberGrpAccess *[]TblAccessControlUserGroup, content_access_id int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("is_deleted = 0 and access_control_id = ? and  tenant_id=?", content_access_id, tenantid).Find(&memberGrpAccess).Error; err != nil {

		return err
	}

	return nil

}

// check access for enteries
func (AccessModel) CheckPresenceOfChannelEntriesInContentAccess(count *int64, accessGroupId, chanId, entryId int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and channel_id = ? and entry_id = ? and  tenant_id=?", accessGroupId, chanId, entryId, tenantid).Count(count).Error; err != nil {

		return err
	}

	return nil

}

// update accesspage
func (AccessModel) UpdateAccessPage(chanAccess *TblAccessControlPages, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and channel_id = ? and entry_id = ? and  tenant_id=?", chanAccess.AccessControlUserGroupId, chanAccess.ChannelId, chanAccess.EntryId, tenantid).UpdateColumns(map[string]interface{}{"modified_on": chanAccess.ModifiedOn, "modified_by": chanAccess.ModifiedBy}).Error; err != nil {

		return err
	}

	return nil
}

// remove membergroup access
func (AccessModel) RemoveMemberGroupsNotUnderContentAccessRights(memgrp_access *TblAccessControlUserGroup, memgrp_array []int, access_id int, DB *gorm.DB) error {

	if err := DB.Debug().Exec(`
		WITH updated_user_groups AS (
			UPDATE tbl_access_control_user_groups
			SET is_deleted = (?),
			deleted_by = (?),
			deleted_on = (?)
			WHERE tbl_access_control_user_groups.IS_DELETED =0 and tbl_access_control_user_groups.access_control_id=? and tbl_access_control_user_groups.member_group_id not in(?)
			RETURNING id
		)
		UPDATE tbl_access_control_pages
		SET is_deleted = (?),
		deleted_by = (?),
		deleted_on = (?)
		FROM updated_user_groups
		WHERE tbl_access_control_pages.access_control_user_group_id = (
			SELECT id
			FROM tbl_access_control_user_groups
			WHERE tbl_access_control_user_groups.id = updated_user_groups.id
		)`, memgrp_access.IsDeleted, memgrp_access.DeletedBy, memgrp_access.DeletedOn, access_id, memgrp_array, memgrp_access.IsDeleted, memgrp_access.DeletedBy, memgrp_access.DeletedOn).Error; err != nil {

		return err
	}

	return nil
}

// remove access for entries
func (AccessModel) RemoveChannelEntriesNotUnderContentAccess(chanAccess *TblAccessControlPages, entryIds []int, DB *gorm.DB, tenantid int) error {

	if err := DB.Debug().Table("tbl_access_control_pages").Where("is_deleted = 0 and access_control_user_group_id = ? and entry_id != 0 and entry_id NOT IN (?) and  tenant_id=?", chanAccess.AccessControlUserGroupId, entryIds, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": chanAccess.IsDeleted, "deleted_on": chanAccess.DeletedOn, "deleted_by": chanAccess.DeletedBy}).Error; err != nil {

		return err
	}

	return nil
}

// Delete Access Control User Group tbl
func (AccessModel) DeleteInAccessUserGroup(accessusergrp *TblAccessControlUserGroup, Id int, DB *gorm.DB, tenantid int) error {

	if err := DB.Table("tbl_access_control_user_groups").Where("access_control_id = ? and  tenant_id=?", Id, tenantid).UpdateColumns(map[string]interface{}{"deleted_by": accessusergrp.DeletedBy, "deleted_on": accessusergrp.DeletedOn, "is_deleted": accessusergrp.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}

// To Get Deleted id in access control user group tbl
func (AccessModel) GetDeleteIdInAccessUserGroup(controlaccessgrp *[]TblAccessControlUserGroup, Id int, DB *gorm.DB, tenantid int) (*[]TblAccessControlUserGroup, error) {

	if err := DB.Table("tbl_access_control_user_groups").Where("access_control_id = ? and  tenant_id=?", Id, tenantid).Find(&controlaccessgrp).Error; err != nil {

		return &[]TblAccessControlUserGroup{}, err
	}

	return controlaccessgrp, nil
}

// Delete Access Control Pages tbl
func (AccessModel) DeleteAccessControlPages(pg_access *TblAccessControlPages, Id []int, DB *gorm.DB, tenantid int) error {
	if err := DB.Table("tbl_access_control_pages").Where("access_control_user_group_id IN ? and  tenant_id=?", Id, tenantid).UpdateColumns(map[string]interface{}{"deleted_by": pg_access.DeletedBy, "deleted_on": pg_access.DeletedOn, "is_deleted": pg_access.IsDeleted}).Error; err != nil {

		return err
	}

	return nil
}
