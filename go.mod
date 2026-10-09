module github.com/BrightMountainMedia/uasurfer

// The minimum is what the code needs, not what it is developed on: a go line
// here is a floor for every module that imports this one, and bs-bidder builds
// on Go 1.24. The toolchain line only applies inside this repository.
go 1.24.0

toolchain go1.25.12
