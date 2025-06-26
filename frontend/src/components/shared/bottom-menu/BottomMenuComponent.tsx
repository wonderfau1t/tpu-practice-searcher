import {List, User, Plus, Star, SignpostBig} from "lucide-react";
import { Link, useLocation } from "react-router-dom";
import "./BottomMenuComponent.scss";
import React, {JSX} from "react";
import {useAuthStore} from "../../../stores/authStore.ts";

const menuConfig: Record<string, { path: string, label: string, icon: JSX.Element}[]> = {
  headHR: [
    { path: '/company/vacanciesList', label: 'Вакансии', icon: <List size="30" color="#6B7280" /> },
    { path: '/company/createVacancy', label: 'Создать вакансию', icon: <Plus size="30" color="#6B7280" /> },
    { path: '/company/companyProfile', label: 'Компания', icon: <User size="30" color="#6B7280" /> },
  ],
  student: [
    { path: '/student/vacanciesList', label: 'Вакансии', icon: <List size="30" color="#6B7280" /> },
    { path: '/student/replies', label: 'Мои отклики', icon: <Star size="30" color="#6B7280" /> },
  ],
  moderator: [
    { path: '/moderator/vacanciesList', label: 'Вакансии', icon: <List size="30" color="#6B7280" /> },
    { path: '/moderator/createVacancy', label: 'Создать вакансию', icon: <Plus size="30" color="#6B7280" /> },
  ],
  HR: [
    { path: '/company/vacanciesList', label: 'Вакансии', icon: <List size="30" color="#6B7280" /> },
    { path: '/company/createVacancy', label: 'Создать вакансию', icon: <Plus size="30" color="#6B7280" /> },
    { path: '/company/companyProfile', label: 'Компания', icon: <User size="30" color="#6B7280" /> },
  ],
  admin: [
    { path: '/admin/companiesApplicationsList', label: 'Заявки', icon: <SignpostBig stroke-width="1.7" size="30" color="#6B7280"/> },
    { path: '/admin/createVacancy', label: 'Создать вакансию', icon: <Plus size="30" color="#6B7280" /> },
    { path: '/admin/vacanciesList', label: 'Вакансии', icon: <List size="30" color="#6B7280" /> },
  ]
}

const BottomMenuComponent = () => {
  const location = useLocation();
  const { role } = useAuthStore();

  const menuItems = menuConfig[role] || menuConfig.student;

  const getActiveTab = () => {
    if (location.pathname.includes('/vacancy/')) {
      return menuItems.findIndex((item) => item.path.includes('vacanciesList'));
    }
    const activeIndex = menuItems.findIndex((item) => location.pathname.includes(item.path));
    return activeIndex !== -1 ? activeIndex : 0;
  };

  const activeTab = getActiveTab();

  const tabCount = menuItems.length;
  const highlightWidth = `${100 / tabCount}%`;
  return (
    <nav
      className="bottom-menu__wrapper"
      data-active-tab={activeTab}
      style={{ '--highlight-width': highlightWidth } as React.CSSProperties}
    >
      { menuItems.map((item, index) => (
        <Link
          key={item.path}
          to={item.path}
          className={`tab bottom-menu__item ${activeTab === index ? "active" : ""}`}
        >
          { item.icon }
          <div>{ item.label }</div>
        </Link>
      )) }
    </nav>
  );
};

export default BottomMenuComponent;