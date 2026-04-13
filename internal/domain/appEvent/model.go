package appEvent

type PublishingEvent struct {
	Source     EventSource        `json:"source"`
	DetailType AppEventDetailType `json:"detailType"`
	Detail     string             `json:"detail"`
}
