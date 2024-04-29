package memberaccess

import (
	"log"
	"strconv"
	"strings"

	"github.com/spurtcms/member"
)

func AccessSetup(config Config) *AccessControl {

	return &AccessControl{
		DB:               config.DB,
		AuthEnable:       config.AuthEnable,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
	}

}

/*List */
func (access *AccessControl) ContentAccessList(limit int, offset int, filter Filter) (accesslist []Tblaccesscontrol, totalCount int64, err error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []Tblaccesscontrol{}, 0, autherr
	}

	contentAccessList, _, errr := Accessmodel.GetContentAccessList(limit, offset, filter, access.DB)

	if errr != nil {

		return []Tblaccesscontrol{}, 0, errr
	}

	var final_content_accesslist []Tblaccesscontrol

	for _, contentAccess := range contentAccessList {

		var access_grant_memgrps []TblAccessControlUserGroup

		Accessmodel.GetAccessGrantedMemberGroups(&access_grant_memgrps, contentAccess.Id, access.DB)

		for _, memgrp := range access_grant_memgrps {

			if memgrp.AccessControlId == contentAccess.Id {

				var memberGroup member.TblMemberGroup

				Accessmodel.GetMemberGroupsByContentAccessMemId(&memberGroup, memgrp.MemberGroupId, access.DB)

				contentAccess.MemberGroups = append(contentAccess.MemberGroups, memberGroup)
			}
		}

		var entriesCount, pageCount int64

		Accessmodel.GetaccessGrantedEntriesCount(&entriesCount, contentAccess.Id, access.DB)

		Accessmodel.GetaccessGrantedPageCount(&pageCount, contentAccess.Id, access.DB)

		if entriesCount > 0 {

			contentAccess.AccessGrantedModules = append(contentAccess.AccessGrantedModules, "Channel")
		}

		if pageCount > 0 {

			contentAccess.AccessGrantedModules = append(contentAccess.AccessGrantedModules, "Space")
		}

		if !contentAccess.ModifiedOn.IsZero() {

			contentAccess.DateString = contentAccess.ModifiedOn.UTC().Format("02 Jan 2006 03:04 PM")

		} else {

			contentAccess.DateString = contentAccess.CreatedOn.UTC().Format("02 Jan 2006 03:04 PM")

		}

		final_content_accesslist = append(final_content_accesslist, contentAccess)

	}

	_, content_access_count, _ := Accessmodel.GetContentAccessList(0, 0, filter, access.DB)

	return final_content_accesslist, content_access_count, nil

}

/*Get Access by id*/
func (access *AccessControl) GetControlAccessById(accessid int) (accesslist Tblaccesscontrol, err error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return Tblaccesscontrol{}, autherr
	}

	AccessControl, _ := Accessmodel.GetContentAccessByAccessId(accessid, access.DB)

	return *AccessControl, nil

}

/**/
func (access *AccessControl) GetselectedPageByAccessControlId(accessid int) ([]Page, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []Page{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetPagesAndPageGroupsInContentAccess(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var pages []Page

	for _, val := range contentAccessPages {

		var pg Page

		pg.Id = strconv.Itoa(val.PageId)

		pg.GroupId = strconv.Itoa(val.PageGroupId)

		pg.SpaceId = strconv.Itoa(val.SpacesId)

		pages = append(pages, pg)
	}

	return pages, nil
}

/**/
func (access *AccessControl) GetselectedGroupByAccessControlId(accessid int) ([]PageGroup, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []PageGroup{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetPageGroupsInContentAccess(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var group []PageGroup

	for _, val := range contentAccessPages {

		var pg PageGroup

		pg.Id = strconv.Itoa(val.PageId)

		pg.SpaceId = strconv.Itoa(val.SpacesId)

		group = append(group, pg)
	}

	return group, nil
}

/**/
func (access *AccessControl) GetselectedSubPageByAccessControlId(accessid int) ([]SubPage, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []SubPage{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetPagesAndPageGroupsInContentAccess(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var pages []SubPage

	for _, val := range contentAccessPages {

		var pg SubPage

		pg.Id = strconv.Itoa(val.PageId)

		pg.GroupId = strconv.Itoa(val.PageGroupId)

		pg.ParentId = strconv.Itoa(val.ParentPageId)

		pg.SpaceId = strconv.Itoa(val.SpacesId)

		pages = append(pages, pg)
	}

	return pages, nil
}

/**/
func (access *AccessControl) GetselectedSpacesByAccessControlId(accessid int) ([]string, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []string{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetSelectedSpaces(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var spaces []string

	for _, val := range contentAccessPages {

		spaces = append(spaces, strconv.Itoa(val.SpacesId))
	}

	return spaces, nil
}

/**/
func (access *AccessControl) GetselectedChannelByAccessControlId(accessid int) ([]string, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []string{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetSelectedSpaces(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var spaces []string

	for _, val := range contentAccessPages {

		spaces = append(spaces, strconv.Itoa(val.SpacesId))
	}

	return spaces, nil
}

/**/
func (access *AccessControl) GetselectedEntiresByAccessControlId(accessid int) ([]string, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []string{}, autherr
	}

	contentAccessPages, err := Accessmodel.GetSelectedSpaces(accessid, access.DB)

	if err != nil {

		log.Println(err)
	}

	var spaces []string

	for _, val := range contentAccessPages {

		spaces = append(spaces, strconv.Itoa(val.SpacesId))
	}

	return spaces, nil
}


func (access *AccessControl) CreateRestrictPage(accessid int, membergroups []int, ids []int, createdBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grps []TblAccessControlUserGroup

	for _, val := range membergroups {

		var membergrp TblAccessControlUserGroup

		membergrp.AccessControlId = accessid

		membergrp.MemberGroupId = val

		membergrp.CreatedBy = createdBy

		membergrp.CreatedOn = CurrentTime

		acces, err := Accessmodel.CreateMemberGroupRestrict(membergrp, access.DB)

		if err != nil {

			log.Println(err)
		}

		grps = append(grps, acces)

	}

	for _, grp := range grps {

		for _, val := range ids {

			var page TblAccessControlPages

			page.AccessControlUserGroupId = grp.Id

			page.PageId = val

			page.CreatedBy = createdBy

			page.CreatedOn = CurrentTime

			Accessmodel.CreatePage(&page, access.DB)

		}

	}

	return nil
}


func (access *AccessControl) CreateRestrictGroup(accessid int, membergroups []int, ids []int, createdBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grps []TblAccessControlUserGroup

	for _, val := range membergroups {

		var membergrp TblAccessControlUserGroup

		membergrp.AccessControlId = accessid

		membergrp.MemberGroupId = val

		membergrp.CreatedBy = createdBy

		membergrp.CreatedOn = CurrentTime

		acces, err := Accessmodel.CreateMemberGroupRestrict(membergrp, access.DB)

		if err != nil {

			log.Println(err)
		}

		grps = append(grps, acces)

	}

	for _, grp := range grps {

		for _, val := range ids {

			var page TblAccessControlPages

			page.AccessControlUserGroupId = grp.Id

			page.PageGroupId = val

			page.CreatedBy = createdBy

			page.CreatedOn = CurrentTime

			Accessmodel.CreatePage(&page, access.DB)

		}

	}

	return nil
}

func (access *AccessControl) CreateRestrictSubPage(accessid int, membergroups []int, ids []int, createdBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grps []TblAccessControlUserGroup

	for _, val := range membergroups {

		var membergrp TblAccessControlUserGroup

		membergrp.AccessControlId = accessid

		membergrp.MemberGroupId = val

		membergrp.CreatedBy = createdBy

		membergrp.CreatedOn = CurrentTime

		acces, err := Accessmodel.CreateMemberGroupRestrict(membergrp, access.DB)

		if err != nil {

			log.Println(err)
		}

		grps = append(grps, acces)

	}

	for _, grp := range grps {

		for _, val := range ids {

			var page TblAccessControlPages

			page.AccessControlUserGroupId = grp.Id

			page.PageId = val

			page.CreatedBy = createdBy

			page.CreatedOn = CurrentTime

			Accessmodel.CreatePage(&page, access.DB)

		}

	}

	return nil

}

/**/
func (access *AccessControl) DeleteSeletedPage(accessid int, ids []int, DeletedBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB)

	if err != nil {

		return err
	}

	for _, val := range grps {

		grpsid = append(grpsid, val.Id)

	}

	var acc TblAccessControlPages

	acc.IsDeleted = 0

	acc.DeletedBy = DeletedBy

	acc.DeletedOn = CurrentTime

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) DeleteSeletedGroup(accessid int, ids []int, DeletedBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB)

	if err != nil {

		return err
	}

	for _, val := range grps {

		grpsid = append(grpsid, val.Id)

	}

	var acc TblAccessControlPages

	acc.IsDeleted = 0

	acc.DeletedBy = DeletedBy

	acc.DeletedOn = CurrentTime

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) DeleteSelectedSpaces(accessid int, ids []int, DeletedBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB)

	if err != nil {

		return err
	}

	for _, val := range grps {

		grpsid = append(grpsid, val.Id)

	}

	var acc TblAccessControlPages

	acc.IsDeleted = 0

	acc.DeletedBy = DeletedBy

	acc.DeletedOn = CurrentTime

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) UpdateAccessControl(accessid int, title string, ModifiedBy int) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var acc TblAccessControl

	acc.Id = accessid

	acc.AccessControlName = title

	acc.AccessControlSlug = strings.ReplaceAll(strings.ToLower(title), " ", "-")

	acc.ModifiedBy = ModifiedBy

	acc.ModifiedOn = CurrentTime

	err := Accessmodel.UpdateContentAccessId(&acc, access.DB)

	if err != nil {

		return err
	}

	return nil
}
