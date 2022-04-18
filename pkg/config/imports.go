package config

import (
	"golang.org/x/tools/go/packages"
	"runtime"
	"strings"
	"sync"
)

type Imports struct {
	// The path to the go package which contains .go file
	PackagePaths []string
	// package -> .go file paths
	PackageLocationMap map[string][]string
}

func (i Imports) GetImports() {
	config := &packages.Config{
		Mode: packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedImports,
	}

	type result struct {
		pkgPath string
		pkgs    []*packages.Package
		err     error
	}

	// Traversing AST and find the package location cause disk I/O block, use goroutine to optimize performance
	results := make(chan result)
	jobs := make(chan string, len(i.PackagePaths))
	quit := make(chan struct{})

	var wg sync.WaitGroup

	worker := func(job chan string, r chan result, quit chan struct{}) {
		for {
			select {
			case p := <-job:
				// Can't use golang.org/x/tools/go/packages because it doesn't support non go modules
				pkgs, err := packages.Load(config, p)
				select {
				case r <- result{p, pkgs, err}:
				case <-quit:
					wg.Done()
					return
				}
			default:
				wg.Done()
				return
			}
		}
	}

	for _, pkgPath := range i.PackagePaths {
		jobs <- pkgPath
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(jobs, results, quit)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			continue
		}
		for _, pkg := range r.pkgs {
			if pkg.Errors != nil {
				//fmt.Fprintf(os.Stderr, "error: %v\n", pkg.Errors)
				close(quit)
				wg.Wait()
			}
			for _, syntax := range pkg.Syntax {
				for _, imp := range syntax.Imports {
					if imp.Path.Value == "" {
						continue
					}
					if _, ok := i.PackageLocationMap[imp.Path.Value]; !ok {
						i.PackageLocationMap[imp.Path.Value] = []string{}
					}
					pkgName := strings.Replace(imp.Path.Value, "\"", "", -1)
					// Package -> file location list
					i.PackageLocationMap[pkgName] = append(i.PackageLocationMap[pkgName], pkg.Fset.File(syntax.Pos()).Name())
				}
			}
		}
	}
}
