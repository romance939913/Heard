import { createSlice } from "@reduxjs/toolkit";

const sessionSlice = createSlice({
    name: 'session',
    initialState: {
        isAuthenticated: null,
        user: null
    },
    reducers: {
        loginReducer: (state, action) => {
            state.isAuthenticated = action.payload;
            state.user = action.payload;
        },
        registerReducer: (state, action) => {
            state.isAuthenticated = action.payload;
            state.user = action.payload;
        },
        logoutReducer: (state) => {
            state.isAuthenticated = null;
            state.user = null;
        },
    }
});

export const { loginReducer, registerReducer, logoutReducer } = sessionSlice.actions;
export default sessionSlice.reducer;