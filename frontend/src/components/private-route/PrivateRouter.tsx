import {useAuthStore} from "../../stores/authStore.ts";
import {Navigate, Outlet} from "react-router-dom";
import Loader from "../shared/loader/Loader.tsx";

interface PrivateRouteProps {
  allowedRoles: string[];
}

const PrivateRoute = ({ allowedRoles }: PrivateRouteProps) => {
  const { role, isLoading } = useAuthStore();

  if (isLoading) {
    return <Loader isLoading={isLoading} delay={300} />
  }

  if (role === 'guest') {
    return <Navigate to="/" replace />;
  }

  if (!role || !allowedRoles.includes(role)) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />
}

export default PrivateRoute;