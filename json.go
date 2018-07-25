package yajl

// #include "api/yajl_tree.h"
// #include "shim-go.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
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
	Inner map[string]JsonObject
}

type JsonArray struct {
	Inner []JsonObject
}

type JsonString struct {
	RealType int
	Inner    string
}

func (jString *JsonString) Compare(other JsonObject, option JsonObjectCompareOption) bool {
	if other == nil {
		return false
	}

	if oString, ok := other.(*JsonString); ok {
		if option&ALL_VALUE_TREAT_AS_STRING > 0 {
			return oString.Inner == jString.Inner
		} else {
			if oString.RealType != jString.RealType {
				return false
			} else {
				return oString.Inner == jString.Inner
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
		if len(oMap.Inner) != len(jMap.Inner) {
			return false
		}

		for k, v := range jMap.Inner {
			if v.Compare(oMap.Inner[k], option) == false {
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
		if len(oArray.Inner) != len(jArray.Inner) {
			return false
		}

		coArray := make([]JsonObject, len(oArray.Inner))
		copy(coArray, oArray.Inner)

		for _, sObject := range jArray.Inner {
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
	csource := C.CString(source)
	defer C.free(unsafe.Pointer(csource))
	err := (*C.char)(C.malloc(1024))
	defer C.free(unsafe.Pointer(err))
	tree := C.yajl_tree_parse(csource, err, 1024)
	defer C.yajl_tree_free(tree)
	if tree == nil {
		return nil, errors.New(C.GoString(err))
	}
	return walkCTree(tree), nil
}

func walkCTree(tree C.yajl_val) JsonObject {
	if tree._type == C.yajl_t_object {
		jm := &JsonMap{}
		jm.Inner = make(map[string]JsonObject)
		for i := C.size_t(0); i < C.yajl_val_get_object_len(tree); i++ {
			key := C.GoString(C.yajl_val_get_object_key(tree, i))
			jm.Inner[key] = walkCTree(C.yajl_val_get_object_value(tree, i))
		}
		return jm
	} else if tree._type == C.yajl_t_array {
		ja := &JsonArray{}
		ja.Inner = make([]JsonObject, 0)
		for i := C.size_t(0); i < C.yajl_val_get_array_len(tree); i++ {
			ja.Inner = append(ja.Inner, walkCTree(C.yajl_val_get_array_value(tree, i)))
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
