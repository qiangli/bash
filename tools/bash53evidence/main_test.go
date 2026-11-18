package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func fakeSuite(t *testing.T, verdicts []verdict, artifact ...string) []byte {
	t.Helper()
	s := suite{SchemaVersion: 1, Verdicts: verdicts}
	s.Infrastructure.Status = "ok"
	s.Context.BashPPMode = "off"
	s.Context.BashPPVerified = true
	if len(artifact) > 0 {
		s.Context.Artifact = artifact[0]
	}
	for _, v := range verdicts {
		switch v.Verdict {
		case "passed":
			s.Summary.Passed++
		case "failed":
			s.Summary.Failed++
		case "timed_out":
			s.Summary.TimedOut++
		case "skipped":
			s.Summary.Skipped++
		}
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestEvidenceWriterAndComparison(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "suite.json")
	output := filepath.Join(dir, "evidence.json")
	args := []string{"-input", input, "-output", output, "-bashy-commit", "abc", "-sh-pin", "def", "-image-id", "sha256:123", "-host-os", "Darwin", "-uid", "1000", "-mode", "0"}
	old := fakeSuite(t, []verdict{{"a", "failed"}, {"b", "passed"}, {"c", "timed_out"}, {"d", "skipped"}})
	if err := os.WriteFile(input, old, 0600); err != nil {
		t.Fatal(err)
	}
	if err := run(args); err != nil {
		t.Fatal(err)
	}
	var prior evidence
	b, _ := os.ReadFile(output)
	if err := json.Unmarshal(b, &prior); err != nil {
		t.Fatal(err)
	}
	if prior.BashyCommit != "abc" || prior.SHPin != "def" || prior.ImageDigest != "sha256:123" || prior.TesteeUID != 1000 || prior.Totals.TimedOut != 1 || !reflect.DeepEqual(prior.Failing, []string{"a", "c"}) || !reflect.DeepEqual(prior.Skips, []string{"d"}) {
		t.Fatalf("bad evidence: %+v", prior)
	}
	baseline := filepath.Join(dir, "baseline.json")
	if err := os.Rename(output, baseline); err != nil {
		t.Fatal(err)
	}
	next := fakeSuite(t, []verdict{{"a", "passed"}, {"b", "failed"}, {"c", "failed"}, {"d", "skipped"}})
	if err := os.WriteFile(input, next, 0600); err != nil {
		t.Fatal(err)
	}
	if err := run(append(args, "-baseline", baseline)); err == nil {
		t.Fatal("NEW failure should fail")
	}
	b, _ = os.ReadFile(output)
	var got evidence
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Comparison, &comparison{New: []string{"b"}, Recovered: []string{"a"}, Unchanged: []string{"c"}}) {
		t.Fatalf("bad comparison: %+v", got.Comparison)
	}
	if err := os.WriteFile(input, old, 0600); err != nil {
		t.Fatal(err)
	}
	if err := run(append(args, "-baseline", baseline)); err != nil {
		t.Fatalf("unchanged baseline: %v", err)
	}
}

func TestEvidenceRejectsMalformedAndIncomplete(t *testing.T) {
	e := evidence{BashPPMode: "0", TesteeUID: 1000}
	for _, tc := range []struct {
		name     string
		verdicts []verdict
		bad      bool
	}{
		{"duplicate", []verdict{{"a", "passed"}, {"a", "failed"}}, true},
		{"unknown", []verdict{{"a", "mystery"}}, true},
		{"good", []verdict{{"a", "passed"}}, false},
	} {
		_, err := collect(fakeSuite(t, tc.verdicts), e)
		if (err != nil) != tc.bad {
			t.Errorf("%s: %v", tc.name, err)
		}
	}
	dir := t.TempDir()
	input := filepath.Join(dir, "in.json")
	output := filepath.Join(dir, "out.json")
	var on suite
	json.Unmarshal(fakeSuite(t, []verdict{{"a", "passed"}}), &on)
	on.Context.BashPPMode = "on"
	onRaw, _ := json.Marshal(on)
	os.WriteFile(input, onRaw, 0600)
	if err := run([]string{"-input", input, "-output", output, "-bashy-commit", "abc", "-sh-pin", "def", "-image-id", "sha256:123", "-host-os", "Darwin", "-uid", "1000", "-mode", "1", "-gate"}); err == nil {
		t.Fatal("incomplete gate passed")
	}
}

func TestSplitArtifactEvidence(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "suite.json")
	output := filepath.Join(dir, "evidence-a.json")
	s := fakeSuite(t, []verdict{{"alpha", "passed"}}, "A")
	if err := os.WriteFile(input, s, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{
		"-input", input,
		"-output", output,
		"-artifact", "A",
		"-testee-sha256", "11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(args); err != nil {
		t.Fatalf("split artifact evidence failed: %v", err)
	}
	b, _ := os.ReadFile(output)
	var got evidence
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got.Artifact != "A" {
		t.Fatalf("artifact = %q, want A", got.Artifact)
	}
	if got.TesteeSHA256 != "11223344556677889900aabbccddeeff11223344556677889900aabbccddeeff" {
		t.Fatalf("testee_sha256 = %q", got.TesteeSHA256)
	}
	if got.SourceCommit != "c0ffee" {
		t.Fatalf("source_commit = %q, want c0ffee", got.SourceCommit)
	}
	if got.BashyCommit != "" || got.SHPin != "" {
		t.Fatalf("legacy fields mislabeled: bashy_commit=%q, sh_pin=%q", got.BashyCommit, got.SHPin)
	}
}

func TestSplitArtifactLowerPinRefusal(t *testing.T) {
	dir := t.TempDir()
	inputB := filepath.Join(dir, "suite-b.json")
	outputB := filepath.Join(dir, "evidence-b.json")
	var s suite
	json.Unmarshal(fakeSuite(t, []verdict{{"alpha", "passed"}}, "B"), &s)
	s.Context.BashPPMode = "on"
	sBytes, _ := json.Marshal(s)
	if err := os.WriteFile(inputB, sBytes, 0600); err != nil {
		t.Fatal(err)
	}

	baseArgs := []string{
		"-input", inputB,
		"-output", outputB,
		"-testee-sha256", "112233",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "1",
	}

	// Artifact B without lower pins must fail:
	argsB := append(append([]string{}, baseArgs...), "-artifact", "B")
	if err := run(argsB); err == nil || !strings.Contains(err.Error(), "requires lower module pins") {
		t.Fatalf("expected error for B without lower pins, got: %v", err)
	}

	// Artifact C without lower pins must fail:
	var sC suite
	json.Unmarshal(sBytes, &sC)
	sC.Context.Artifact = "C"
	sCBytes, _ := json.Marshal(sC)
	inputC := filepath.Join(dir, "suite-c.json")
	os.WriteFile(inputC, sCBytes, 0600)
	outputC := filepath.Join(dir, "evidence-c.json")
	argsC := []string{
		"-input", inputC,
		"-output", outputC,
		"-artifact", "C",
		"-testee-sha256", "112233",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "1",
	}
	if err := run(argsC); err == nil || !strings.Contains(err.Error(), "requires lower module pins") {
		t.Fatalf("expected error for C without lower pins, got: %v", err)
	}

	// Artifact B with lower pins succeeds:
	argsBWithPins := append(argsB, "-lower-pins", "github.com/foo/bar@v1.0.0,mvdan.cc/sh/v3@v3.10.0")
	if err := run(argsBWithPins); err != nil {
		t.Fatalf("B with lower pins failed: %v", err)
	}
	b, _ := os.ReadFile(outputB)
	var got evidence
	json.Unmarshal(b, &got)
	wantPins := []string{"github.com/foo/bar@v1.0.0", "mvdan.cc/sh/v3@v3.10.0"}
	if !reflect.DeepEqual(got.LowerModulePins, wantPins) {
		t.Fatalf("got lower module pins %v, want %v", got.LowerModulePins, wantPins)
	}

	// Artifact A without lower pins succeeds:
	inputA := filepath.Join(dir, "suite-a.json")
	outputA := filepath.Join(dir, "evidence-a.json")
	os.WriteFile(inputA, fakeSuite(t, []verdict{{"alpha", "passed"}}, "A"), 0600)
	argsA := []string{
		"-input", inputA,
		"-output", outputA,
		"-artifact", "A",
		"-testee-sha256", "112233",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsA); err != nil {
		t.Fatalf("A without lower pins should succeed, got: %v", err)
	}
}

func TestPayloadIdentity(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "suite.json")
	output := filepath.Join(dir, "evidence.json")
	os.WriteFile(input, fakeSuite(t, []verdict{{"alpha", "passed"}}, "A"), 0600)
	args := []string{
		"-input", input,
		"-output", output,
		"-artifact", "A",
		"-testee-sha256", "testee-hash-123",
		"-payload-sha256", "payload-hash-456",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(args); err != nil {
		t.Fatalf("run with payload-sha256 failed: %v", err)
	}
	b, _ := os.ReadFile(output)
	var got evidence
	json.Unmarshal(b, &got)
	if got.PayloadSHA256 != "payload-hash-456" {
		t.Fatalf("payload_sha256 = %q, want payload-hash-456", got.PayloadSHA256)
	}
	if got.TesteeSHA256 != "testee-hash-123" {
		t.Fatalf("testee_sha256 = %q, want testee-hash-123", got.TesteeSHA256)
	}

	// Without payload-sha256:
	argsNoPayload := []string{
		"-input", input,
		"-output", output,
		"-artifact", "A",
		"-testee-sha256", "testee-hash-123",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsNoPayload); err != nil {
		t.Fatalf("run without payload-sha256 failed: %v", err)
	}
	b, _ = os.ReadFile(output)
	var gotNoPayload evidence
	json.Unmarshal(b, &gotNoPayload)
	if gotNoPayload.PayloadSHA256 != "" {
		t.Fatalf("payload_sha256 = %q, want empty", gotNoPayload.PayloadSHA256)
	}
}

func TestLegacyToSplitComparison(t *testing.T) {
	legacy := evidence{
		SchemaVersion: 1,
		BashPPMode:    "0",
		Fixtures:      []string{"a", "b"},
		Failing:       []string{"a"},
	}
	// split A in same mode 0 and same inventory:
	splitA := evidence{
		SchemaVersion: 1,
		Artifact:      "A",
		BashPPMode:    "0",
		Fixtures:      []string{"a", "b"},
		Failing:       []string{"a"},
	}
	if err := compare(&splitA, legacy); err != nil {
		t.Fatalf("legacy baseline to split A comparison failed: %v", err)
	}

	// split B in same mode 1:
	legacyMode1 := evidence{
		SchemaVersion: 1,
		BashPPMode:    "1",
		Fixtures:      []string{"a", "b"},
		Failing:       []string{"a"},
	}
	splitB := evidence{
		SchemaVersion: 1,
		Artifact:      "B",
		BashPPMode:    "1",
		Fixtures:      []string{"a", "b"},
		Failing:       []string{"a"},
	}
	if err := compare(&splitB, legacyMode1); err != nil {
		t.Fatalf("legacy baseline to split B comparison failed: %v", err)
	}

	// Comparing different explicit artifacts (A vs B) must reject:
	explicitA := evidence{
		SchemaVersion: 1,
		Artifact:      "A",
		BashPPMode:    "0",
		Fixtures:      []string{"a", "b"},
	}
	if err := compare(&splitB, explicitA); err == nil || !strings.Contains(err.Error(), "baseline artifact differs") {
		t.Fatalf("explicit artifact A to B comparison should reject with 'baseline artifact differs', got: %v", err)
	}

	// Comparing same explicit artifact (A vs A) must succeed:
	anotherA := evidence{
		SchemaVersion: 1,
		Artifact:      "A",
		BashPPMode:    "0",
		Fixtures:      []string{"a", "b"},
	}
	if err := compare(&anotherA, explicitA); err != nil {
		t.Fatalf("explicit artifact A to A comparison failed: %v", err)
	}
}

func TestReportArtifactMismatch(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "suite.json")
	output := filepath.Join(dir, "evidence.json")
	os.WriteFile(input, fakeSuite(t, []verdict{{"alpha", "passed"}}, "A"), 0600)
	// Requested artifact B vs report artifact A:
	argsMismatch := []string{
		"-input", input,
		"-output", output,
		"-artifact", "B",
		"-testee-sha256", "112233",
		"-lower-pins", "pin1@v1",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsMismatch); err == nil || !strings.Contains(err.Error(), "does not match requested artifact") {
		t.Fatalf("expected artifact mismatch error, got: %v", err)
	}

	// Requested artifact A vs report with no artifact (legacy report):
	inputLegacy := filepath.Join(dir, "suite-legacy.json")
	os.WriteFile(inputLegacy, fakeSuite(t, []verdict{{"alpha", "passed"}}), 0600)
	argsLegacyMismatch := []string{
		"-input", inputLegacy,
		"-output", output,
		"-artifact", "A",
		"-testee-sha256", "112233",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsLegacyMismatch); err == nil || !strings.Contains(err.Error(), "does not match requested artifact") {
		t.Fatalf("expected legacy report vs requested artifact mismatch error, got: %v", err)
	}
}

func TestMissingExplicitArtifactFields(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "suite.json")
	output := filepath.Join(dir, "evidence.json")
	os.WriteFile(input, fakeSuite(t, []verdict{{"alpha", "passed"}}, "A"), 0600)

	// Missing testee sha256:
	argsNoSHA := []string{
		"-input", input,
		"-output", output,
		"-artifact", "A",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsNoSHA); err == nil {
		t.Fatal("expected error on missing testee-sha256")
	}

	// Missing source commit:
	argsNoCommit := []string{
		"-input", input,
		"-output", output,
		"-artifact", "A",
		"-testee-sha256", "112233",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsNoCommit); err == nil {
		t.Fatal("expected error on missing source-commit")
	}

	// Invalid artifact label:
	argsBadArt := []string{
		"-input", input,
		"-output", output,
		"-artifact", "Z",
		"-testee-sha256", "112233",
		"-source-commit", "c0ffee",
		"-image-id", "sha256:img123",
		"-host-os", "Linux",
		"-uid", "1000",
		"-mode", "0",
	}
	if err := run(argsBadArt); err == nil || !strings.Contains(err.Error(), "invalid artifact") {
		t.Fatalf("expected invalid artifact error, got: %v", err)
	}
}
