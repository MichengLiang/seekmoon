package model

// ProjectIdentity identifies the local MoonBit project under inspection.
type ProjectIdentity struct {
	Root   string `json:"root"`
	Module string `json:"module,omitempty"`
}

// ProjectContext captures Moon project configuration evidence.
type ProjectContext struct {
	Identity             ProjectIdentity `json:"identity"`
	ModuleConfig         EvidenceObject  `json:"module_config"`
	PackageConfig        EvidenceObject  `json:"package_config"`
	DeclaredTarget       EvidenceString  `json:"declared_target"`
	ExistingDependencies EvidenceObject  `json:"existing_dependencies"`
}
