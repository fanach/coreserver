package db

import (
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	// DefaultDBFile is the default path of LevelDB database file
	DefaultDBFile = "./db/default.leveldb"
	// UserDBFile is the path of user database file
	UserDBFile = "./db/user.leveldb"
	// SessDBFile is the path of session database file
	SessDBFile = "./db/sess.leveldb"
)

// LevelDB is a simple wrapper for leveldb.DB
type LevelDB struct {
	DB            *leveldb.DB
	LevelDBConfig *LevelDBConfig
}

// LevelDBConfig contains some options of leveldb
type LevelDBConfig struct {
	DBFile       string
	Options      *opt.Options
	ReadOptions  *opt.ReadOptions
	WriteOptions *opt.WriteOptions
}

// NewLevelDB opens LevelDB database file and init LevelDB instance
func NewLevelDB(config LevelDBConfig) (ldb *LevelDB, err error) {
	ldb = &LevelDB{}
	if strings.TrimSpace(config.DBFile) == "" {
		config.DBFile = DefaultDBFile
	}
	db, err := leveldb.OpenFile(config.DBFile, config.Options)
	if err != nil {
		log.Printf("open db file error: %s\n", err)
		return
	}
	if db != nil {
		ldb.DB = db
	}
	return
}

// Use usesLevel given LevelDB database file
func (ldb *LevelDB) Use(file string) (db *leveldb.DB, err error) {
	return leveldb.OpenFile(file, ldb.LevelDBConfig.Options)
}
