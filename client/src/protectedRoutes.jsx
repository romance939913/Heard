import React from 'react';
import { Outlet, Navigate } from 'react-router-dom';
import { connect } from 'react-redux';

const ProtectedRoutes = (props) => {
    return props.isAuthenticated ? <Outlet/> : <Navigate to="/login"/>
} 

const mapStateToProps = state => ({
    isAuthenticated: state.session.isAuthenticated
})

export default connect(mapStateToProps, null)(ProtectedRoutes);
