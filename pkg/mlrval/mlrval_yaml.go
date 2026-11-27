// YAML decode/encode for Mlrval and Mlrmap.
// Converts between YAML native types (from gopkg.in/yaml.v3) and Miller's
// record model. YAML maps become Mlrmap; keys are stringified (YAML allows
// non-string keys). Used by the YAML record reader and writer.

package mlrval

import (
	"fmt"
	"io"
	"sort"

	"go.yaml.in/yaml/v3"
)

// MlrvalDecodeFromYAML decodes one YAML document from the decoder into an
// *Mlrval. Returns (nil, true, nil) on EOF. The decoded value can be a map
// (one record), array (array of records when elements are maps), or scalar.
//
// The document is decoded into a yaml.Node as well as into native Go values,
// since Go maps do not keep the key order of the input: the node is used to put
// map keys in the order they appear in the document.
func MlrvalDecodeFromYAML(decoder *yaml.Decoder) (*Mlrval, bool, error) {
	var node yaml.Node
	err := decoder.Decode(&node)
	if err == io.EOF {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, err
	}
	var doc interface{}
	err = node.Decode(&doc)
	if err != nil {
		return nil, false, err
	}
	if doc == nil {
		return NULL, false, nil
	}
	mv, err := mlrvalFromYAMLNativeOrdered(doc, &node)
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

// mlrvalFromYAMLNativeOrdered is like mlrvalFromYAMLNative, but takes map key
// order from the corresponding YAML node rather than sorting the keys.
func mlrvalFromYAMLNativeOrdered(v interface{}, node *yaml.Node) (*Mlrval, error) {
	node = resolveYAMLNode(node)
	if node == nil {
		return mlrvalFromYAMLNative(v)
	}
	switch val := v.(type) {
	case map[string]interface{}:
		if node.Kind == yaml.MappingNode {
			return mlrvalFromYAMLMapOrdered(val, node)
		}
	case map[interface{}]interface{}:
		if node.Kind == yaml.MappingNode {
			lookup := make(map[string]interface{}, len(val))
			for k, kv := range val {
				lookup[yamlKeyString(k)] = kv
			}
			return mlrvalFromYAMLMapOrdered(lookup, node)
		}
	case []interface{}:
		if node.Kind == yaml.SequenceNode && len(node.Content) == len(val) {
			out := FromEmptyArray()
			for i, elem := range val {
				mv, err := mlrvalFromYAMLNativeOrdered(elem, node.Content[i])
				if err != nil {
					return nil, err
				}
				out.ArrayAppend(mv)
			}
			return out, nil
		}
	}
	return mlrvalFromYAMLNative(v)
}

// mlrvalFromYAMLMapOrdered builds a map from the decoded values in m, with keys
// in the order of the mapping node. Any keys not found in the node (which
// should not happen) are appended in sorted order.
func mlrvalFromYAMLMapOrdered(m map[string]interface{}, node *yaml.Node) (*Mlrval, error) {
	out := FromEmptyMap()
	seen := make(map[string]bool, len(m))
	for _, pair := range yamlMappingPairs(node) {
		v, ok := m[pair.key]
		if !ok || seen[pair.key] {
			continue
		}
		seen[pair.key] = true
		valMv, err := mlrvalFromYAMLNativeOrdered(v, pair.value)
		if err != nil {
			return nil, err
		}
		out.MapPut(FromString(pair.key), valMv)
	}
	if len(seen) < len(m) {
		keys := make([]string, 0, len(m)-len(seen))
		for k := range m {
			if !seen[k] {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			valMv, err := mlrvalFromYAMLNative(m[k])
			if err != nil {
				return nil, err
			}
			out.MapPut(FromString(k), valMv)
		}
	}
	return out, nil
}

type yamlKeyValueNodes struct {
	key   string
	value *yaml.Node
}

// yamlMappingPairs returns the keys of a mapping node in document order, with
// the node holding each key's value. Merge keys ("<<") are expanded in place;
// as in yaml.v3, explicit keys take precedence over merged ones, and earlier
// merged mappings take precedence over later ones.
func yamlMappingPairs(node *yaml.Node) []yamlKeyValueNodes {
	explicit := make(map[string]*yaml.Node)
	for i := 0; i+1 < len(node.Content); i += 2 {
		if !isYAMLMergeKey(node.Content[i]) {
			explicit[yamlNodeKeyString(node.Content[i])] = node.Content[i+1]
		}
	}

	pairs := make([]yamlKeyValueNodes, 0, len(node.Content)/2)
	seen := make(map[string]bool)
	add := func(key string, value *yaml.Node) {
		if !seen[key] {
			seen[key] = true
			pairs = append(pairs, yamlKeyValueNodes{key, value})
		}
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if !isYAMLMergeKey(keyNode) {
			key := yamlNodeKeyString(keyNode)
			add(key, explicit[key])
			continue
		}
		merged := resolveYAMLNode(node.Content[i+1])
		var sources []*yaml.Node
		if merged != nil && merged.Kind == yaml.SequenceNode {
			sources = merged.Content
		} else {
			sources = []*yaml.Node{merged}
		}
		for _, source := range sources {
			source = resolveYAMLNode(source)
			if source == nil || source.Kind != yaml.MappingNode {
				continue
			}
			for _, pair := range yamlMappingPairs(source) {
				if _, isExplicit := explicit[pair.key]; !isExplicit {
					add(pair.key, pair.value)
				}
			}
		}
	}
	return pairs
}

func isYAMLMergeKey(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.Value == "<<" && node.ShortTag() == "!!merge"
}

// yamlNodeKeyString stringifies a mapping-key node the same way yamlKeyString
// does for decoded keys.
func yamlNodeKeyString(node *yaml.Node) string {
	var k interface{}
	if err := node.Decode(&k); err != nil {
		return node.Value
	}
	return yamlKeyString(k)
}

// resolveYAMLNode skips over document and alias nodes.
func resolveYAMLNode(node *yaml.Node) *yaml.Node {
	for node != nil {
		switch {
		case node.Kind == yaml.DocumentNode && len(node.Content) == 1:
			node = node.Content[0]
		case node.Kind == yaml.AliasNode:
			node = node.Alias
		default:
			return node
		}
	}
	return nil
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

// MlrmapToYAMLNative converts an Mlrmap to a *yaml.Node that preserves
// key insertion order. The returned node is suitable for yaml.Marshal.
func MlrmapToYAMLNative(mlrmap *Mlrmap) (*yaml.Node, error) {
	if mlrmap == nil {
		return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, nil
	}
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	for pe := mlrmap.Head; pe != nil; pe = pe.Next {
		keyNode, err := encodeScalarNode(pe.Key)
		if err != nil {
			return nil, err
		}
		valNode, err := mlrvalToYAMLNode(pe.Value)
		if err != nil {
			return nil, err
		}
		node.Content = append(node.Content, keyNode, valNode)
	}
	return node, nil
}

// encodeScalarNode creates a yaml.Node by delegating to yaml.v3's Encode,
// which handles edge cases like NaN/Inf floats and non-UTF-8 strings.
func encodeScalarNode(v interface{}) (*yaml.Node, error) {
	node := &yaml.Node{}
	if err := node.Encode(v); err != nil {
		return nil, err
	}
	return node, nil
}

// mlrvalToYAMLNode converts *Mlrval to a *yaml.Node for yaml.Marshal.
func mlrvalToYAMLNode(mv *Mlrval) (*yaml.Node, error) {
	if mv == nil {
		return encodeScalarNode(nil)
	}
	switch mv.Type() {
	case MT_ABSENT, MT_NULL:
		return encodeScalarNode(nil)
	case MT_VOID:
		return encodeScalarNode("")
	case MT_STRING:
		s, _ := mv.GetStringValue()
		return encodeScalarNode(s)
	case MT_INT:
		i, _ := mv.GetIntValue()
		return encodeScalarNode(i)
	case MT_FLOAT:
		f, _ := mv.GetFloatValue()
		return encodeScalarNode(f)
	case MT_BOOL:
		b, _ := mv.GetBoolValue()
		return encodeScalarNode(b)
	case MT_ARRAY:
		arr := mv.GetArray()
		seqNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, elem := range arr {
			v, err := mlrvalToYAMLNode(elem)
			if err != nil {
				return nil, err
			}
			seqNode.Content = append(seqNode.Content, v)
		}
		return seqNode, nil
	case MT_MAP:
		m := mv.GetMap()
		return MlrmapToYAMLNative(m)
	case MT_BYTES:
		return encodeScalarNode(mv.String())
	case MT_ERROR, MT_PENDING:
		return encodeScalarNode(mv.String())
	default:
		return encodeScalarNode(mv.String())
	}
}
