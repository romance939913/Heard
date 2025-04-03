import React, { Fragment } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Login from './components/userAuth/login';
import Signin from './components/userAuth/signin';
import Feed from './components/feed/feed';
import ProtectedRoute from './protectedRoutes';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Fragment>
          <Route element={<Login />} path="/login" />
          <Route element={<Signin />} path="/signin" />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Feed />
              </ProtectedRoute>
            }
          />
        </Fragment>
      </Routes>
    </BrowserRouter>
  )
}

export default App;