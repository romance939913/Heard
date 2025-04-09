import React, { useState, useContext } from 'react';
import { AuthContext } from '../../context/authProvider';
import { login } from '../../util/sessionUtil';
import { Link, useNavigate, useLocation } from 'react-router-dom';

function Login() {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const { auth, setAuth } = useContext(AuthContext);
  const navigate = useNavigate();
  const location = useLocation();
  const from = location.state?.from?.pathname || "/";

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prevData => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    login(formData)
      .then((response) => {
        const email = response?.data?.email;
        const id = response?.data?.id;
        const username = response?.data?.username;
        const access_token = response?.data?.access_token;
        debugger
        setAuth({ email, id, username, access_token});
        navigate(from, { replace: true });
      })
      .catch((err) => {
        console.log(err)
      })
  }

  return (
    <div>
      <h1>this is the login page</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          placeholder="email"
        />
        <input
          type="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          placeholder="password"
        />
        <button type="submit">login</button>
      </form>
    </div>
  );
}

export default Login;
