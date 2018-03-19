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

	printMessageExtensions(pdep, fdep.FIELD_OPTION.MessageName())

	printOptionType(pdep)
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
		case *fproto.OneOfFieldElement:
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

	for _, fld := range message_element.Fields {
		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fmt.Printf("* FIELD: %s - TYPE: %s\n", xfld.Name, xfld.Type)
			Util_PrintType("", tp_print, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: map<%s, %s>\n", xfld.Name, xfld.KeyType, xfld.Type)
			Util_PrintType("KEY ", tp_print, xfld.KeyType)
			Util_PrintType("", tp_print, xfld.Type)
		case *fproto.OneOfFieldElement:
			fmt.Printf("* FIELD: %s - TYPE: oneof\n", xfld.Name)
			for _, oofld := range xfld.Fields {
				fmt.Printf("\tONEOF FIELD: %s\n", oofld.FieldName())
			}
		}
	}
}

func printMessageExtensions(pdep *fdep.Dep, typeName string) {
	fmt.Printf("%s PRINT MESSAGE EXTENSIONS: %s %s\n", lines, typeName, lines)

	tp_print, err := pdep.GetType(typeName)
	if err != nil {
		log.Fatal(err)
	}

	extensions, err := tp_print.GetTypeExtensions()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message type name: %s\n", tp_print.FullOriginalName())

	for extpkg, extType := range extensions {
		fmt.Printf("* EXTENSION: %s [package: %s] [proto path: %s]\n", extType.FullOriginalName(), extpkg, extType.FileDep.FilePath)
		if em, emok := extType.Item.(*fproto.MessageElement); emok {
			for _, emf := range em.Fields {
				switch xfld := emf.(type) {
				case *fproto.FieldElement:
					fmt.Printf("\tFIELD: %s [type: %s]\n", xfld.Name, xfld.Type)
				case *fproto.MapFieldElement:
					fmt.Printf("\tMAP FIELD: %s [type: map<%s, %s>]\n", xfld.Name, xfld.KeyType, xfld.Type)
				case *fproto.OneOfFieldElement:
					fmt.Printf("\tFIELD: %s [type: oneof]\n", xfld.Name)
				}
			}
		}
	}
}

func printOptionType(pdep *fdep.Dep) {
	fmt.Printf("%s PRINT OPTION TYPE %s\n", lines, lines)

	o, err := pdep.GetOption(fdep.FIELD_OPTION, "validate.field")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OPTION: %s\n", o.OptionName)

	fmt.Printf("Source option type: %s\n", o.SourceOption.FullOriginalName())
	if o.Option != nil {
		fmt.Printf("Option type: %s\n", o.Option.FullOriginalName())
	}
	if o.Name != "" {
		fmt.Printf("Option name: %s\n", o.Name)
	}
	if o.FieldItem != nil {
		parent_dt := pdep.DepTypeFromElement(o.FieldItem.ParentElement())

		switch xfld := o.FieldItem.(type) {
		case *fproto.FieldElement:
			fmt.Printf("Field item fieldname: %s [type: %s]\n", xfld.Name, xfld.Type)
			Util_PrintType("", parent_dt, xfld.Type)
		case *fproto.MapFieldElement:
			fmt.Printf("Field item fieldname: %s [type: map<%s, %s>]\n", xfld.Name, xfld.KeyType, xfld.Type)
			Util_PrintType("KEY ", parent_dt, xfld.KeyType)
			Util_PrintType("", parent_dt, xfld.Type)
		case *fproto.OneOfFieldElement:
			fmt.Printf("Field item fieldname: %s [type: oneof]\n", xfld.Name)
		}
	}
}

func Util_PrintType(desc string, parent *fdep.DepType, typeName string) {
	tp_fld, err := parent.MustGetType(typeName)
	if err != nil {
		log.Fatal(err)
	}
	if tp_fld.IsScalar() {
		fmt.Printf("\t%sTYPE: SCALAR %s\n", desc, tp_fld.ScalarType.ProtoType())
	} else {
		fmt.Printf("\t%sTYPE: %s [%s] - from file: %s\n", desc, tp_fld.FullOriginalName(), tp_fld.Item.ElementTypeName(), tp_fld.FileDep.FilePath)
	}
}
