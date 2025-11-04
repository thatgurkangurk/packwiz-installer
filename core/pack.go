package core

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Side string
type DLType string

const (
	Side_Client = Side("client")
	Side_Server = Side("server")
	Side_Both   = Side("both")

	DL_Url        = DLType("url")
	DL_Curseforge = DLType("curseforge")
)

type Download struct {
	Type DLType `json:"type"`
	Data string `json:"data"`
}

type Mod struct {
	Path       string    `json:"path"`
	Hash       string    `json:"hash"`
	HashFormat string    `json:"hashFormat"`
	Side       Side      `json:"side,omitempty"`
	Downloads  *Download `json:"download"`
}

type Pack struct {
	Name    string `json:"name"`
	Author  string `json:"author,omitempty"`
	Version string `json:"version,omitempty"`
	Mods    []*Mod `json:"files,omitempty"`
}

type CurseforgeData struct {
	ProjectID int `json:"projectId"`
	FileID    int `json:"fileId"`
}

func (d *CurseforgeData) String() string {
	return fmt.Sprintf("%d:%d", d.ProjectID, d.FileID)
}

func ParseCfData(s string) (*CurseforgeData, error) {
	a := strings.Split(s, ":")
	if len(a) < 2 {
		return nil, fmt.Errorf("invalid curseforge data")
	} else {
		pid, err := strconv.Atoi(a[0])
		if err != nil {
			return nil, err
		}
		fid, err := strconv.Atoi(a[1])
		if err != nil {
			return nil, err
		}

		return &CurseforgeData{
			ProjectID: pid,
			FileID:    fid,
		}, nil
	}
}

func tomlToPack(
	packUrl *url.URL,
	pack *PackToml,
	index *IndexToml,
	metafiles []*MetafileToml,
) (*Pack, error) {
	var ppack = &Pack{
		Name:    pack.Name,
		Author:  pack.Author,
		Version: pack.Version,
	}

	var mods = make([]*Mod, 0, len(index.Files))
	for _, f := range index.Files {
		if f.Metafile {
			i := slices.IndexFunc(metafiles, func(m *MetafileToml) bool {
				return m.IndexName == f.File
			})
			if i == -1 {
				return nil, fmt.Errorf("metafile not found: %s", f.File)
			}

			metafile := metafiles[i]
			var dl = &Download{}
			switch metafile.Download.Mode {
			// use url when mode is empty
			// https://github.com/packwiz/packwiz/blob/7545d9a777739655de749dedcd383dee6bbfd2e2/core/mod.go#L39
			case "":
				dl = &Download{
					Type: DL_Url,
					Data: metafile.Download.Url,
				}
			case "metadata:curseforge":
				cfData := &CurseforgeData{
					ProjectID: metafile.Update.CurseForge.ProjectId,
					FileID:    metafile.Update.CurseForge.FileId,
				}
				dl = &Download{
					Type: DL_Curseforge,
					Data: cfData.String(),
				}
			}

			modDir := filepath.ToSlash(filepath.Join(filepath.Dir(pack.Index.File), filepath.Dir(f.File)))
			modPath := filepath.ToSlash(filepath.Join(modDir, metafile.Filename))
			m := &Mod{
				Path:       modPath,
				Hash:       metafile.Download.Hash,
				HashFormat: metafile.Download.HashFormat,
				Side:       Side(metafile.Side),
				Downloads:  dl,
			}

			mods = append(mods, m)
		} else {
			hashFmt := f.HashFormat
			if hashFmt == "" {
				hashFmt = index.HashFormat
			}
			modPath := filepath.ToSlash(filepath.Join(filepath.Dir(pack.Index.File), f.File))
			modUrl := packUrl.JoinPath("..", modPath)
			dl := &Download{
				Type: DL_Url,
				Data: modUrl.String(),
			}

			m := &Mod{
				Path:       f.File,
				Hash:       f.Hash,
				HashFormat: hashFmt,
				Side:       Side_Both,
				Downloads:  dl,
			}
			mods = append(mods, m)
		}
	}

	ppack.Mods = mods
	return ppack, nil
}

func NewPack(r *Repository) (*Pack, error) {
	return tomlToPack(r.Url, r.Pack, r.Index, r.Metafiles)
}
