export type Folder = {
  id: string;
  databaseId: string;
  parentId: string;
  name: string;
};

export type FolderList = {
  items: Folder[];
};
