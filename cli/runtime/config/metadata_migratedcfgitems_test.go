// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetMigratedCfgItems(t *testing.T) {
	// setup
	func() {
		LocalDirName = TestLocalDirName
	}()

	defer func() {
		cleanupDir(LocalDirName)
	}()

	tests := []struct {
		name   string
		in     string
		out    []*MigratedCfgItem
		errStr string
	}{
		{
			name: "success empty items",
			in:   ``,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
				{Value: KeyCurrentContext, Type: yaml.MappingNode},
			},
		},
		{
			name: "success with patch strategies",
			in: `configMetadata:
    migratedCfgItems:
        - contexts
        - currentContext
        - clientOptions
`,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
				{Value: KeyCurrentContext, Type: yaml.MappingNode},
				{Value: KeyClientOptions, Type: yaml.MappingNode},
			},
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			// Setup
			var node yaml.Node
			err := yaml.Unmarshal([]byte(spec.in), &node)
			assert.NoError(t, err)

			err = persistConfigMetadata(&node)
			assert.NoError(t, err)

			//Test case
			c, err := GetMigratedCfgItems()
			assert.Equal(t, c, spec.out)
			assert.NoError(t, err)
		})
	}
}

func TestSetMigratedCfgItem(t *testing.T) {
	// Setup config data
	f1, err := os.CreateTemp("", "tanzu_config")
	assert.Nil(t, err)
	err = os.WriteFile(f1.Name(), []byte(""), 0644)
	assert.Nil(t, err)

	err = os.Setenv(EnvConfigKey, f1.Name())
	assert.NoError(t, err)

	f2, err := os.CreateTemp("", "tanzu_config_v2")
	assert.Nil(t, err)
	err = os.WriteFile(f2.Name(), []byte(""), 0644)
	assert.Nil(t, err)

	err = os.Setenv(EnvConfigV2Key, f2.Name())
	assert.NoError(t, err)

	//Setup metadata
	fMeta, err := os.CreateTemp("", "tanzu_config_metadata")
	assert.Nil(t, err)
	err = os.WriteFile(fMeta.Name(), []byte(""), 0644)
	assert.Nil(t, err)

	err = os.Setenv(EnvConfigMetadataKey, fMeta.Name())
	assert.NoError(t, err)

	// Cleanup
	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f1.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f2.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(fMeta.Name())

	tests := []struct {
		name   string
		key    string
		out    []*MigratedCfgItem
		errStr string
	}{
		{
			name: "success add new item",
			key:  KeyContexts,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
			},
		},
		{
			name: "success add another item",
			key:  KeyCurrentContext,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
				{Value: KeyCurrentContext, Type: yaml.MappingNode},
			},
		},
		{
			name: "success add existing item",
			key:  KeyCurrentContext,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
				{Value: KeyCurrentContext, Type: yaml.MappingNode},
			},
		},
		{
			name: "success add another item",
			key:  KeyClientOptions,
			out: []*MigratedCfgItem{
				{Value: KeyContexts, Type: yaml.SequenceNode},
				{Value: KeyCurrentContext, Type: yaml.MappingNode},
				{Value: KeyClientOptions, Type: yaml.MappingNode},
			},
		},
	}
	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			err := SetMigratedCfgItem(spec.key)
			if spec.errStr != "" {
				assert.Equal(t, err.Error(), spec.errStr)
				c, err := GetConfigMetadataPatchStrategy()
				assert.NoError(t, err)
				assert.Equal(t, "", c[spec.key])
			} else {
				assert.NoError(t, err)
				c, err := GetMigratedCfgItems()
				assert.NoError(t, err)
				assert.Equal(t, spec.out, c)
			}
		})
	}
}
