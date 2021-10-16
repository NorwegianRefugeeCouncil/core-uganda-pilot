import {createAsyncThunk, createEntityAdapter, createSlice, EntityState} from "@reduxjs/toolkit";
import {IdentityProvider} from "../types/types";
import {RootState} from "../app/store";
import {defaultClient, IdentityProviderListRequest, IdentityProviderListResponse} from "../data/client";

export interface IdentityProviderState extends EntityState<IdentityProvider> {

}

const adapter = createEntityAdapter<IdentityProvider>({
    // Assume IDs are stored in a field other than `book.id`
    selectId: (identityProvider) => identityProvider.id,
    // Keep the "all IDs" array sorted based on book titles
    sortComparer: (a, b) => a.id.localeCompare(b.id),
})

export const fetchIdentityProviders = createAsyncThunk<IdentityProviderListResponse, IdentityProviderListRequest>(
    'identityProviders/fetch',
    async (arg, thunkAPI) => {
        try {
            const response = await defaultClient.listIdentityProviders({organizationId: arg.organizationId})
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

export const identityProvidersSlice = createSlice({
    name: 'identityProviders',
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
        builder.addCase(fetchIdentityProviders.pending, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = true
        })
        builder.addCase(fetchIdentityProviders.rejected, (state, action) => {
            state.fetchSuccess = false
            state.fetchPending = false
            state.fetchError = action.payload
        })
        builder.addCase(fetchIdentityProviders.fulfilled, (state, action) => {
            state.fetchSuccess = true
            state.fetchPending = false
            state.fetchError = undefined
            if (action.payload.response?.items) {
                adapter.setAll(state, action.payload.response.items)
            }
        })
    }
})

export const identityProviderActions = identityProvidersSlice.actions

let selectors = adapter.getSelectors();
let globalSelectors = adapter.getSelectors<RootState>(state => state.identityProviders);

function selectForOrganization(state: IdentityProviderState, organizationId: string) {
    return selectors.selectAll(state).filter(i => i.organizationId === organizationId)
}

export const identityProviderSelectors = {
    ...selectors,
    selectForOrganization
}

export const identityProviderGlobalSelectors = {
    ...globalSelectors,
    selectForOrganization: (state: RootState, organizationId: string) => selectForOrganization(state.identityProviders, organizationId)
}

export default identityProvidersSlice.reducer

