package pack

type Report struct {
	Repaired int
	Warnings []Issue
	Errors   []Issue
}

func (r *Report) merge(other Report) {
	r.Repaired += other.Repaired
	r.Warnings = append(r.Warnings, other.Warnings...)
	r.Errors = append(r.Errors, other.Errors...)
}

type Issue struct {
	Message string

	Level IssueLevel

	Kind IssueKind

	Path string

	// Stable identifier for the entity involved (question / pack)
	Ref string

	// Optional extra context (free-form)
	Meta map[string]any
}

type IssueLevel int

// Simple. No “info”, no “fatal”. Keep it binary.
const (
	IssueWarning IssueLevel = iota
	IssueError
)

type IssueKind int

const (
	IssueMissingID IssueKind = iota + 1
	IssueDuplicateID
	IssueMissingField
	IssueInvalidDifficulty
	IssueInvalidAnswer
	IssueInvalidFormat
)

func NewIssue(
	level IssueLevel,
	kind IssueKind,
	message string,
	path string,
	ref string,
) Issue {
	return Issue{
		Level:   level,
		Kind:    kind,
		Message: message,
		Path:    path,
		Ref:     ref,
	}
}

func NewError(kind IssueKind, message, path, ref string) Issue {
	return NewIssue(IssueError, kind, message, path, ref)
}

func NewWarning(kind IssueKind, message, path, ref string) Issue {
	return NewIssue(IssueWarning, kind, message, path, ref)
}
