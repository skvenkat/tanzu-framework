// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndDeleteConfigMetadataFeatureFlags(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tests := []struct {
		name  string
		key   string
		value bool
	}{
		{
			name:  "success context-aware-cli-for-plugins",
			key:   "migrateToNewConfig",
			value: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetConfigMetadataFeatureFlag(tc.key, strconv.FormatBool(tc.value))
			assert.NoError(t, err)

			migrate, err := ShouldMigrateToNewConfig()
			assert.NoError(t, err)
			assert.Equal(t, tc.value, migrate)

			err = DeleteConfigMetadataFeatureFlag(tc.key)
			assert.NoError(t, err)

			migrate, err = ShouldMigrateToNewConfig()
			assert.Equal(t, "not found", err.Error())
			assert.Equal(t, tc.value, migrate)

			err = SetConfigMetadataFeatureFlag(tc.key, strconv.FormatBool(!tc.value))
			assert.NoError(t, err)

			migrate, err = ShouldMigrateToNewConfig()
			assert.NoError(t, err)
			assert.Equal(t, !tc.value, migrate)
		})
	}
}

func TestSetConfigMetadataFeatureFlag(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()
	defer func() {
		cleanupDir(LocalDirName)
	}()
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "success disable migrateToNewConfig",
			key:   "migrateToNewConfig",
			value: "false",
		},
		{
			name:  "success enable migrateToNewConfig",
			key:   "migrateToNewConfig",
			value: "true",
		},
		{
			name:  "success disable migrateToNewConfig",
			key:   "migrateToNewConfig",
			value: "false",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := SetConfigMetadataFeatureFlag(tc.key, tc.value)
			assert.NoError(t, err)

			migrate, err := ShouldMigrateToNewConfig()
			assert.NoError(t, err)
			expected, err := strconv.ParseBool(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, expected, migrate)
		})
	}
}
