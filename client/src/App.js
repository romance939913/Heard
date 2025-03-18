import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Login from './components/userAuth/login';
import Signin from './components/userAuth/signin';
import Feed from './components/feed/feed';
import ProtectedRoutes from './util/protectedRoutes';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Login/>} path="/login"/>
        <Route element={<Signin/>} path="/signin"/>
        <Route element={<ProtectedRoutes/>}>
          <Route element={<Feed/>} path="/"/>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App