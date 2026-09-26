package lib

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Minimal TZif v1 data with explicit offsets and optional transitions.
func timezoneTestData(offsets []int32, transitions []int32, indices []byte) []byte {
	data := make([]byte, 44)
	copy(data, "TZif")
	binary.BigEndian.PutUint32(data[32:36], uint32(len(transitions)))
	binary.BigEndian.PutUint32(data[36:40], uint32(len(offsets)))
	binary.BigEndian.PutUint32(data[40:44], 8)
	for _, transition := range transitions {
		data = binary.BigEndian.AppendUint32(data, uint32(transition))
	}
	data = append(data, indices...)
	for i, offset := range offsets {
		data = binary.BigEndian.AppendUint32(data, uint32(offset))
		data = append(data, byte(i), byte(i*4))
	}
	return append(data, "STD\x00DST\x00"...)
}

func writeTimezoneTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSetTZFromEnv(t *testing.T) {
	savedLocal := time.Local
	t.Cleanup(func() { time.Local = savedLocal })
	directory := t.TempDir()
	zoneFile := filepath.Join(directory, "zone.tzif")
	spacedFile := filepath.Join(directory, "zone with spaces-时区.tzif")
	badFile := filepath.Join(directory, "bad.tzif")
	emptyFile := filepath.Join(directory, "empty.tzif")
	data := timezoneTestData([]int32{5400}, nil, nil)
	writeTimezoneTestFile(t, zoneFile, data)
	writeTimezoneTestFile(t, spacedFile, data)
	writeTimezoneTestFile(t, badFile, []byte("not timezone data"))
	writeTimezoneTestFile(t, emptyFile, nil)

	for _, test := range []struct {
		name       string
		tz         string
		wantOffset int
		wantError  bool
		unchanged  bool
	}{
		{"UTC", "UTC", 0, false, false},
		{"colon UTC", ":UTC", 0, false, false},
		{"colon only", ":", 0, false, false},
		{"Local", "Local", 9000, false, true},
		{"colon Local", ":Local", 9000, false, true},
		{"empty", "", 9000, false, true},
		{"absolute file", zoneFile, 5400, false, false},
		{"colon file", ":" + zoneFile, 5400, false, false},
		{"spaces and Unicode", spacedFile, 5400, false, false},
		{"colon spaces and Unicode", ":" + spacedFile, 5400, false, false},
		{"missing file", filepath.Join(directory, "missing.tzif"), 0, true, true},
		{"colon missing file", ":" + filepath.Join(directory, "missing.tzif"), 0, true, true},
		{"invalid file", badFile, 0, true, true},
		{"empty file", ":" + emptyFile, 0, true, true},
		{"directory", ":" + directory, 0, true, true},
		{"invalid name", "this-is-not-a-valid-timezone-name", 0, true, true},
		{"colon invalid name", ":this-is-not-a-valid-timezone-name", 0, true, true},
		{"extra colon", "::UTC", 0, true, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			before := time.FixedZone("before", 9000)
			time.Local = before
			t.Setenv("TZ", test.tz)
			err := SetTZFromEnv()
			if (err != nil) != test.wantError {
				t.Fatalf("error = %v, wantError = %v", err, test.wantError)
			}
			if test.unchanged && time.Local != before {
				t.Fatal("the previous timezone changed")
			}
			if err != nil {
				if want := fmt.Sprintf("TZ environment variable appears malformed: \"%s\"", test.tz); err.Error() != want {
					t.Errorf("error = %q, want %q", err, want)
				}
				return
			}
			if _, offset := time.Unix(0, 0).Local().Zone(); offset != test.wantOffset {
				t.Errorf("offset = %d, want %d", offset, test.wantOffset)
			}
		})
	}
}

func TestSetTZFromEnvTransitions(t *testing.T) {
	savedLocal := time.Local
	t.Cleanup(func() { time.Local = savedLocal })
	path := filepath.Join(t.TempDir(), "transitions.tzif")
	writeTimezoneTestFile(t, path, timezoneTestData([]int32{3600, 7200}, []int32{100, 200}, []byte{1, 0}))
	t.Setenv("TZ", ":"+path)
	if err := SetTZFromEnv(); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ timestamp, offset int64 }{{0, 3600}, {150, 7200}, {250, 3600}} {
		if _, offset := time.Unix(test.timestamp, 0).Local().Zone(); int64(offset) != test.offset {
			t.Errorf("timestamp %d: offset = %d, want %d", test.timestamp, offset, test.offset)
		}
	}
}

func TestSetTZFromEnvReload(t *testing.T) {
	savedLocal := time.Local
	t.Cleanup(func() { time.Local = savedLocal })
	path := filepath.Join(t.TempDir(), "reload.tzif")
	t.Setenv("TZ", ":"+path)
	for _, offset := range []int32{5400, -19800} {
		writeTimezoneTestFile(t, path, timezoneTestData([]int32{offset}, nil, nil))
		if err := SetTZFromEnv(); err != nil {
			t.Fatal(err)
		}
		if _, got := time.Unix(0, 0).Local().Zone(); got != int(offset) {
			t.Errorf("reloaded offset = %d, want %d", got, offset)
		}
	}
	t.Setenv("TZ", "UTC")
	if err := SetTZFromEnv(); err != nil {
		t.Fatal(err)
	}
	if time.Local != time.UTC {
		t.Fatal("switching back to UTC did not take effect")
	}
}

func TestSetTZFromEnvFileSize(t *testing.T) {
	savedLocal := time.Local
	t.Cleanup(func() { time.Local = savedLocal })
	for _, size := range []int64{10 << 20, (10 << 20) + 1} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "large.tzif")
			writeTimezoneTestFile(t, path, timezoneTestData([]int32{5400}, nil, nil))
			if err := os.Truncate(path, size); err != nil {
				t.Fatal(err)
			}
			t.Setenv("TZ", ":"+path)
			before := time.FixedZone("before", 9000)
			time.Local = before
			err := SetTZFromEnv()
			if size > 10<<20 {
				if err == nil || time.Local != before {
					t.Fatalf("oversized file changed the timezone or had no error: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			} else if _, offset := time.Unix(0, 0).Local().Zone(); offset != 5400 {
				t.Errorf("offset = %d, want 5400", offset)
			}
		})
	}
}
