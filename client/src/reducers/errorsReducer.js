import { createSlice } from "@reduxjs/toolkit";

const errorsSlice = createSlice({
    name: 'errors',
    initialState: {
        error: null
    },
    reducers: {
        errorReducer: (state, action) => {
            state.error = action.payload
        },
        removeErrorReducer: (state) => {
            state.error = null
        },
    }
});

export const { errorReducer, removeErrorReducer } = errorsSlice.actions;
export default errorsSlice.reducer;