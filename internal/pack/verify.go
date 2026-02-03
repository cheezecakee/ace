package pack

import (
	"fmt"

	"github.com/cheezecakee/ace/internal/engine"
)

// Verify runs a quick verification throughout
// the packs to see if everything is correct

type Verify interface {
	Verify() Report
}

// Repair runs a repair if verification returns
// any errors
type Repair interface {
	Repair(issue Issue) Report
}

type VerifyRepair interface {
	Verify
	Repair
}

func (q *RawChoiceQuestion) Verify() Report {
	var report Report

	if q.ID == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingID,
			"Missing question ID",
			"choice",
			"",
		))
	}

	if q.Prompt == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing prompt",
			"choice",
			q.ID,
		))
	}

	if engine.ParseDifficulty(q.Difficulty) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueInvalidDifficulty,
			fmt.Sprintf("Invalid difficulty: %s", q.Difficulty),
			"choice",
			q.ID,
		))
	}

	if len(q.Options) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"No options provided",
			"choice",
			q.ID,
		))
	}

	if q.Answer < 0 || q.Answer >= len(q.Options) {
		report.Errors = append(report.Errors, NewError(
			IssueInvalidAnswer,
			fmt.Sprintf("Answer index %d out of range (0-%d)", q.Answer, len(q.Options)-1),
			"choice",
			q.ID,
		))
	}

	return report
}

func (q *RawChoiceQuestion) Repair(issue Issue) Report {
	var report Report

	return report
}

func (q *RawMultiQuestion) Verify() Report {
	var report Report

	if q.ID == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingID,
			"Missing question ID",
			"multi",
			"",
		))
	}

	if q.Prompt == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing prompt",
			"multi",
			q.ID,
		))
	}

	if engine.ParseDifficulty(q.Difficulty) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueInvalidDifficulty,
			fmt.Sprintf("Invalid difficulty: %s", q.Difficulty),
			"multi",
			q.ID,
		))
	}

	if len(q.Options) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"No options provided",
			"multi",
			q.ID,
		))
	}

	if len(q.Answer) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"No correct answers specified",
			"multi",
			q.ID,
		))
	}

	// Check all answer indices are valid
	for _, ans := range q.Answer {
		if ans < 0 || ans >= len(q.Options) {
			report.Errors = append(report.Errors, NewError(
				IssueInvalidAnswer,
				fmt.Sprintf("Answer index %d out of range (0-%d)", ans, len(q.Options)-1),
				"multi",
				q.ID,
			))
			break
		}
	}

	return report
}

func (q *RawMultiQuestion) Repair(issue Issue) Report {
	var report Report

	return report
}

func (q *RawBoolQuestion) Verify() Report {
	var report Report

	if q.ID == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingID,
			"Missing question ID",
			"bool",
			"",
		))
	}

	if q.Prompt == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing prompt",
			"bool",
			q.ID,
		))
	}

	if engine.ParseDifficulty(q.Difficulty) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueInvalidDifficulty,
			fmt.Sprintf("Invalid difficulty: %s", q.Difficulty),
			"bool",
			q.ID,
		))
	}

	return report
}

func (q *RawBoolQuestion) Repair(issue Issue) Report {
	var report Report

	return report
}

func (q *RawTextQuestion) Verify() Report {
	var report Report

	if q.ID == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingID,
			"Missing question ID",
			"text",
			"",
		))
	}

	if q.Prompt == "" {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Missing prompt",
			"text",
			q.ID,
		))
	}

	if engine.ParseDifficulty(q.Difficulty) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueInvalidDifficulty,
			fmt.Sprintf("Invalid difficulty: %s", q.Difficulty),
			"text",
			q.ID,
		))
	}

	if q.Expected == "" && len(q.Keywords) == 0 {
		report.Errors = append(report.Errors, NewError(
			IssueMissingField,
			"Must have either expected answer or keywords",
			"text",
			q.ID,
		))
	}

	return report
}

func (q *RawTextQuestion) Repair(issue Issue) Report {
	var report Report

	return report
}
