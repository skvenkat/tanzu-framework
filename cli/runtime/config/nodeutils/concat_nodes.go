// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package nodeutils

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ConcatNodes to merge two yaml nodes src(source) to dst(destination) node
func ConcatNodes(src, dst *yaml.Node) error {
	if dst == nil || dst.Content == nil {
		*dst = *src
		return nil
	}

	err := checkErrors(src, dst)
	if err != nil {
		return err
	}

	switch src.Kind {
	case yaml.MappingNode:
		concatMappingNode(src, dst)
	case yaml.SequenceNode:
		concatSequenceNodes(src, dst)
	case yaml.DocumentNode:
		err := ConcatNodes(src.Content[0], dst.Content[0])
		if err != nil {
			return errors.New("at key " + src.Content[0].Value + ": " + err.Error())
		}
	case yaml.ScalarNode:
		concatScalarNodes(src, dst)
	default:
		return errors.New("can only merge mapping and sequence nodes")
	}
	return nil
}

func concatMappingNode(src, dst *yaml.Node) {
	for i := 0; i < len(src.Content); i += 2 {
		found := false
		for j := 0; j < len(dst.Content); j += 2 {
			if ok, _ := equalScalars(src.Content[i], dst.Content[j]); ok {
				found = true
				break
			}
		}
		if !found {
			dst.Content = append(dst.Content, src.Content[i:i+2]...)
		}
	}
}

func concatSequenceNodes(src, dst *yaml.Node) {
	if dst.Content == nil && len(dst.Content) == 0 {
		dst.Content = src.Content
	}
}

func concatScalarNodes(src, dst *yaml.Node) {
	if dst.Value != "" {
		dst.Value = src.Value
	}
}
