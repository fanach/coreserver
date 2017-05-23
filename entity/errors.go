package entity

// Err is the self-defined error
type Err struct {
	ErrNo  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}
