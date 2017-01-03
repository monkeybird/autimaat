// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Package dictionary provides a custom dictionary.
// It allows users to define ad lookup definitions for specific terms.
package dictionary

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/monkeybird/autimaat/app/util"
	"github.com/monkeybird/autimaat/irc"
	"github.com/monkeybird/autimaat/irc/cmd"
	"github.com/monkeybird/autimaat/irc/proto"
	"github.com/monkeybird/autimaat/plugins"
)

func init() { plugins.Register(&plugin{}) }

type plugin struct {
	m           sync.RWMutex
	cmd         *cmd.Set
	file        string
	terms       map[string][]int
	definitions []string
}

// Load initializes the module and loads any internal resources
// which may be required.
func (p *plugin) Load(prof irc.Profile) error {
	p.m.Lock()

	p.file = filepath.Join(prof.Root(), "dictionary.txt")
	p.terms = make(map[string][]int)
	p.cmd = cmd.New(
		prof.CommandPrefix(),
		prof.IsWhitelisted,
	)

	p.cmd.Bind(TextDefineName, false, p.cmdDefine).
		Add(TextDefineTermName, true, cmd.RegAny)
	p.cmd.Bind(TextDefinitionsName, false, p.cmdDefinitions)

	p.m.Unlock()
	return p.loadFile()
}

// Unload cleans the module up and unloads any internal resources.
func (p *plugin) Unload(prof irc.Profile) error {
	return nil
}

// Dispatch sends the given, incoming IRC message to the plugin for
// processing as it sees fit.
func (p *plugin) Dispatch(w irc.ResponseWriter, r *irc.Request) {
	p.cmd.Dispatch(w, r)
}

// cmdDefine yields the definition of a given term, if found.
func (p *plugin) cmdDefine(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.RLock()
	defer p.m.RUnlock()

	key := strings.ToLower(params.String(0))
	indices, ok := p.terms[key]
	if !ok {
		proto.PrivMsg(w, r.Target, TextDefineNotFound, r.SenderName, util.Bold(params.String(0)))
		return
	}

	for _, index := range indices {
		proto.PrivMsg(w, r.Target, TextDefineDisplay, r.SenderName, p.definitions[index])
	}
}

// cmdDefinitions presents the user with a list of all defined terms,
// minus their definitions.
func (p *plugin) cmdDefinitions(w irc.ResponseWriter, r *irc.Request, params cmd.ParamList) {
	p.m.RLock()
	defer p.m.RUnlock()

	set := make([]string, 0, len(p.terms))
	for key := range p.terms {
		set = append(set, key)
	}

	sort.Strings(set)

	proto.PrivMsg(w, r.SenderName, TextDefinitionsDisplay, util.Bold("%d", len(set)))

	// We want to send this list in chunks. Else it will be cut
	// off early and most of it is lost.
	for {
		if len(set) > 30 {
			proto.PrivMsg(w, r.SenderName, strings.Join(set[:30], ", "))
			set = set[30:]
		} else {
			proto.PrivMsg(w, r.SenderName, strings.Join(set, ", "))
			break
		}
	}
}

// loadFile loads dictionary contents from disk.
func (p *plugin) loadFile() error {
	p.m.Lock()
	defer p.m.Unlock()

	fd, err := os.Open(p.file)
	if err != nil {
		return err
	}

	defer fd.Close()

	var terms []string
	var indices []int

	scn := bufio.NewScanner(fd)
	for scn.Scan() {
		line := strings.TrimSpace(scn.Text())
		if len(line) == 0 {
			continue
		}

		// New definition for currently active term?
		// These lines start with >
		if strings.HasPrefix(line, ">") {
			line = strings.TrimSpace(line[1:])
			if len(line) == 0 {
				continue
			}

			idx := indexOf(p.definitions, line)
			if idx > -1 {
				// no need to append duplicate definition
				indices = append(indices, idx)
				continue
			}

			p.definitions = append(p.definitions, line)
			indices = append(indices, len(p.definitions)-1)
			continue
		}

		// Store indices for currently active terms, if applicable.
		if len(terms) > 0 && len(indices) > 0 {
			for _, t := range terms {
				p.terms[t] = indices
			}
		}

		// We have a new set of terms to be defined.
		terms = split(line, ",")
		indices = nil
	}

	// Store indices for last terms in the file, if applicable.
	if len(terms) > 0 && len(indices) > 0 {
		for _, t := range terms {
			p.terms[t] = indices
		}
	}

	return scn.Err()
}

// indexOf returns the index of v in set. Returns -1 if not found.
func indexOf(set []string, v string) int {
	for i, sv := range set {
		if strings.EqualFold(sv, v) {
			return i
		}
	}
	return -1
}

// split splits v, using delimiter d. It filters out empty entries
// and transforms all resulting values to lower case.
func split(v, d string) []string {
	fields := strings.Split(v, d)
	out := make([]string, 0, len(fields))

	for _, f := range fields {
		f = strings.TrimSpace(f)
		if len(f) > 0 {
			out = append(out, strings.ToLower(f))
		}
	}

	return out
}
