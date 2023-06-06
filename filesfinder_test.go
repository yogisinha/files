package filesfinder_test

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/yogisinha/filesfinder"
)

// test File system. file entries are sorted by mod time descending
var testFileSystem = fstest.MapFS{
	"a/b/c": &fstest.MapFile{
		ModTime: time.Now().AddDate(0, -2, 0),
	},
	"a/b/d": &fstest.MapFile{
		ModTime: time.Now().AddDate(0, -7, 0),
	},
	"a/b/e": &fstest.MapFile{
		ModTime: time.Now().AddDate(0, -11, 0),
	},
	"a/b/f": &fstest.MapFile{
		ModTime: time.Now().AddDate(-1, 0, 0),
	},
	"a/b/g": &fstest.MapFile{
		ModTime: time.Now().AddDate(-1, -1, 0),
	},
	"a/b/h": &fstest.MapFile{
		ModTime: time.Now().AddDate(-1, -6, 0),
	},
	"a/b/i": &fstest.MapFile{
		ModTime: time.Now().AddDate(-1, -9, 0),
	},
	"a/b/j": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, 0, 0),
	},
	"a/b/k": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -3, 0),
	},
	"a/b/l": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -6, 0),
	},
	"a/b/l1": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -6, -1),
	},
	"a/b/l2": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -6, -4),
	},
	"a/b/l3": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -6, -10),
	},
	"a/b/m": &fstest.MapFile{
		ModTime: time.Now().AddDate(-2, -10, 0),
	},
	"a/b/n": &fstest.MapFile{
		ModTime: time.Now().AddDate(-3, 0, 0),
	},
}

type testCase struct {
	options []filesfinder.Option
	want    int
}

func TestFindFilesOlderThanDifferentCombinations(t *testing.T) {
	testCases := []testCase{
		{
			options: []filesfinder.Option{
				filesfinder.WithMonths(2), filesfinder.WithDays(2),
			},
			want: 14,
		},
		{
			options: []filesfinder.Option{
				filesfinder.WithYears(2), filesfinder.WithMonths(6), filesfinder.WithDays(2),
			},
			want: 4,
		},
		{
			options: []filesfinder.Option{
				filesfinder.WithYears(2), filesfinder.WithMonths(6),
				filesfinder.WithDays(4), filesfinder.WithHours(2), filesfinder.WithMinutes(30),
			},
			want: 3,
		},
		{
			options: []filesfinder.Option{
				filesfinder.WithYears(2), filesfinder.WithMonths(6), filesfinder.WithHours(241),
			},
			want: 2,
		},
	}

	for _, tt := range testCases {
		f, err := filesfinder.New(testFileSystem, tt.options...)
		if err != nil {
			t.Fatal(err)
		}
		results, err := f.OlderThan()
		if err != nil {
			t.Logf("Error %s", err)
		}
		got := len(results)
		if tt.want != got {
			t.Errorf("want %d, got %d", tt.want, got)
		}
	}

}
