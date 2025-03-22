package memberaccess

import (
	"log"
	"strconv"
	"strings"
	"time"

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
func (access *AccessControl) ContentAccessList(limit int, offset int, filter Filter, tenantid string) (accesslist []Tblaccesscontrol, totalCount int64, err error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []Tblaccesscontrol{}, 0, autherr
	}

	Accessmodel.DataAccess = access.DataAccess
	Accessmodel.UserId = access.UserId

	contentAccessList, _, errr := Accessmodel.GetContentAccessList(limit, offset, filter, access.DB, tenantid)

	if errr != nil {

		return []Tblaccesscontrol{}, 0, errr
	}

	var final_content_accesslist []Tblaccesscontrol

	for _, contentAccess := range contentAccessList {

		var access_grant_memgrps []TblAccessControlUserGroup

		Accessmodel.GetAccessGrantedMemberGroups(&access_grant_memgrps, contentAccess.Id, access.DB, tenantid)

		for _, memgrp := range access_grant_memgrps {

			if memgrp.AccessControlId == contentAccess.Id {

				var memberGroup member.TblMemberGroup

				Accessmodel.GetMemberGroupsByContentAccessMemId(&memberGroup, memgrp.MemberGroupId, access.DB, tenantid)

				contentAccess.MemberGroups = append(contentAccess.MemberGroups, memberGroup)
			}
		}

		var entriesCount, pageCount int64

		Accessmodel.GetaccessGrantedEntriesCount(&entriesCount, contentAccess.Id, access.DB, tenantid)

		Accessmodel.GetaccessGrantedPageCount(&pageCount, contentAccess.Id, access.DB, tenantid)

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

	_, content_access_count, _ := Accessmodel.GetContentAccessList(0, 0, filter, access.DB, tenantid)

	return final_content_accesslist, content_access_count, nil

}

/*Get Access by id*/
func (access *AccessControl) GetControlAccessById(accessid int, tenantid string) (accesslist Tblaccesscontrol, err error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return Tblaccesscontrol{}, autherr
	}

	AccessControl, _ := Accessmodel.GetContentAccessByAccessId(accessid, access.DB, tenantid)

	return *AccessControl, nil

}

/**/
// func (access *AccessControl) GetselectedPageByAccessControlId(accessid int, tenantid string) ([]Page, error) {

// 	autherr := AuthandPermission(access)

// 	if autherr != nil {

// 		return []Page{}, autherr
// 	}

// 	contentAccessPages, err := Accessmodel.GetPagesAndPageGroupsInContentAccess(accessid, access.DB, tenantid)

// 	if err != nil {

// 		log.Println(err)
// 	}

// 	var pages []Page

// 	for _, val := range contentAccessPages {

// 		var pg Page

// 		pg.Id = strconv.Itoa(val.PageId)

// 		pg.GroupId = strconv.Itoa(val.PageGroupId)

// 		pg.SpaceId = strconv.Itoa(val.SpacesId)

// 		pages = append(pages, pg)
// 	}

// 	return pages, nil
// }

/**/
// func (access *AccessControl) GetselectedGroupByAccessControlId(accessid int, tenantid string) ([]PageGroup, error) {

// 	autherr := AuthandPermission(access)

// 	if autherr != nil {

// 		return []PageGroup{}, autherr
// 	}

// 	contentAccessPages, err := Accessmodel.GetPageGroupsInContentAccess(accessid, access.DB, tenantid)

// 	if err != nil {

// 		log.Println(err)
// 	}

// 	var group []PageGroup

// 	for _, val := range contentAccessPages {

// 		var pg PageGroup

// 		pg.Id = strconv.Itoa(val.PageId)

// 		pg.SpaceId = strconv.Itoa(val.SpacesId)

// 		group = append(group, pg)
// 	}

// 	return group, nil
// }

/**/
// func (access *AccessControl) GetselectedSubPageByAccessControlId(accessid int, tenantid string) ([]SubPage, error) {

// 	autherr := AuthandPermission(access)

// 	if autherr != nil {

// 		return []SubPage{}, autherr
// 	}

// 	contentAccessPages, err := Accessmodel.GetPagesAndPageGroupsInContentAccess(accessid, access.DB, tenantid)

// 	if err != nil {

// 		log.Println(err)
// 	}

// 	var pages []SubPage

// 	for _, val := range contentAccessPages {

// 		var pg SubPage

// 		pg.Id = strconv.Itoa(val.PageId)

// 		pg.GroupId = strconv.Itoa(val.PageGroupId)

// 		pg.ParentId = strconv.Itoa(val.ParentPageId)

// 		pg.SpaceId = strconv.Itoa(val.SpacesId)

// 		pages = append(pages, pg)
// 	}

// 	return pages, nil
// }

/**/
// func (access *AccessControl) GetselectedSpacesByAccessControlId(accessid int, tenantid string) ([]string, error) {

// 	autherr := AuthandPermission(access)

// 	if autherr != nil {

// 		return []string{}, autherr
// 	}

// 	contentAccessPages, err := Accessmodel.GetSelectedSpaces(accessid, access.DB, tenantid)

// 	if err != nil {

// 		log.Println(err)
// 	}

// 	var spaces []string

// 	for _, val := range contentAccessPages {

// 		spaces = append(spaces, strconv.Itoa(val.SpacesId))
// 	}

// 	return spaces, nil
// }

/**/
// func (access *AccessControl) GetselectedChannelByAccessControlId(accessid int, tenantid string) ([]string, error) {

// 	autherr := AuthandPermission(access)

// 	if autherr != nil {

// 		return []string{}, autherr
// 	}

// 	contentAccessPages, err := Accessmodel.GetSelectedSpaces(accessid, access.DB, tenantid)

// 	if err != nil {

// 		log.Println(err)
// 	}

// 	var spaces []string

// 	for _, val := range contentAccessPages {

// 		spaces = append(spaces, strconv.Itoa(val.SpacesId))
// 	}

// 	return spaces, nil
// }

/**/
func (access *AccessControl) GetselectedEntiresByAccessControlId(accessid int, tenantid string) ([]int, []Entry, error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []int{}, []Entry{}, autherr
	}

	var channelEntries []Entry

	var contentAccessEntries []TblAccessControlPages

	Accessmodel.GetAccessGrantedEntries(&contentAccessEntries, accessid, access.DB, tenantid)

	channelMap := make(map[int][]TblAccessControlPages)

	for _, accessEntry := range contentAccessEntries {

		chanEntry := Entry{Id: strconv.Itoa(accessEntry.EntryId), ChannelId: strconv.Itoa(accessEntry.ChannelId), TenantId: accessEntry.TenantId}

		channelEntries = append(channelEntries, chanEntry)

		if _, exists := channelMap[accessEntry.ChannelId]; !exists {

			channelMap[accessEntry.ChannelId] = []TblAccessControlPages{}

		}

		channelMap[accessEntry.ChannelId] = append(channelMap[accessEntry.ChannelId], accessEntry)

	}

	var channelIds []int

	for channelId, entriesArr := range channelMap {

		var entriesCountInChannel int64

		Accessmodel.GetEntriesCountUnderChannel(&entriesCountInChannel, channelId, access.DB, tenantid)

		if int(entriesCountInChannel) == len(entriesArr) {

			channelIds = append(channelIds, channelId)
		}
	}

	return channelIds, channelEntries, nil
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

func (access *AccessControl) CreateRestrictGroup(accessid int, membergroups []int, ids []int, createdBy int, tenantid string) error {

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

		membergrp.TenantId = tenantid

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

			page.TenantId = tenantid

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
func (access *AccessControl) DeleteSeletedPage(accessid int, ids []int, DeletedBy int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB, tenantid)

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

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB, tenantid)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) DeleteSeletedGroup(accessid int, ids []int, DeletedBy int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB, tenantid)

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

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB, tenantid)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) DeleteSelectedSpaces(accessid int, ids []int, DeletedBy int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var grpsid []int

	grps, err := Accessmodel.GetGroupsByAccessId(accessid, access.DB, tenantid)

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

	er := Accessmodel.DeletePage(&acc, grpsid, ids, access.DB, tenantid)

	if er != nil {

		return er
	}

	return nil
}

func (access *AccessControl) UpdateAccessControl(accessid int, title string, ModifiedBy int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var acc TblAccessControl

	acc.Id = accessid

	acc.AccessControlName = title

	acc.AccessControlSlug = strings.ReplaceAll(strings.ToLower(title), " ", "-")

	acc.ModifiedBy = ModifiedBy

	acc.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Accessmodel.UpdateContentAccessId(&acc, access.DB, tenantid)

	if err != nil {

		return err
	}

	return nil
}

// Create Accesscontrol
func (access *AccessControl) CreateAccessControl(title string, ModifiedBy int, tenantid string) (accessdata TblAccessControl, aerr error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return accessdata, autherr
	}

	var acc TblAccessControl

	acc.AccessControlName = title

	acc.AccessControlSlug = strings.ReplaceAll(strings.ToLower(title), " ", "-")

	acc.CreatedBy = ModifiedBy

	acc.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	acc.TenantId = tenantid

	err := Accessmodel.NewContentAccessEntry(&acc, access.DB)

	if err != nil {

		return TblAccessControl{}, err
	}

	return acc, nil
}

// Delete Accesscontrol
func (access *AccessControl) DeleteMemberAccessControl(accessid int, ModifiedBy int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	var accesscontrol TblAccessControl

	accesscontrol.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	accesscontrol.DeletedBy = ModifiedBy

	accesscontrol.IsDeleted = 1

	err := Accessmodel.DeleteControlAccess(&accesscontrol, accessid, access.DB, tenantid)

	var acusergrp TblAccessControlUserGroup

	acusergrp.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	acusergrp.DeletedBy = ModifiedBy

	acusergrp.IsDeleted = 1

	Accessmodel.DeleteInAccessUserGroup(&acusergrp, accessid, access.DB, tenantid)

	var accessgrp []TblAccessControlUserGroup

	Accessmodel.GetDeleteIdInAccessUserGroup(&accessgrp, accessid, access.DB, tenantid)

	var pgid []int

	for _, v := range accessgrp {

		pgid = append(pgid, v.Id)

	}

	var accesscontrolpg TblAccessControlPages

	accesscontrolpg.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	accesscontrolpg.DeletedBy = ModifiedBy

	accesscontrolpg.IsDeleted = 1

	Accessmodel.DeleteAccessControlPages(&accesscontrolpg, pgid, access.DB, tenantid)

	if err != nil {

		return err

	}

	return nil
}

func (access *AccessControl) CreateRestrictEntries(accessid int, membergroups []int, entryids []Entry, createdBy int, tenantid string) error {

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

		membergrp.TenantId = tenantid

		acces, err := Accessmodel.CreateMemberGroupRestrict(membergrp, access.DB)

		if err != nil {

			log.Println(err)
		}

		grps = append(grps, acces)

	}

	for _, grp := range grps {

		for _, val := range entryids {

			var page TblAccessControlPages

			page.AccessControlUserGroupId = grp.Id

			page.ChannelId, _ = strconv.Atoi(val.ChannelId)

			page.EntryId, _ = strconv.Atoi(val.Id)

			page.CreatedBy = createdBy

			page.CreatedOn = CurrentTime

			page.TenantId = tenantid

			Accessmodel.CreatePage(&page, access.DB)

		}

	}

	return nil
}

// function used to retrieve the access granted member group list
func (access *AccessControl) GetaccessMemberGroup(accessid int, tenantid string) (group []int, err error) {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return []int{}, autherr
	}
	var accessGrantedMemgrps []int

	gerr := Accessmodel.GetAccessGrantedMemberGroupsList(&accessGrantedMemgrps, accessid, access.DB, tenantid)

	if gerr != nil {

		return []int{}, gerr
	}
	return accessGrantedMemgrps, nil
}

//function helps to update a member restrict access

func (access *AccessControl) UpdateRestrictEntries(accessid int, membergroups []int, entryids []Entry, userid int, tenantid string) error {

	autherr := AuthandPermission(access)

	if autherr != nil {

		return autherr
	}

	for _, memgrp_id := range membergroups {

		var access_count int64

		err := Accessmodel.CheckPresenceOfAccessGrantedMemberGroups(&access_count, memgrp_id, accessid, access.DB, tenantid)

		if err != nil {

			log.Println(err)

		}

		var memberGrpAccess TblAccessControlUserGroup

		memberGrpAccess.AccessControlId = accessid

		memberGrpAccess.MemberGroupId = memgrp_id

		if access_count == 0 {

			memberGrpAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			memberGrpAccess.CreatedBy = userid

			memberGrpAccess.IsDeleted = 0
			memberGrpAccess.TenantId = tenantid

			err = Accessmodel.GrantAccessToMemberGroups(&memberGrpAccess, access.DB, tenantid)

			if err != nil {

				log.Println(err)

			}

		} else if access_count == 1 {

			memberGrpAccess.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

			memberGrpAccess.ModifiedBy = userid

			err = Accessmodel.UpdateContentAccessMemberGroup(&memberGrpAccess, access.DB, tenantid)

			if err != nil {

				log.Println(err)

			}

		}

	}
	var MemGrpAccess []TblAccessControlUserGroup

	Accessmodel.GetMemberGrpByAccessControlId(&MemGrpAccess, accessid, access.DB, tenantid)

	var entryIds []int

	seen_entry := make(map[int]bool)

	for _, memgrp := range MemGrpAccess {
		for _, entry := range entryids {

			chanId, _ := strconv.Atoi(entry.ChannelId)

			entryId, _ := strconv.Atoi(entry.Id)

			var entryCount int64

			err := Accessmodel.CheckPresenceOfChannelEntriesInContentAccess(&entryCount, memgrp.Id, chanId, entryId, access.DB, tenantid)

			if err != nil {

				log.Println(err)
			}

			var channelAccess TblAccessControlPages

			channelAccess.AccessControlUserGroupId = memgrp.Id

			channelAccess.ChannelId = chanId

			channelAccess.EntryId = entryId

			channelAccess.TenantId = tenantid

			if entryCount == 0 {

				channelAccess.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				channelAccess.CreatedBy = userid

				channelAccess.IsDeleted = 0

				err = Accessmodel.CreatePage(&channelAccess, access.DB)

				if err != nil {

					log.Println(err)
				}

			} else if entryCount == 1 {

				channelAccess.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

				channelAccess.ModifiedBy = userid

				err = Accessmodel.UpdateAccessPage(&channelAccess, access.DB, tenantid)

				if err != nil {

					log.Println(err)
				}

			}

			if !seen_entry[entryId] {

				entryIds = append(entryIds, entryId)

				seen_entry[entryId] = true

			}
		}
	}
	var memgrp_access1 TblAccessControlUserGroup

	memgrp_access1.IsDeleted = 1

	memgrp_access1.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	memgrp_access1.DeletedBy = userid

	Accessmodel.RemoveMemberGroupsNotUnderContentAccessRights(&memgrp_access1, membergroups, accessid, access.DB)

	for _, memgrp := range MemGrpAccess {

		var pg_access1 TblAccessControlPages

		pg_access1.AccessControlUserGroupId = memgrp.Id

		pg_access1.IsDeleted = 1

		pg_access1.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

		pg_access1.DeletedBy = userid

		err := Accessmodel.RemoveChannelEntriesNotUnderContentAccess(&pg_access1, entryIds, access.DB, tenantid)

		if err != nil {

			log.Println(err)
		}

	}

	return nil
}
