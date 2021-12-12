package integrations

// func run() (err error) {
// 	// 	ctx := context.Background()

// 	// 	// if err = cache.Clear(); err != nil {
// 	// 	// 	return fmt.Errorf("error in cache clear: %w", err)
// 	// 	// }

// 	// 	applier.NewApplier(cache.New)
// 	// 	analyzer := analyzer.NewAnalyzer([]analyzer.Type{})
// 	// 	analyzer.AnalyzeFile(ctx, )
// 	// 	// rc, err := openStream(*tarPath)
// 	// 	// if err != nil {
// 	// 	// 	return err
// 	// 	// }

// 	// 	// files, err = analyzer.AnalyzeFromFile(ctx, rc)
// 	// 	// if err != nil {
// 	// 	// 	return err
// 	// 	// }

// 	// 	os, err := analyzer.GetOS(files)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	fmt.Printf("%+v\n", os)

// 	// 	pkgs, err := analyzer.GetPackages(files)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	fmt.Printf("Packages: %d\n", len(pkgs))

// 	// 	libs, err := analyzer.GetLibraries(files)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	for filepath, libList := range libs {
// 	// 		fmt.Printf("%s: %d\n", filepath, len(libList))
// 	// 	}
// 	// 	return nil
// 	// }

// 	// func openStream(path string) (*os.File, error) {
// 	// 	if path == "-" {
// 	// 		if terminal.IsTerminal(0) {
// 	// 			flag.Usage()
// 	// 			os.Exit(64)
// 	// 		} else {
// 	// 			return os.Stdin, nil
// 	// 		}
// 	// 	}
// 	// 	return os.Open(path)
// }
