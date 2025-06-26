import {createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";
import SelectRolePage from "./pages/SelectRolePage/SelectRolePage.tsx";
import NotFoundPage from "./pages/NotFoundPage/NotFoundPage.tsx";
import { useEffect } from "react";
import { useAuthStore } from "./stores/authStore.ts";
import { tg } from "./lib/telegram.ts";
import PrivateRoute from "./components/private-route/PrivateRouter.tsx";
import PublicRoute from "./components/public-route/PublicRoute.tsx";
import '../src/styles/global.scss';
import CompanyRegistrationPage from "./pages/company-registration-page/CompanyRegistrationPage.tsx";
import TelegramLayout from "./components/layouts/TelegramLayout.tsx";
import CreateVacancyPage from "./pages/create-vacancy-page/CreateVacancyPage.tsx";
import CompanyVacanciesPage from "./pages/company-vacancies-page/CompanyVacanciesPage.tsx";
import BottomMenuLayout from "./components/layouts/bottom-menu-layout/BottomMenuLayout.tsx";
import MainHRPInfoPage from "./pages/profile-pages/main-hr/main-hr-profile-info-page/MainHRPInfoPage.tsx";
import MainHRProfileLayout from "./components/layouts/profile-layout/MainHRProfileLayout.tsx";
import MainHRHRsPage from "./pages/profile-pages/main-hr/main-hr-hrs-page/MainHRHRsPage.tsx";
import MainHRCreateHRPage from "./pages/profile-pages/main-hr/main-hr-create-hr-page/MainHRCreateHRPage.tsx";
import VacancyPage from "./pages/vacancy-page/VacancyPage.tsx";
import StudentVacanciesListPage from "./pages/student/vacancies/StudentVacanciesListPage.tsx";
import StudentRepliesPage from "./pages/student/replies/StudentRepliesPage.tsx";
import CompanyProfilePage from "./pages/company-profile-page/CompanyProfilePage.tsx";
import ModeratorVacanciesListPage from "./pages/moderator/vacancies/ModeratorVacanciesListPage.tsx";
import {ToastContainer} from "react-toastify";
import CompaniesApplicationListPage from "./pages/companies-application-list-page/CompaniesApplicationListPage.tsx";
import CompanyApplicationPage from "./pages/company-application-page/CompanyApplicationPage.tsx";
import CompanyUnderReviewPage from "./pages/company-under-review-page/CompanyUnderReviewPage.tsx";

const router = createBrowserRouter([
  {
    path: '/',
    element: <TelegramLayout />,
    errorElement: <NotFoundPage />,
    children: [
      {
        path: '/',
        element: <PublicRoute />,
        children: [
          {
            path: '/',
            element: <SelectRolePage />,
          },
          {
            path: 'companyRegistration',
            element: <CompanyRegistrationPage />,
          },
          {
            path: 'companyUnderReview',
            element: <CompanyUnderReviewPage />,
          }
        ],
      },
      {
        path: '',
        element: <BottomMenuLayout />,
        children: [
          {
            path: 'company',
            element: <PrivateRoute allowedRoles={['headHR', 'HR']} />,
            children: [
              {
                path: '',
                element: <Navigate to="vacanciesList" replace />,
              },
              {
                path: 'vacanciesList',
                element: <CompanyVacanciesPage />,
              },
              {
                path: 'createVacancy',
                element: <CreateVacancyPage />,
              },
              {
                path: 'companyProfile',
                element: <MainHRProfileLayout />,
                children: [
                  {
                    path: '',
                    element: <Navigate to="info" replace />,
                  },
                  {
                    path: 'info',
                    element: <MainHRPInfoPage />,
                  },
                  {
                    path: 'hrs',
                    element: <MainHRHRsPage />,
                  },
                  {
                    path: 'createHR',
                    element: <MainHRCreateHRPage />,
                  },
                ],
              },
            ],
          },
          {
            path: 'moderator',
            element: <PrivateRoute allowedRoles={['moderator']} />,
            children: [
              {
                path: 'vacanciesList',
                element: <ModeratorVacanciesListPage />
              },
              {
                path: 'createVacancy',
                element: <CreateVacancyPage />
              },
            ]
          },
          {
            path: 'admin',
            element: <PrivateRoute allowedRoles={['admin']} />,
            children: [
              {
                path: 'companiesApplicationsList',
                element: <CompaniesApplicationListPage />,
              },
              {
                path: 'createVacancy',
                element: <CreateVacancyPage />,
              },
              {
                path: 'vacanciesList',
                element: <StudentVacanciesListPage />
              },
            ],
          },
          {
            path: 'student',
            element: <PrivateRoute allowedRoles={['student']} />,
            children: [
              {
                path: 'vacanciesList',
                element: <StudentVacanciesListPage />
              },
              {
                path: 'replies',
                element: <StudentRepliesPage />
              },
            ],
          },
          {
            path: 'vacancy/:id',
            element: <VacancyPage />,
          },
          {
            path: 'currentCompany/:id',
            element: <CompanyProfilePage />
          },
          {
            path: 'currentCompanyApplication/:id',
            element: <CompanyApplicationPage />
          }
        ],
      },
    ],
  },
]);

function App() {
  const { verifyUser } = useAuthStore();
  useEffect(() => {
    tg.ready();
    tg.expand();
    const initData = tg.initData;
    verifyUser(initData);
  }, [verifyUser]);

  return (
    <>
      <RouterProvider router={router} />
      <ToastContainer position="top-right" autoClose={2000} hideProgressBar={false} />
    </>
  );
}

export default App;
