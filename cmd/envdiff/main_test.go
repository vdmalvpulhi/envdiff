package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "envdiff")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoDiff(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	cmd := exec.Command(bin, a, b)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "No differences") {
		t.Errorf("expected no-differences message, got: %q", string(out))
	}
}

func TestMain_WithDiff_ExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\nONLY_A=1\n")
	b := writeTempEnv(t, "FOO=different\nONLY_B=2\n")
	cmd := exec.Command(bin, a, b)
	out, _ := cmd.CombinedOutput()
	outStr := string(out)
	if !strings.Contains(outStr, "CONFLICT") {
		t.Errorf("expected CONFLICT in output, got: %q", outStr)
	}
}

func TestMain_JSONFormat(t *testing.T) {
	bin := buildBinary(t)
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=baz\n")
	cmd := exec.Command(bin, "-format", "json", a, b)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "\"Key\"") {
		t.Errorf("expected JSON output with Key field, got: %q", string(out))
	}
}

func TestMain_MissingArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing args")
	}
}
