package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func containsKey(oi *ast.ObjectItem, name string) bool {
	for _, k := range oi.Keys {
		if k.Token.Text == name {
			return true
		}
	}
	return false
}

func findModuleSource(module *ast.ObjectItem) (string, string) {
	modName := *&module.Keys[1].Token.Text

	ot, ok := module.Val.(*ast.ObjectType)
	if !ok {
		log.Fatalf("Expected ast.ObjectType")
	}

	for _, ent := range ot.List.Items {
		if !containsKey(ent, "source") {
			continue
		}

		lit, ok := ent.Val.(*ast.LiteralType)
		if !ok {
			log.Fatalf("Expected source Val to be *ast.LiteralType")
		}

		return modName, lit.Token.Text
	}

	log.Fatalf("No source field found on module resource")
	return "", ""
}

func findModuleReferencs(file *ast.File) map[string]string {
	refs := make(map[string]string)

	ol, ok := file.Node.(*ast.ObjectList)
	if !ok {
		log.Fatalln("no objectlist")
	}

	for _, oi := range ol.Items {
		for _, key := range oi.Keys {
			if key.Token.Text == "module" {
				modName, modSource := findModuleSource(oi)
				refs[strings.Trim(modName, "\"")] = strings.Trim(modSource, "\"")
			}
		}
	}

	return refs
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
}

func main() {
	files := make([]string, 0)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".tf" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	checkError(err)

	fmt.Printf("Files %v\n", files)

	for _, filename := range files {
		content, err := ioutil.ReadFile(filename)
		checkError(err)

		fileast, err := hcl.ParseBytes(content)
		checkError(err)

		refs := findModuleReferencs(fileast)
		fmt.Printf("%s -> %s\n", filename, refs)

		abs, err := filepath.Abs(filename)
		checkError(err)
		fmt.Printf("%s\n", abs)

		repodir := filepath.Dir(abs)
		fmt.Printf("%s\n", repodir)

		for k, v := range refs {
			fmt.Printf("%s %s\n", k, filepath.Join(abs, v))
		}
	}
}
