package memberaccess

import (
	"strconv"

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

	var pages []Page

	var subpages []SubPage

	var pagegroups []PageGroup

	contentAccessPages, err := Accessmodel.GetPagesAndPageGroupsInContentAccess(accessid, access.DB)

	var pageArrContainer [][]TblPage

	var pgArrContainer []TblAccessControlPages

	seen1 := make(map[int]bool)

	seen2 := make(map[int]bool)

	seen3 := make(map[int]bool)

	seen4 := make(map[int]bool)

	for _, pagz := range contentAccessPages {

		if pagz.ParentPageId == 0 {

			if !seen1[pagz.PageId] {

				var pg Page

				pg.Id = strconv.Itoa(pagz.PageId)

				pg.GroupId = strconv.Itoa(pagz.PageGroupId)

				pg.SpaceId = strconv.Itoa(pagz.SpacesId)

				pages = append(pages, pg)

				seen1[pagz.PageId] = true

			}
		}

		if pagz.ParentPageId != 0 {

			if !seen2[pagz.PageId] {

				var spg SubPage

				spg.Id = strconv.Itoa(pagz.PageId)

				spg.GroupId = strconv.Itoa(pagz.PageGroupId)

				spg.ParentId = strconv.Itoa(pagz.ParentPageId)

				spg.SpaceId = strconv.Itoa(pagz.SpacesId)

				subpages = append(subpages, spg)

				seen2[pagz.PageId] = true
			}

		}

		if pagz.PageGroupId != 0 {

			if !seen3[pagz.PageGroupId] {

				var pagesinPgg []TblPage

				Accessmodel.GetPagesUnderPageGroup(&pagesinPgg, pagz.PageGroupId, access.DB)

				pageArrContainer = append(pageArrContainer, pagesinPgg)

				seen3[pagz.PageGroupId] = true

			}

			if !seen4[pagz.PageId] {

				pgArrContainer = append(pgArrContainer, pagz)

				seen4[pagz.PageId] = true

			}

		}

	}

	// log.Println("orgpgg", pageArrContainer)

	// log.Println("pggchk", pgArrContainer)

	groupedObjects := make(map[int][]TblAccessControlPages)

	OriginalPgg := make(map[int][]TblPage)

	for _, pgz := range pgArrContainer {

		if _, exists := groupedObjects[pgz.PageGroupId]; !exists {

			groupedObjects[pgz.PageGroupId] = []TblAccessControlPages{}

		}

		groupedObjects[pgz.PageGroupId] = append(groupedObjects[pgz.PageGroupId], pgz)
	}

	for _, pggArr := range pageArrContainer {

		for _, pages := range pggArr {

			OriginalPgg[pages.PageGroupId] = pggArr

		}

	}

	for pggId, array := range OriginalPgg {

		for pggid := range groupedObjects {

			if len(OriginalPgg[pggId]) == len(groupedObjects[pggid]) && pggId == pggid {

				for index, result := range array {

					if index == 0 {

						var pgg PageGroup

						pgg.Id = strconv.Itoa(result.PageGroupId)

						pgg.SpaceId = strconv.Itoa(result.SpacesId)

						pagegroups = append(pagegroups, pgg)

						break
					}

				}

			}
		}
	}

	// var access_grant_memgrps_list []int

	var accessGrantedMemgrps []int

	Accessmodel.GetAccessGrantedMemberGroupsList(&accessGrantedMemgrps, accessid, access.DB)

	var spaceIds []int

	Accessmodel.GetContentAccessSpaces(&spaceIds, accessid, access.DB)

	var accessSpaceIds []int

	for _, spaceId := range spaceIds {

		var contentAccessPages []int

		Accessmodel.GetcontentAccessPagesBySpaceId(&contentAccessPages, spaceId, accessid, access.DB)

		var tblPageData []TblPage

		Accessmodel.GetPagesUnderSpaces(&tblPageData, spaceId, access.DB)

		if len(tblPageData) == len(contentAccessPages) {

			accessSpaceIds = append(accessSpaceIds, spaceId)
		}

	}

	var channelEntries []Entry

	var contentAccessEntries []TblAccessControlPages

	Accessmodel.GetAccessGrantedEntries(&contentAccessEntries, accessid, access.DB)

	channelMap := make(map[int][]TblAccessControlPages)

	for _, accessEntry := range contentAccessEntries {

		chanEntry := Entry{Id: strconv.Itoa(accessEntry.EntryId), ChannelId: strconv.Itoa(accessEntry.ChannelId)}

		channelEntries = append(channelEntries, chanEntry)

		if _, exists := channelMap[accessEntry.ChannelId]; !exists {

			channelMap[accessEntry.ChannelId] = []TblAccessControlPages{}

		}

		channelMap[accessEntry.ChannelId] = append(channelMap[accessEntry.ChannelId], accessEntry)

	}

	var channelIds []int

	for channelId, entriesArr := range channelMap {

		var entriesCountInChannel int64

		Accessmodel.GetEntriesCountUnderChannel(&entriesCountInChannel, channelId, access.DB)

		if int(entriesCountInChannel) == len(entriesArr) {

			channelIds = append(channelIds, channelId)
		}
	}

	return AccessControl, pages, subpages, pagegroups, accessGrantedMemgrps, accessSpaceIds, channelIds, channelEntries, nil

}

/**/
func (access *AccessControl) GetselectedPageAccessControl() {

}

