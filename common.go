package memberaccess

import (
	"errors"
	"time"
)

var (
	ErrorAuth       = errors.New("auth enabled not initialised")
	ErrorPermission = errors.New("permissions enabled not initialised")
	CurrentTime, _  = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	Empty           string
)

func TruncateDescription(description string, limit int) string {
	if len(description) <= limit {
		return description
	}

	truncated := description[:limit] + "..."
	return truncated
}


func AuthandPermission(accessControl *AccessControl) error {

	//check auth enable if enabled, use auth pkg otherwise it will return error
	if accessControl.AuthEnable && !accessControl.Auth.AuthFlg {

		return ErrorAuth
	}
	//check permission enable if enabled, use team-role pkg otherwise it will return error
	if accessControl.PermissionEnable && !accessControl.Auth.PermissionFlg {

		return ErrorPermission

	}

	return nil
}