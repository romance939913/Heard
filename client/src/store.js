import { configureStore } from '@reduxjs/toolkit';
import sessionReducer from './reducers/sessionReducer';


const store = configureStore({
    reducer: {
        session: sessionReducer,
    },
    middleware: (getDefaultMiddleware) => getDefaultMiddleware(),
});

export default store;