import {createAsyncThunk, createEntityAdapter, createSlice} from "@reduxjs/toolkit";
import {Folder} from "../types/types";
import {RootState} from "../app/store";
import {defaultClient, FolderListResponse} from "../data/client";

const adapter = createEntityAdapter<Folder>({
    // Assume IDs are stored in a field other than `book.id`
    selectId: (folder) => folder.id,
    // Keep the "all IDs" array sorted based on book titles
    sortComparer: (a, b) => a.name.localeCompare(b.name),
})

export const fetchFolders = createAsyncThunk<FolderListResponse>(
    'folders/fetch',
    async (_, thunkAPI) => {
        try {
            const response = await defaultClient.listFolders({})
            if (response.success) {
                return response
            } else {
                return thunkAPI.rejectWithValue(response)
            }
        } catch (err) {
            return thunkAPI.rejectWithValue(err)
        }
    }
)

export const foldersSlice = createSlice({
    name: 'folders',
    initialState: {
        ...adapter.getInitialState(),
        fetchPending: false,
        fetchError: undefined as any,
        fetchSuccess: true
    },
    reducers: {
        addOne: adapter.addOne,
        addMany: adapter.addMany,
        removeAll: adapter.removeAll,
        removeMany: adapter.removeMany,
        removeOne: adapter.removeOne,
        updateMany: adapter.updateMany,
        updateOne: adapter.updateOne,
        upsertOne: adapter.upsertOne,
        upsertMany: adapter.upsertMany,
        setOne: adapter.setOne,
        setMany: adapter.setMany,
        setAll: adapter.setAll,
    },
    extraReducers: builder => {
        builder.addCase(fetchFolders.pending, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = true
        })
        builder.addCase(fetchFolders.rejected, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = false
            state.fetchError = action.payload
        })
        builder.addCase(fetchFolders.fulfilled, (state, action) => {
            state.fetchSuccess = true
            state.fetchPending = false
            state.fetchError = undefined
            if (action.payload.response?.items) {
                adapter.setAll(state, action.payload.response.items)
            }
        })
    }
})

export const folderActions = foldersSlice.actions

export const folderSelectors = adapter.getSelectors()
export const folderGlobalSelectors = adapter.getSelectors<RootState>(state => state.folders)

export const selectByParentId = (s: RootState, parentId: string | undefined) => {
    return folderSelectors.selectAll(s.folders).filter(folder => folder.parentId === parentId)
}

export const selectDatabaseRootFolders = (s: RootState, databaseId: string) => {
    return folderSelectors.selectAll(s.folders).filter(folder => !folder.parentId && folder.databaseId === databaseId)
}

export const selectParent = (s: RootState, childId: string | undefined) => {
    if (!childId) {
        return undefined
    }
    const child = folderGlobalSelectors.selectById(s, childId)
    if (!child) {
        return undefined
    }
    return folderGlobalSelectors.selectById(s, child.parentId)
}

export const selectParents = (s: RootState, folderId: string | undefined, includeSelf = false): Folder[] => {
    let result: Folder[] = []
    if (!folderId) {
        return result
    }
    let walk = folderId
    while (true) {
        const parent = selectParent(s, walk)
        if (parent) {
            result.push(parent)
            walk = parent.id
        } else {
            break
        }
    }
    if (includeSelf) {
        const folder = folderGlobalSelectors.selectById(s, folderId)
        if (folder) {
            result = [folder, ...result]
        }
    }
    return result
}


export default foldersSlice.reducer

