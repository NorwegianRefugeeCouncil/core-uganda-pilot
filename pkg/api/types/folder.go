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

func (f *FolderList) HasFolderWithID(id string) bool {
	for _, item := range f.Items {
		if item.ID == id {
			return true
		}
	}
	return false
}
