package session

type State int

const (
	NotStarted State = iota + 1
	Running
	Completed
	Failed
	TimeExpired
)

func (s State) String() string {
	switch s {
	case NotStarted:
		return "not started"
	case Running:
		return "running"
	case Completed:
		return "completed"
	case Failed:
		return "failed"
	case TimeExpired:
		return "time expired"
	default:
		return "unknown"
	}
}

func (s State) IsTerminal() bool {
	return s == Completed || s == Failed || s == TimeExpired
}
