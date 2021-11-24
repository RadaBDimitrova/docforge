// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0
// Code generated by counterfeiter. DO NOT EDIT.
package gitfakes

import (
	"sync"

	"github.com/gardener/docforge/pkg/git"
	gita "github.com/go-git/go-git/v5"
)

type FakeRepositoryWorktree struct {
	CheckoutStub        func(*gita.CheckoutOptions) error
	checkoutMutex       sync.RWMutex
	checkoutArgsForCall []struct {
		arg1 *gita.CheckoutOptions
	}
	checkoutReturns struct {
		result1 error
	}
	checkoutReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRepositoryWorktree) Checkout(arg1 *gita.CheckoutOptions) error {
	fake.checkoutMutex.Lock()
	ret, specificReturn := fake.checkoutReturnsOnCall[len(fake.checkoutArgsForCall)]
	fake.checkoutArgsForCall = append(fake.checkoutArgsForCall, struct {
		arg1 *gita.CheckoutOptions
	}{arg1})
	stub := fake.CheckoutStub
	fakeReturns := fake.checkoutReturns
	fake.recordInvocation("Checkout", []interface{}{arg1})
	fake.checkoutMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeRepositoryWorktree) CheckoutCallCount() int {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	return len(fake.checkoutArgsForCall)
}

func (fake *FakeRepositoryWorktree) CheckoutCalls(stub func(*gita.CheckoutOptions) error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = stub
}

func (fake *FakeRepositoryWorktree) CheckoutArgsForCall(i int) *gita.CheckoutOptions {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	argsForCall := fake.checkoutArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeRepositoryWorktree) CheckoutReturns(result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	fake.checkoutReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepositoryWorktree) CheckoutReturnsOnCall(i int, result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	if fake.checkoutReturnsOnCall == nil {
		fake.checkoutReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkoutReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepositoryWorktree) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeRepositoryWorktree) recordInvocation(key string, args []interface{}) {
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

var _ git.RepositoryWorktree = new(FakeRepositoryWorktree)