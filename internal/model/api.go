package model

import (
	"fmt"
	"strings"
)

type ModuleIndexTree struct {
	Name    string            `json:"name"`
	Package *PackageSummary   `json:"package,omitempty"`
	Childs  []ModuleIndexTree `json:"childs,omitempty"`
}

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

type APISummary struct {
	Name string `json:"name"`
}

type TypeSummary struct {
	Name    string       `json:"name"`
	Methods []APISummary `json:"methods,omitempty"`
}

type PackageData struct {
	Name      string     `json:"name"`
	Traits    []APIEntry `json:"traits,omitempty"`
	Errors    []APIEntry `json:"errors,omitempty"`
	Types     []APIEntry `json:"types,omitempty"`
	TypeAlias []APIEntry `json:"typealias,omitempty"`
	Values    []APIEntry `json:"values,omitempty"`
	Misc      []APIEntry `json:"misc,omitempty"`
}

type APIEntry struct {
	Name           string         `json:"name"`
	Docstring      EvidenceString `json:"docstring"`
	Signature      string         `json:"signature"`
	PlainSignature EvidenceString `json:"plain_signature"`
	Loc            EvidenceObject `json:"loc"`
	Methods        []APIEntry     `json:"methods,omitempty"`
	Impls          []APIEntry     `json:"impls,omitempty"`
}

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
