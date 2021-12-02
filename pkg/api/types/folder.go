package types

// Folder is a namespace for Forms.
// Folder can be used to group FormDefinitions together.
type Folder struct {
	// ID represents the ID of the Folder
	ID string `json:"id"`
	// DatabaseID represents the ID of the Database
	DatabaseID string `json:"databaseId"`
	// ParentID represents the ID of the parent Folder
	// If the ParentID is empty, this means that the Folder exists
	// at the root of the Database
	ParentID string `json:"parentId,omitempty"`
	// Name of the folder
	Name string `json:"name"`
}

// FolderList represents a list of Folder
type FolderList struct {
	Items []*Folder `json:"items"`
}

// HasFolderWithID returns whether the FolderList contains a
// Folder with the given Folder.ID
func (f *FolderList) HasFolderWithID(id string) bool {
	for _, item := range f.Items {
		if item.ID == id {
			return true
		}
	}
	return false
}
