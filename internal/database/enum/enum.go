package enum

type RoleID int16

const (
	RoleOwner RoleID = iota + 1
	RoleManager
)
