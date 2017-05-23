package conf

import (
	"errors"
	"os"
	"strconv"
)

const (
	// keys of environment variables

	// LEVELDB_DIR Directory of LevelDB data files
	LEVELDB_DIR Env = "LEVELDB_DIR"
)

var (
	// ErrNotSet returned when the key of env is not set
	ErrNotSet = errors.New("env not set")
)

// Env is struct of environment variable
type Env string

// ToInt parse value to int
func (e *Env) ToInt() (int, error) {
	val, found := lookup(e)
	if !found {
		return 0, ErrNotSet
	}

	i, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ToBool parse value to bool
func (e *Env) ToBool() (bool, error) {
	val, found := lookup(e)
	if !found {
		return false, ErrNotSet
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return b, nil
}

// String convert Env to string
func (e *Env) String() string {
	return string(*e)
}

func lookup(e *Env) (val string, found bool) {
	return os.LookupEnv(e.String())
}
