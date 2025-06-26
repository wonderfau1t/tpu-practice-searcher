import './MainHRProfileLayout.scss';
import ProfileNavComponent from "../../profile-nav-component/ProfileNavComponent.tsx";
import {Outlet} from "react-router-dom";
import {useAuthStore} from "../../../stores/authStore.ts";

const MainHRProfileLayout = () => {
  const { role } = useAuthStore();
  return (
    <div className="company-profile">
      <div className="company-profile__nav-tabs">
        <ProfileNavComponent role={role} />
      </div>

      <div className="company-profile__content">
        <Outlet />
      </div>
    </div>
  );
}

export default MainHRProfileLayout;