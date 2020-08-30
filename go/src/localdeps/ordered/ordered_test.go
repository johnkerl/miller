// ================================================================
// FOUND AT https://gitlab.com/c0b/go-ordered-json
// ================================================================

package ordered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestMarshalOrderedMap(t *testing.T) {
	om := NewOrderedMap()
	om.Set("a", 34)
	om.Set("b", []int{3, 4, 5})
	b, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("Marshal OrderedMap: %v", err)
	}
	// fmt.Printf("%q\n", b)
	const expected = "{\"a\":34,\"b\":[3,4,5]}"
	if !bytes.Equal(b, []byte(expected)) {
		t.Errorf("Marshal OrderedMap: %q not equal to expected %q", b, expected)
	}
}

func ExampleOrderedMap_UnmarshalJSON() {
	const jsonStream = `{
  "country"     : "United States",
  "countryCode" : "US",
  "region"      : "CA",
  "regionName"  : "California",
  "city"        : "Mountain View",
  "zip"         : "94043",
  "lat"         : 37.4192,
  "lon"         : -122.0574,
  "timezone"    : "America/Los_Angeles",
  "isp"         : "Google Cloud",
  "org"         : "Google Cloud",
  "as"          : "AS15169 Google Inc.",
  "mobile"      : true,
  "proxy"       : false,
  "query"       : "35.192.xx.xxx"
}`

	// compare with if using a regular generic map, the unmarshalled result
	//  is a map with unpredictable order of keys
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStream), &m)
	if err != nil {
		fmt.Println("error:", err)
	}
	for key := range m {
		// fmt.Printf("%-12s: %v\n", key, m[key])
		_ = key
	}

	// use the OrderedMap to Unmarshal from JSON object
	var om *OrderedMap = NewOrderedMap()
	err = json.Unmarshal([]byte(jsonStream), om)
	if err != nil {
		fmt.Println("error:", err)
	}

	// use an iterator func to loop over all key-value pairs,
	// it is ok to call Set append-modify new key-value pairs,
	// but not safe to call Delete during iteration.
	iter := om.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		fmt.Printf("%-12s: %v\n", pair.Key, pair.Value)
		if pair.Key == "city" {
			om.Set("mobile", false)
			om.Set("extra", 42)
		}
	}

	// Output:
	// country     : United States
	// countryCode : US
	// region      : CA
	// regionName  : California
	// city        : Mountain View
	// zip         : 94043
	// lat         : 37.4192
	// lon         : -122.0574
	// timezone    : America/Los_Angeles
	// isp         : Google Cloud
	// org         : Google Cloud
	// as          : AS15169 Google Inc.
	// mobile      : false
	// proxy       : false
	// query       : 35.192.xx.xxx
	// extra       : 42
}

func TestUnmarshalOrderedMapFromInvalid(t *testing.T) {
	om := NewOrderedMap()

	om.Set("m", math.NaN())
	b, err := json.Marshal(om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", b, err)
	}
	// fmt.Println(om, b, err)
	om.Delete("m")

	err = json.Unmarshal([]byte("[]"), om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error")
	}

	err = json.Unmarshal([]byte("["), om)
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(nil))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{}3"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte("{]"))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(`{"a": 3, "b": [{`))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}

	err = om.UnmarshalJSON([]byte(`{"a": 3, "b": [}`))
	if err == nil {
		t.Fatal("Unmarshal OrderedMap: expecting error:", om)
	}
	// fmt.Println("error:", om, err)
}

func TestUnmarshalOrderedMap(t *testing.T) {
	var (
		data  = []byte(`{"as":"AS15169 Google Inc.","city":"Mountain View","country":"United States","countryCode":"US","isp":"Google Cloud","lat":37.4192,"lon":-122.0574,"org":"Google Cloud","query":"35.192.25.53","region":"CA","regionName":"California","status":"success","timezone":"America/Los_Angeles","zip":"94043"}`)
		pairs = []*KVPair{
			{"as", "AS15169 Google Inc."},
			{"city", "Mountain View"},
			{"country", "United States"},
			{"countryCode", "US"},
			{"isp", "Google Cloud"},
			{"lat", 37.4192},
			{"lon", -122.0574},
			{"org", "Google Cloud"},
			{"query", "35.192.25.53"},
			{"region", "CA"},
			{"regionName", "California"},
			{"status", "success"},
			{"timezone", "America/Los_Angeles"},
			{"zip", "94043"},
		}
		obj = NewOrderedMapFromKVPairs(pairs)
	)

	om := NewOrderedMap()
	err := json.Unmarshal(data, om)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}

	// fix number type for deepequal test
	for _, key := range []string{"lat", "lon"} {
		numf, _ := om.Get(key).(json.Number).Float64()
		om.Set(key, numf)
	}

	// check by Has and GetValue
	for _, kv := range pairs {
		if !om.Has(kv.Key) {
			t.Fatalf("expect key %q exists in Unmarshaled OrderedMap")
		}
		value, ok := om.GetValue(kv.Key)
		if !ok || value != kv.Value {
			t.Fatalf("expect for key %q: the value %v should equal to %v, in Unmarshaled OrderedMap", kv.Key, value, kv.Value)
		}
	}

	b, err := json.MarshalIndent(om, "", "  ")
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}
	const expected = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "org": "Google Cloud",
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043"
}`
	if !bytes.Equal(b, []byte(expected)) {
		t.Fatalf("Unmarshal OrderedMap marshal indent from %#v not equal to expected: %q\n", om, expected)
	}

	if !reflect.DeepEqual(om, obj) {
		t.Fatalf("Unmarshal OrderedMap not deeply equal: %#v %#v", om, obj)
	}

	val, ok := om.Delete("org")
	if !ok {
		t.Fatalf("org should exist")
	}
	om.Set("org", val)
	b, err = json.MarshalIndent(om, "", "  ")
	// fmt.Println("after delete", om, string(b), err)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}
	const expected2 = `{
  "as": "AS15169 Google Inc.",
  "city": "Mountain View",
  "country": "United States",
  "countryCode": "US",
  "isp": "Google Cloud",
  "lat": 37.4192,
  "lon": -122.0574,
  "query": "35.192.25.53",
  "region": "CA",
  "regionName": "California",
  "status": "success",
  "timezone": "America/Los_Angeles",
  "zip": "94043",
  "org": "Google Cloud"
}`
	if !bytes.Equal(b, []byte(expected2)) {
		t.Fatalf("Unmarshal OrderedMap marshal indent from %#v not equal to expected: %s\n", om, expected2)
	}
}

func TestUnmarshalNestedOrderedMap(t *testing.T) {
	var (
		data = []byte(`{"a": true, "b": [3, 4, { "b": "3", "d": [] }]}`)
		obj  = NewOrderedMapFromKVPairs([]*KVPair{
			{"a", true},
			{"b", []interface{}{3, 4, NewOrderedMapFromKVPairs([]*KVPair{
				{"b", "3"},
				{"d", []interface{}{}},
			})}},
		})
	)

	om := NewOrderedMap()
	err := json.Unmarshal(data, om)
	if err != nil {
		t.Fatalf("Unmarshal OrderedMap: %v", err)
	}

	// b, err := json.MarshalIndent(om, "", "  ")
	// fmt.Println(om, string(b), err, obj)

	// fix number type for deepequal test
	elearr := om.Get("b").([]interface{})
	for i, v := range elearr {
		if num, ok := v.(json.Number); ok {
			numi, _ := num.Int64()
			elearr[i] = int(numi)
		}
	}

	if !reflect.DeepEqual(om, obj) {
		t.Fatalf("Unmarshal OrderedMap not deeply equal: %#v expected %#v", om, obj)
	}
}

func ExampleOrderedMap_EntriesReverseIter() {
	// initialize from a list of key-value pairs
	om := NewOrderedMapFromKVPairs([]*KVPair{
		{"country", "United States"},
		{"countryCode", "US"},
		{"region", "CA"},
		{"regionName", "California"},
		{"city", "Mountain View"},
		{"zip", "94043"},
		{"lat", 37.4192},
		{"lon", -122.0574},
		{"timezone", "America/Los_Angeles"},
		{"isp", "Google Cloud"},
		{"org", "Google Cloud"},
		{"as", "AS15169 Google Inc."},
		{"mobile", true},
		{"proxy", false},
		{"query", "35.192.xx.xxx"},
	})

	iter := om.EntriesReverseIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		fmt.Printf("%-12s: %v\n", pair.Key, pair.Value)
	}

	// Output:
	// query       : 35.192.xx.xxx
	// proxy       : false
	// mobile      : true
	// as          : AS15169 Google Inc.
	// org         : Google Cloud
	// isp         : Google Cloud
	// timezone    : America/Los_Angeles
	// lon         : -122.0574
	// lat         : 37.4192
	// zip         : 94043
	// city        : Mountain View
	// regionName  : California
	// region      : CA
	// countryCode : US
	// country     : United States
}

var unmarshalTests = []struct {
	in        string
	new       func() interface{}
	out       interface{}
	err       error
	useNumber bool
	golden    bool
}{
	{in: `true`, new: func() interface{} { return new(bool) }, out: true},
	{in: `1`, new: func() interface{} { return new(int) }, out: 1},
	{in: `1.2`, new: func() interface{} { return new(float64) }, out: 1.2},
	{in: `-5`, new: func() interface{} { return new(int16) }, out: int16(-5)},
	{in: `2`, new: func() interface{} { return new(json.Number) }, out: json.Number("2"), useNumber: true},
	{in: `2`, new: func() interface{} { return new(json.Number) }, out: json.Number("2")},
	{in: `2`, new: func() interface{} { return new(interface{}) }, out: float64(2.0)},
	{in: `2`, new: func() interface{} { return new(interface{}) }, out: json.Number("2"), useNumber: true},
	{in: `"a\u1234"`, new: func() interface{} { return new(string) }, out: "a\u1234"},
	{in: `"http:\/\/"`, new: func() interface{} { return new(string) }, out: "http://"},
	{in: `"g-clef: \uD834\uDD1E"`, new: func() interface{} { return new(string) }, out: "g-clef: \U0001D11E"},
	{in: `"invalid: \uD834x\uDD1E"`, new: func() interface{} { return new(string) }, out: "invalid: \uFFFDx\uFFFD"},
	{in: "null", new: func() interface{} { return new(interface{}) }, out: nil},
	{in: "{}", new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs([]*KVPair{})},
	{in: `{"a": 3}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}})},
	{in: `{"a": 3, "b": true}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}, {"b", true}})},
	{in: `{"a": 3, "b": true, "c": null}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}, {"b", true}, {"c", nil}})},
	{in: `{"a": 3, "c": null, "d": []}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}, {"c", nil}, {"d", []interface{}{}}})},
	{in: `{"a": 3, "c": null, "d": [3,4,true]}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}, {"c", nil}, {"d", []interface{}{
			json.Number("3"), json.Number("4"), true,
		}}})},
	{in: `{"a": 3, "c": null, "d": [3,4,true, { "inner": "abc" }]}`, new: func() interface{} { return NewOrderedMap() }, out: *NewOrderedMapFromKVPairs(
		[]*KVPair{{"a", json.Number("3")}, {"c", nil}, {"d", []interface{}{
			json.Number("3"), json.Number("4"), true, NewOrderedMapFromKVPairs([]*KVPair{{"inner", "abc"}}),
		}}})},
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range unmarshalTests {
		in := []byte(tt.in)
		if tt.new == nil {
			continue
		}

		// v = new(right-type)
		v := tt.new() // reflect.New(reflect.TypeOf(tt.ptr).Elem())
		dec := json.NewDecoder(bytes.NewReader(in))
		if tt.useNumber {
			dec.UseNumber()
		}
		if err := dec.Decode(v); !reflect.DeepEqual(err, tt.err) {
			t.Errorf("#%d: %v, want %v", i, err, tt.err)
			continue
		} else if err != nil {
			continue
		}
		if !reflect.DeepEqual(reflect.ValueOf(v).Elem().Interface(), tt.out) {
			t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, tt.out)
			data, _ := json.Marshal(v)
			println(string(data))
			data, _ = json.Marshal(tt.out)
			println(string(data))
			continue
		}

		// Check round trip also decodes correctly.
		if tt.err == nil {
			enc, err := json.Marshal(v)
			if err != nil {
				t.Errorf("#%d: error re-marshaling: %v", i, err)
				continue
			}
			if tt.golden && !bytes.Equal(enc, in) {
				t.Errorf("#%d: remarshal mismatch:\nhave: %s\nwant: %s", i, enc, in)
			}
			vv := tt.new() // reflect.New(reflect.TypeOf(tt.ptr).Elem())
			dec = json.NewDecoder(bytes.NewReader(enc))
			if tt.useNumber {
				dec.UseNumber()
			}
			if err := dec.Decode(vv); err != nil {
				t.Errorf("#%d: error re-unmarshaling %#q: %v", i, enc, err)
				continue
			}
			if !reflect.DeepEqual(v, vv) {
				t.Errorf("#%d: mismatch\nhave: %#+v\nwant: %#+v", i, v, vv)
				t.Errorf("     In: %q", strings.Map(noSpace, string(in)))
				t.Errorf("Marshal: %q", strings.Map(noSpace, string(enc)))
				continue
			}
		}
	}
}

func noSpace(c rune) rune {
	if isSpace(byte(c)) { //only used for ascii
		return -1
	}
	return c
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
