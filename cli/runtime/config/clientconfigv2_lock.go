// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//nolint:dupl
package config

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/juju/fslock"
)

const (
	LocalTanzuConfigV2FileLock = ".tanzu-v2.lock"
	// DefaultConfigV2LockTimeout is the default time waiting on the filelock
	DefaultConfigV2LockTimeout = 10 * time.Minute
)

var tanzuConfigV2LockFile string

// tanzuConfigV2Lock used as a static lock variable that stores fslock
// This is used for interprocess locking of the config file
var tanzuConfigV2Lock *fslock.Lock

// v2Mutex is used to handle the locking behavior between concurrent calls
// within the existing process trying to acquire the lock
var v2Mutex sync.Mutex

// AcquireTanzuConfigV2Lock tries to acquire lock to update tanzu config file with timeout
func AcquireTanzuConfigV2Lock() {
	var err error

	if tanzuConfigV2LockFile == "" {
		path, err := ClientConfigV2Path()
		if err != nil {
			panic(fmt.Sprintf("cannot get config path while acquiring lock on tanzu config file, reason: %v", err))
		}
		tanzuConfigV2LockFile = filepath.Join(filepath.Dir(path), LocalTanzuConfigV2FileLock)
	}

	// using fslock to handle interprocess locking
	lock, err := getFileLockWithTimeOut(tanzuConfigV2LockFile, DefaultConfigV2LockTimeout)
	if err != nil {
		panic(fmt.Sprintf("cannot acquire lock for tanzu config file, reason: %v", err))
	}

	// Lock the mutex to prevent concurrent calls to acquire and configure the tanzuConfigLock
	v2Mutex.Lock()
	tanzuConfigV2Lock = lock
}

// ReleaseTanzuConfigV2Lock releases the lock if the tanzuConfigLock was acquired
func ReleaseTanzuConfigV2Lock() {
	if tanzuConfigV2Lock == nil {
		return
	}
	if errUnlock := tanzuConfigV2Lock.Unlock(); errUnlock != nil {
		panic(fmt.Sprintf("cannot release lock for tanzu config file, reason: %v", errUnlock))
	}

	tanzuConfigV2Lock = nil
	// Unlock the mutex to allow other concurrent calls to acquire and configure the tanzuConfigLock
	v2Mutex.Unlock()
}
