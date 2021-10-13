package attachments

type Attachment struct {
	ID           string `json:"id" bson:"id"`
	AttachedToID string `json:"attachedToId" bson:"attachedToId"`
	Body         string `json:"body" bson:"body"`
}

type AttachmentList struct {
	Items []*Attachment `json:"items"`
}
