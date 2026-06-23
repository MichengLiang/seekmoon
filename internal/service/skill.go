package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/yumiaura/seekmoon/internal/model"
	"github.com/yumiaura/seekmoon/internal/source"
	"github.com/yumiaura/seekmoon/internal/store"
)

type SkillFlow struct {
	Skills   source.SkillsClient
	Sessions store.SessionStore
}

func (s SkillFlow) Search(ctx context.Context, input SkillSearchInput) ([]model.SkillEntry, error) {
	sourceResult := s.Skills.FetchSkills(ctx)
	if sourceResult.Status != model.StatePresent || sourceResult.Value == nil {
		return nil, sourceFailure("skill search", eraseSourceValue(sourceResult), "retry skill search")
	}
	query := strings.ToLower(input.Query)
	var out []model.SkillEntry
	candidates := store.SessionCandidates{SessionID: defaultSessionID, SnapshotID: "skills", Candidates: map[int]model.CandidateRef{}}
	for _, entry := range *sourceResult.Value {
		if !skillMatches(entry, query) {
			continue
		}
		out = append(out, entry)
		candidates.Candidates[len(out)] = model.CandidateRef{Kind: "skill", Module: entry.Module, Version: entry.Version, Package: entry.Package, Name: entry.Name}
	}
	if err := s.Sessions.WriteCandidates(ctx, candidates); err != nil {
		return nil, err
	}
	return out, nil
}

func (s SkillFlow) View(ctx context.Context, input SkillViewInput) (model.SkillProfile, error) {
	ref, err := candidateFromRequest(ctx, s.Sessions, input.Entry)
	if err != nil {
		return model.SkillProfile{}, err
	}
	entryID := ref.Module
	if input.Entry.Module != "" {
		entryID = input.Entry.Module
	}
	sourceResult := s.Skills.FetchSkill(ctx, entryID)
	if sourceResult.Status != model.StatePresent || sourceResult.Value == nil {
		return model.SkillProfile{}, sourceFailure(entryID, eraseSourceValue(sourceResult), "run skill search or check the skill entry")
	}
	entry := *sourceResult.Value
	profile := model.SkillProfile{
		Entry:             entry,
		SkillMD:           model.Unknown[string](),
		WasmAsset:         s.Skills.FetchSkillAsset(ctx, entry.WasmURL),
		ChecksumAsset:     s.Skills.FetchSkillAsset(ctx, entry.ChecksumURL),
		RunwasmCoordinate: model.Derived(model.RunwasmCoordinate(entry), string(model.SourceDerived)),
	}
	if entry.DetailURL != "" {
		profile.SkillMD = model.Present(entry.DetailURL, string(model.SourceSkillsAPI))
	}
	return profile, nil
}

func skillMatches(entry model.SkillEntry, query string) bool {
	if query == "" {
		return true
	}
	haystack := strings.ToLower(strings.Join([]string{entry.Module, entry.Package, entry.Name, evidenceStringValue(entry.Repository)}, " "))
	for _, value := range entry.Metadata {
		haystack += " " + strings.ToLower(valueString(value))
	}
	return strings.Contains(haystack, query)
}

func valueString(value any) string {
	if value == nil {
		return ""
	}
	out := fmt.Sprint(value)
	out = strings.ReplaceAll(out, "\n", " ")
	out = strings.ReplaceAll(out, "\t", " ")
	return strings.TrimSpace(out)
}
