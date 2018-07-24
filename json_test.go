package yajl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ = fmt.Println

func TestParse(t *testing.T) {
	var err error
	_, err = ParseJson(`{"a":"b"}`)
	assert.Nil(t, err)

	_, err = ParseJson(`{"a":"b",}`)
	assert.NotNil(t, err)
}

func TestCompare(t *testing.T) {
	j1, _ := ParseJson(`{"a":"b"}`)
	j2, _ := ParseJson(`{"a":"b"}`)
	assert.True(t, j1.Compare(j2, DEFAULT))

	j3, _ := ParseJson(`{"a":"c"}`)
	j4, _ := ParseJson(`{"a":"b"}`)
	assert.False(t, j3.Compare(j4, DEFAULT))

	j5, _ := ParseJson(`{"namelist":[{"name":"jack"},{"name":"rose"}]}`)
	j6, _ := ParseJson(`{"namelist":[{"name":"rose"},{"name":"jack"}]}`)
	assert.True(t, j5.Compare(j6, DEFAULT))

	j7, _ := ParseJson(`{"namelist":[{"name":"jack"},{"name":"ross"}]}`)
	j8, _ := ParseJson(`{"namelist":[{"name":"rose"},{"name":"jack"}]}`)
	assert.False(t, j7.Compare(j8, DEFAULT))

	j9, _ := ParseJson(`{"boolean":"true"}`)
	j10, _ := ParseJson(`{"boolean":true}`)
	assert.True(t, j9.Compare(j10, ALL_VALUE_TREAT_AS_STRING))

	j11, _ := ParseJson(`{"boolean":"true"}`)
	j12, _ := ParseJson(`{"boolean":true}`)
	assert.False(t, j11.Compare(j12, DEFAULT))

	j13, _ := ParseJson(`{"nil":null}`)
	j14, _ := ParseJson(`{"nil":"null"}`)
	assert.True(t, j13.Compare(j14, ALL_VALUE_TREAT_AS_STRING))

	j15, _ := ParseJson(`{"nil":null}`)
	j16, _ := ParseJson(`{"nil":"null"}`)
	assert.False(t, j15.Compare(j16, DEFAULT))

	j17, _ := ParseJson(`[{"nil":null}]`)
	j18, _ := ParseJson(`[{"nil":"null"}]`)
	assert.True(t, j17.Compare(j18, ALL_VALUE_TREAT_AS_STRING))

	j19, _ := ParseJson(`[{"nil":null},{"nil":null}]`)
	j20, _ := ParseJson(`[{"nil":"null"}]`)
	assert.False(t, j19.Compare(j20, ALL_VALUE_TREAT_AS_STRING))

	j21, _ := ParseJson(`[{"nil":null}]`)
	j22, _ := ParseJson(`[{"nil":"null"},{"diff":null}]`)
	assert.False(t, j21.Compare(j22, ALL_VALUE_TREAT_AS_STRING))
}
