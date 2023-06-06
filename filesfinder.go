package filesfinder

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"
)

type finder struct {
	fsys        fs.FS
	currentTime time.Time
}

const (
	minsinsecs = 60 * time.Second
	hrstosecs  = 60 * minsinsecs
)

type Option func(f *finder) error

// New provides the construction of finder instance.
// first argument is fs.FS type
// second is the list of time options letting the user specify the time
// as a combination of years, months, days, hours and minutes
func New(fsys fs.FS, opts ...Option) (finder, error) {
	f := finder{
		fsys:        fsys,
		currentTime: time.Now(),
	}
	for _, opt := range opts {
		if err := opt(&f); err != nil {
			return finder{}, err
		}
	}
	return f, nil

}

func WithYears(years int) Option {
	return func(f *finder) error {
		f.currentTime = f.currentTime.AddDate(-years, 0, 0)
		return nil
	}
}

func WithMonths(months int) Option {
	return func(f *finder) error {
		f.currentTime = f.currentTime.AddDate(0, -months, 0)
		return nil
	}
}

func WithDays(days int) Option {
	return func(f *finder) error {
		f.currentTime = f.currentTime.AddDate(0, 0, -days)
		return nil
	}
}

func WithHours(hours int) Option {
	return func(f *finder) error {
		f.currentTime = f.currentTime.Add(time.Duration(hours) * -hrstosecs)
		return nil
	}
}

func WithMinutes(minutes int) Option {
	return func(f *finder) error {
		f.currentTime = f.currentTime.Add(time.Duration(minutes) * -minsinsecs)
		return nil
	}
}

type Result struct {
	Path    string
	ModTime time.Time
	Error   error
}
type Results []Result

// Finds files older than the time options specified
// in local file system.
func (f finder) OlderThan() (Results, error) {

	var results Results
	walkerr := fs.WalkDir(f.fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			results = append(results, Result{Path: path, Error: err})
			return fs.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		finfo, err := d.Info()
		if err != nil { // handle other way. add to struct
			results = append(results, Result{Path: path, Error: err})
			return fs.SkipDir
		}
		if finfo.ModTime().Before(f.currentTime) {
			results = append(results, Result{Path: path, ModTime: finfo.ModTime()})
		}
		return nil

	})

	return results, walkerr
}

func makeFinder(inputArgs []string) (finder, error) {
	if inputArgs == nil || len(inputArgs) == 0 {
		return finder{}, fmt.Errorf("No input provided")
	}
	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	years := fset.Int("y", 0, "Years option.")
	months := fset.Int("mo", 0, "Months option.")
	days := fset.Int("d", 0, "Days option.")
	hours := fset.Int("h", 0, "Hours option.")
	minutes := fset.Int("mi", 0, "Minutes option.")
	fset.SetOutput(os.Stdout)
	err := fset.Parse(inputArgs)
	if err != nil {
		return finder{}, err
	}

	args := fset.Args()
	if len(args) < 1 {
		return finder{}, fmt.Errorf("no folder path provided to search under")
	}

	fs, err := New(os.DirFS(args[0]),
		WithYears(*years),
		WithMonths(*months),
		WithDays(*days),
		WithHours(*hours),
		WithMinutes(*minutes))

	return fs, err
}

func RunCLI() {
	finder, err := makeFinder(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	results, err := finder.OlderThan()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, res := range results {
		fmt.Printf("%s		%s\n", res.Path, res.ModTime)
	}
	fmt.Println("Erros:")
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("%s\n", res.Error)
		}
	}

}
