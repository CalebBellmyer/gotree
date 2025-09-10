package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// opts holds CLI options
type opts struct {
	depth     int
	all       bool
	dirsOnly  bool
	showSize  bool
}

func main() {
	var o opts
	flag.IntVar(&o.depth, "depth", 0, "max depth to display (0 = unlimited)")
	flag.BoolVar(&o.all, "all", false, "include hidden files and directories (those starting with .)")
	flag.BoolVar(&o.dirsOnly, "dirs-only", false, "show directories only")
	flag.BoolVar(&o.showSize, "size", false, "show file sizes")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: gotree [options] [PATH]\n\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Renders a text-based directory tree. Defaults to current directory.")
		fmt.Fprintln(flag.CommandLine.Output(), "\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	info, err := os.Stat(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gotree: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "gotree: %s is not a directory\n", root)
		os.Exit(1)
	}

	// Normalize root to clean path
	root = filepath.Clean(root)
	absRoot, err := filepath.Abs(root)
	if err != nil {
    	absRoot = root // fallback
	}
	fmt.Println(filepath.Base(absRoot))


	if err := walk(root, &o, 0, nil, ""); err != nil {
		fmt.Fprintf(os.Stderr, "gotree: %v\n", err)
		os.Exit(1)
	}
}

// walk prints the directory tree for path p.
// depth=0 is unlimited; otherwise stop when level == depth.
// prefix carries the drawing guide (e.g., "│   ").
func walk(p string, o *opts, level int, guides []bool, prefix string) error {
	if o.depth > 0 && level >= o.depth {
		return nil
	}

	entries, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	// Filter hidden unless -all
	filtered := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !o.all && strings.HasPrefix(name, ".") {
			continue
		}
		if o.dirsOnly && !e.IsDir() {
			continue
		}
		filtered = append(filtered, e)
	}

	// Sort: directories first, then files; alphabetical within groups
	sort.Slice(filtered, func(i, j int) bool {
		ei, ej := filtered[i], filtered[j]
		if ei.IsDir() != ej.IsDir() {
			return ei.IsDir() && !ej.IsDir()
		}
		return strings.ToLower(ei.Name()) < strings.ToLower(ej.Name())
	})

	for idx, e := range filtered {
		isLast := idx == len(filtered)-1
		var branch string
		if isLast {
			branch = "└── "
		} else {
			branch = "├── "
		}

		line := prefix + branch + e.Name()
		if o.showSize && !e.IsDir() {
			if info, err := e.Info(); err == nil {
				line += fmt.Sprintf(" (%s)", humanSize(info.Size()))
			}
		}
		fmt.Println(line)

		if e.IsDir() {
			nextPrefix := prefix
			if isLast {
				nextPrefix += "    "
			} else {
				nextPrefix += "│   "
			}
			if err := walk(filepath.Join(p, e.Name()), o, level+1, append(guides, !isLast), nextPrefix); err != nil {
				return err
			}
		}
	}
	return nil
}

// humanSize converts bytes to a readable string.
func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	d := float64(n)
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	for i := 0; i < len(units); i++ {
		if d < unit {
			return fmt.Sprintf("%.1f%s", d, units[i])
		}
		d /= unit
	}
	return fmt.Sprintf("%.1fPiB", d)
}
