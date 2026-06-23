package model

import "testing"

func TestParseModuleCoordinate(t *testing.T) {
	got, err := ParseModuleCoordinate("moonbitlang/core")
	if err != nil {
		t.Fatalf("ParseModuleCoordinate: %v", err)
	}
	if got.Owner != "moonbitlang" || got.Name != "core" || got.String() != "moonbitlang/core" {
		t.Fatalf("unexpected coordinate: %#v", got)
	}

	for _, input := range []string{"", "core", "a/b/c", "/core", "moonbitlang/"} {
		if _, err := ParseModuleCoordinate(input); err == nil {
			t.Fatalf("ParseModuleCoordinate(%q) should fail", input)
		}
	}
}

func TestPackageRelPath(t *testing.T) {
	got, err := PackageRelPath("moonbitlang/core", "moonbitlang/core/argparse")
	if err != nil {
		t.Fatalf("PackageRelPath: %v", err)
	}
	if got != "argparse" {
		t.Fatalf("PackageRelPath = %q, want argparse", got)
	}

	root, err := PackageRelPath("moonbitlang/core", "moonbitlang/core")
	if err != nil {
		t.Fatalf("root PackageRelPath: %v", err)
	}
	if root != "" {
		t.Fatalf("root PackageRelPath = %q, want empty", root)
	}

	if _, err := PackageRelPath("moonbitlang/core", "other/core/argparse"); err == nil {
		t.Fatal("PackageRelPath should reject packages outside module")
	}
}

func TestAdoptionConclusionValidation(t *testing.T) {
	if _, err := ParseAdoptionConclusion(string(ConclusionContinueVerification)); err != nil {
		t.Fatalf("ParseAdoptionConclusion: %v", err)
	}
	if _, err := ParseAdoptionConclusion("maybe"); err == nil {
		t.Fatal("unknown adoption conclusion should fail")
	}
}

func TestRunwasmCoordinate(t *testing.T) {
	got := RunwasmCoordinate(SkillEntry{Module: "Yoorkin/cowsay", Version: "0.1.0", Package: "cowsay"})
	if got != "Yoorkin/cowsay/cowsay@0.1.0" {
		t.Fatalf("RunwasmCoordinate = %q", got)
	}
}

func FuzzParseModuleCoordinate(f *testing.F) {
	f.Add("moonbitlang/core")
	f.Add("not-a-coordinate")
	f.Fuzz(func(t *testing.T, value string) {
		coord, err := ParseModuleCoordinate(value)
		if err != nil {
			return
		}
		if coord.String() != value {
			t.Fatalf("round trip coordinate = %q, want %q", coord.String(), value)
		}
	})
}

func FuzzPackageRelPath(f *testing.F) {
	f.Add("moonbitlang/core", "moonbitlang/core/argparse")
	f.Fuzz(func(t *testing.T, module, pkg string) {
		rel, err := PackageRelPath(module, pkg)
		if err != nil {
			return
		}
		if rel != "" && len(rel) >= len(pkg) {
			t.Fatalf("relpath %q should be shorter than package %q", rel, pkg)
		}
	})
}
