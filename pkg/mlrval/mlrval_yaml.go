// YAML decode/encode for Mlrval and Mlrmap.
// Converts between YAML native types (from gopkg.in/yaml.v3) and Miller's
// record model. YAML maps become Mlrmap; keys are stringified (YAML allows
// non-string keys). Used by the YAML record reader and writer.

package mlrval

import (
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"
)

// MlrvalDecodeFromYAML decodes one YAML document from the decoder into an
// *Mlrval. Returns (nil, true, nil) on EOF. The decoded value can be a map
// (one record), array (array of records when elements are maps), or scalar.
func MlrvalDecodeFromYAML(decoder *yaml.Decoder) (*Mlrval, bool, error) {
	var doc interface{}
	err := decoder.Decode(&doc)
	if err == io.EOF {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, err
	}
	if doc == nil {
		return NULL, false, nil
	}
	mv, err := mlrvalFromYAMLNative(doc)
	if err != nil {
		return nil, false, err
	}
	return mv, false, nil
}

// mlrvalFromYAMLNative converts a YAML-decoded value (map[interface{}]interface{},
// []interface{}, or scalar) into *Mlrval.
func mlrvalFromYAMLNative(v interface{}) (*Mlrval, error) {
	if v == nil {
		return NULL, nil
	}
	switch val := v.(type) {
	case map[interface{}]interface{}:
		return mlrvalFromYAMLMap(val)
	case map[string]interface{}:
		return mlrvalFromYAMLStringMap(val)
	case []interface{}:
		return mlrvalFromYAMLArray(val)
	case string:
		return FromString(val), nil
	case bool:
		return FromBool(val), nil
	case int:
		return FromInt(int64(val)), nil
	case int64:
		return FromInt(val), nil
	case uint64:
		if val <= 1<<63-1 {
			return FromInt(int64(val)), nil
		}
		return FromFloat(float64(val)), nil
	case float64:
		return FromFloat(val), nil
	case float32:
		return FromFloat(float64(val)), nil
	default:
		return FromString(fmt.Sprint(val)), nil
	}
}

func mlrvalFromYAMLMap(m map[interface{}]interface{}) (*Mlrval, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, yamlKeyString(k))
	}
	sort.Strings(keys)
	out := FromEmptyMap()
	for _, keyStr := range keys {
		var v interface{}
		for k, val := range m {
			if yamlKeyString(k) == keyStr {
				v = val
				break
			}
		}
		valMv, err := mlrvalFromYAMLNative(v)
		if err != nil {
			return nil, err
		}
		out.MapPut(FromString(keyStr), valMv)
	}
	return out, nil
}

func mlrvalFromYAMLStringMap(m map[string]interface{}) (*Mlrval, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := FromEmptyMap()
	for _, k := range keys {
		valMv, err := mlrvalFromYAMLNative(m[k])
		if err != nil {
			return nil, err
		}
		out.MapPut(FromString(k), valMv)
	}
	return out, nil
}

func yamlKeyString(k interface{}) string {
	switch t := k.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case float64:
		return fmt.Sprintf("%g", t)
	default:
		return fmt.Sprint(k)
	}
}

func mlrvalFromYAMLArray(a []interface{}) (*Mlrval, error) {
	out := FromEmptyArray()
	for _, elem := range a {
		mv, err := mlrvalFromYAMLNative(elem)
		if err != nil {
			return nil, err
		}
		out.ArrayAppend(mv)
	}
	return out, nil
}

// MlrmapToYAMLNative converts an Mlrmap to a Go value suitable for
// yaml.Marshal: map[string]interface{} with nested maps/arrays.
func MlrmapToYAMLNative(mlrmap *Mlrmap) (interface{}, error) {
	if mlrmap == nil {
		return nil, nil
	}
	out := make(map[string]interface{})
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		v, err := mlrvalToYAMLNative(pe.Value)
		if err != nil {
			return nil, err
		}
		out[pe.Key] = v
	}
	return out, nil
}

// mlrvalToYAMLNative converts *Mlrval to a Go value for yaml.Marshal.
func mlrvalToYAMLNative(mv *Mlrval) (interface{}, error) {
	if mv == nil {
		return nil, nil
	}
	switch mv.Type() {
	case MT_ABSENT, MT_VOID, MT_NULL:
		return nil, nil
	case MT_STRING:
		s, _ := mv.GetStringValue()
		return s, nil
	case MT_INT:
		i, _ := mv.GetIntValue()
		return i, nil
	case MT_FLOAT:
		f, _ := mv.GetFloatValue()
		return f, nil
	case MT_BOOL:
		b, _ := mv.GetBoolValue()
		return b, nil
	case MT_ARRAY:
		arr := mv.GetArray()
		if arr == nil {
			return []interface{}{}, nil
		}
		out := make([]interface{}, 0, len(arr))
		for _, elem := range arr {
			v, err := mlrvalToYAMLNative(elem)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case MT_MAP:
		m := mv.GetMap()
		return MlrmapToYAMLNative(m)
	case MT_ERROR, MT_PENDING:
		return mv.String(), nil
	default:
		return mv.String(), nil
	}
}
