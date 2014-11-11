package regnet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Regnet struct {
	Patterns map[string]Pattern
}

type Pattern struct {
	raw      string
	Compiled *regexp.Regexp
}

const (
	blockIdent string = "REGNET_BLOCK"
	blockKey   string = "REGNET_KEY"
)

//
func New() (r *Regnet, err error) {
	regentBlock, err := regexp.Compile(`\%{([^}]+)\}`)
	if err != nil {
		return nil, err
	}
	blockPattern := Pattern{blockIdent, regentBlock}
	patterns := make(map[string]Pattern)
	patterns[blockIdent] = blockPattern
	compiledPattern, err := regexp.Compile(`[\w]+`)
	if err != nil {
		return nil, err
	}
	keyPattern := Pattern{blockKey, compiledPattern}
	patterns[blockKey] = keyPattern
	return &Regnet{patterns}, nil
}

//
func (regnet *Regnet) AddPattern(name string, pattern string) (err error) {
	if _, present := regnet.GetPattern(name); present == true {
		return errors.New("regnet: pattern " + name + " already exists.")
	}

	r := regnet.Patterns[blockIdent].Compiled
	slices := r.FindAllString(pattern, -1)
	for indx := range slices {
		key := regnet.Patterns[blockKey].Compiled.FindString(slices[indx])
		value, present := regnet.GetPattern(key)
		if present == false {
			return errors.New("regnet: pattern " + key + " not found. Define it before " + name + " regnet.")
		} else {
			//replace regent this its derefrenced pattern string
			pattern = strings.Replace(pattern, "%{" + key + "}", value.Compiled.String(), -1)
		}
	}
	//  contains only Regnet, so get the value and compile it
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	patternCompiled := Pattern{name, compiled}
	regnet.Patterns[name] = patternCompiled
	return nil
}

type Match struct {
	Ident   string
	results []string
}

func (match *Match) String() string {
	return fmt.Sprintf("%s : %s", match.Ident, match.results)
}

func (match *Match) Step() []string {
	return match.results
}

func (regnet *Regnet) MatchInText(text, patterns string) (match *Match, err error) {
	regnets := regnet.Patterns[blockIdent].Compiled.FindAllString(patterns, -1)
	if regnets != nil {
		stripped := regnet.Patterns[blockKey].Compiled.FindString(regnets[0])
		pattern, present := regnet.GetPattern(stripped)
		if present {
			matched := pattern.Compiled.FindAllString(text, -1)
			return &Match{stripped, matched}, nil
		} else {
			return nil, errors.New("regnet: pattern " + stripped + " not found.")
		}
	} else {
		return nil, errors.New("regnet: invalid pattern definition. Format: %{insert_regent_name_here}")
	}
	return nil, nil
}

func (regnet *Regnet) GetPattern(name string) (pattern Pattern, present bool) {
	pattern, present = regnet.Patterns[name]
	return pattern, present
}

func (regnet *Regnet) AddPatternsFromFile(path string) (err error) {
	files, err := filepath.Glob(path)
	if err == nil {
		for file := range files {
			fmt.Println(files[file])
			if patternFile, err := os.Open(files[file]); err == nil {
				defer patternFile.Close()
				reader := bufio.NewReader(patternFile)
				eof := false
				for !eof {
					var line string
					line, err = reader.ReadString('\n')
					if len(line) > 1 && strings.HasPrefix(line, "#") == false {
						index := subStr(line, 0, strings.Index(line, " "))
						pattern := subStr(line, strings.Index(line, " "), len(line))
						error := regnet.AddPattern(strings.TrimSpace(index), strings.TrimSpace(pattern))
						if error != nil {
							return error
						}
					}
					if err == io.EOF {
						eof = true
					}
				}
			}
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}