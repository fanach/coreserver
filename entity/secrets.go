package entity

const (
	// HidenString ...
	HidenString = "***"
)

// Password is wrapper for string, display "***" when print
type Password string

func (p Password) String() string {
	if len(p) == 0 {
		return ""
	}
	return HidenString
}
