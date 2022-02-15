// Code generated by entc, DO NOT EDIT.

package ent

import (
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent/permission"
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
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.DefaultName holds the default value on creation for the name field.
	user.DefaultName = userDescName.Default.(string)
}
