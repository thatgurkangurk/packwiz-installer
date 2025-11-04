// Copyright 2025 Gurkan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package core

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Updates tracks changes to be made during installation/update
type Updates struct {
	Added     []*Mod
	Removed   []*Mod
	Unchanged []*Mod
}

func (u *Updates) String() string {
	var s string
	s += "Added:\n"
	for _, m := range u.Added {
		s += fmt.Sprintf("  %s\n", m.Path)
	}
	s += "Removed:\n"
	for _, m := range u.Removed {
		s += fmt.Sprintf("  %s\n", m.Path)
	}
	s += "Unchanged:\n"
	for _, m := range u.Unchanged {
		s += fmt.Sprintf("  %s\n", m.Path)
	}
	return s
}

// LocalInstaller manages installation and updates of mods in a local directory
type LocalInstaller struct {
	BaseDir    string
	Pack       *Pack
	GameSide   Side
	httpClient *http.Client
}

// NewLocalInstaller creates a new installer for the given pack in the specified directory
func NewLocalInstaller(p *Pack, dir string, gameSide Side) (*LocalInstaller, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	return &LocalInstaller{
		BaseDir:    abs,
		Pack:       p,
		GameSide:   gameSide,
		httpClient: http.DefaultClient,
	}, nil
}

func (i *LocalInstaller) saveCache(name string, v any) error {
	p := filepath.Join(i.BaseDir, ".pw-install", fmt.Sprintf("%s.json", name))
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, os.ModePerm)
}

func (i *LocalInstaller) restoreCache(name string, v any) error {
	p := filepath.Join(i.BaseDir, ".pw-install", fmt.Sprintf("%s.json", name))
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (i *LocalInstaller) setInstalledMods() error {
	return i.saveCache("installed", i.Pack.Mods)
}

func (i *LocalInstaller) getInstalledMods() ([]*Mod, error) {
	var mods []*Mod
	err := i.restoreCache("installed", &mods)
	if err != nil {
		return nil, err
	}
	return mods, nil
}

func (i *LocalInstaller) checkIntegrity(m *Mod) (bool, error) {
	// existence
	p := filepath.Join(i.BaseDir, m.Path)
	stat, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if stat.IsDir() {
			return false, nil
		}
		return false, err
	}

	// hash
	data, err := os.ReadFile(p)
	if err != nil {
		return false, err
	}
	valid, err := MatchHash(data, m.HashFormat, m.Hash)
	if err != nil {
		return false, err
	}

	return valid, nil
}

// GetUpdates determines which mods need to be added, removed, or are unchanged
func (i *LocalInstaller) GetUpdates() (*Updates, error) {
	installed, err := i.getInstalledMods()
	if err != nil {
		return nil, err
	}

	// Filter mods based on game side
	filteredMods := make([]*Mod, 0, len(i.Pack.Mods))
	for _, m := range i.Pack.Mods {
		if i.GameSide.ShouldInstall(m.Side) {
			filteredMods = append(filteredMods, m)
		}
	}

	a, r, u := diffSliceFunc(installed, filteredMods, func(a, b *Mod) int {
		res := cmp.Compare(a.Path, b.Path)
		if res == 0 && a.Hash != b.Hash {
			res = -1
		}
		return res
	})
	return &Updates{
		Added:     a,
		Removed:   r,
		Unchanged: u,
	}, nil
}

// InstallMod installs or updates a single mod
func (i *LocalInstaller) InstallMod(ctx context.Context, m *Mod) error {
	var (
		data []byte
		err  error
	)
	switch m.Downloads.Type {
	case DL_Url:
		data, err = httpGetValidBytes(ctx, i.httpClient, m.Downloads.Data, m.HashFormat, m.Hash)
		if err != nil {
			return err
		}
	case DL_Curseforge:
		cfData, err := ParseCfData(m.Downloads.Data)
		if err != nil {
			return err
		}
		u, err := DefaultCurseClient.GetDownloadUrl(ctx, cfData)
		if err != nil {
			return err
		}
		data, err = httpGetValidBytes(ctx, i.httpClient, u, m.HashFormat, m.Hash)
		if err != nil {
			return err
		}
	}

	p := filepath.Join(i.BaseDir, m.Path)
	err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, os.ModePerm)
}

// Install executes installation and update of the modpack
func (i *LocalInstaller) Install(ctx context.Context) (*Updates, error) {
	var result = &Updates{}
	update, err := i.GetUpdates()
	if err != nil {
		return nil, fmt.Errorf("check updates: %w", err)
	}

	mut := sync.Mutex{}
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())

	for _, m := range update.Unchanged {
		m := m // capture for closure
		eg.Go(func() error {
			ok, err := i.checkIntegrity(m)
			if err != nil {
				return fmt.Errorf("check integrity: %w", err)
			}
			mut.Lock()
			if ok {
				result.Unchanged = append(result.Unchanged, m)
			} else {
				update.Added = append(update.Added, m)
			}
			mut.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for _, m := range update.Added {
		m := m // capture for closure
		eg.Go(func() error {
			err := i.InstallMod(ctx, m)
			if err != nil {
				return fmt.Errorf("install mod: %w", err)
			}
			mut.Lock()
			result.Added = append(result.Added, m)
			mut.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for _, m := range update.Removed {
		m := m // capture for closure
		eg.Go(func() error {
			p := filepath.Join(i.BaseDir, m.Path)
			err := os.Remove(p)
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("remove mod: %w", err)
				}
			}
			mut.Lock()
			result.Removed = append(result.Removed, m)
			mut.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	err = i.setInstalledMods()
	if err != nil {
		return nil, fmt.Errorf("save cache: %w", err)
	}
	return result, nil
}
