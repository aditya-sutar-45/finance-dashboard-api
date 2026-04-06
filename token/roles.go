package token

type Role string

const (
	RoleViewer  Role = "viewer"
	RoleAnalyst Role = "analyst"
	RoleAdmin   Role = "admin"
)

var RoleHierarchy = map[Role]int{
	RoleViewer:  1,
	RoleAnalyst: 2,
	RoleAdmin:   3,
}
