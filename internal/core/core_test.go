package core

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func TestNaturalLess(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		{"item_1.txt", "item_2.txt", true},
		{"item_2.txt", "item_10.txt", true},
		{"item_10.txt", "item_2.txt", false},
		{"a.txt", "b.txt", true},
		{"photo_001.jpg", "photo_002.jpg", true},
	}

	for _, tt := range tests {
		got := NaturalLess(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("NaturalLess(%q, %q) = %v; want %v", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestPlanExtChange(t *testing.T) {
	tmpDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tmpDir, "foto1.jpeg"), []byte("test"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "foto2.JPEG"), []byte("test"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "doc.pdf"), []byte("test"), 0644)

	plan, err := PlanExtChange(ExtChangeOptions{
		Directory:  tmpDir,
		CurrentExt: ".jpeg",
		NewExt:     ".jpg",
		Recursive:  false,
	})
	if err != nil {
		t.Fatalf("PlanExtChange failed: %v", err)
	}

	if plan.Total() != 2 {
		t.Fatalf("expected 2 files in plan, got %d", plan.Total())
	}

	for _, c := range plan.Changes {
		if filepath.Ext(c.NewName) != ".jpg" {
			t.Errorf("expected extension .jpg, got %s", filepath.Ext(c.NewName))
		}
	}
}

func TestPlanNumbering(t *testing.T) {
	tmpDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tmpDir, "nota-b.txt"), []byte("1"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "nota-a.txt"), []byte("2"), 0644)

	plan, err := PlanNumbering(NumbererOptions{
		Directory: tmpDir,
		Extension: ".txt",
		Start:     big.NewInt(1),
		Width:     3,
		InPlace:   true,
	})
	if err != nil {
		t.Fatalf("PlanNumbering failed: %v", err)
	}

	if plan.Total() != 2 {
		t.Fatalf("expected 2 files in plan, got %d", plan.Total())
	}

	// Orden natural: nota-a precede a nota-b
	if plan.Changes[0].OldName != "nota-a.txt" || plan.Changes[0].NewName != "001.txt" {
		t.Errorf("first change should be nota-a.txt -> 001.txt, got %s -> %s", plan.Changes[0].OldName, plan.Changes[0].NewName)
	}
	if plan.Changes[1].OldName != "nota-b.txt" || plan.Changes[1].NewName != "002.txt" {
		t.Errorf("second change should be nota-b.txt -> 002.txt, got %s -> %s", plan.Changes[1].OldName, plan.Changes[1].NewName)
	}
}

func TestPlanRename(t *testing.T) {
	tmpDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tmpDir, "IMG_001_old.jpg"), []byte("test"), 0644)

	plan, err := PlanRename(RenamerOptions{
		Directory: tmpDir,
		Search:    "_old",
		Replace:   "_new",
		Prefix:    "vacaciones_",
		Suffix:    "_hq",
	})
	if err != nil {
		t.Fatalf("PlanRename failed: %v", err)
	}

	if plan.Total() != 1 {
		t.Fatalf("expected 1 file in plan, got %d", plan.Total())
	}

	expected := "vacaciones_IMG_001_new_hq.jpg"
	if plan.Changes[0].NewName != expected {
		t.Errorf("expected %q, got %q", expected, plan.Changes[0].NewName)
	}
}

func TestConvertCase(t *testing.T) {
	tests := []struct {
		input    string
		style    CaseStyle
		expected string
	}{
		{"mi documento_lindo", CaseKebab, "mi-documento-lindo"},
		{"MiDocumentoLindo", CaseSnake, "mi_documento_lindo"},
		{"mi-documento-lindo", CaseCamel, "miDocumentoLindo"},
		{"mi_documento_lindo", CasePascal, "MiDocumentoLindo"},
	}

	for _, tt := range tests {
		got := ConvertCase(tt.input, tt.style)
		if got != tt.expected {
			t.Errorf("ConvertCase(%q, %q) = %q; want %q", tt.input, tt.style, got, tt.expected)
		}
	}
}

func TestExecutePlan(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	_ = os.WriteFile(file1, []byte("contenido"), 0644)

	plan := NewFilePlan(tmpDir)
	plan.Add(FileChange{
		SourcePath: file1,
		TargetPath: filepath.Join(tmpDir, "renamed.txt"),
		OldName:    "file1.txt",
		NewName:    "renamed.txt",
	})

	res, err := Execute(plan, ModeRename)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if res.Processed != 1 {
		t.Errorf("expected 1 processed, got %d", res.Processed)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "renamed.txt")); err != nil {
		t.Errorf("renamed file does not exist: %v", err)
	}
}
