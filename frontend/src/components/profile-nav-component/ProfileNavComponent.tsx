import {Link, useLocation} from "react-router-dom";
import './ProfileNavComponent.scss';

interface ProfileNavProps {
  role: "guest" | "student" | "moderator" | "headHR" | "HR" | "admin";
}

const tabConfig: Record<ProfileNavProps['role'], { path: string; label: string }[]> = {
  headHR: [
    { path: 'info', label: 'Инфо' },
    { path: 'hrs', label: 'Все HR' },
    { path: 'createHR', label: 'Добавить HR' },
  ],
  HR: [
    { path: 'info', label: 'Информация' },
    { path: 'hrs', label: 'Все HR' },
    { path: 'createHR', label: 'Добавить HR' },
  ],
  student: [
    { path: 'info', label: 'Инфо студента' },
  ],
  moderator: [
    { path: 'info', label: 'Инфо модератора' },
  ],
  guest: [],
  admin: [],
};

const ProfileNavComponent = ({ role }: ProfileNavProps) => {
  const location = useLocation();

  const activeTab = tabConfig[role].findIndex((tab) =>
    location.pathname.includes(tab.path)
  );

  return (
    <nav className="profile__nav-tabs" data-active-tab={activeTab}>
      { tabConfig[role].map((tab, index) => (
        <Link
          key={tab.path}
          to={tab.path}
          className={`tab ${index === activeTab ? "active" : ""}`}
        >
          { tab.label }
        </Link>
      )) }
    </nav>
  )
}

export default ProfileNavComponent;