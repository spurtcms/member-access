package memberaccess

import (
	"fmt"

	"github.com/spurtcms/auth"
	"gorm.io/gorm"
)

func Testing() {

	Auth := auth.AuthSetup(auth.Config{
		UserId:     1,
		ExpiryTime: 2,
		SecretKey:  "Secret123",
	})

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, "Secret123")

	permisison, _ := Auth.IsGranted("Member Restrict", auth.CRUD, TenantId)

	accesscontroller := AccessSetup(Config{
		DB:               &gorm.DB{},
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})

	//contentaccess
	if permisison {

		//list contentaccess
		contentaccesslist, count, err := accesscontroller.ContentAccessList(10, 0, Filter{}, TenantId)
		fmt.Println(contentaccesslist, count, err)

		//create contentaccess
		_, cerr := accesscontroller.CreateAccessControl("Demo Entries Access", 1, 1)

		if cerr != nil {

			fmt.Println(cerr)
		}

		// update contentaccess
		uerr := accesscontroller.UpdateAccessControl(1, "Entries Access", 1, 1)

		if uerr != nil {

			fmt.Println(uerr)
		}

		// delete contentaccess
		derr := accesscontroller.DeleteMemberAccessControl(2, 1, 1)

		if derr != nil {

			fmt.Println(derr)
		}

		// list selectedpage
		// pagelist, lerr := accesscontroller.GetselectedPageByAccessControlId(1, 1)
		// fmt.Println(pagelist, lerr)

		// create restrictgroup
		gerr := accesscontroller.CreateRestrictGroup(1, []int{1, 2}, []int{1, 2, 3}, 1, TenantId)

		if gerr != nil {

			fmt.Println(gerr)
		}
	}

}
