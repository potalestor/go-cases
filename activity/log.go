package activity

import (
	"sync"
)

var (
	// TheConfig contains activity settings
	TheConfig *Config
	once      sync.Once
	h         handler
)

// Log is used for writting events
//
//  - activity is a combination of EventID and status.
//	  For example, "activity = activity.Auth | activity.AuthUserLocked" means
//	  that the authentication is failed because the user was locked
//    Or "activity = activity.Auth" means that the authentication is success
//
//  - data is an optional argument. It can consist of a tuple of arbitrary primitive types
//    and / or activity.Param structures
//
//  Examples:
//   - activity.Log(activity.Auth, activity.NewParam("UserID", "user___012345"), activity.NewParam("Ip", "192.168.0.12"))
//   - activity.Log(activity.Auth | activity.AuthUserLocked, activity.NewParam("UserID", "user___012345"), activity.NewParam("Ip", "192.168.0.12"))
func Log(activity Activity, data ...interface{}) {
	// initialize logger at once
	once.Do(func() {
		h = handlerFactory(TheConfig)
	})
	if h == nil {
		return
	}
	h.log(activity, data)
}
