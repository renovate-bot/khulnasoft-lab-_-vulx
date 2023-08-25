// Code generated by mockery v1.0.0. DO NOT EDIT.

package cache

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/khulnasoft-lab/vul/pkg/fanal/types"
)

// MockLocalArtifactCache is an autogenerated mock type for the LocalArtifactCache type
type MockLocalArtifactCache struct {
	mock.Mock
}

type LocalArtifactCacheClearReturns struct {
	Err error
}

type LocalArtifactCacheClearExpectation struct {
	Returns LocalArtifactCacheClearReturns
}

func (_m *MockLocalArtifactCache) ApplyClearExpectation(e LocalArtifactCacheClearExpectation) {
	var args []interface{}
	_m.On("Clear", args...).Return(e.Returns.Err)
}

func (_m *MockLocalArtifactCache) ApplyClearExpectations(expectations []LocalArtifactCacheClearExpectation) {
	for _, e := range expectations {
		_m.ApplyClearExpectation(e)
	}
}

// Clear provides a mock function with given fields:
func (_m *MockLocalArtifactCache) Clear() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type LocalArtifactCacheCloseReturns struct {
	Err error
}

type LocalArtifactCacheCloseExpectation struct {
	Returns LocalArtifactCacheCloseReturns
}

func (_m *MockLocalArtifactCache) ApplyCloseExpectation(e LocalArtifactCacheCloseExpectation) {
	var args []interface{}
	_m.On("Close", args...).Return(e.Returns.Err)
}

func (_m *MockLocalArtifactCache) ApplyCloseExpectations(expectations []LocalArtifactCacheCloseExpectation) {
	for _, e := range expectations {
		_m.ApplyCloseExpectation(e)
	}
}

// Close provides a mock function with given fields:
func (_m *MockLocalArtifactCache) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type LocalArtifactCacheGetArtifactArgs struct {
	ArtifactID         string
	ArtifactIDAnything bool
}

type LocalArtifactCacheGetArtifactReturns struct {
	ArtifactInfo types.ArtifactInfo
	Err          error
}

type LocalArtifactCacheGetArtifactExpectation struct {
	Args    LocalArtifactCacheGetArtifactArgs
	Returns LocalArtifactCacheGetArtifactReturns
}

func (_m *MockLocalArtifactCache) ApplyGetArtifactExpectation(e LocalArtifactCacheGetArtifactExpectation) {
	var args []interface{}
	if e.Args.ArtifactIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.ArtifactID)
	}
	_m.On("GetArtifact", args...).Return(e.Returns.ArtifactInfo, e.Returns.Err)
}

func (_m *MockLocalArtifactCache) ApplyGetArtifactExpectations(expectations []LocalArtifactCacheGetArtifactExpectation) {
	for _, e := range expectations {
		_m.ApplyGetArtifactExpectation(e)
	}
}

// GetArtifact provides a mock function with given fields: artifactID
func (_m *MockLocalArtifactCache) GetArtifact(artifactID string) (types.ArtifactInfo, error) {
	ret := _m.Called(artifactID)

	var r0 types.ArtifactInfo
	if rf, ok := ret.Get(0).(func(string) types.ArtifactInfo); ok {
		r0 = rf(artifactID)
	} else {
		r0 = ret.Get(0).(types.ArtifactInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(artifactID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type LocalArtifactCacheGetBlobArgs struct {
	BlobID         string
	BlobIDAnything bool
}

type LocalArtifactCacheGetBlobReturns struct {
	BlobInfo types.BlobInfo
	Err      error
}

type LocalArtifactCacheGetBlobExpectation struct {
	Args    LocalArtifactCacheGetBlobArgs
	Returns LocalArtifactCacheGetBlobReturns
}

func (_m *MockLocalArtifactCache) ApplyGetBlobExpectation(e LocalArtifactCacheGetBlobExpectation) {
	var args []interface{}
	if e.Args.BlobIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.BlobID)
	}
	_m.On("GetBlob", args...).Return(e.Returns.BlobInfo, e.Returns.Err)
}

func (_m *MockLocalArtifactCache) ApplyGetBlobExpectations(expectations []LocalArtifactCacheGetBlobExpectation) {
	for _, e := range expectations {
		_m.ApplyGetBlobExpectation(e)
	}
}

// GetBlob provides a mock function with given fields: blobID
func (_m *MockLocalArtifactCache) GetBlob(blobID string) (types.BlobInfo, error) {
	ret := _m.Called(blobID)

	var r0 types.BlobInfo
	if rf, ok := ret.Get(0).(func(string) types.BlobInfo); ok {
		r0 = rf(blobID)
	} else {
		r0 = ret.Get(0).(types.BlobInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(blobID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
