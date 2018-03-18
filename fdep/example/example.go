package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
)

var (
	lines = strings.Repeat("=", 20)
)

// This directory must be the working directory
func main() {
	// load test files
	pdep := fdep.NewDep()

	// add include path
	err := pdep.AddIncludeDir("proto_test/include")
	if err != nil {
		log.Fatal(err)
	}

	// add application files, with root as "app"
	// we use the "WithRoot" version, because if we added "proto_test" directly,
	// the "include" directory would be also included, but with a wrong path
	// (prepended by "include")
	err = pdep.AddPathWithRoot("app", "proto_test/app", fdep.DepType_Own)
	if err != nil {
		log.Fatal(err)
	}

	printFiles(pdep)

	printPackages(pdep)

	printExtensions(pdep)

	printTypes(pdep)

	printFields(pdep)

	printFieldTypes(pdep, "app.core.SendMailAttach")

	printFieldTypes(pdep, "fproto_wrap_headers.Headers")
}

func printFiles(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT FILES %s\n", lines, lines)
	for filepath, file := range pdep.Files {
		fmt.Printf("File: %s (%s) [package: %s]\n", filepath, file.DepType.String(), file.ProtoFile.PackageName)
	}
}

func printPackages(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT PACKAGES %s\n", lines, lines)
	for pkg, filelist := range pdep.Packages {
		fmt.Printf("Package: %s [files: %s]\n", pkg, strings.Join(filelist, ", "))
	}
}

func printExtensions(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT EXTENSIONS %s\n", lines, lines)
	for ext, pkglist := range pdep.Extensions {
		fmt.Printf("Extension: %s [packages: %s]\n", ext, strings.Join(pkglist, ", "))
	}
}

func printTypes(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT TYPES %s\n", lines, lines)

	//
	// app.core.User
	//
	tp_user, err := pdep.GetType("app.core.User")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type 'app.core.User' is in file '%s', package '%s' [name: %s]\n",
		tp_user.FileDep.FilePath, tp_user.OriginalAlias, tp_user.Name)

	//
	// app.core.SendMail.Body
	//
	tp_sendmail_body, err := pdep.GetType("app.core.SendMail.Body")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Type 'app.core.SendMail.Body' is in file '%s', package '%s' [name: %s, alias: %s]\n",
		tp_sendmail_body.FileDep.FilePath, tp_sendmail_body.OriginalAlias, tp_sendmail_body.Name, tp_sendmail_body.Alias)

	//
	// app.core.SendMail.Body in the context of app.Core.Sendmail
	//
	tp_sendmail, err := pdep.GetType("app.core.SendMail")
	if err != nil {
		log.Fatal(err)
	}

	tp_sendmail_body2, err := tp_sendmail.GetType("Body")
	if err != nil {
		log.Fatal(err)
	}

	// When getting a type in the context of other type, the alias may be blank if the type is on the same file.
	fmt.Printf("Type 'Body' in the context of 'app.core.SendMail' is in file '%s', package '%s' [name: %s, alias: %s]\n",
		tp_sendmail_body2.FileDep.FilePath, tp_sendmail_body2.OriginalAlias, tp_sendmail_body2.Name, tp_sendmail_body2.Alias)
}

func printFields(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT FIELDS %s\n", lines, lines)

	//
	// app.core.Sendmail
	//
	tp_sendmail, err := pdep.GetType("app.core.SendMail")
	if err != nil {
		log.Fatal(err)
	}

	message_element, is_message_element := tp_sendmail.Item.(*fproto.MessageElement)
	if !is_message_element {
		log.Fatal("Should be a *fproto.MessageElement")
	}

	fmt.Printf("MESSAGE ELEMENT NAME: %s\n", message_element.Name)
	for _, fld := range message_element.Fields {
		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fmt.Printf("* FIELD: %s - TYPE: %s\n", xfld.Name, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: map<%s, %s>\n", xfld.Name, xfld.KeyType, xfld.Type)
		case *fproto.OneofFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: oneof\n", xfld.Name)
		}
	}
}

func printFieldTypes(pdep *fdep.Dep, typeName string) {
	fmt.Printf("%s PRINT FIELD TYPES: %s %s\n", lines, typeName, lines)

	tp_print, err := pdep.GetType(typeName)
	if err != nil {
		log.Fatal(err)
	}

	message_element, is_message_element := tp_print.Item.(*fproto.MessageElement)
	if !is_message_element {
		log.Fatal("Should be a *fproto.MessageElement")
	}

	fmt.Printf("Message type name: %s\n", tp_print.FullOriginalName())

	print_type := func(desc string, typeName string) {
		tp_fld, err := tp_print.MustGetType(typeName)
		if err != nil {
			log.Fatal(err)
		}
		if tp_fld.IsScalar() {
			fmt.Printf("\t%sTYPE: SCALAR %s\n", desc, tp_fld.ScalarType.ProtoType())
		} else {
			fmt.Printf("\t%sTYPE: %s [%s] - from file: %s\n", desc, tp_fld.FullOriginalName(), tp_fld.Item.ElementTypeName(), tp_fld.FileDep.FilePath)
		}
	}

	for _, fld := range message_element.Fields {
		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fmt.Printf("* FIELD: %s - TYPE: %s\n", xfld.Name, xfld.Type)
			print_type("", xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: map<%s, %s>\n", xfld.Name, xfld.KeyType, xfld.Type)
			print_type("KEY ", xfld.KeyType)
			print_type("", xfld.Type)
		case *fproto.OneofFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: oneof\n", xfld.Name)
			for _, oofld := range xfld.Fields {
				fmt.Printf("\tONEOF FIELD: %s\n", oofld.FieldName())
			}
		}
	}
}
