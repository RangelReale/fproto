package fproto

import (
	"strings"
	"testing"
)

var (
	testfile = `
syntax = "proto3";
package p_user;
option go_package = "myapp/proto/p_user";

message User {
	message Address {
		string address = 1;
		string city = 2;
	}

    int32 id = 1;
    string name = 2;
    string email = 3;

	Address address = 4;
}

	`
)

func TestFindProtoFileName(t *testing.T) {
	r := strings.NewReader(testfile)

	pfile, err := Parse(r)
	if err != nil {
		t.Fatalf("Error parsing proto file: %c", err)
	}

	if pfile.PackageName != "p_user" {
		t.Fatalf("Package name should be 'p_user', but '%s' found", pfile.PackageName)
	}

	user_address := pfile.FindName("User.Address")
	if len(user_address) != 1 {
		t.Fatalf("Error finding User.Address, expected 1 item got %d", len(user_address))
	}

	user_address_item, ok := user_address[0].(*MessageElement)
	if !ok {
		t.Fatalf("User.Address should be a MessageElement")
	}

	if user_address_item.Name != "Address" {
		t.Fatalf("User.Address' name should be a 'Address' but is %s", user_address_item.Name)
	}

	user_item, ok := user_address_item.Parent.(*MessageElement)
	if !ok {
		t.Fatalf("User.Address' parent should be a MessageElement")
	}

	if user_item.Name != "User" {
		t.Fatalf("User.Address' parent name should be a 'User' but is %s", user_item.Name)
	}
}
