package service

import (
	"errors"
	"log"
	"time"

	"github.com/dchest/uniuri"
	"github.com/fanach/coreserver/entity"
	"github.com/fanach/coreserver/util"
)

var (
	// SessionAliveDuration is a period sessions last for
	SessionAliveDuration = 24 * time.Hour
	// SessionGCDuration is a period sessions will be cleaned up
	SessionGCDuration = 1 * time.Hour
	// ErrIncorrectLogin returns when login failed with wrong username or password
	ErrIncorrectLogin = errors.New("incorrect username or password")
)

// CreateSess creates session
func CreateSess(username, password string) (sess *entity.Session, err error) {
	// validate username and password
	err = validate(username, password)
	if err != nil {
		log.Printf("validate username and password failed: %v\n", err)
		return nil, err
	}

	// new session
	sess = &entity.Session{}
	sess.SessID = uniuri.New()

	return
}

// DeleteSess deletes session
func DeleteSess(sessID string) (err error) {
	// do nothing now
	return nil
}

func validate(username, password string) (err error) {
	user, err := queryUserByName(username)
	if err != nil {
		err = ErrUserNotFound
		return
	}
	if user != nil && compareBCryptPassword(string(user.Password), password) == nil {
		// ok
		return
	}
	err = ErrIncorrectLogin
	return
}

func compareBCryptPassword(hash, password string) (err error) {
	return util.CompareBCryptPassword(hash, password)
}
