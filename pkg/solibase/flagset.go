package solibase

import (
	"flag"
)

type CustomFlagSet struct {
	prefix string
	fs     *flag.FlagSet
}

func (f *CustomFlagSet) StringVar(p *string, name, value, usage string) {
	f.fs.StringVar(p, f.prefix+"-"+name, value, usage)
}

func (f *CustomFlagSet) IntVar(i *int, name string, value int, usage string) {
	f.fs.IntVar(i, f.prefix+"-"+name, value, usage)
}

type FlagSetGenerator struct {
	FS *flag.FlagSet
}

func (g *FlagSetGenerator) New(prefix string) FlagSet {
	return &CustomFlagSet{prefix: prefix, fs: g.FS}
}
