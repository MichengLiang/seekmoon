package model

import (
	"fmt"
	"strings"
)

// ModuleIndexTree mirrors module_index.json package hierarchy.
type ModuleIndexTree struct {
	Name    string            `json:"name"`
	Package *PackageSummary   `json:"package,omitempty"`
	Childs  []ModuleIndexTree `json:"childs,omitempty"`
}

// PackageSummary summarizes one package node from a module index.
type PackageSummary struct {
	Path      string         `json:"path"`
	RelPath   EvidenceString `json:"relpath"`
	Traits    []APISummary   `json:"traits,omitempty"`
	Errors    []APISummary   `json:"errors,omitempty"`
	Types     []TypeSummary  `json:"types,omitempty"`
	TypeAlias []APISummary   `json:"typealias,omitempty"`
	Values    []APISummary   `json:"values,omitempty"`
	Misc      []APISummary   `json:"misc,omitempty"`
}

// APISummary names a lightweight API item from module_index.json.
type APISummary struct {
	Name string `json:"name"`
}

// TypeSummary names a type and its method summaries.
type TypeSummary struct {
	Name    string       `json:"name"`
	Methods []APISummary `json:"methods,omitempty"`
}

// PackageData contains detailed API entries from package_data.json.
type PackageData struct {
	Name      string     `json:"name"`
	Traits    []APIEntry `json:"traits,omitempty"`
	Errors    []APIEntry `json:"errors,omitempty"`
	Types     []APIEntry `json:"types,omitempty"`
	TypeAlias []APIEntry `json:"typealias,omitempty"`
	Values    []APIEntry `json:"values,omitempty"`
	Misc      []APIEntry `json:"misc,omitempty"`
}

// APIEntry describes one detailed API symbol.
type APIEntry struct {
	Name           string         `json:"name"`
	Docstring      EvidenceString `json:"docstring"`
	Signature      string         `json:"signature"`
	PlainSignature EvidenceString `json:"plain_signature"`
	Loc            EvidenceObject `json:"loc"`
	Methods        []APIEntry     `json:"methods,omitempty"`
	Impls          []APIEntry     `json:"impls,omitempty"`
}

// PackageRelPath derives the asset path for a package under its module.
func PackageRelPath(module, packagePath string) (string, error) {
	if module == "" || packagePath == "" {
		return "", fmt.Errorf("module and package path are required")
	}
	if packagePath == module {
		return "", nil
	}
	prefix := module + "/"
	if !strings.HasPrefix(packagePath, prefix) {
		return "", fmt.Errorf("package path %q is not under module %q", packagePath, module)
	}
	return strings.TrimPrefix(packagePath, prefix), nil
}
