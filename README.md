# Member-access Package

The 'Member-access' package grants website admin the authority to provide access to the members.This streamlined process enables effortless audience management and curation, ensuring a personalized experience for each member within the website. 

## Features



# Installation

``` bash
go get github.com/spurtcms/Member-access
```


# Usage example

``` bash
import (
	"github.com/spurtcms/auth"
	"github.com/spurtcms/member-access"
)

func main() {

	Auth := auth.AuthSetup(auth.Config{
		UserId:     1,
		ExpiryTime: 2,
		SecretKey:  "SecretKey@123",
		DB: &gorm.DB{},
		RoleId: 1,
	})

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Member Restrict", auth.CRUD)

	accesscontroller := access.AccessSetup(access.Config{
		DB:               &gorm.DB{},
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})

	//contentaccess
	if permisison {

		//list contentaccess
		contentaccesslist, count, err := accesscontroller.ContentAccessList(10, 0, access.Filter{})
		fmt.Println(contentaccesslist, count, err)

		//create contentaccess
		cerr := accesscontroller.CreateAccessControl("Demo Entries Access", 1)

		if cerr != nil {

			fmt.Println(cerr)
		}

		// update contentaccess
		uerr := accesscontroller.UpdateAccessControl(1, "Entries Access", 1)

		if uerr != nil {

			fmt.Println(uerr)
		}

		// delete contentaccess
		derr := accesscontroller.DeleteMemberAccessControl(2, 1)

		if derr != nil {

			fmt.Println(derr)
		}

		// list selectedpage
		pagelist, lerr := accesscontroller.GetselectedPageByAccessControlId(1)
		fmt.Println(pagelist, lerr)

		// create restrictgroup
		gerr := accesscontroller.CreateRestrictGroup(1, []int{1, 2}, []int{1, 2, 3}, 1)

		if gerr != nil {

			fmt.Println(gerr)
		}
	}

}
```
# Getting help
If you encounter a problem with the package,please refer [Please refer [(https://www.spurtcms.com/documentation/cms-admin)] or you can create a new Issue in this repo[https://github.com/spurtcms/member-access/issues]. 

