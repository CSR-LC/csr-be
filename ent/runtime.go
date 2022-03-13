// Code generated by entc, DO NOT EDIT.

package ent

import (
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/permission"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/role"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/schema"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	permissionFields := schema.Permission{}.Fields()
	_ = permissionFields
	// permissionDescName is the schema descriptor for name field.
	permissionDescName := permissionFields[0].Descriptor()
	// permission.DefaultName holds the default value on creation for the name field.
	permission.DefaultName = permissionDescName.Default.(string)
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescName is the schema descriptor for name field.
	roleDescName := roleFields[0].Descriptor()
	// role.NameValidator is a validator for the "name" field. It is called by the builders before save.
	role.NameValidator = roleDescName.Validators[0].(func(string) error)
	// roleDescSlug is the schema descriptor for slug field.
	roleDescSlug := roleFields[1].Descriptor()
	// role.SlugValidator is a validator for the "slug" field. It is called by the builders before save.
	role.SlugValidator = roleDescSlug.Validators[0].(func(string) error)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.DefaultName holds the default value on creation for the name field.
	user.DefaultName = userDescName.Default.(string)
}
