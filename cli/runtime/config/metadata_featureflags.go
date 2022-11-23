// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strings"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/tanzu-framework/cli/runtime/config/nodeutils"
)

const (
	FeatureMigrateToNewConfig = "migrateToNewConfig"
)

// GetConfigMetadataFeatureFlags retrieves feature flags
func GetConfigMetadataFeatureFlags() (map[string]string, error) {
	// Retrieve Metadata config node
	node, err := getMetadataNode()
	if err != nil {
		return nil, err
	}
	return getConfigMetadataFeatureFlags(node)
}

func getConfigMetadataFeatureFlags(node *yaml.Node) (map[string]string, error) {
	cfgMetadata, err := convertNodeToMetadata(node)
	if err != nil {
		return nil, err
	}
	if cfgMetadata != nil && cfgMetadata.ConfigMetadata != nil &&
		cfgMetadata.ConfigMetadata.FeatureFlags != nil {
		return cfgMetadata.ConfigMetadata.FeatureFlags, nil
	}
	return nil, nil
}

// IsConfigMetadataFeatureEnabled checks and returns whether specific plugin and key is true
func IsConfigMetadataFeatureEnabled(key string) (bool, error) {
	node, err := getMetadataNode()
	if err != nil {
		return false, err
	}
	val, err := getConfigMetadataFeature(node, key)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(val, "true"), nil
}

// ShouldMigrateToNewConfig checks migrateToNewConfig feature flag
func ShouldMigrateToNewConfig() (bool, error) {
	return IsConfigMetadataFeatureEnabled(FeatureMigrateToNewConfig)
}

func getConfigMetadataFeature(node *yaml.Node, key string) (string, error) {
	cfgMetadata, err := convertNodeToMetadata(node)
	if err != nil {
		return "", err
	}

	if cfgMetadata == nil || cfgMetadata.ConfigMetadata == nil ||
		cfgMetadata.ConfigMetadata.FeatureFlags == nil {
		return "", errors.New("not found")
	}

	if val, ok := cfgMetadata.ConfigMetadata.FeatureFlags[key]; ok {
		return val, nil
	}
	return "", errors.New("not found")
}

// DeleteConfigMetadataFeatureFlag delete the env entry of specified key
func DeleteConfigMetadataFeatureFlag(key string) error {
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getMetadataNode()

	if err != nil {
		return err
	}
	err = deleteFeatureFlag(node, key)
	if err != nil {
		return err
	}
	return persistConfigMetadata(node)
}

func deleteFeatureFlag(node *yaml.Node, key string) (err error) {
	// find feature flags node
	keys := []nodeutils.Key{
		{Name: KeyConfigMetadata, Type: yaml.MappingNode},
		{Name: KeyFeatureFlags, Type: yaml.MappingNode},
	}
	featureFlagsNode := nodeutils.FindNode(node.Content[0], nodeutils.WithKeys(keys))
	if featureFlagsNode == nil {
		return err
	}

	// convert env nodes to map
	featureFlags, err := nodeutils.ConvertNodeToMap(featureFlagsNode)
	if err != nil {
		return err
	}

	// delete the specified entry in the map
	delete(featureFlags, key)

	// convert updated map to env node
	newFeatureFlagNode, err := nodeutils.ConvertMapToNode(featureFlags)
	if err != nil {
		return err
	}
	featureFlagsNode.Content = newFeatureFlagNode.Content[0].Content
	return nil
}

// SetConfigMetadataFeatureFlag add or update a env key and value
func SetConfigMetadataFeatureFlag(key, value string) (err error) {
	AcquireTanzuConfigLock()
	defer ReleaseTanzuConfigLock()
	node, err := getMetadataNode()
	if err != nil {
		return err
	}
	persist, err := setFeatureFlag(node, key, value)
	if err != nil {
		return err
	}
	if persist {
		return persistConfigMetadata(node)
	}
	return err
}

//nolint:dupl
func setFeatureFlag(node *yaml.Node, key, value string) (persist bool, err error) {
	// find feature flags stanza node
	keys := []nodeutils.Key{
		{Name: KeyConfigMetadata, Type: yaml.MappingNode},
		{Name: KeyFeatureFlags, Type: yaml.MappingNode},
	}
	featureFlagsNode := nodeutils.FindNode(node.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys(keys))
	if featureFlagsNode == nil {
		return persist, err
	}
	// convert env node to map
	featureFlags, err := nodeutils.ConvertNodeToMap(featureFlagsNode)
	if err != nil {
		return persist, err
	}
	// add or update the envs map per specified key value pair
	if len(featureFlags) == 0 || featureFlags[key] != value {
		featureFlags[key] = value
		persist = true
	}
	// convert map to yaml node
	newFeatureFlagsNode, err := nodeutils.ConvertMapToNode(featureFlags)
	if err != nil {
		return persist, err
	}
	featureFlagsNode.Content = newFeatureFlagsNode.Content[0].Content
	return persist, err
}
