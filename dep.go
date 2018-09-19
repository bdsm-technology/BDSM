package main

import (
	"debug/elf"
	"io/ioutil"
)

func getSoDeps(filepath string) ([]string, error) {
	file, err := elf.Open("mods/" + filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file.ImportedLibraries()
}

type scanResult map[string]map[string]bool

func (modsList scanResult) filterDeps(deps []string) (ret []string) {
	for _, dep := range deps {
		if _, ok := modsList[dep]; ok {
			ret = append(ret, dep)
		}
	}
	return
}

func (modsList scanResult) getModDeps(file string) ([]string, error) {
	list, err := getSoDeps(file)
	if err != nil {
		return nil, err
	}
	return modsList.filterDeps(list), nil
}

func scanMods() (modsList scanResult, err error) {
	list, err := ioutil.ReadDir("mods")
	if err != nil {
		return
	}
	modsList = make(scanResult)
	for _, info := range list {
		modsList[info.Name()] = nil
	}
	for name := range modsList {
		deps, err := modsList.getModDeps(name)
		if err != nil {
			return nil, err
		}
		for _, item := range deps {
			modsList[name] = make(map[string]bool)
			modsList[name][item] = true
		}
	}
	for name, deps := range modsList {
		changed := true
		for changed {
			changed = false
			var ndeps []string
			for dep := range deps {
				subdeps := modsList[dep]
				for subdep := range subdeps {
					if _, ok := deps[subdep]; !ok {
						changed = true
						ndeps = append(ndeps, subdep)
					}
				}
			}
			if changed {
				for _, n := range ndeps {
					modsList[name][n] = true
				}
			}
		}
	}
	return
}
