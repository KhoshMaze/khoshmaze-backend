package model

type UserRoles = Authority

const (
	BranchAccountant UserRoles = 1 << 3
	Waiter           UserRoles = 1 << 5
	Cashier          UserRoles = 1 << 8
	BranchManager    UserRoles = 1 << 10
	RestaurantOwner  UserRoles = 1 << 12

	Support    UserRoles = 1 << 27
	Accountant UserRoles = 1 << 28
	SuperAdmin UserRoles = 1 << 30
	Founder    UserRoles = 1 << 62
)

// ------------------------------------------------------------

type UserPermissions = Authority

const (
	Read               UserPermissions = 1 << 0
	Update             UserPermissions = 1 << 1
	Create             UserPermissions = 1 << 2
	Delete             UserPermissions = 1 << 3
	BanUser            UserPermissions = 1 << 4
	ChangeBranchStatus UserPermissions = 1 << 5
	SeeFinancialReport UserPermissions = 1 << 6
	UpdateMenuStock    UserPermissions = 1 << 7
	CanGiveDiscount    UserPermissions = 1 << 8
)
