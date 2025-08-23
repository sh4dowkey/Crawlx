package main

func init() {
	// The `flag.BoolVar` call should be in a separate `init` function or `main` function.
	// Since you've split the files, you need to make sure the flag is defined before `flag.Parse()`.
	// For this example, let's assume `main.go` is where the flags are defined.
	// We'll keep the `verbose` variable as a global and access it directly.
}
