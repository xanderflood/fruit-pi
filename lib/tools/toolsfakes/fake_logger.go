// Code generated by counterfeiter. DO NOT EDIT.
package toolsfakes

import (
	"sync"

	"github.com/xanderflood/fruit-pi/lib/tools"
)

type FakeLogger struct {
	ErrorStub        func(...interface{})
	errorMutex       sync.RWMutex
	errorArgsForCall []struct {
		arg1 []interface{}
	}
	ErrorfStub        func(string, ...interface{})
	errorfMutex       sync.RWMutex
	errorfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	InfoStub        func(...interface{})
	infoMutex       sync.RWMutex
	infoArgsForCall []struct {
		arg1 []interface{}
	}
	InfofStub        func(string, ...interface{})
	infofMutex       sync.RWMutex
	infofArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	DetailStub        func(...interface{})
	detailMutex       sync.RWMutex
	detailArgsForCall []struct {
		arg1 []interface{}
	}
	DetailfStub        func(string, ...interface{})
	detailfMutex       sync.RWMutex
	detailfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	DebugStub        func(...interface{})
	debugMutex       sync.RWMutex
	debugArgsForCall []struct {
		arg1 []interface{}
	}
	DebugfStub        func(string, ...interface{})
	debugfMutex       sync.RWMutex
	debugfArgsForCall []struct {
		arg1 string
		arg2 []interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLogger) Error(arg1 ...interface{}) {
	fake.errorMutex.Lock()
	fake.errorArgsForCall = append(fake.errorArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Error", []interface{}{arg1})
	fake.errorMutex.Unlock()
	if fake.ErrorStub != nil {
		fake.ErrorStub(arg1...)
	}
}

func (fake *FakeLogger) ErrorCallCount() int {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return len(fake.errorArgsForCall)
}

func (fake *FakeLogger) ErrorArgsForCall(i int) []interface{} {
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	return fake.errorArgsForCall[i].arg1
}

func (fake *FakeLogger) Errorf(arg1 string, arg2 ...interface{}) {
	fake.errorfMutex.Lock()
	fake.errorfArgsForCall = append(fake.errorfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Errorf", []interface{}{arg1, arg2})
	fake.errorfMutex.Unlock()
	if fake.ErrorfStub != nil {
		fake.ErrorfStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) ErrorfCallCount() int {
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	return len(fake.errorfArgsForCall)
}

func (fake *FakeLogger) ErrorfArgsForCall(i int) (string, []interface{}) {
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	return fake.errorfArgsForCall[i].arg1, fake.errorfArgsForCall[i].arg2
}

func (fake *FakeLogger) Info(arg1 ...interface{}) {
	fake.infoMutex.Lock()
	fake.infoArgsForCall = append(fake.infoArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Info", []interface{}{arg1})
	fake.infoMutex.Unlock()
	if fake.InfoStub != nil {
		fake.InfoStub(arg1...)
	}
}

func (fake *FakeLogger) InfoCallCount() int {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return len(fake.infoArgsForCall)
}

func (fake *FakeLogger) InfoArgsForCall(i int) []interface{} {
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	return fake.infoArgsForCall[i].arg1
}

func (fake *FakeLogger) Infof(arg1 string, arg2 ...interface{}) {
	fake.infofMutex.Lock()
	fake.infofArgsForCall = append(fake.infofArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Infof", []interface{}{arg1, arg2})
	fake.infofMutex.Unlock()
	if fake.InfofStub != nil {
		fake.InfofStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) InfofCallCount() int {
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	return len(fake.infofArgsForCall)
}

func (fake *FakeLogger) InfofArgsForCall(i int) (string, []interface{}) {
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	return fake.infofArgsForCall[i].arg1, fake.infofArgsForCall[i].arg2
}

func (fake *FakeLogger) Detail(arg1 ...interface{}) {
	fake.detailMutex.Lock()
	fake.detailArgsForCall = append(fake.detailArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Detail", []interface{}{arg1})
	fake.detailMutex.Unlock()
	if fake.DetailStub != nil {
		fake.DetailStub(arg1...)
	}
}

func (fake *FakeLogger) DetailCallCount() int {
	fake.detailMutex.RLock()
	defer fake.detailMutex.RUnlock()
	return len(fake.detailArgsForCall)
}

func (fake *FakeLogger) DetailArgsForCall(i int) []interface{} {
	fake.detailMutex.RLock()
	defer fake.detailMutex.RUnlock()
	return fake.detailArgsForCall[i].arg1
}

func (fake *FakeLogger) Detailf(arg1 string, arg2 ...interface{}) {
	fake.detailfMutex.Lock()
	fake.detailfArgsForCall = append(fake.detailfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Detailf", []interface{}{arg1, arg2})
	fake.detailfMutex.Unlock()
	if fake.DetailfStub != nil {
		fake.DetailfStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) DetailfCallCount() int {
	fake.detailfMutex.RLock()
	defer fake.detailfMutex.RUnlock()
	return len(fake.detailfArgsForCall)
}

func (fake *FakeLogger) DetailfArgsForCall(i int) (string, []interface{}) {
	fake.detailfMutex.RLock()
	defer fake.detailfMutex.RUnlock()
	return fake.detailfArgsForCall[i].arg1, fake.detailfArgsForCall[i].arg2
}

func (fake *FakeLogger) Debug(arg1 ...interface{}) {
	fake.debugMutex.Lock()
	fake.debugArgsForCall = append(fake.debugArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Debug", []interface{}{arg1})
	fake.debugMutex.Unlock()
	if fake.DebugStub != nil {
		fake.DebugStub(arg1...)
	}
}

func (fake *FakeLogger) DebugCallCount() int {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return len(fake.debugArgsForCall)
}

func (fake *FakeLogger) DebugArgsForCall(i int) []interface{} {
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	return fake.debugArgsForCall[i].arg1
}

func (fake *FakeLogger) Debugf(arg1 string, arg2 ...interface{}) {
	fake.debugfMutex.Lock()
	fake.debugfArgsForCall = append(fake.debugfArgsForCall, struct {
		arg1 string
		arg2 []interface{}
	}{arg1, arg2})
	fake.recordInvocation("Debugf", []interface{}{arg1, arg2})
	fake.debugfMutex.Unlock()
	if fake.DebugfStub != nil {
		fake.DebugfStub(arg1, arg2...)
	}
}

func (fake *FakeLogger) DebugfCallCount() int {
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	return len(fake.debugfArgsForCall)
}

func (fake *FakeLogger) DebugfArgsForCall(i int) (string, []interface{}) {
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	return fake.debugfArgsForCall[i].arg1, fake.debugfArgsForCall[i].arg2
}

func (fake *FakeLogger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.errorMutex.RLock()
	defer fake.errorMutex.RUnlock()
	fake.errorfMutex.RLock()
	defer fake.errorfMutex.RUnlock()
	fake.infoMutex.RLock()
	defer fake.infoMutex.RUnlock()
	fake.infofMutex.RLock()
	defer fake.infofMutex.RUnlock()
	fake.detailMutex.RLock()
	defer fake.detailMutex.RUnlock()
	fake.detailfMutex.RLock()
	defer fake.detailfMutex.RUnlock()
	fake.debugMutex.RLock()
	defer fake.debugMutex.RUnlock()
	fake.debugfMutex.RLock()
	defer fake.debugfMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLogger) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ tools.Logger = new(FakeLogger)
