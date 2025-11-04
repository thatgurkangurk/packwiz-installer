package core

import (
	"fmt"
	"net/url"

	toml "github.com/pelletier/go-toml/v2"
)

// PackToml represents the top-level pack.toml manifest.
type PackToml struct {
	Name        string            `toml:"name"`
	Author      string            `toml:"author,omitempty"`
	Version     string            `toml:"version,omitempty"`
	Description string            `toml:"description,omitempty"`
	PackFormat  string            `toml:"pack-format"`
	Index       PackIndexSection  `toml:"index"`
	Versions    map[string]string `toml:"versions,omitempty"`
}

// PackIndexSection is the `index` table inside pack.toml.
type PackIndexSection struct {
	File       string `toml:"file"`
	HashFormat string `toml:"hash-format,omitempty"`
	Hash       string `toml:"hash,omitempty"`
}

// IndexToml represents the index.toml format that lists files.
type IndexToml struct {
	HashFormat string            `toml:"hash-format,omitempty"`
	Files      []IndexedFileToml `toml:"files,omitempty"`
}

// IndexedFileToml is an entry inside index.toml's [[files]] array.
type IndexedFileToml struct {
	File       string `toml:"file"`
	Hash       string `toml:"hash,omitempty"`
	Alias      string `toml:"alias,omitempty"`
	HashFormat string `toml:"hash-format,omitempty"`
	Metafile   bool   `toml:"metafile,omitempty"`
	Preserve   bool   `toml:"preserve,omitempty"`
}

// UpdateModrinth describes a modrinth-style update reference.
type UpdateModrinth struct {
	ModId   string `toml:"mod-id,omitempty"`
	Version string `toml:"version,omitempty"`
}

// UpdateCurseForge describes a curseforge update reference.
type UpdateCurseForge struct {
	FileId    int `toml:"file-id,omitempty"`
	ProjectId int `toml:"project-id,omitempty"`
}

// MetafileUpdate contains update references for a metafile.
type MetafileUpdate struct {
	Modrinth   *UpdateModrinth   `toml:"modrinth,omitempty,nullable"`
	CurseForge *UpdateCurseForge `toml:"curseforge,omitempty,nullable"`
}

// MetafileDownload describes where to download a metafile and its hash info.
type MetafileDownload struct {
	HashFormat string `toml:"hash-format,omitempty"`
	Hash       string `toml:"hash,omitempty"`
	Url        string `toml:"url,omitempty"`
	Mode       string `toml:"mode,omitempty"`
}

// MetafileOption contains optional flagging and defaults for a metafile.
type MetafileOption struct {
	Optional    bool   `toml:"optional,omitempty"`
	Default     bool   `toml:"default,omitempty"`
	Description string `toml:"description,omitempty"`
}

// MetafileToml represents a single metafile (a single file definition TOML).
type MetafileToml struct {
	Filename  string            `toml:"filename"`
	Name      string            `toml:"name"`
	Side      string            `toml:"side,omitempty"`
	Download  *MetafileDownload `toml:"download,omitempty"`
	Update    *MetafileUpdate   `toml:"update,omitempty"`
	Option    *MetafileOption   `toml:"option,omitempty"`
	IndexName string
}

// ==============================
// Parsing helpers
// ==============================

// ParsePackToml parses pack.toml bytes into a PackToml and performs light validation.
// `source` is an optional context string (filename/URL) to improve error messages.
func ParsePackToml(data []byte, source string) (*PackToml, error) {
	var p PackToml
	if err := toml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pack.toml (%s): %w", source, err)
	}

	// minimal validation
	if p.Name == "" {
		return nil, fmt.Errorf("pack.toml (%s): missing required field 'name'", source)
	}
	if p.PackFormat == "" {
		return nil, fmt.Errorf("pack.toml (%s): missing required field 'pack-format'", source)
	}

	// ensure map is non-nil
	if p.Versions == nil {
		p.Versions = make(map[string]string)
	}

	return &p, nil
}

// ParseIndexToml parses index.toml bytes into IndexToml and performs light validation.
func ParseIndexToml(data []byte, source string) (*IndexToml, error) {
	var idx IndexToml
	if err := toml.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal index.toml (%s): %w", source, err)
	}

	// Ensure slice non-nil for easier iteration by callers.
	if idx.Files == nil {
		idx.Files = make([]IndexedFileToml, 0)
	}

	return &idx, nil
}

// ParseMetafileToml parses an individual metafile TOML document into MetafileToml.
// Performs light validation (name/filename present).
func ParseMetafileToml(data []byte, source string) (*MetafileToml, error) {
	var m MetafileToml
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metafile (%s): %w", source, err)
	}

	// Basic required fields
	if m.Filename == "" {
		return nil, fmt.Errorf("metafile (%s): missing required field 'filename'", source)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("metafile (%s): missing required field 'name'", source)
	}

	// Validate download URL if present
	if m.Download != nil && m.Download.Url != "" {
		if err := validateURL(m.Download.Url); err != nil {
			return nil, fmt.Errorf("metafile (%s): invalid download.url: %w", source, err)
		}
	}

	// Normalize option/update pointers to avoid nil checks by callers
	if m.Option == nil {
		m.Option = &MetafileOption{}
	}
	if m.Update == nil {
		m.Update = nil // keep nil when absent (caller may check)
	}

	return &m, nil
}

// ==============================
// Utility / validation methods
// ==============================

// IsOptional returns true if the metafile is marked optional.
func (m *MetafileToml) IsOptional() bool {
	if m == nil || m.Option == nil {
		return false
	}
	return m.Option.Optional
}

// HasDownload checks if metafile contains a download block with a URL.
func (m *MetafileToml) HasDownload() bool {
	return m != nil && m.Download != nil && m.Download.Url != ""
}

// GetDownloadURL returns the download URL (empty string if none).
func (m *MetafileToml) GetDownloadURL() string {
	if m == nil || m.Download == nil {
		return ""
	}
	return m.Download.Url
}

// validateURL does a minimal sanity check for a URL.
func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url must be absolute with scheme and host")
	}
	return nil
}
