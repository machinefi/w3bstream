package client_test

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/client"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

var t = &client.TypeInfo{
	Expose: "TypeName",
	Import: "path/to/repo/pkg",
}

func ExampleTypeInfo_AsAlias() {
	ti := t.AsAlias()
	output(ti.Snippet(f)...)

	ti = t.SetAliasLevel(1).AsAlias()
	output(ti.Snippet(f)...)

	ti = t.SetAliasLevel(2).AsAlias()
	output(ti.Snippet(f)...)

	ti = t.SetAliasLevel(3).AsAlias()
	output(ti.Snippet(f)...)

	ti = t.SetAliasLevel(4).AsAlias()
	output(ti.Snippet(f)...)

	ti = t.SetAliasLevel(5).AsAlias()
	output(ti.Snippet(f)...)

	// Output:
	// type TypeName = pkg.TypeName
	// type PkgTypeName = pkg.TypeName
	// type RepoPkgTypeName = pkg.TypeName
	// type ToRepoPkgTypeName = pkg.TypeName
	// type PathToRepoPkgTypeName = pkg.TypeName
	// type PathToRepoPkgTypeName = pkg.TypeName
}

func ExampleTypeInfo_AsEnum() {
	strOpts := enumgen.Options{
		{Label: "OPTION1", Str: ptrx.Ptr("OPTION1")},
		{Label: "OPTION2", Str: ptrx.Ptr("OPTION2")},
	}
	intOpts := enumgen.Options{
		{Label: "OPTION1", Int: ptrx.Ptr(int64(1))},
		{Label: "OPTION2", Int: ptrx.Ptr(int64(2))},
	}
	floatOpts := enumgen.Options{
		{Label: "OPTION1", Float: ptrx.Float64(40000001.9999)},
		{Label: "OPTION2", Float: ptrx.Float64(40000002.9999)},
	}

	ti := t.AsEnum(strOpts)
	output(ti.Snippet(f)...)

	ti = ti.AsAlias().SetAliasLevel(1)
	output(ti.Snippet(f)...)

	ti = t.AsEnum(intOpts)
	output(ti.Snippet(f)...)

	ti = ti.AsAlias().SetAliasLevel(2)
	output(ti.Snippet(f)...)

	ti = t.AsEnum(floatOpts)
	output(ti.Snippet(f)...)

	ti = ti.AsAlias().SetAliasLevel(3)
	output(ti.Snippet(f)...)

	// Output:
	// type TypeName string
	// const (
	// TYPE_NAME__OPTION1 TypeName = "OPTION1" // OPTION1
	// TYPE_NAME__OPTION2 TypeName = "OPTION2" // OPTION2
	// )
	//
	// type PkgTypeName = pkg.TypeName
	// const (
	// PKG_TYPE_NAME__OPTION1 PkgTypeName = "OPTION1" // OPTION1
	// PKG_TYPE_NAME__OPTION2 PkgTypeName = "OPTION2" // OPTION2
	// )
	//
	// type TypeName int64
	// const (
	// TYPE_NAME__OPTION1 TypeName = 1 // OPTION1
	// TYPE_NAME__OPTION2 TypeName = 2 // OPTION2
	// )
	//
	// type RepoPkgTypeName = pkg.TypeName
	// const (
	// REPO_PKG_TYPE_NAME__OPTION1 RepoPkgTypeName = 1 // OPTION1
	// REPO_PKG_TYPE_NAME__OPTION2 RepoPkgTypeName = 2 // OPTION2
	// )
	//
	// type TypeName float64
	// const (
	// TYPE_NAME__OPTION1 TypeName = 40000001.9999 // OPTION1
	// TYPE_NAME__OPTION2 TypeName = 40000002.9999 // OPTION2
	// )
	//
	// type ToRepoPkgTypeName = pkg.TypeName
	// const (
	// TO_REPO_PKG_TYPE_NAME__OPTION1 ToRepoPkgTypeName = 40000001.9999 // OPTION1
	// TO_REPO_PKG_TYPE_NAME__OPTION2 ToRepoPkgTypeName = 40000002.9999 // OPTION2
	// )
}
