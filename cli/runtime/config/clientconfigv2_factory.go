// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Skip duplicate that matching with metadata config file

// Package config Provide API methods to Read/Write specific stanza of config file
//
//nolint:dupl
package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// getClientConfigV2Node retrieves the config from the local directory with file lock
func getClientConfigV2Node() (*yaml.Node, error) {
	// Acquire tanzu config v2 lock
	AcquireTanzuConfigV2Lock()
	defer ReleaseTanzuConfigV2Lock()
	return getClientConfigV2NodeNoLock()
}

// getClientConfigV2NodeNoLock retrieves the config from the local directory without acquiring the lock
func getClientConfigV2NodeNoLock() (*yaml.Node, error) {
	cfgPath, err := ClientConfigV2Path()
	if err != nil {
		return nil, errors.Wrap(err, "failed getting client config path")
	}
	bytes, err := os.ReadFile(cfgPath)
	if err != nil || len(bytes) == 0 {
		node, err := newClientConfigNode()
		if err != nil {
			return nil, errors.Wrap(err, " failed to create new client config")
		}
		return node, nil
	}
	var node yaml.Node
	err = yaml.Unmarshal(bytes, &node)
	if err != nil {
		return nil, errors.Wrap(err, "getClientConfigNoLock: failed to construct struct from config data")
	}
	node.Content[0].Style = 0
	return &node, nil
}

func persistClientConfigV2(node *yaml.Node) error {
	path, err := ClientConfigV2Path()
	if err != nil {
		return errors.Wrap(err, "could not find config path")
	}
	return persistNode(node, WithCfgPath(path))
}
