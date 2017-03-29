package rethinkdb

import (
	"encoding/json"
	"fmt"
	"strings"

	r "gopkg.in/dancannon/gorethink.v2"
)

// DB implemtation
type DB struct{}

// Update rethinkdb by mongodb selector and document
func (db *DB) Update(selector, document []byte) (int32, error) {
	d, err := parse(selector, document)
	if err != nil {
		return 0, fmt.Errorf("Parsing error invalid string %s:%s", string(selector), string(document))
	}
	affected, err := build(d)
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func parse(selector, document []byte) (map[string]map[string]interface{}, error) {
	d := make(map[string]map[string]interface{})
	var operators map[string]*json.RawMessage
	err := json.Unmarshal([]byte(document), &operators)
	if err != nil {
		return nil, err
	}
	for operator, raw := range operators {
		switch operator {
		case "$push":
			operatePush(d, raw)
		}
	}
	return d, nil
}

func operatePush(
	d map[string]map[string]interface{},
	raw *json.RawMessage,
) {
	var fields map[string]interface{}
	json.Unmarshal(*raw, &fields)
	for field, raw := range fields {
		switch raw := raw.(type) {
		case string:
			d[field] = make(map[string]interface{})
			d[field]["$default"] = raw
		case bool:
			d[field] = make(map[string]interface{})
			d[field]["$default"] = raw
		case float64:
			d[field] = make(map[string]interface{})
			d[field]["$default"] = raw
		case map[string]interface{}:
			d[field] = make(map[string]interface{})
			for k, v := range raw {
				if runes := []rune(k); runes[0] == '$' {
					d[field][k] = v
				} else {
					d[field]["$default"] = raw
				}
			}
		case []interface{}:
			d[field] = make(map[string]interface{})
			d[field]["$default"] = raw
		}
	}
}

func build(d map[string]map[string]interface{}) (int32, error) {
	session, _ := r.Connect(r.ConnectOpts{
		Address: "external-rethinkdb:28015",
	})
	selected := r.DB("2910eb12_d64a_49cc_b2be_54201441e27b").Table("test_mongoql")
	query := selected.Update(buildUpdateFunc(d))
	fmt.Println(query.String())
	result, err := query.Run(session)
	if err != nil {
		return 0, err
	}
	var a []map[string]int
	result.All(&a)
	return int32(a[0]["replaced"]), nil
}

func buildUpdateFunc(d map[string]map[string]interface{}) func(r.Term) interface{} {
	return func(t r.Term) interface{} {
		ret := make(map[string]interface{})
		for field, modifiers := range d {
			modifiedTerm := modifyTerm(t, modifiers, field)
			buildHierarchy(ret, field, modifiedTerm)
		}
		return ret
	}
}

func buildHierarchy(ret map[string]interface{}, field string, t r.Term) {
	splitedField := strings.Split(field, ".")
	ref := ret
	for i, f := range splitedField {
		if i == len(splitedField)-1 {
			ref[f] = t
		} else {
			ref[f] = make(map[string]interface{})
			ref = ref[f].(map[string]interface{})
		}
	}
}

func modifyTerm(
	t r.Term,
	modifiers map[string]interface{},
	field string,
) r.Term {
	splitedField := strings.Split(field, ".")
	for _, f := range splitedField {
		t = t.Field(f)
	}
	if v, ok := modifiers["$default"]; ok {
		return t.Append(v)
	} else if v, ok := modifiers["$each"]; ok {
		t = t.Add(v)
		if v, ok := modifiers["$sort"]; ok {
			modifierSort(&t, v)
		}
		if v, ok := modifiers["$slice"]; ok {
			modifierSlice(&t, v)
		}
		return t
	}
	return t
}

func modifierSort(t *r.Term, v interface{}) {
	if fieldsSorted, ok := v.(map[string]interface{}); ok {
		for field, power := range fieldsSorted {
			intpower := int(power.(float64))
			if intpower == -1 {
				*t = t.OrderBy(r.Desc(field))
			} else if intpower == 1 {
				*t = t.OrderBy(field)
			}
		}
	}
}

func modifierSlice(t *r.Term, v interface{}) {
	if n, ok := v.(float64); ok {
		if n >= 0 {
			*t = t.Slice(0, n)
		} else {
			*t = t.Slice(n, -1, r.SliceOpts{RightBound: "closed"})
		}
	}
}
