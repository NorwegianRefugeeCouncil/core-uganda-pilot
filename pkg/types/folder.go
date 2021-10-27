package types

type Folder struct {
	ID         string `json:"id"`
	DatabaseID string `json:"databaseId"`
	ParentID   string `json:"parentId,omitempty"`
	Name       string `json:"name"`
}

type FolderList struct {
	Items []*Folder `json:"items"`
}
