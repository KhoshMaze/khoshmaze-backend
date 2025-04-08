package model

/*

	How permission system works:

	Each user has a permission attribute which is a number. This number is sum of all of user permissions.
	Since the sum of any K power of 2 numbers is unique, system finds each user permissions (The binary
	version 1 bits represent the permissions of user)


	This system can be mixed with Roles for more specific and stronger permission system.
*/

type Authority int64

// HasSpecific checks if the user has a specific permission
func (p Authority) HasSpecific(permission Authority) bool {
	return uint64(p)&uint64(permission) == uint64(permission)
}

func (p Authority) HasAll(permissions ...Authority) bool {
	for _, permission := range permissions {
		if !p.HasSpecific(permission) {
			return false
		}
	}
	return true
}

func (p *Authority) Add(permission Authority) {
	*p |= Authority(permission)
}

func (p *Authority) Remove(permission Authority) {
	*p &= Authority(^permission)
}

// Deprecated
func (p Authority) GetAll() []Authority {
	var permissions []Authority
	allPermissions := []Authority{}

	for _, permission := range allPermissions {
		if p.HasSpecific(permission) {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
