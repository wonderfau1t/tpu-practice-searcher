import CompanyHRsComponent from "../../../../components/company-hrs-component/CompanyHRsComponent.tsx";
import {useProfileStore} from "../../../../stores/profileStore.ts";
import {useEffect} from "react";
import './MainHRHRsPage.scss';
import Loader from "../../../../components/shared/loader/Loader.tsx";

const MainHRHRsPage = () => {
  const { hrs, isLoading, error, fetchHRs } = useProfileStore();

  useEffect(() => {
    fetchHRs();
  }, [fetchHRs]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: {error}</div>;

  return (
    <CompanyHRsComponent hrs={ hrs }/>
  )
}

export default MainHRHRsPage;