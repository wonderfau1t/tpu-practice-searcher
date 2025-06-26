import {useAuthStore} from "../../stores/authStore.ts";
import {Navigate, Outlet, useLocation} from "react-router-dom";
import Loader from "../shared/loader/Loader.tsx";

const roleRedirects: Record<string, string> = {
  student: '/student/vacanciesList',
  moderator: '/moderator/vacanciesList',
  headHR: '/company/vacanciesList',
  HR: '/company/vacanciesList',
  guest: '/',
  admin: '/admin/companiesApplicationsList',
}

const PublicRoute = () => {
  const { role, isForbidden, isAuthChecked, isLoading } = useAuthStore();
  const location = useLocation();

  console.log('PublicRoute: ', { role, isForbidden, isAuthChecked, isLoading, pathname: location.pathname });
  if (isLoading || !isAuthChecked) return <Loader isLoading={isLoading} delay={300} />

  if (isForbidden && location.pathname !== '/companyUnderReview') {
    return <Navigate to="/companyUnderReview" replace />;
  }

  if (role !== "guest") {
    const redirectTo = roleRedirects[role];
    return <Navigate to={redirectTo} replace />;
  }

  return <Outlet />
}

export default PublicRoute;