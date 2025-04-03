import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { loginReducer } from '../../reducers/sessionReducer';
import { errorReducer, removeErrorReducer } from '../../reducers/errorsReducer';
import { login } from '../../util/sessionUtil';


function Login() {
  const dispatch = useDispatch()
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });

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
        dispatch(loginReducer(response.data))
      })
      .catch((e) => {
        dispatch(errorReducer(e))
        setTimeout(() => { removeErrorReducer() }, 3000);
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
