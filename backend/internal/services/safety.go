package services

import (
	"regexp"
	"strings"
)

// Result of a safety scan.
type SafetyResult struct {
	Status string // safe | review | blocked
	Reason string
}

// Block hard: anything indicating minors or non-consensual context.
var blockPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(child|children|kid|kids|minor|minors|underage|under\s*age|preteen|pre[-\s]?teen|teen(?:ie)?\s*sex|loli|shota|cp\b|baby|infant|toddler)\b`),
	regexp.MustCompile(`(?i)\b(school\s*girl|school\s*boy|elementary|middle\s*school|high\s*school|jr\.?\s*high|kindergarten)\b`),
	regexp.MustCompile(`(?i)\b(1[0-7]\s*(yo|y/o|years?\s*old))\b`),
	regexp.MustCompile(`(?i)\b(rape|raped|non[-\s]?consensual|forced|kidnap|abducted|trafficking|incest|bestiality|zoophilia|necro)\b`),
	regexp.MustCompile(`(?i)\b(hidden\s*cam(era)?|spy\s*cam|upskirt|voyeur\s*illegal|blackmail|revenge\s*porn|leaked\s*private)\b`),
	regexp.MustCompile(`(?i)\b(snuff|gore\s*porn)\b`),
}

// Mark for human review (suggestive but ambiguous).
var reviewPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(young|petite|tiny|small|first[-\s]?time|barely\s*legal|18\+?\s*only|just\s*turned\s*18)\b`),
	regexp.MustCompile(`(?i)\b(uniform|step[-\s]?(sister|brother|mom|dad|daughter|son))\b`),
	regexp.MustCompile(`(?i)\b(amateur\s*leaked|hidden|secret\s*recording)\b`),
}

// Scan inspects title, excerpt, content, and tags for unsafe terms.
func ScanSafety(title, excerpt, content string, tagsAndCats ...string) SafetyResult {
	hay := strings.ToLower(strings.Join(append([]string{title, excerpt, content}, tagsAndCats...), " \n "))

	for _, re := range blockPatterns {
		if m := re.FindString(hay); m != "" {
			return SafetyResult{Status: "blocked", Reason: "matched blocklist: " + m}
		}
	}
	for _, re := range reviewPatterns {
		if m := re.FindString(hay); m != "" {
			return SafetyResult{Status: "review", Reason: "matched reviewlist: " + m}
		}
	}
	return SafetyResult{Status: "safe"}
}
