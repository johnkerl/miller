# Proposal: YAML I/O format support

**Status: Implemented.** YAML I/O, conversion shortcuts, list-wrap flags, unit tests, and regression tests are in place.

Miller already handles a subset of JSON for input and output. This document proposed adding **YAML** as an additional I/O format, reusing the same record model (stream of `Mlrmap` records) and the existing `mlrval` type system.

## Goals

- **Input**: Read YAML files as a stream of records. Support:
  - Single document: one object → one record; one array of objects → one record per element.
  - Multiple documents (separated by `---`) → one or more records per document, same rules.
- **Output**: Write records as YAML. Support:
  - One YAML document per record (with `---` between), or
  - One top-level YAML document that is an array of objects (like `--json` with list wrap).
- **Consistency**: Same semantics as JSON where applicable (flatten/nested keys, type preservation), and same CLI pattern (`--yaml`, `--iyaml`, `--oyaml`, `-i yaml`, `-o yaml`).

## Non-goals (initial version)

- YAML-specific features (tags, anchors/aliases, custom types) can be left as best-effort (e.g. round-trip as strings or omit).
- No separate “YAML Lines” format; YAML’s multi-document and array forms are enough.

---

## 1. Dependencies

- **gopkg.in/yaml.v3** is already in `go.mod` as an indirect dependency. Make it a **direct** dependency and use it for decode/encode.

```go
// go.mod
require (
	// ... existing
	gopkg.in/yaml.v3 v3.0.1
)
```

Remove the `// indirect` comment so the dependency is explicit.

---

## 2. New / modified files

### 2.1 `pkg/mlrval/mlrval_yaml.go` (new)

- **Decode (YAML → mlrval)**  
  - Use `yaml.Unmarshal` or `yaml.NewDecoder` into `interface{}`.  
  - Implement a converter from generic YAML-decoded values to `*Mlrval`:
    - `map[interface{}]interface{}` → `Mlrmap` (keys must be stringifiable; same as JSON “object keys must be strings”).
    - `[]interface{}` → `Mlrval` array.
    - Scalars → `FromString`, `FromInt`, `FromFloat`, `FromBool`, `NULL` as appropriate.
  - Expose something like:
    - `MlrvalDecodeFromYAML(decoder *yaml.Decoder) (*Mlrval, bool, error)` (mirroring `MlrvalDecodeFromJSON` with eof).
  - YAML maps use `map[interface{}]interface{}`; key comparison must normalize to string (e.g. all keys converted to string for Miller’s record model).

- **Encode (mlrval → YAML)**  
  - Convert `*Mlrval` / `*Mlrmap` to a Go type that `yaml.Marshal` accepts:
    - `Mlrmap` → `map[string]interface{}` (values recursively converted).
    - Arrays → `[]interface{}`; scalars → `string`, `int64`, `float64`, `bool`, `nil`.
  - Add:
    - `(mv *Mlrval) MarshalYAML() (interface{}, error)` if using the `yaml.Marshaler` interface, **or**
    - Helper that returns `interface{}` for a single record and call `yaml.Marshal` in the writer.
  - No need to support `yaml.Unmarshaler` on `Mlrval` for now unless we want YAML→YAML round-trip of arbitrary nodes.

### 2.2 `pkg/input/record_reader_yaml.go` (new)

- Mirror structure of `record_reader_json.go`:
  - `RecordReaderYAML` with `Read(filenames, context, readerChannel, errorChannel, downstreamDoneChannel)`.
  - Per file: `yaml.NewDecoder(handle)`, then loop:
    - `decoder.Decode(&doc)` into `interface{}`.
    - If document is a map → convert to one `Mlrmap`, send one record.
    - If document is a slice → for each element that is a map, convert to `Mlrmap` and send; non-maps same as JSON (error or skip with clear message).
    - EOF from `Decode` → end of stream for that file.
  - Optionally support comment stripping (similar to `JSONCommentEnabledReader`) if we want `#` comments at line-start; can be a follow-up.
  - No need to set `context.JSONHadBrackets`; YAML writer can have its own option for “wrap in top-level array” (see below).

### 2.3 `pkg/output/record_writer_yaml.go` (new)

- Mirror structure of `record_writer_json_jsonl.go`:
  - **RecordWriterYAML**:
    - Option: “wrap in outer list” (like `WrapJSONOutputInOuterList`): single top-level YAML document as an array of objects.
    - Else: one YAML document per record, separated by `---`.
  - Each record is an `*mlrval.Mlrmap`; convert to `map[string]interface{}` via the new mlrval→YAML helper, then `yaml.Marshal` and write.
  - Use a small indent (e.g. 2 spaces) for readability when writing a single document; multi-doc style can be one document per line or pretty-printed.

### 2.4 Format registration and CLI

- **pkg/cli/separators.go**  
  Add `"yaml"` to:
  - `defaultFSes`  → `"N/A"`
  - `defaultPSes` → `"N/A"`
  - `defaultRSes` → `"N/A"`
  - `defaultAllowRepeatIFSes` → `false`  
  (same as `"json"`.)

- **pkg/input/record_reader_factory.go**  
  In `Create()`, add:
  - `case "yaml": return NewRecordReaderYAML(readerOptions, recordsPerBatch)`

- **pkg/output/record_writer_factory.go**  
  In `Create()`, add:
  - `case "yaml": return NewRecordWriterYAML(writerOptions)`

- **pkg/cli/option_parse.go**
  - Add long-form flags:
    - `--yaml` (or `--y2y`): input and output format `yaml`; set writer options for “wrap in outer list” if desired (e.g. `WrapYAMLOutputInOuterList = true` by default for symmetry with `--json`).
    - `--iyaml`: input format `yaml`.
    - `--oyaml`: output format `yaml` (and set wrap option).
  - Ensure `-i yaml` and `-o yaml` work (they will once `"yaml"` is in `defaultFSes` and the factories are updated; the existing “unrecognized format” check uses `defaultFSes`).

- **Writer options (pkg/cli/option_types.go)**  
  Add:
  - `WrapYAMLOutputInOuterList bool` (default `true` for `--yaml`/`--oyaml`, so one array document by default).

### 2.5 Conversion shortcuts (optional but consistent with JSON)

- Add the same style of shortcuts as for JSON, e.g.:
  - `--y2c`, `--y2t`, `--y2d`, `--y2j`, `--y2l`, `--y2p`, `--y2m`, … (YAML in → CSV/TSV/DKVP/JSON/JSONL/PPRINT/Markdown out).
  - `--c2y`, `--j2y`, etc. (other formats in → YAML out).
  These can be added in one pass or in a follow-up.

### 2.6 Help and docs

- **pkg/terminals/help/entry.go**  
  In `helpFileFormats()`, add a short “YAML” section describing:
  - Single document (object or array of objects) and multiple documents (`---`).
  - That record shape matches JSON (nested objects → dotted keys when flattened, etc.).

- **docs/src/file-formats.md**  
  Add a “YAML” subsection under file formats, with small example (one document = one record, array of objects = multiple records, multi-doc example).

- **Man page / main help**  
  Ensure “yaml” appears in the list of formats where other formats are enumerated (e.g. in `option_parse.go` or wherever the format list is printed).

---

## 3. Edge cases and behavior

- **Empty input / no records**  
  Same as JSON: if “wrap in outer list” is true, emit a single YAML document `[]`; otherwise emit nothing (or no documents).

- **Non-map elements in an array**  
  Match JSON reader: error with a clear message (“expected map (object); got …”).

- **YAML map keys that are not strings**  
  Convert to string (e.g. integer or float key `1` → `"1"`) for Miller’s record model; optionally document that keys are stringified.

- **Streaming**  
  Use `yaml.Decoder` so we don’t load the whole file; multiple documents and single-document arrays can be processed incrementally.

- **Comments**  
  YAML `#` comments are not part of the data model; we can either strip them at read (e.g. only at line start, like JSON comments) or leave as future work.

---

## 4. Testing

- Add tests in `pkg/mlrval` for:
  - `MlrvalDecodeFromYAML` (and any helper) with single map, array of maps, scalars, nested map.
  - Marshal path: Mlrmap → YAML string → decode again and compare (round-trip of records).
- Add tests for `record_reader_yaml` and `record_writer_yaml` (e.g. in `pkg/input` and `pkg/output` or in a test that runs `mlr --iyaml --oyaml` on a small fixture).
- Add a couple of regression tests in the main test suite (e.g. `mlr --yaml cat`, `mlr --icsv --oyaml cat`, `mlr --iyaml --ocsv cat`).

---

## 5. Summary of code touch points

| Area              | Change |
|-------------------|--------|
| go.mod            | Add direct `gopkg.in/yaml.v3` |
| pkg/mlrval        | New `mlrval_yaml.go` (decode + encode helpers) |
| pkg/input         | New `record_reader_yaml.go`; extend `record_reader_factory.go` |
| pkg/output        | New `record_writer_yaml.go`; extend `record_writer_factory.go` |
| pkg/cli            | `separators.go` (yaml maps), `option_types.go` (WrapYAML…), `option_parse.go` (--yaml, --iyaml, --oyaml, conversion flags) |
| pkg/terminals/help | `helpFileFormats()` YAML section |
| docs              | `file-formats.md` YAML subsection |

This keeps YAML aligned with Miller’s existing JSON subset and record model while reusing the same CLI and option patterns.

---

## Implementation summary (done)

* **Dependency:** `gopkg.in/yaml.v3` is a direct requirement in `go.mod`.
* **Core:** `pkg/mlrval/mlrval_yaml.go` — decode (including `map[string]interface{}` from the library) with **deterministic key order** (sorted on decode so YAML→CSV/TSV etc. have stable schema) and encode; `pkg/input/record_reader_yaml.go`; `pkg/output/record_writer_yaml.go` (list-wrap and multi-doc modes).
* **CLI:** `--yaml`, `--iyaml`, `--oyaml`, `-i yaml`, `-o yaml`; `--ylistwrap` / `--no-ylistwrap`; conversion shortcuts in the format matrix (`--c2y`, `--y2c`, etc.).
* **Docs:** YAML in `helpFileFormats()`, `docs/src/file-formats.md.in`, and the separators table in `reference-main-separators.md.in`.
* **Tests:** `pkg/mlrval/mlrval_yaml_test.go` (scalars, maps, arrays, multi-doc, round-trip, non-string keys); `test/cases/io-yaml-io/0001` (YAML→CSV), `0002` (YAML→YAML multi-doc); `test/cases/io-format-conversion-keystroke-savers`: c2y, j2y, y2j plus **y2c, y2t, y2d, y2n, y2p, y2x, y2m** and **t2y, d2y, n2y, p2y, x2y, m2y** for full conversion coverage.
* **Regtest:** Directory and case-path iteration in `pkg/terminals/regtest/regtester.go` is **sorted** for deterministic run order.

No open TODOs; proposal is complete.
