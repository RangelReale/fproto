package fdep

import (
	"strings"
	"testing"
)

func TestDep(t *testing.T) {

	dep := NewDep()
	dep.IgnoreNotFoundDependencies = true
	err := dep.AddReader("myapp/proto/p_user/user.proto", strings.NewReader(testfile_user), DepType_Own)
	if err != nil {
		t.Fatalf("Error parsing test user proto: %v", err)
	}

	err = dep.AddReader("google/protobuf/empty.proto", strings.NewReader(testfile_google_empty), DepType_Imported)
	if err != nil {
		t.Fatalf("Error parsing test user proto: %v", err)
	}

	// p_user.User
	user_type, err := dep.GetType("p_user.User")
	if err != nil {
		t.Fatalf("Error getting type p_user.User: %v", err)
	}

	if user_type.Alias != "p_user" {
		t.Fatalf("p_user.User alias should be 'p_user', but is '%s'", user_type.Alias)
	}

	if user_type.Name != "User" {
		t.Fatalf("p_user.User name should be 'User', but is '%s'", user_type.Name)
	}

	// p_user.User.Address
	user_address_type, err := dep.GetType("p_user.User.Address")
	if err != nil {
		t.Fatalf("Error getting type p_user.User.Address: %v", err)
	}

	if user_address_type.Alias != "p_user" {
		t.Fatalf("p_user.User.Address alias should be 'p_user', but is '%s'", user_address_type.Alias)
	}

	if user_address_type.Name != "User.Address" {
		t.Fatalf("p_user.User.Address name should be 'User.Address', but is '%s'", user_address_type.Name)
	}

	// find the User type using the user.proto file as base
	f_user_type, err := dep.Files["myapp/proto/p_user/user.proto"].GetType("User")
	if err != nil {
		t.Fatalf("Error getting type User from user.proto: %v", err)
	}

	if f_user_type.Alias != "" {
		t.Fatalf("User from user.proto's alias should be blank, but is '%s'", f_user_type.Alias)
	}

	if user_type.Name != "User" {
		t.Fatalf("User User from user.proto's name should be 'User', but is '%s'", f_user_type.Name)
	}

	// google.protobuf.Empty
	empty_type, err := dep.GetType("google.protobuf.Empty")
	if err != nil {
		t.Fatalf("Error getting type google.protobuf.Empty: %v", err)
	}

	if empty_type.Alias != "google.protobuf" {
		t.Fatalf("google.protobuf.Empty alias should be 'google.protobuf', but is '%s'", empty_type.Alias)
	}

	if empty_type.Name != "Empty" {
		t.Fatalf("google.protobuf.Empty name should be 'Empty', but is '%s'", empty_type.Name)
	}
}
