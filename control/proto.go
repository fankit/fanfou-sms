package control


const (
	StateOk       State = "ok"
	StatePaused   State = "paused"
	StateAlerting State = "alerting"
	StatePending  State = "pending"
	StateNoData   State = "no_data"
)

type State string

type Body struct {
	Title       string                   `json:"title"`
	RuleID      int                      `json:"ruleId"`
	RuleName    string                   `json:"ruleName"`
	RuleURL     string                   `json:"ruleUrl"`
	State       State                    `json:"state"`
	ImageURL    string                   `json:"imageUrl"`
	Message     string                   `json:"message"`
	EvalMatches []map[string]interface{} `json:"evalMatches"`
}
