package yajl

// #include "api/yajl_tree.h"
// #include "shim-go.h"
import "C"

import (
	"errors"
)

type JsonObjectCompareOption int

const (
	DEFAULT                   = 0
	ALL_VALUE_TREAT_AS_STRING = (1 << 0)
)

type JsonObject interface {
	Compare(other JsonObject, option JsonObjectCompareOption) bool
}

type JsonMap struct {
	inner map[string]JsonObject
}

type JsonArray struct {
	inner []JsonObject
}

type JsonString struct {
	realtype int
	inner    string
}

func (jString *JsonString) Compare(other JsonObject, option JsonObjectCompareOption) bool {
	if other == nil {
		return false
	}

	if oString, ok := other.(*JsonString); ok {
		if option&ALL_VALUE_TREAT_AS_STRING > 0 {
			return oString.inner == jString.inner
		} else {
			if oString.realtype != jString.realtype {
				return false
			} else {
				return oString.inner == jString.inner
			}
		}
	} else {
		return false
	}
}

func (jMap *JsonMap) Compare(other JsonObject, option JsonObjectCompareOption) bool {
	if other == nil {
		return false
	}

	if oMap, ok := other.(*JsonMap); ok {
		if len(oMap.inner) != len(jMap.inner) {
			return false
		}

		for k, v := range jMap.inner {
			if v.Compare(oMap.inner[k], option) == false {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func (jArray *JsonArray) Compare(other JsonObject, option JsonObjectCompareOption) bool {
	if other == nil {
		return false
	}

	if oArray, ok := other.(*JsonArray); ok {
		if len(oArray.inner) != len(jArray.inner) {
			return false
		}

		coArray := make([]JsonObject, len(oArray.inner))
		copy(coArray, oArray.inner)

		for _, sObject := range jArray.inner {
			found := false
			for j, oObject := range coArray {
				if sObject.Compare(oObject, option) {
					coArray = append(coArray[:j], coArray[j+1:]...)
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func ParseJson(source string) (JsonObject, error) {
	err := make([]byte, 1024)
	tree := C.yajl_tree_parse(C.CString(source), C.CString(string(err)), 1024)
	defer C.yajl_tree_free(tree)
	if tree == nil {
		return nil, errors.New(string(err))
	}
	return walkCTree(tree), nil
}

func walkCTree(tree C.yajl_val) JsonObject {
	if tree._type == C.yajl_t_object {
		jm := &JsonMap{}
		jm.inner = make(map[string]JsonObject)
		for i := C.size_t(0); i < C.yajl_val_get_object_len(tree); i++ {
			key := C.GoString(C.yajl_val_get_object_key(tree, i))
			jm.inner[key] = walkCTree(C.yajl_val_get_object_value(tree, i))
		}
		return jm
	} else if tree._type == C.yajl_t_array {
		ja := &JsonArray{}
		ja.inner = make([]JsonObject, 0)
		for i := C.size_t(0); i < C.yajl_val_get_array_len(tree); i++ {
			ja.inner = append(ja.inner, walkCTree(C.yajl_val_get_array_value(tree, i)))
		}
		return ja
	} else if tree._type == C.yajl_t_string {
		v := C.GoString(C.yajl_val_get_string(tree))
		//surroud quto around a string
		return &JsonString{C.yajl_t_string, v}
	} else if tree._type == C.yajl_t_number {
		v := C.GoString(C.yajl_val_get_number(tree))
		return &JsonString{C.yajl_t_number, v}
	} else if tree._type == C.yajl_t_true {
		return &JsonString{C.yajl_t_true, "true"}
	} else if tree._type == C.yajl_t_false {
		return &JsonString{C.yajl_t_false, "false"}
	}
	//ensure at least return JsonNull
	return &JsonString{C.yajl_t_null, "null"}
}
