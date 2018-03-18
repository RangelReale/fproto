package main

import (
	"fmt"
	"log"
	"strings"

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
