import React from 'react';
import { Outlet, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';

const ProtectedRoute = ({ children }) => {
    const isAuthenticated = useSelector((state) => state.session.isAuthenticated);
    console.log(`The user is auhenticated: ${isAuthenticated}`)
    return isAuthenticated ? children : <Navigate to="/login"/>
} 


export default ProtectedRoute;
