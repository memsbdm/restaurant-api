package enum

type RoleID int16

const (
	RoleAdmin RoleID = iota + 1
	RoleUser
)
