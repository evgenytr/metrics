// Package main implements custom multichecker staticlint.
//
// staticlint is a tool for static analysis of Go programs.
//
// staticlint examines Go source code and reports suspicious constructs,
// such as Printf calls whose arguments do not align with the format
// string. It uses heuristics that do not guarantee all reports are
// genuine problems, but it can find errors not caught by the compilers.
//
// By default, all analyzers are run.
// To select specific analyzers, use the -NAME flag for each one,
// or -NAME=false to run all analyzers not explicitly disabled.
//
// For detailed info please run
//
//	staticlint help
//
// To see details and flags of a specific analyzer, run 'staticlint help name'.
//
// Included are:
//   - all analyzers from staticcheck.io:
//     SA1, SA2, SA3, SA4, SA5, SA6, SA9, S1, ST1, QF1
//   - durationcheck: linter to detect cases where two time.Duration values are being multiplied in possibly erroneous ways
//     github.com/charithe/durationcheck
//   - wastedassign: finds wasted assignment statements
//     reassigned, but never used afterward
//     reassigned, but reassigned without using the value
//     github.com/sanposhiho/wastedassign
//   - osexitmain checker to detect os.Exit direct calls in main function of package main
//   - all analyzers from x/analysis/passes:
//     asmdecl
//     assign
//     atomic
//     atomicalign
//     bools
//     buildssa
//     buildtag
//     cgocall
//     composite
//     copylock
//     ctrlflow
//     deepequalerrors
//     errorsas
//     fieldalignment
//     findcall
//     framepointer
//     httpresponse
//     ifaceassert
//     inspect
//     loopclosure
//     lostcancel
//     nilfunc
//     nilness
//     pkgfact
//     printf
//     reflectvaluecompare
//     shadow
//     shift
//     sigchanyzer
//     sortslice
//     stdmethods
//     stringintconv
//     structtag
//     tests
//     unmarshal
//     unreachable
//     unsafeptr
//     unusedresult
//     unusedwrite
package main
