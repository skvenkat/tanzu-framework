// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/vmware-tanzu/tanzu-framework/cli/runtime/config/nodeutils"
)

type MigratedCfgItem struct {
	Type  yaml.Kind
	Value string
}

// GetMigratedCfgItems retrieves migratedCfgItems
func GetMigratedCfgItems() ([]*MigratedCfgItem, error) {
	// Retrieve config metadata node
	node, err := getMetadataNode()
	if err != nil {
		return nil, err
	}

	defaultMigratedItems := []*MigratedCfgItem{
		{Value: KeyContexts, Type: yaml.SequenceNode},
		{Value: KeyCurrentContext, Type: yaml.MappingNode},
	}

	items, err := getMigratedCfgItems(node)
	if err != nil {
		return defaultMigratedItems, nil
	}

	if len(items) == 0 {
		return defaultMigratedItems, nil
	}

	return items, nil
}

// SetMigratedCfgItem add migrated cfg item
func SetMigratedCfgItem(value string) error {
	// Retrieve config metadata node
	AcquireTanzuMetadataLock()
	defer ReleaseTanzuMetadataLock()
	node, err := getMetadataNodeNoLock()
	if err != nil {
		return err
	}

	// Add or update patch strategy
	err = setMigratedCfgItem(node, value)
	if err != nil {
		return err
	}
	return persistConfigMetadata(node)
}

// SetDefaultMigratedCfgItems set default migrated cfg item
func SetDefaultMigratedCfgItems(values []string) error {
	// Retrieve config metadata node
	AcquireTanzuMetadataLock()
	defer ReleaseTanzuMetadataLock()
	node, err := getMetadataNodeNoLock()
	if err != nil {
		return err
	}

	for _, val := range values {
		// Add or update patch strategy
		err = setMigratedCfgItem(node, val)
		if err != nil {
			return err
		}
	}

	return persistConfigMetadata(node)
}

func getMigratedCfgItems(node *yaml.Node) ([]*MigratedCfgItem, error) {
	// convert yaml node to metadata struct
	metadata, err := convertNodeToMetadata(node)
	if err != nil {
		return nil, err
	}

	if metadata != nil && metadata.ConfigMetadata != nil && metadata.ConfigMetadata.MigratedCfgItems != nil {
		mapped := make([]*MigratedCfgItem, len(metadata.ConfigMetadata.MigratedCfgItems))

		for i, val := range metadata.ConfigMetadata.MigratedCfgItems {
			mapped[i] = generateMigratedCfgItem(val)
		}

		return mapped, nil
	}

	return nil, errors.New("migrated config items not found")
}

func generateMigratedCfgItem(cfgItem string) *MigratedCfgItem {
	switch cfgItem {
	case KeyContexts:
		return &MigratedCfgItem{
			Type:  yaml.SequenceNode,
			Value: cfgItem,
		}
	case KeyServers:
		return &MigratedCfgItem{
			Type:  yaml.SequenceNode,
			Value: cfgItem,
		}
	case KeyCurrentContext:
		return &MigratedCfgItem{
			Type:  yaml.MappingNode,
			Value: cfgItem,
		}
	case KeyCurrentServer:
		return &MigratedCfgItem{
			Type:  yaml.ScalarNode,
			Value: cfgItem,
		}
	case KeyClientOptions:
		return &MigratedCfgItem{
			Type:  yaml.MappingNode,
			Value: cfgItem,
		}
	default:
		return &MigratedCfgItem{
			Type:  yaml.ScalarNode,
			Value: cfgItem,
		}
	}
}
func setMigratedCfgItem(node *yaml.Node, value string) error {
	// find migratedCfgItem node
	keys := []nodeutils.Key{
		{Name: KeyConfigMetadata, Type: yaml.MappingNode},
		{Name: KeyMigratedCfgItems, Type: yaml.SequenceNode},
	}
	migratedCfgItemNode := nodeutils.FindNode(node.Content[0], nodeutils.WithForceCreate(), nodeutils.WithKeys(keys))
	if migratedCfgItemNode == nil {
		return nodeutils.ErrNodeNotFound
	}
	newItemNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
		Style: 0,
		Tag:   "!!str",
	}

	exists := false
	var result []*yaml.Node
	for _, item := range migratedCfgItemNode.Content {
		if item.Value == value {
			exists = true
			result = append(result, newItemNode)
			continue
		}
		result = append(result, item)
	}

	if !exists {
		result = append(result, newItemNode)
	}

	migratedCfgItemNode.Content = result

	return nil
}
