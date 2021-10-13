package attachments

import "net/url"

/*type Interface interface {
	Attachments() AttachmentClient
}*/

type AttachmentListOptions struct {
	AttachedToID string
}

func (a *AttachmentListOptions) MarshalQueryParameters() (url.Values, error) {
	ret := url.Values{}
	ret.Set("attachedToId", a.AttachedToID)
	return ret, nil
}

func (a *AttachmentListOptions) UnmarshalQueryParameters(values url.Values) error {
	a.AttachedToID = values.Get("attachedToId")
	return nil
}
