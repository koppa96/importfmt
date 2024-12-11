package main

import (
	"go/token"
	"io"
	"slices"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func format(f *dst.File, writer io.Writer) error {
	var (
		stdSpecs []*dst.ImportSpec
		extSpecs []*dst.ImportSpec
	)
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*dst.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok != token.IMPORT {
			continue
		}

		for _, spec := range genDecl.Specs {
			importSpec := spec.(*dst.ImportSpec)
			importSpec.Decs.Before = dst.None
			importSpec.Decs.After = dst.NewLine

			if strings.ContainsRune(importSpec.Path.Value, '.') {
				extSpecs = append(extSpecs, importSpec)
			} else {
				stdSpecs = append(stdSpecs, importSpec)
			}
		}
	}

	if len(stdSpecs) == 0 && len(extSpecs) == 0 {
		return decorator.Fprint(writer, f)
	}

	slices.SortFunc(stdSpecs, func(a, b *dst.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	slices.SortFunc(extSpecs, func(a, b *dst.ImportSpec) int {
		return strings.Compare(a.Path.Value, b.Path.Value)
	})

	if len(extSpecs) > 0 && len(stdSpecs) > 0 {
		stdSpecs[len(stdSpecs)-1].Decs.After = dst.EmptyLine
	}

	specs := make([]dst.Spec, 0, len(stdSpecs)+len(extSpecs))
	for _, spec := range stdSpecs {
		specs = append(specs, spec)
	}

	for _, spec := range extSpecs {
		specs = append(specs, spec)
	}

	f.Decls = preserveFirstImportDeclWithSpecs(f.Decls, specs)

	return decorator.Fprint(writer, f)
}

func preserveFirstImportDeclWithSpecs(decls []dst.Decl, specs []dst.Spec) []dst.Decl {
	result := make([]dst.Decl, 0, len(decls))
	firstFound := false

	for _, decl := range decls {
		genDecl, ok := decl.(*dst.GenDecl)
		if !ok {
			result = append(result, decl)
			continue
		}

		if genDecl.Tok == token.IMPORT {
			if firstFound {
				continue
			}

			firstFound = true
			genDecl.Lparen = len(specs) == 1
			genDecl.Rparen = len(specs) == 1
			genDecl.Specs = specs
		}

		result = append(result, decl)
	}

	return result
}
