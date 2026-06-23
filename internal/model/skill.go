package model

type SkillEntry struct {
	Module       string         `json:"module"`
	Author       string         `json:"author"`
	AuthorAvatar EvidenceString `json:"author_avatar,omitempty"`
	Version      string         `json:"version"`
	Package      string         `json:"package"`
	Name         string         `json:"name"`
	DetailURL    string         `json:"detail_url"`
	WasmURL      string         `json:"wasm_url"`
	ChecksumURL  string         `json:"checksum_url"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Repository   EvidenceString `json:"repository"`
	CreatedAt    string         `json:"created_at"`
}

type SkillProfile struct {
	Entry             SkillEntry     `json:"entry"`
	SkillMD           EvidenceString `json:"skill_md"`
	WasmAsset         EvidenceObject `json:"wasm_asset"`
	ChecksumAsset     EvidenceObject `json:"checksum_asset"`
	RunwasmCoordinate EvidenceString `json:"runwasm_coordinate"`
}

func RunwasmCoordinate(entry SkillEntry) string {
	if entry.Package == "" {
		return entry.Module + "@" + entry.Version
	}
	return entry.Module + "/" + entry.Package + "@" + entry.Version
}
