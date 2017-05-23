package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/fanach/coreserver/db"
	"github.com/fanach/coreserver/entity"
	"github.com/fanach/coreserver/util"
)

var (
	userdb *leveldb.DB

	// ErrUsernameConflict returns when username already exist while registering
	ErrUsernameConflict = errors.New("duplicated username")
	// ErrUserNotFound returns when user does not exist
	ErrUserNotFound = errors.New("user not found")
)

func initUserDB() error {
	cfg := db.LevelDBConfig{
		DBFile: db.UserDBFile,
	}

	ldb, err := db.NewLevelDB(cfg)
	if err != nil {
		return err
	}
	userdb = ldb.DB
	return nil
}

// CreateUser creates a new user
func CreateUser(user entity.ReqPostUser) (newUser *entity.User, err error) {
	newUser = &entity.User{}
	newUser.Username = user.Username

	bcryptPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("hash password error: %v\n", err)
		return
	}
	newUser.Password = bcryptPassword
	newUser.WeChatID = user.WeChatID
	newUser.Type = user.Type
	newUser.Email = user.Email

	newUser.ID = genUserID(newUser.Username)
	newUser.RegTime = time.Now()

	// check if duplicated
	foundUser, err := queryUserByName(newUser.Username)
	if err != nil && err != leveldb.ErrNotFound {
		log.Printf("query user by name error: %v\n", err)
		return
	}

	if foundUser != nil && len(foundUser.ID) > 0 {
		// log.Printf("duplicated username: %s\n", user.Username)
		return nil, ErrUsernameConflict
	}

	// save
	err = saveUser(*newUser)
	if err != nil {
		log.Printf("save user error: %v\n", err)
		return
	}

	return
}

// GetUser returns user by ID
// return ErrUserNotFound if user does not exist
func GetUser(userID string) (user *entity.User, err error) {
	return getUser(userID)
}

// GetUsers return all users
func GetUsers() (user *[]entity.User, err error) {
	return getUsers()
}

// UpdateUser updates an user
func UpdateUser(userID string, user entity.ReqPutUser) (modifiedUser *entity.User, err error) {
	modifiedUser, err = getUser(userID)
	if err != nil {
		log.Printf("get user by id %s error: %v", userID, err)
		return
	}

	modifiedUser.UpdateTime = time.Now()
	if len(user.Email) > 0 {
		modifiedUser.Email = user.Email
	}
	if len(user.WeChatID) > 0 {
		modifiedUser.WeChatID = user.WeChatID
	}
	if len(user.Type) > 0 {
		modifiedUser.Type = user.Type
	}
	if len(user.Password) > 0 {
		bcryptPassword, err2 := hashPassword(user.Password)
		if err2 != nil {
			log.Printf("hash password error: %v\n", err2)
			return
		}
		modifiedUser.Password = bcryptPassword
	}

	err = saveUser(*modifiedUser)
	if err != nil {
		log.Printf("update user error: %v\n", err)
		return
	}

	return
}

// DeleteUser deletes a user
func DeleteUser(userID string) (err error) {
	return deleteUser(userID)
}

func saveUser(user entity.User) (err error) {
	v, err := json.Marshal(user)
	if err != nil {
		return
	}
	return userdb.Put([]byte(user.ID), v, nil)
}

func deleteUser(userID string) (err error) {
	return userdb.Delete([]byte(userID), nil)
}

func getUser(userID string) (user *entity.User, err error) {
	data, err := userdb.Get([]byte(userID), nil)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return
	}
	return
}

func getUsers() (pUsers *[]entity.User, err error) {
	iter := userdb.NewIterator(nil, nil)
	user := entity.User{}
	users := []entity.User{}
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.

		// key := iter.Key()
		// value := iter.Value()

		if err = json.Unmarshal(iter.Value(), &user); err != nil {
			log.Printf("unmarshal user error: %v\n", err)
			continue
		}

		users = append(users, user)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Printf("iterator over userdb error: %v\n", err)
		return
	}

	pUsers = &users
	return
}

func queryUserByName(username string) (user *entity.User, err error) {
	return getUser(genUserID(username))
}

func hashPassword(password entity.Password) (hashPassword entity.Password, err error) {
	bcryptPassword, err := util.BCrypt(string(password))
	return entity.Password(bcryptPassword), err
}

func genUserID(username string) (id string) {
	return util.MD5sum(username)
}
