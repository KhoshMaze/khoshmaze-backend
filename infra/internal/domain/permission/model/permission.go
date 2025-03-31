package model

/*

	How permission system works:

	Each user has a permission attribute which is a number. This number is sum of all of user permissions.
	Since the sum of any K power of 2 numbers is unique, system finds each user permissions (The binary
	version 1 bits represent the permissions of user)


	This system can be mixed with Roles for more specific and stronger permission system.
*/

type Permission uint64

const (
	ReadUser    Permission = 1 << 0 // 1
	WriteUser   Permission = 1 << 1 // 2
	DeleteUser  Permission = 1 << 2 // 4
	ReadAdmin   Permission = 1 << 3 // 8
	WriteAdmin  Permission = 1 << 4 // 16
	DeleteAdmin Permission = 1 << 5 // 32
)

// UserPermissions represents all permissions a user has
type UserPermissions uint64

// HasPermission checks if the user has a specific permission
func (p UserPermissions) HasPermission(permission Permission) bool {
	return uint64(p)&uint64(permission) == uint64(permission)
}

// HasAllPermissions checks if the user has all specified permissions
func (p UserPermissions) HasAllPermissions(permissions ...Permission) bool {
	for _, permission := range permissions {
		if !p.HasPermission(permission) {
			return false
		}
	}
	return true
}

// AddPermission adds a permission to the user's permissions
func (p *UserPermissions) AddPermission(permission Permission) {
	*p |= UserPermissions(permission)
}

// RemovePermission removes a permission from the user's permissions
func (p *UserPermissions) RemovePermission(permission Permission) {
	*p &= UserPermissions(^permission)
}

// GetPermissions returns all permissions as a slice
func (p UserPermissions) GetPermissions() []Permission {
	var permissions []Permission
	allPermissions := []Permission{
		ReadUser, WriteUser, DeleteUser,
		ReadAdmin, WriteAdmin, DeleteAdmin,
	}

	for _, permission := range allPermissions {
		if p.HasPermission(permission) {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
