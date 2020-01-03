// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package geom

import (
	"fmt"
	"math"
)

// Point defines a location.
type Point struct {
	X, Y float64
}

// NewPoint creates a new Point.
func NewPoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

// NewPointPtr creates a new *Point.
func NewPointPtr(x, y float64) *Point {
	p := NewPoint(x, y)
	return &p
}

// Align modifies this Point to align with integer coordinates.
func (p *Point) Align() {
	p.X = math.Floor(p.X)
	p.Y = math.Floor(p.Y)
}

// Add modifies this Point by adding the supplied coordinates.
func (p *Point) Add(pt Point) {
	p.X += pt.X
	p.Y += pt.Y
}

// Subtract modifies this Point by subtracting the supplied coordinates.
func (p *Point) Subtract(pt Point) {
	p.X -= pt.X
	p.Y -= pt.Y
}

// String implements the fmt.Stringer interface.
func (p Point) String() string {
	return fmt.Sprintf("%v, %v", p.X, p.Y)
}
