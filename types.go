/* types.go - Go wrapper for Emacs module API.

Copyright (C) 2016 Yann Hodique <yann.hodique@gmail.com>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or (at
your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.  */

package goemacs

/*
#include "include/wrapper.h"
*/
import "C"

type Value interface {
	getVal() C.emacs_value
}

type baseValue struct {
	val C.emacs_value
	env *Environment
}

func (v baseValue) getVal() C.emacs_value {
	return v.val
}

type String interface {
	Value
	String() string
}

type stringValue struct {
	baseValue
}

func (s stringValue) String() string {
	res, _ := s.env.GoString(s)
	return res
}

type Callable interface {
	Value
	Callable() bool
}

type Symbol interface {
	Callable
	makeCallable()
}

type symbolValue struct {
	baseValue
	callable bool
}

func (s symbolValue) Callable() bool {
	return s.callable
}

func (s symbolValue) makeCallable() {
	s.callable = true
}

type Function interface {
	Callable
}

type functionValue struct {
	baseValue
}

func (f functionValue) Callable() bool {
	return true
}
