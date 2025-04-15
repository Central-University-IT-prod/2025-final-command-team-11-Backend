package models

type Role string

const (
	RoleUser       Role = "USER"
	RoleAdmin      Role = "ADMIN"
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleSupport    Role = "SUPPORT"
)

func ParseRole(role string) (Role, error) {
	switch role {
	case string(RoleUser):
		return RoleUser, nil
	case string(RoleAdmin):
		return RoleAdmin, nil
	case string(RoleSuperAdmin):
		return RoleSuperAdmin, nil
	case string(RoleSupport):
		return RoleSupport, nil
	default:
		return "", ErrInvalidRole
	}
}
