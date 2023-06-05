package filesfinder

import "os"

type finder struct {
	folder string
}

type option func(f *finder) error

func New(opts ...option) (finder, error) {
	f := finder{
		folder: ".",
	}
	for _, opt := range opts {
		if err := opt(&f); err != nil {
			return finder{}, err
		}
	}
	return f, nil

}

func WithPath(folder string) option {
	return func(f *finder) error {
		f.folder = folder
		return nil
	}

}

func (f finder) FilesUnder() (int, error) {
	entries, err := os.ReadDir(f.folder)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}
