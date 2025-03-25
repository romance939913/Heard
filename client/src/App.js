import React from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { connect } from 'react-redux';
import Login from './components/userAuth/login';
import Signin from './components/userAuth/signin';
import Feed from './components/feed/feed';
import ProtectedRoutes from './protectedRoutes';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Login/>} path="/login"/>
        <Route element={<Signin/>} path="/signin"/>
        <Route element={<ProtectedRoutes/>}>
          <Route element={<Feed />} path="/" />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

const mapStateToProps = state => ({
  loggedin: state.session.id
})

export default connect(mapStateToProps, null)(App);