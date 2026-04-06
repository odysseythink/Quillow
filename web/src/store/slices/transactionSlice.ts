import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getTransactions } from '../../api/transactions';
import { TransactionGroup } from '../../types/models';
import { PaginationMeta } from '../../types/api';

interface TransactionState {
  items: TransactionGroup[];
  pagination: PaginationMeta | null;
  loading: boolean;
  error: string | null;
}

const initialState: TransactionState = {
  items: [],
  pagination: null,
  loading: false,
  error: null,
};

export const fetchTransactions = createAsyncThunk(
  'transactions/fetch',
  async ({ page, limit, type, start, end }: { page?: number; limit?: number; type?: string; start?: string; end?: string }) => {
    const res = await getTransactions(page, limit, type, start, end);
    return res.data;
  }
);

const transactionSlice = createSlice({
  name: 'transactions',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchTransactions.pending, (state) => { state.loading = true; })
      .addCase(fetchTransactions.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload.data.map((r) => ({ id: r.id, ...r.attributes }));
        state.pagination = action.payload.meta.pagination;
      })
      .addCase(fetchTransactions.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch transactions';
      });
  },
});

export default transactionSlice.reducer;
