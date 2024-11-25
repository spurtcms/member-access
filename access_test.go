package memberaccess

import (
	"fmt"
	"log"
	"testing"

	"github.com/spurtcms/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var SecretKey = "Secret123"

// Db connection
func DBSetup() (*gorm.DB, error) {

	dbConfig := map[string]string{
		"username": "postgres",
		"password": "postgres",
		"host":     "localhost",
		"port":     "5432",
		"dbname":   "nov_16",
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user=" + dbConfig["username"] + " password=" + dbConfig["password"] +
			" dbname=" + dbConfig["dbname"] + " host=" + dbConfig["host"] +
			" port=" + dbConfig["port"] + " sslmode=disable TimeZone=Asia/Kolkata",
	}), &gorm.Config{})

	if err != nil {

		log.Fatal("Failed to connect to database:", err)

	}
	if err != nil {

		return nil, err

	}

	return db, nil
}

// test contentaccesslist function
func TestContentAccessList(t *testing.T) {

	db, _ := DBSetup()

	config := auth.Config{
		UserId: 1,
		// ExpiryTime: 2,
		ExpiryFlg: false,
		SecretKey: "Secret123",
		DB:        db,
		RoleId:    1,
	}

	Auth := auth.AuthSetup(config)

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Member Restrict", auth.CRUD, 1)

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		contentaccesslist, count, err := accesscontroller.ContentAccessList(10, 0, Filter{}, 1)

		if err != nil {

			panic(err)
		}

		fmt.Println(contentaccesslist, count)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test getcontrolaccessbyid function
func TestGetControlAccessById(t *testing.T) {

	db, _ := DBSetup()

	config := auth.Config{
		UserId: 1,
		// ExpiryTime: 2,
		ExpiryFlg: false,
		SecretKey: "Secret123",
		DB:        db,
		RoleId:    1,
	}

	Auth := auth.AuthSetup(config)

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Member Restrict", auth.CRUD, 1)

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		contentaccess, err := accesscontroller.GetControlAccessById(1, 1)

		if err != nil {

			panic(err)
		}

		fmt.Println(contentaccess)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test getselectedentriesbyaccesscontrolid function
func TestGetselectedEntiresByAccessControlId(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})
	channelIds, entries, err := accesscontroller.GetselectedEntiresByAccessControlId(1, 1)

	if err != nil {

		panic(err)
	}

	fmt.Println("entries", entries)
	fmt.Println("channelIds", channelIds)

}

// test createrestrictpage function
func TestCreateRestrictPage(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := accesscontroller.CreateRestrictPage(1, []int{1, 2}, []int{1, 2, 3}, 1)

	if err != nil {

		panic(err)
	}

}

// test createrestrictgroup function
func TestCreateRestrictGroup(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := accesscontroller.CreateRestrictGroup(1, []int{1, 2}, []int{1, 2, 3}, 1, 1)

	if err != nil {

		panic(err)
	}

}

// test createrestrictsubpage function
func TestCreateRestrictSubPage(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := accesscontroller.CreateRestrictSubPage(1, []int{1, 2}, []int{1, 2, 3}, 1)

	if err != nil {

		panic(err)
	}

}

// test deleteselectedpage function
func TestDeleteSeletedPage(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := accesscontroller.DeleteSeletedPage(1, []int{1, 2, 3}, 1, 1)

	if err != nil {

		panic(err)
	}

}

// test deleteselectedgroup function
func TestDeleteSeletedGroup(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

		err := accesscontroller.DeleteSeletedGroup(1,[]int{1,2},1,1)

		if err != nil {

			panic(err)
		}

}

// test updateaccesscontrol function
func TestUpdateAccessControl(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

		err := accesscontroller.UpdateAccessControl(1,"sports",1,1)

		if err != nil {

			panic(err)
		}

}

// test createaccesscontrol function
func TestCreateAccessControl(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

		_,err := accesscontroller.CreateAccessControl("Demo Entries Access",1,1)

		if err != nil {

			panic(err)
		}

}

// test deleteaccesscontrol function
func TestDeleteMemberAccessControl(t *testing.T) {

	db, _ := DBSetup()

	accesscontroller := AccessSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

		err := accesscontroller.DeleteMemberAccessControl(2,1,1)

		if err != nil {

			panic(err)
		}

}
