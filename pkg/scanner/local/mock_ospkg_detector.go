// Code generated by mockery v1.0.0. DO NOT EDIT.

package local

import mock "github.com/stretchr/testify/mock"
import pkgtypes "github.com/khulnasoft-lab/vul/pkg/types"
import time "time"
import types "github.com/aquasecurity/fanal/types"

// MockOspkgDetector is an autogenerated mock type for the OspkgDetector type
type MockOspkgDetector struct {
	mock.Mock
}

type OspkgDetectorDetectArgs struct {
	ImageName         string
	ImageNameAnything bool
	OsFamily          string
	OsFamilyAnything  bool
	OsName            string
	OsNameAnything    bool
	Created           time.Time
	CreatedAnything   bool
	Pkgs              []types.Package
	PkgsAnything      bool
}

type OspkgDetectorDetectReturns struct {
	DetectedVulns []pkgtypes.DetectedVulnerability
	Eosl          bool
	Err           error
}

type OspkgDetectorDetectExpectation struct {
	Args    OspkgDetectorDetectArgs
	Returns OspkgDetectorDetectReturns
}

func (_m *MockOspkgDetector) ApplyDetectExpectation(e OspkgDetectorDetectExpectation) {
	var args []interface{}
	if e.Args.ImageNameAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ImageName)
	}
	if e.Args.OsFamilyAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.OsFamily)
	}
	if e.Args.OsNameAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.OsName)
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
	_m.On("Detect", args...).Return(e.Returns.DetectedVulns, e.Returns.Eosl, e.Returns.Err)
}

func (_m *MockOspkgDetector) ApplyDetectExpectations(expectations []OspkgDetectorDetectExpectation) {
	for _, e := range expectations {
		_m.ApplyDetectExpectation(e)
	}
}

// Detect provides a mock function with given fields: imageName, osFamily, osName, created, pkgs
func (_m *MockOspkgDetector) Detect(imageName string, osFamily string, osName string, created time.Time, pkgs []types.Package) ([]pkgtypes.DetectedVulnerability, bool, error) {
	ret := _m.Called(imageName, osFamily, osName, created, pkgs)

	var r0 []pkgtypes.DetectedVulnerability
	if rf, ok := ret.Get(0).(func(string, string, string, time.Time, []types.Package) []pkgtypes.DetectedVulnerability); ok {
		r0 = rf(imageName, osFamily, osName, created, pkgs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pkgtypes.DetectedVulnerability)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string, string, string, time.Time, []types.Package) bool); ok {
		r1 = rf(imageName, osFamily, osName, created, pkgs)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, string, time.Time, []types.Package) error); ok {
		r2 = rf(imageName, osFamily, osName, created, pkgs)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
