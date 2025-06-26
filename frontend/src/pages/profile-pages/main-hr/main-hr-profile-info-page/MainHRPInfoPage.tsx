import './MainHRInfoPage.scss';
import CompanyProfileInfoComponent
  from "../../../../components/company-profile-info-component/CompanyProfileInfoComponent.tsx";
import {useProfileStore} from "../../../../stores/profileStore.ts";
import {useEffect} from "react";
import Loader from "../../../../components/shared/loader/Loader.tsx";
// import EditCompanyInfoComponent from "../../../../components/edit-company-info/EditCompanyInfoComponent.tsx";

const MainHRPInfoPage = () => {
  const { isLoading, error, fetchCompanyInfo } = useProfileStore();

  useEffect(() => {
    fetchCompanyInfo();
  }, [fetchCompanyInfo]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: { error }</div>;

  return (
    <div className="main-hr-info__wrapper">
      <CompanyProfileInfoComponent />
    </div>
  );
}

export default MainHRPInfoPage;