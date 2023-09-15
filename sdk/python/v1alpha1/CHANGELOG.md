# CHANGELOG

This is a manually generated log to track changes to the repository for each release.
Each section should include general headers such as **Implemented enhancements**
and **Merged pull requests**. Critical items to know are:

 - renamed commands
 - deprecated / removed commands
 - changed defaults
 - backward incompatible changes (recipe file format? image file format?)
 - migration guidance (how to convert images?)
 - changed behaviour (recipe sections work differently)

The versions coincide with releases on pip. Only major versions will be released as tags on Github.

## [0.0.x](https://github.com/converged-computing/metrics-operator/tree/main) (0.0.x)
 - Allow getting raw logs for any metric (without parser) (0.0.19)
 - Refactor of structure of Operator and addition of metrics (0.0.18)
 - Add wait for delete function to python parser (0.0.17)
 - LAMMPS python parser (0.0.16)
   - custom flags and multi example for osu-benchmarks (0.0.16.2)
   - fixes to metric osu-benchmarks (0.0.16.1)
 - add FIO storage (IO) metric (0.0.15)
 - resources specification added and tweaks to perf-sysstat (0.0.14)
 - pidstat python parser and better support for metric in Go (0.0.13)
 - Separation of parsing logs into separate metric module functions (0.0.12.1)
 - Support for Netmark parser and plotting example (0.0.12)
 - Support for OSU Benchmarks parser (and plotting) (0.0.11)
 - First release with support for parsing io-sysstat output (0.0.1)
 - Skeleton release (0.0.0)
