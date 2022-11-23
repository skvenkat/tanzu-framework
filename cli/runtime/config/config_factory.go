// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package config Provide API methods to Read/Write specific stanza of config file
package config

import (
	"os"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/tanzu-framework/cli/runtime/config/nodeutils"
)

func persistConfig(node *yaml.Node) error {
	migrate, err := ShouldMigrateToNewConfig()
	if err != nil {
		migrate = false
	}

	if migrate {
		return persistClientConfigV2(node)
	}

	cfgNode, err := getClientConfigNodeNoLock()
	if err != nil {
		return err
	}

	// deep copy of change node
	var cfgNodeToPersist yaml.Node
	data, err := yaml.Marshal(node)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &cfgNodeToPersist)
	if err != nil {
		return err
	}

	cfgV2Node, err := getClientConfigV2NodeNoLock()
	if err != nil {
		return err
	}

	migratedItems, err := GetMigratedCfgItems()
	if err != nil {
		return err
	}

	for _, migratedItem := range migratedItems {
		// Find the migrated node from the updated node
		itemNode := nodeutils.FindNode(cfgNodeToPersist.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys([]nodeutils.Key{
			{Name: migratedItem.Value, Type: migratedItem.Type},
		}))

		// Find the migrated node from config.yaml
		itemCfgNode := nodeutils.FindNode(cfgNode.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys([]nodeutils.Key{
			{Name: migratedItem.Value, Type: migratedItem.Type},
		}))

		// Find the migrated node from config-v2.yaml
		itemV2Node := nodeutils.FindNode(cfgV2Node.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys([]nodeutils.Key{
			{Name: migratedItem.Value, Type: migratedItem.Type},
		}))

		*itemV2Node = *itemNode

		// Reset migrated node in config.yaml
		*itemNode = *itemCfgNode
	}

	// Store the not migrated config data to config.yaml
	err = persistClientConfig(&cfgNodeToPersist)
	if err != nil {
		return err
	}

	// Store the migrated config data to config-v2.yaml
	err = persistClientConfigV2(cfgV2Node)
	if err != nil {
		return err
	}

	// Store the config data to legacy client config file/location
	err = persistLegacyClientConfig(node)
	if err != nil {
		return err
	}

	return nil
}

// persistNode stores/writes the yaml node to config.yaml
func persistNode(node *yaml.Node, opts ...CfgOpts) error {
	configurations := &CfgOptions{}
	for _, opt := range opts {
		opt(configurations)
	}
	cfgPathExists, err := fileExists(configurations.CfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to check config path existence")
	}
	if !cfgPathExists {
		localDir, err := LocalDir()
		if err != nil {
			return errors.Wrap(err, "could not find local tanzu dir for OS")
		}
		if err := os.MkdirAll(localDir, 0755); err != nil {
			return errors.Wrap(err, "could not make local tanzu directory")
		}
	}
	data, err := yaml.Marshal(node)
	if err != nil {
		return errors.Wrap(err, "failed to marshal nodeutils")
	}
	err = os.WriteFile(configurations.CfgPath, data, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write the config to file")
	}
	return nil
}

func getClientConfig() (*yaml.Node, error) {
	migrate, err := ShouldMigrateToNewConfig()
	if err != nil {
		migrate = false
	}

	if migrate {
		return getClientConfigV2Node()
	}
	return getMultiConfig()
}

// Use for CREATE/Update/Delete apis
// getClientConfig retrieve config data from config.yaml, config-alt.yaml based on feature flag with no file lock
func getClientConfigNoLock() (*yaml.Node, error) {
	// Check config migration feature flag
	migrate, err := ShouldMigrateToNewConfig()
	if err != nil {
		migrate = false
	}

	if migrate {
		return getClientConfigV2NodeNoLock()
	}
	return getMultiConfigNoLock()
}

// getMultiConfig retrieves combined config.yaml and config.alt.yaml
func getMultiConfig() (*yaml.Node, error) {
	node1, err := getClientConfigNode()
	if err != nil {
		return node1, err
	}
	node2, err := getClientConfigV2Node()
	if err != nil {
		return node2, err
	}
	// Merge node1 and node2
	err = nodeutils.ConcatNodes(node1, node2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge node1 and node2")
	}
	return node2, err
}

// getMultiConfigNoLock retrieves combined config.yaml and config.alt.yaml
func getMultiConfigNoLock() (*yaml.Node, error) {
	cfgNode, err := getClientConfigNodeNoLock()
	if err != nil {
		return cfgNode, err
	}
	cfgV2Node, err := getClientConfigV2NodeNoLock()
	if err != nil {
		return cfgV2Node, err
	}
	// Merge node1 and node2
	err = nodeutils.ConcatNodes(cfgNode, cfgV2Node)
	if err != nil {
		return nil, errors.Wrap(err, "failed to merge node1 and node2")
	}
	return cfgV2Node, err
}
