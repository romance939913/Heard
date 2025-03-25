import axios from 'axios';

export const register = (formUser) => (
    axios.post('/api/register', {
        'username': formUser.username,
        'email': formUser.email,
        'password': formUser.password
    })
);

export const login = (formUser) => (
    axios.post('/api/login', {
        'email': formUser.email,
        'password': formUser.password
    })
);

export const logout = () => (
    axios.post('/api/logout')
);