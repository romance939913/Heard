import { useContext } from "react";
import { AuthContext } from "../../context/authProvider";
import { useLocation, Navigate, Outlet } from "react-router-dom";

const RequireAuth = () => {
    const { auth } = useContext(AuthContext);
    const location = useLocation();

    return (
        auth?.user
            ? <Outlet />
            : <Navigate to="/login" state={{ from: location }} replace />
    )
}

export default RequireAuth