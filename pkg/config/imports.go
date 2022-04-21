package config

import (
	"fmt"
	"os"
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

	type result struct {
		astFile AstFile
		err     error
	}

	// Traversing AST and find the package location cause disk I/O block, use goroutine to optimize performance
	results := make(chan result)
	jobs := make(chan AstFile, len(i.PackagePaths))
	quit := make(chan struct{})

	var wg sync.WaitGroup

	worker := func(job chan AstFile, r chan result, quit chan struct{}) {
		for {
			select {
			case p := <-job:
				err := p.Load()
				select {
				case r <- result{p, err}:
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
		job := AstFile{FilePath: pkgPath}
		jobs <- job
	}

	for idx := 0; idx < runtime.NumCPU(); idx++ {
		wg.Add(1)
		go worker(jobs, results, quit)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", r.err)
			continue
		}
		for _, imp := range r.astFile.Imports {
			pkgName := strings.Replace(imp.Path.Value, "\"", "", -1)
			// Package -> file location list
			i.PackageLocationMap[pkgName] = append(i.PackageLocationMap[pkgName], r.astFile.FilePath)
		}
	}
}
