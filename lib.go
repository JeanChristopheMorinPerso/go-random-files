package randomfiles

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
)

type Options struct {
	Out    io.Writer // output progress
	Source io.Reader // randomness source

	FileSize     int // the size per file.
	FilenameSize int
	MaxFiles     int64  // max number of files to create
	Alphabet     []rune // for filenames

	FanoutDepth int // how deep the hierarchy goes
	FanoutFiles int // how many files per dir
	FanoutDirs  int // how many dirs per dir

	RandomSeed   int64 // use a random seed. if 0, use a random seed
	RandomSize   bool  // randomize file sizes
	RandomFanout bool  // randomize fanout numbers
}

var RunesEasy = []rune("abcdefghijklmnopqrstuvwxyz01234567890")
var RunesHard = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890!@#$%^&*()-_+= ;.,<>'\"[]{}() ")

var filesCreated = int64(0)
func WriteRandomFiles(root string, depth int, opts *Options) error {

	numfiles := opts.FanoutFiles
	if opts.RandomFanout {
		numfiles = rand.Intn(numfiles) + 1
	}

	for i := 0; i < numfiles; i++ {
		if filesCreated >= opts.MaxFiles {
			return nil
		}
		if err := WriteRandomFile(root, opts); err != nil {
			return err
		}
	}

	if depth+1 <= opts.FanoutDepth {
		numdirs := opts.FanoutDirs
		if opts.RandomFanout {
			numdirs = rand.Intn(numdirs) + 1
		}

		for i := 0; i < numdirs; i++ {
			if filesCreated >= opts.MaxFiles {
				return nil
			}
			if err := WriteRandomDir(root, depth+1, opts); err != nil {
				return err
			}
		}
	}

	return nil
}

func RandomFilename(length int, alphabet []rune) string {
	b := make([]rune, length)
	alphabetLength := len(alphabet)

	for i := range b {
		b[i] = alphabet[rand.Intn(alphabetLength)]
	}
	return string(b)
}

func WriteRandomFile(root string, opts *Options) error {
	filesize := int64(opts.FileSize)
	if opts.RandomSize {
		filesize = rand.Int63n(filesize) + 1
	}

	n := rand.Intn(opts.FilenameSize-4) + 4
	name := RandomFilename(n, opts.Alphabet)
	filepath := path.Join(root, name)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	if opts.Out != nil {
		fmt.Fprintln(opts.Out, filepath)
	}

	filesCreated++

	return f.Close()
}

func WriteRandomDir(root string, depth int, opts *Options) error {
	if depth > opts.FanoutDepth {
		return nil
	}

	n := rand.Intn(opts.FilenameSize-4) + 4
	name := RandomFilename(n, opts.Alphabet)
	root = path.Join(root, name)
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	if opts.Out != nil {
		fmt.Fprintln(opts.Out, root)
	}

	return WriteRandomFiles(root, depth, opts)
}
