import {createAsyncThunk, createEntityAdapter, createSlice} from "@reduxjs/toolkit";
import {Organization} from "../types/types";
import {RootState} from "../app/store";
import {defaultClient, OrganizationListResponse} from "../data/client";

const adapter = createEntityAdapter<Organization>({
    // Assume IDs are stored in a field other than `book.id`
    selectId: (organization) => organization.id,
    // Keep the "all IDs" array sorted based on book titles
    sortComparer: (a, b) => a.name.localeCompare(b.name),
})

export const fetchOrganizations = createAsyncThunk<OrganizationListResponse>(
    'organizations/fetch',
    async (_, thunkAPI) => {
        try {
            const response = await defaultClient.listOrganizations()
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

export const organizationsSlice = createSlice({
    name: 'organizations',
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
        builder.addCase(fetchOrganizations.pending, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = true
        })
        builder.addCase(fetchOrganizations.rejected, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = false
            state.fetchError = action.payload
        })
        builder.addCase(fetchOrganizations.fulfilled, (state, action) => {
            state.fetchSuccess = true
            state.fetchPending = false
            state.fetchError = undefined
            if (action.payload.response?.items) {
                adapter.setAll(state, action.payload.response.items)
            }
        })
    }
})

export const organizationActions = organizationsSlice.actions
export const organizationSelectors = adapter.getSelectors()
export const organizationGlobalSelectors = adapter.getSelectors<RootState>(state => state.organizations)
export default organizationsSlice.reducer

