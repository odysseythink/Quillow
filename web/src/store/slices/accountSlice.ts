import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getAccounts } from '../../api/accounts';
import { Account } from '../../types/models';
import { PaginationMeta } from '../../types/api';

interface AccountState {
  items: Account[];
  pagination: PaginationMeta | null;
  loading: boolean;
  error: string | null;
}

const initialState: AccountState = {
  items: [],
  pagination: null,
  loading: false,
  error: null,
};

export const fetchAccounts = createAsyncThunk(
  'accounts/fetch',
  async ({ page, limit, type }: { page?: number; limit?: number; type?: string }) => {
    const res = await getAccounts(page, limit, type);
    return res.data;
  }
);

const accountSlice = createSlice({
  name: 'accounts',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchAccounts.pending, (state) => { state.loading = true; })
      .addCase(fetchAccounts.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload.data.map((r) => ({ ...r.attributes, id: r.id }));
        state.pagination = action.payload.meta.pagination;
      })
      .addCase(fetchAccounts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch accounts';
      });
  },
});

export default accountSlice.reducer;
