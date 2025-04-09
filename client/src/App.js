import Login from './components/userAuth/login';
import Signin from './components/userAuth/register';
import Feed from './components/feed/feed';
import Layout from './Layout';
import { Routes, Route } from 'react-router-dom';
import RequireAuth from './components/RequireAuth/RequireAuth';

function App() {
  return (
      <Routes>
        <Route element={<Layout />}>
          {/*public routes*/}
          <Route element={<Login />} path="/login" />
          <Route element={<Signin />} path="/register" />

          {/*protected routes*/}
          <Route element={<RequireAuth />}>
            <Route element={<Feed />} path="/" />
          </Route>

          {/*Catch all*/}
        </Route>
      </Routes>
  )
}

export default App;