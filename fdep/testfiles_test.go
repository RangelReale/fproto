package fdep

var (
	testfile_user = `
syntax = "proto3";
package p_user;
option go_package = "myapp/proto/p_user";

import "google/protobuf/empty.proto";

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

message UserListResponse {
    repeated User list = 1;
}

service UserSvc {
    rpc List(google.protobuf.Empty) returns (UserListResponse);
    rpc Add(User) returns (google.protobuf.Empty);
}
	`

	testfile_google_empty = `
syntax = "proto3";

package google.protobuf;

option csharp_namespace = "Google.Protobuf.WellKnownTypes";
option go_package = "github.com/golang/protobuf/ptypes/empty";
option java_package = "com.google.protobuf";
option java_outer_classname = "EmptyProto";
option java_multiple_files = true;
option objc_class_prefix = "GPB";
option cc_enable_arenas = true;

// A generic empty message that you can re-use to avoid defining duplicated
// empty messages in your APIs. A typical example is to use it as the request
// or the response type of an API method. For instance:
//
//     service Foo {
//       rpc Bar(google.protobuf.Empty) returns (google.protobuf.Empty);
//     }
//
// The JSON representation for "Empty"" is empty JSON object ""{}"".
message Empty {}
`
)
