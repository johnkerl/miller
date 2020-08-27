// Package ordered provided a type OrderedMap for use in JSON handling
// although JSON spec says the keys order of an object should not matter
// but sometimes when working with particular third-party proprietary code
// which has incorrect using the keys order, we have to maintain the object keys
// in the same order of incoming JSON object, this package is useful for these cases.
//
// Disclaimer:
// same as Go's default [map](https://blog.golang.org/go-maps-in-action),
// this OrderedMap is not safe for concurrent use, if need atomic access, may use a sync.Mutex to synchronize.
package ordered

// Refers
//  JSON and Go        https://blog.golang.org/json-and-go
//  Go-Ordered-JSON    https://github.com/virtuald/go-ordered-json
//  Python OrderedDict https://github.com/python/cpython/blob/2.7/Lib/collections.py#L38
//  port OrderedDict   https://github.com/cevaris/ordered_map

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
)

// the key-value pair type, for initializing from a list of key-value pairs, or for looping entries in the same order
type KVPair struct {
	Key   string
	Value interface{}
}

type m map[string]interface{}

// the OrderedMap type, has similar operations as the default map, but maintained
// the keys order of inserted; similar to map, all single key operations (Get/Set/Delete) runs at O(1).
type OrderedMap struct {
	m
	l    *list.List
	keys map[string]*list.Element // the double linked list for delete and lookup to be O(1)
}

// Create a new OrderedMap
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		m:    make(map[string]interface{}),
		l:    list.New(),
		keys: make(map[string]*list.Element),
	}
}

// Create a new OrderedMap and populate from a list of key-value pairs
func NewOrderedMapFromKVPairs(pairs []*KVPair) *OrderedMap {
	om := NewOrderedMap()
	for _, pair := range pairs {
		om.Set(pair.Key, pair.Value)
	}
	return om
}

// return all keys
// func (om *OrderedMap) Keys() []string { return om.keys }

// set value for particular key, this will remember the order of keys inserted
// but if the key already exists, the order is not updated.
func (om *OrderedMap) Set(key string, value interface{}) {
	if _, ok := om.m[key]; !ok {
		om.keys[key] = om.l.PushBack(key)
	}
	om.m[key] = value
}

// Check if value exists
func (om *OrderedMap) Has(key string) bool {
	_, ok := om.m[key]
	return ok
}

// Get value for particular key, or nil if not exist; but don't rely on nil for non-exist; should check by Has or GetValue
func (om *OrderedMap) Get(key string) interface{} {
	return om.m[key]
}

// Get value and exists together
func (om *OrderedMap) GetValue(key string) (value interface{}, ok bool) {
	value, ok = om.m[key]
	return
}

// deletes the element with the specified key (m[key]) from the map. If there is no such element, this is a no-op.
func (om *OrderedMap) Delete(key string) (value interface{}, ok bool) {
	value, ok = om.m[key]
	if ok {
		om.l.Remove(om.keys[key])
		delete(om.keys, key)
		delete(om.m, key)
	}
	return
}

// Iterate all key/value pairs in the same order of object constructed
func (om *OrderedMap) EntriesIter() func() (*KVPair, bool) {
	e := om.l.Front()
	return func() (*KVPair, bool) {
		if e != nil {
			key := e.Value.(string)
			e = e.Next()
			return &KVPair{key, om.m[key]}, true
		}
		return nil, false
	}
}

// Iterate all key/value pairs in the reverse order of object constructed
func (om *OrderedMap) EntriesReverseIter() func() (*KVPair, bool) {
	e := om.l.Back()
	return func() (*KVPair, bool) {
		if e != nil {
			key := e.Value.(string)
			e = e.Prev()
			return &KVPair{key, om.m[key]}, true
		}
		return nil, false
	}
}

// this implements type json.Marshaler interface, so can be called in json.Marshal(om)
func (om *OrderedMap) MarshalJSON() (res []byte, err error) {
	res = append(res, '{')
	front, back := om.l.Front(), om.l.Back()
	for e := front; e != nil; e = e.Next() {
		k := e.Value.(string)
		res = append(res, fmt.Sprintf("%q:", k)...)
		var b []byte
		b, err = json.Marshal(om.m[k])
		if err != nil {
			return
		}
		res = append(res, b...)
		if e != back {
			res = append(res, ',')
		}
	}
	res = append(res, '}')
	// fmt.Printf("marshalled: %v: %#v\n", res, res)
	return
}

// this implements type json.Unmarshaler interface, so can be called in json.Unmarshal(data, om)
func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	err = om.parseobject(dec)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("expect end of JSON object but got more token: %T: %v or err: %v", t, t, err)
	}

	return nil
}

func (om *OrderedMap) parseobject(dec *json.Decoder) (err error) {
	var t json.Token
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		var value interface{}
		value, err = handledelim(t, dec)
		if err != nil {
			return err
		}

		// om.keys = append(om.keys, key)
		om.keys[key] = om.l.PushBack(key)
		om.m[key] = value
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func parsearray(dec *json.Decoder) (arr []interface{}, err error) {
	var t json.Token
	arr = make([]interface{}, 0)
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return
		}

		var value interface{}
		value, err = handledelim(t, dec)
		if err != nil {
			return
		}
		arr = append(arr, value)
	}
	t, err = dec.Token()
	if err != nil {
		return
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		err = fmt.Errorf("expect JSON array close with ']'")
		return
	}

	return
}

func handledelim(t json.Token, dec *json.Decoder) (res interface{}, err error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			om2 := NewOrderedMap()
			err = om2.parseobject(dec)
			if err != nil {
				return
			}
			return om2, nil
		case '[':
			var value []interface{}
			value, err = parsearray(dec)
			if err != nil {
				return
			}
			return value, nil
		default:
			return nil, fmt.Errorf("Unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}
