// Code generated by counterfeiter. DO NOT EDIT.
package encryptionfakes

import (
	"crypto/cipher"
	"sync"

	"code.cloudfoundry.org/bbs/encryption"
)

type FakeKey struct {
	BlockStub        func() cipher.Block
	blockMutex       sync.RWMutex
	blockArgsForCall []struct {
	}
	blockReturns struct {
		result1 cipher.Block
	}
	blockReturnsOnCall map[int]struct {
		result1 cipher.Block
	}
	LabelStub        func() string
	labelMutex       sync.RWMutex
	labelArgsForCall []struct {
	}
	labelReturns struct {
		result1 string
	}
	labelReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeKey) Block() cipher.Block {
	fake.blockMutex.Lock()
	ret, specificReturn := fake.blockReturnsOnCall[len(fake.blockArgsForCall)]
	fake.blockArgsForCall = append(fake.blockArgsForCall, struct {
	}{})
	fake.recordInvocation("Block", []interface{}{})
	blockStubCopy := fake.BlockStub
	fake.blockMutex.Unlock()
	if blockStubCopy != nil {
		return blockStubCopy()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.blockReturns
	return fakeReturns.result1
}

func (fake *FakeKey) BlockCallCount() int {
	fake.blockMutex.RLock()
	defer fake.blockMutex.RUnlock()
	return len(fake.blockArgsForCall)
}

func (fake *FakeKey) BlockCalls(stub func() cipher.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = stub
}

func (fake *FakeKey) BlockReturns(result1 cipher.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = nil
	fake.blockReturns = struct {
		result1 cipher.Block
	}{result1}
}

func (fake *FakeKey) BlockReturnsOnCall(i int, result1 cipher.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = nil
	if fake.blockReturnsOnCall == nil {
		fake.blockReturnsOnCall = make(map[int]struct {
			result1 cipher.Block
		})
	}
	fake.blockReturnsOnCall[i] = struct {
		result1 cipher.Block
	}{result1}
}

func (fake *FakeKey) Label() string {
	fake.labelMutex.Lock()
	ret, specificReturn := fake.labelReturnsOnCall[len(fake.labelArgsForCall)]
	fake.labelArgsForCall = append(fake.labelArgsForCall, struct {
	}{})
	fake.recordInvocation("Label", []interface{}{})
	labelStubCopy := fake.LabelStub
	fake.labelMutex.Unlock()
	if labelStubCopy != nil {
		return labelStubCopy()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.labelReturns
	return fakeReturns.result1
}

func (fake *FakeKey) LabelCallCount() int {
	fake.labelMutex.RLock()
	defer fake.labelMutex.RUnlock()
	return len(fake.labelArgsForCall)
}

func (fake *FakeKey) LabelCalls(stub func() string) {
	fake.labelMutex.Lock()
	defer fake.labelMutex.Unlock()
	fake.LabelStub = stub
}

func (fake *FakeKey) LabelReturns(result1 string) {
	fake.labelMutex.Lock()
	defer fake.labelMutex.Unlock()
	fake.LabelStub = nil
	fake.labelReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeKey) LabelReturnsOnCall(i int, result1 string) {
	fake.labelMutex.Lock()
	defer fake.labelMutex.Unlock()
	fake.LabelStub = nil
	if fake.labelReturnsOnCall == nil {
		fake.labelReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.labelReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeKey) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.blockMutex.RLock()
	defer fake.blockMutex.RUnlock()
	fake.labelMutex.RLock()
	defer fake.labelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeKey) recordInvocation(key string, args []interface{}) {
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

var _ encryption.Key = new(FakeKey)
