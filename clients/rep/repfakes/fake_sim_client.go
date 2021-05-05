// Code generated by counterfeiter. DO NOT EDIT.
package repfakes

import (
	"net/http"
	"sync"
	"time"

	"code.cloudfoundry.org/bbs/clients/rep"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/lager"
)

type FakeSimClient struct {
	CancelTaskStub        func(lager.Logger, string) error
	cancelTaskMutex       sync.RWMutex
	cancelTaskArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	cancelTaskReturns struct {
		result1 error
	}
	cancelTaskReturnsOnCall map[int]struct {
		result1 error
	}
	PerformStub        func(lager.Logger, rep.Work) (rep.Work, error)
	performMutex       sync.RWMutex
	performArgsForCall []struct {
		arg1 lager.Logger
		arg2 rep.Work
	}
	performReturns struct {
		result1 rep.Work
		result2 error
	}
	performReturnsOnCall map[int]struct {
		result1 rep.Work
		result2 error
	}
	ResetStub        func() error
	resetMutex       sync.RWMutex
	resetArgsForCall []struct {
	}
	resetReturns struct {
		result1 error
	}
	resetReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateClientStub        func(*http.Client)
	setStateClientMutex       sync.RWMutex
	setStateClientArgsForCall []struct {
		arg1 *http.Client
	}
	StateStub        func(lager.Logger) (rep.CellState, error)
	stateMutex       sync.RWMutex
	stateArgsForCall []struct {
		arg1 lager.Logger
	}
	stateReturns struct {
		result1 rep.CellState
		result2 error
	}
	stateReturnsOnCall map[int]struct {
		result1 rep.CellState
		result2 error
	}
	StateClientTimeoutStub        func() time.Duration
	stateClientTimeoutMutex       sync.RWMutex
	stateClientTimeoutArgsForCall []struct {
	}
	stateClientTimeoutReturns struct {
		result1 time.Duration
	}
	stateClientTimeoutReturnsOnCall map[int]struct {
		result1 time.Duration
	}
	StopLRPInstanceStub        func(lager.Logger, models.ActualLRPKey, models.ActualLRPInstanceKey) error
	stopLRPInstanceMutex       sync.RWMutex
	stopLRPInstanceArgsForCall []struct {
		arg1 lager.Logger
		arg2 models.ActualLRPKey
		arg3 models.ActualLRPInstanceKey
	}
	stopLRPInstanceReturns struct {
		result1 error
	}
	stopLRPInstanceReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSimClient) CancelTask(arg1 lager.Logger, arg2 string) error {
	fake.cancelTaskMutex.Lock()
	ret, specificReturn := fake.cancelTaskReturnsOnCall[len(fake.cancelTaskArgsForCall)]
	fake.cancelTaskArgsForCall = append(fake.cancelTaskArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	stub := fake.CancelTaskStub
	fakeReturns := fake.cancelTaskReturns
	fake.recordInvocation("CancelTask", []interface{}{arg1, arg2})
	fake.cancelTaskMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSimClient) CancelTaskCallCount() int {
	fake.cancelTaskMutex.RLock()
	defer fake.cancelTaskMutex.RUnlock()
	return len(fake.cancelTaskArgsForCall)
}

func (fake *FakeSimClient) CancelTaskCalls(stub func(lager.Logger, string) error) {
	fake.cancelTaskMutex.Lock()
	defer fake.cancelTaskMutex.Unlock()
	fake.CancelTaskStub = stub
}

func (fake *FakeSimClient) CancelTaskArgsForCall(i int) (lager.Logger, string) {
	fake.cancelTaskMutex.RLock()
	defer fake.cancelTaskMutex.RUnlock()
	argsForCall := fake.cancelTaskArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeSimClient) CancelTaskReturns(result1 error) {
	fake.cancelTaskMutex.Lock()
	defer fake.cancelTaskMutex.Unlock()
	fake.CancelTaskStub = nil
	fake.cancelTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) CancelTaskReturnsOnCall(i int, result1 error) {
	fake.cancelTaskMutex.Lock()
	defer fake.cancelTaskMutex.Unlock()
	fake.CancelTaskStub = nil
	if fake.cancelTaskReturnsOnCall == nil {
		fake.cancelTaskReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.cancelTaskReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) Perform(arg1 lager.Logger, arg2 rep.Work) (rep.Work, error) {
	fake.performMutex.Lock()
	ret, specificReturn := fake.performReturnsOnCall[len(fake.performArgsForCall)]
	fake.performArgsForCall = append(fake.performArgsForCall, struct {
		arg1 lager.Logger
		arg2 rep.Work
	}{arg1, arg2})
	stub := fake.PerformStub
	fakeReturns := fake.performReturns
	fake.recordInvocation("Perform", []interface{}{arg1, arg2})
	fake.performMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSimClient) PerformCallCount() int {
	fake.performMutex.RLock()
	defer fake.performMutex.RUnlock()
	return len(fake.performArgsForCall)
}

func (fake *FakeSimClient) PerformCalls(stub func(lager.Logger, rep.Work) (rep.Work, error)) {
	fake.performMutex.Lock()
	defer fake.performMutex.Unlock()
	fake.PerformStub = stub
}

func (fake *FakeSimClient) PerformArgsForCall(i int) (lager.Logger, rep.Work) {
	fake.performMutex.RLock()
	defer fake.performMutex.RUnlock()
	argsForCall := fake.performArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeSimClient) PerformReturns(result1 rep.Work, result2 error) {
	fake.performMutex.Lock()
	defer fake.performMutex.Unlock()
	fake.PerformStub = nil
	fake.performReturns = struct {
		result1 rep.Work
		result2 error
	}{result1, result2}
}

func (fake *FakeSimClient) PerformReturnsOnCall(i int, result1 rep.Work, result2 error) {
	fake.performMutex.Lock()
	defer fake.performMutex.Unlock()
	fake.PerformStub = nil
	if fake.performReturnsOnCall == nil {
		fake.performReturnsOnCall = make(map[int]struct {
			result1 rep.Work
			result2 error
		})
	}
	fake.performReturnsOnCall[i] = struct {
		result1 rep.Work
		result2 error
	}{result1, result2}
}

func (fake *FakeSimClient) Reset() error {
	fake.resetMutex.Lock()
	ret, specificReturn := fake.resetReturnsOnCall[len(fake.resetArgsForCall)]
	fake.resetArgsForCall = append(fake.resetArgsForCall, struct {
	}{})
	stub := fake.ResetStub
	fakeReturns := fake.resetReturns
	fake.recordInvocation("Reset", []interface{}{})
	fake.resetMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSimClient) ResetCallCount() int {
	fake.resetMutex.RLock()
	defer fake.resetMutex.RUnlock()
	return len(fake.resetArgsForCall)
}

func (fake *FakeSimClient) ResetCalls(stub func() error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = stub
}

func (fake *FakeSimClient) ResetReturns(result1 error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = nil
	fake.resetReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) ResetReturnsOnCall(i int, result1 error) {
	fake.resetMutex.Lock()
	defer fake.resetMutex.Unlock()
	fake.ResetStub = nil
	if fake.resetReturnsOnCall == nil {
		fake.resetReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.resetReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) SetStateClient(arg1 *http.Client) {
	fake.setStateClientMutex.Lock()
	fake.setStateClientArgsForCall = append(fake.setStateClientArgsForCall, struct {
		arg1 *http.Client
	}{arg1})
	stub := fake.SetStateClientStub
	fake.recordInvocation("SetStateClient", []interface{}{arg1})
	fake.setStateClientMutex.Unlock()
	if stub != nil {
		fake.SetStateClientStub(arg1)
	}
}

func (fake *FakeSimClient) SetStateClientCallCount() int {
	fake.setStateClientMutex.RLock()
	defer fake.setStateClientMutex.RUnlock()
	return len(fake.setStateClientArgsForCall)
}

func (fake *FakeSimClient) SetStateClientCalls(stub func(*http.Client)) {
	fake.setStateClientMutex.Lock()
	defer fake.setStateClientMutex.Unlock()
	fake.SetStateClientStub = stub
}

func (fake *FakeSimClient) SetStateClientArgsForCall(i int) *http.Client {
	fake.setStateClientMutex.RLock()
	defer fake.setStateClientMutex.RUnlock()
	argsForCall := fake.setStateClientArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSimClient) State(arg1 lager.Logger) (rep.CellState, error) {
	fake.stateMutex.Lock()
	ret, specificReturn := fake.stateReturnsOnCall[len(fake.stateArgsForCall)]
	fake.stateArgsForCall = append(fake.stateArgsForCall, struct {
		arg1 lager.Logger
	}{arg1})
	stub := fake.StateStub
	fakeReturns := fake.stateReturns
	fake.recordInvocation("State", []interface{}{arg1})
	fake.stateMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSimClient) StateCallCount() int {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	return len(fake.stateArgsForCall)
}

func (fake *FakeSimClient) StateCalls(stub func(lager.Logger) (rep.CellState, error)) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = stub
}

func (fake *FakeSimClient) StateArgsForCall(i int) lager.Logger {
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	argsForCall := fake.stateArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeSimClient) StateReturns(result1 rep.CellState, result2 error) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = nil
	fake.stateReturns = struct {
		result1 rep.CellState
		result2 error
	}{result1, result2}
}

func (fake *FakeSimClient) StateReturnsOnCall(i int, result1 rep.CellState, result2 error) {
	fake.stateMutex.Lock()
	defer fake.stateMutex.Unlock()
	fake.StateStub = nil
	if fake.stateReturnsOnCall == nil {
		fake.stateReturnsOnCall = make(map[int]struct {
			result1 rep.CellState
			result2 error
		})
	}
	fake.stateReturnsOnCall[i] = struct {
		result1 rep.CellState
		result2 error
	}{result1, result2}
}

func (fake *FakeSimClient) StateClientTimeout() time.Duration {
	fake.stateClientTimeoutMutex.Lock()
	ret, specificReturn := fake.stateClientTimeoutReturnsOnCall[len(fake.stateClientTimeoutArgsForCall)]
	fake.stateClientTimeoutArgsForCall = append(fake.stateClientTimeoutArgsForCall, struct {
	}{})
	stub := fake.StateClientTimeoutStub
	fakeReturns := fake.stateClientTimeoutReturns
	fake.recordInvocation("StateClientTimeout", []interface{}{})
	fake.stateClientTimeoutMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSimClient) StateClientTimeoutCallCount() int {
	fake.stateClientTimeoutMutex.RLock()
	defer fake.stateClientTimeoutMutex.RUnlock()
	return len(fake.stateClientTimeoutArgsForCall)
}

func (fake *FakeSimClient) StateClientTimeoutCalls(stub func() time.Duration) {
	fake.stateClientTimeoutMutex.Lock()
	defer fake.stateClientTimeoutMutex.Unlock()
	fake.StateClientTimeoutStub = stub
}

func (fake *FakeSimClient) StateClientTimeoutReturns(result1 time.Duration) {
	fake.stateClientTimeoutMutex.Lock()
	defer fake.stateClientTimeoutMutex.Unlock()
	fake.StateClientTimeoutStub = nil
	fake.stateClientTimeoutReturns = struct {
		result1 time.Duration
	}{result1}
}

func (fake *FakeSimClient) StateClientTimeoutReturnsOnCall(i int, result1 time.Duration) {
	fake.stateClientTimeoutMutex.Lock()
	defer fake.stateClientTimeoutMutex.Unlock()
	fake.StateClientTimeoutStub = nil
	if fake.stateClientTimeoutReturnsOnCall == nil {
		fake.stateClientTimeoutReturnsOnCall = make(map[int]struct {
			result1 time.Duration
		})
	}
	fake.stateClientTimeoutReturnsOnCall[i] = struct {
		result1 time.Duration
	}{result1}
}

func (fake *FakeSimClient) StopLRPInstance(arg1 lager.Logger, arg2 models.ActualLRPKey, arg3 models.ActualLRPInstanceKey) error {
	fake.stopLRPInstanceMutex.Lock()
	ret, specificReturn := fake.stopLRPInstanceReturnsOnCall[len(fake.stopLRPInstanceArgsForCall)]
	fake.stopLRPInstanceArgsForCall = append(fake.stopLRPInstanceArgsForCall, struct {
		arg1 lager.Logger
		arg2 models.ActualLRPKey
		arg3 models.ActualLRPInstanceKey
	}{arg1, arg2, arg3})
	stub := fake.StopLRPInstanceStub
	fakeReturns := fake.stopLRPInstanceReturns
	fake.recordInvocation("StopLRPInstance", []interface{}{arg1, arg2, arg3})
	fake.stopLRPInstanceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeSimClient) StopLRPInstanceCallCount() int {
	fake.stopLRPInstanceMutex.RLock()
	defer fake.stopLRPInstanceMutex.RUnlock()
	return len(fake.stopLRPInstanceArgsForCall)
}

func (fake *FakeSimClient) StopLRPInstanceCalls(stub func(lager.Logger, models.ActualLRPKey, models.ActualLRPInstanceKey) error) {
	fake.stopLRPInstanceMutex.Lock()
	defer fake.stopLRPInstanceMutex.Unlock()
	fake.StopLRPInstanceStub = stub
}

func (fake *FakeSimClient) StopLRPInstanceArgsForCall(i int) (lager.Logger, models.ActualLRPKey, models.ActualLRPInstanceKey) {
	fake.stopLRPInstanceMutex.RLock()
	defer fake.stopLRPInstanceMutex.RUnlock()
	argsForCall := fake.stopLRPInstanceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSimClient) StopLRPInstanceReturns(result1 error) {
	fake.stopLRPInstanceMutex.Lock()
	defer fake.stopLRPInstanceMutex.Unlock()
	fake.StopLRPInstanceStub = nil
	fake.stopLRPInstanceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) StopLRPInstanceReturnsOnCall(i int, result1 error) {
	fake.stopLRPInstanceMutex.Lock()
	defer fake.stopLRPInstanceMutex.Unlock()
	fake.StopLRPInstanceStub = nil
	if fake.stopLRPInstanceReturnsOnCall == nil {
		fake.stopLRPInstanceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.stopLRPInstanceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeSimClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cancelTaskMutex.RLock()
	defer fake.cancelTaskMutex.RUnlock()
	fake.performMutex.RLock()
	defer fake.performMutex.RUnlock()
	fake.resetMutex.RLock()
	defer fake.resetMutex.RUnlock()
	fake.setStateClientMutex.RLock()
	defer fake.setStateClientMutex.RUnlock()
	fake.stateMutex.RLock()
	defer fake.stateMutex.RUnlock()
	fake.stateClientTimeoutMutex.RLock()
	defer fake.stateClientTimeoutMutex.RUnlock()
	fake.stopLRPInstanceMutex.RLock()
	defer fake.stopLRPInstanceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSimClient) recordInvocation(key string, args []interface{}) {
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

var _ rep.SimClient = new(FakeSimClient)
