package match

import (
	"sort"
	"strconv"
	"strings"
)

// Match finds the best MinGW build for the given architecture.
// targetArch is the MinGW arch label ("x86_64", "i686", "aarch64") —
// already mapped from Go runtime arch by the detect package.
func Match(targetArch string, builds []Build, rules Rules, overrides Overrides) MatchResult {
	if overrides.Arch != "" {
		targetArch = overrides.Arch
	}

	// Filter by architecture.
	var candidates []Build
	for _, b := range builds {
		if b.Arch == targetArch {
			candidates = append(candidates, b)
		}
	}

	if len(candidates) == 0 {
		return MatchResult{
			Explanation: []DimensionChoice{
				{Dimension: "arch", Choice: targetArch, Reason: "no builds found for this architecture", Manual: overrides.Arch != ""},
			},
		}
	}

	// Score each candidate.
	type scored struct {
		Build Build
		Score int
	}
	var scoredCandidates []scored
	for _, b := range candidates {
		s := 0
		s += scoreDim(b.Thread, prefList(overrides.Thread, rules.ThreadPreference, "thread"))
		s += scoreDim(b.Exception, prefList(overrides.Exception, rules.ExceptionPreference[targetArch], "exception"))
		s += scoreDim(b.CRT, prefList(overrides.CRT, rules.CRTPreference, "crt"))
		scoredCandidates = append(scoredCandidates, scored{Build: b, Score: s})
	}

	// Sort by score desc, then GCC version desc, then source priority desc (tiebreaker).
	sort.Slice(scoredCandidates, func(i, j int) bool {
		if scoredCandidates[i].Score != scoredCandidates[j].Score {
			return scoredCandidates[i].Score > scoredCandidates[j].Score
		}
		if c := compareGCC(scoredCandidates[i].Build.GCC, scoredCandidates[j].Build.GCC); c != 0 {
			return c > 0
		}
		return scoredCandidates[i].Build.Priority > scoredCandidates[j].Build.Priority
	})

	best := scoredCandidates[0]
	var alternatives []Build
	for i := 1; i < len(scoredCandidates); i++ {
		alternatives = append(alternatives, scoredCandidates[i].Build)
	}

	return MatchResult{
		Build:        best.Build,
		Alternatives: alternatives,
		Explanation:  buildExplanation(targetArch, best.Build, overrides, rules),
		Score:        best.Score,
	}
}

// prefList returns the effective preference list for a dimension.
// If the user provided an override, it becomes the only option (and scores 3).
func prefList(override string, defaults []string, dim string) []string {
	if override != "" {
		return []string{override}
	}
	return defaults
}

// scoreDim returns points for a build's value based on its position in the preference list.
// First preference = 3 points, second = 2, third = 1, not found = 0.
func scoreDim(value string, prefs []string) int {
	for i, p := range prefs {
		if value == p {
			return 3 - i
		}
	}
	return 0
}

// compareGCC compares two GCC version strings like "14.2.0".
// Returns positive if a > b, negative if a < b, 0 if equal.
func compareGCC(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		ai, _ := strconv.Atoi(aParts[i])
		bi, _ := strconv.Atoi(bParts[i])
		if ai != bi {
			return ai - bi
		}
	}
	return len(aParts) - len(bParts)
}

// reasonForArch returns a human-readable reason for the arch choice.
func reasonForArch(arch string) string {
	switch arch {
	case "x86_64":
		return "your CPU is 64-bit"
	case "i686":
		return "your CPU is 32-bit"
	case "aarch64":
		return "your CPU is ARM64"
	default:
		return "detected architecture"
	}
}

// reasonForThread returns a human-readable reason for the thread model choice.
func reasonForThread(model string) string {
	switch model {
	case "posix":
		return "best C++11 std::thread support, wider compatibility"
	case "win32":
		return "native Windows threading, lighter weight"
	default:
		return "selected thread model"
	}
}

// reasonForException returns a human-readable reason for the exception handling choice.
func reasonForException(model string) string {
	switch model {
	case "seh":
		return "optimal exception handling performance on x86_64"
	case "dwarf":
		return "best exception handling for 32-bit targets"
	case "sjlj":
		return "broad compatibility across architectures"
	default:
		return "selected exception handling model"
	}
}

// reasonForCRT returns a human-readable reason for the CRT choice.
func reasonForCRT(crt string) string {
	switch crt {
	case "ucrt":
		return "modern Windows C runtime, recommended by Microsoft"
	case "msvcrt":
		return "compatible with older Windows versions (pre-Win10)"
	default:
		return "selected C runtime"
	}
}

func buildExplanation(targetArch string, best Build, overrides Overrides, rules Rules) []DimensionChoice {
	threadDim := DimensionChoice{
		Dimension: "thread",
		Choice:    best.Thread,
		Reason:    reasonForThread(best.Thread),
		Manual:    overrides.Thread != "",
	}
	if overrides.Thread != "" {
		threadDim.Reason = "[manual override] " + threadDim.Reason
	}

	excDim := DimensionChoice{
		Dimension: "exception",
		Choice:    best.Exception,
		Reason:    reasonForException(best.Exception),
		Manual:    overrides.Exception != "",
	}
	if overrides.Exception != "" {
		excDim.Reason = "[manual override] " + excDim.Reason
	}

	crtDim := DimensionChoice{
		Dimension: "crt",
		Choice:    best.CRT,
		Reason:    reasonForCRT(best.CRT),
		Manual:    overrides.CRT != "",
	}
	if overrides.CRT != "" {
		crtDim.Reason = "[manual override] " + crtDim.Reason
	}

	return []DimensionChoice{
		{Dimension: "arch", Choice: targetArch, Reason: reasonForArch(targetArch), Manual: overrides.Arch != ""},
		threadDim,
		excDim,
		crtDim,
	}
}
