package activity

import (
	"log/syslog"
)

// EventIDs
const (
	Auth Activity = 1 << iota
	AddUser
)

// Auth statuses
const (
	AuthUserNotRegistered Activity = 1 << (iota + 32)
	AuthUserLocked
	AuthUserOutsider
	AuthOrgLocked
	AuthCertError
	AuthAnotherError
)

// AddUser statuses
const (
	AuthUserNotRegistered Activity = 1 << (iota + 32)
	AuthUserLocked
	AuthUserOutsider
	AuthOrgLocked
	AuthCertError
	AuthAnotherError
)

// Activity ...
type Activity uint64

func (a Activity) String() string {
	switch a.eventID() {
	case Auth:
		return "Authentication " + a.authStatus()
	}
	return ""
}

func (a Activity) authStatus() string {
	if a.status() == 0 {
		return "OK"
	}
	s := "failed: "
	switch a.status() {
	case AuthUserNotRegistered:
		s += "user is not registered"
	case AuthUserLocked:
		s += "user is locked"
	case AuthUserOutsider:
		s += "user is outsider"
	case AuthOrgLocked:
		s += "user company is locked"
	case AuthCertError:
		s += "user certificate error"
	default:
		s += "another reason"
	}
	return s
}

func (a Activity) severity() syslog.Priority {
	if a.status() == 0 {
		return syslog.LOG_INFO
	}
	return syslog.LOG_ERR
}

func (a Activity) eventID() Activity {
	return a & 0xFFFFFFFF
}

func (a Activity) status() Activity {
	return a & 0xFFFFFFFF00000000
}
