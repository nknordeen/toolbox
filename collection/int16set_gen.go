// Code created from "set.go.tmpl" - don't edit by hand
//
// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package collection

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Int16Set holds a set of int16 values.
type Int16Set map[int16]bool

// NewInt16Set creates a new set from its input values.
func NewInt16Set(values ...int16) Int16Set {
	s := Int16Set{}
	s.Add(values...)
	return s
}

// Empty returns true if there are no values in the set.
func (s Int16Set) Empty() bool {
	return len(s) == 0
}

// Clear the set.
func (s Int16Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Add values to the set.
func (s Int16Set) Add(values ...int16) {
	for _, v := range values {
		s[v] = true
	}
}

// Contains returns true if the value exists within the set.
func (s Int16Set) Contains(value int16) bool {
	_, ok := s[value]
	return ok
}

// Clone returns a copy of the set.
func (s Int16Set) Clone() Int16Set {
	if s == nil {
		return nil
	}
	clone := Int16Set{}
	for value := range s {
		clone[value] = true
	}
	return clone
}

// Values returns all values in the set.
func (s Int16Set) Values() []int16 {
	values := make([]int16, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return values
}

// MarshalJSON implements the json.Marshaler interface.
func (s Int16Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s Int16Set) UnmarshalJSON(data []byte) error {
	var values []int16
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (s Int16Set) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(s.Values())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s Int16Set) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var values []int16
	if err := unmarshal(&values); err != nil {
		return err
	}
	s.Clear()
	s.Add(values...)
	return nil
}
