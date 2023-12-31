// Code generated by mockery v1.0.0. DO NOT EDIT.

package library

import mock "github.com/stretchr/testify/mock"
import pkgtypes "github.com/khulnasoft-lab/vul/pkg/types"
import time "time"
import types "github.com/aquasecurity/fanal/types"

// MockOperation is an autogenerated mock type for the Operation type
type MockOperation struct {
	mock.Mock
}

type OperationDetectArgs struct {
	ImageName         string
	ImageNameAnything bool
	FilePath          string
	FilePathAnything  bool
	Created           time.Time
	CreatedAnything   bool
	Pkgs              []types.Package
	PkgsAnything      bool
}

type OperationDetectReturns struct {
	Vulns []pkgtypes.DetectedVulnerability
	Err   error
}

type OperationDetectExpectation struct {
	Args    OperationDetectArgs
	Returns OperationDetectReturns
}

func (_m *MockOperation) ApplyDetectExpectation(e OperationDetectExpectation) {
	var args []interface{}
	if e.Args.ImageNameAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ImageName)
	}
	if e.Args.FilePathAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.FilePath)
	}
	if e.Args.CreatedAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Created)
	}
	if e.Args.PkgsAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Pkgs)
	}
	_m.On("Detect", args...).Return(e.Returns.Vulns, e.Returns.Err)
}

func (_m *MockOperation) ApplyDetectExpectations(expectations []OperationDetectExpectation) {
	for _, e := range expectations {
		_m.ApplyDetectExpectation(e)
	}
}

// Detect provides a mock function with given fields: imageName, filePath, created, pkgs
func (_m *MockOperation) Detect(imageName string, filePath string, created time.Time, pkgs []types.Package) ([]pkgtypes.DetectedVulnerability, error) {
	ret := _m.Called(imageName, filePath, created, pkgs)

	var r0 []pkgtypes.DetectedVulnerability
	if rf, ok := ret.Get(0).(func(string, string, time.Time, []types.Package) []pkgtypes.DetectedVulnerability); ok {
		r0 = rf(imageName, filePath, created, pkgs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pkgtypes.DetectedVulnerability)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, time.Time, []types.Package) error); ok {
		r1 = rf(imageName, filePath, created, pkgs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
