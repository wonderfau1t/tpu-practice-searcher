import { useParams } from "react-router-dom";
import {useProfileStore} from "../../stores/profileStore.ts";
import {useEffect} from "react";
import Loader from "../../components/shared/loader/Loader.tsx";
import CompanyProfileInfoComponent from "../../components/company-profile-info-component/CompanyProfileInfoComponent.tsx";

const CompanyProfilePage = () => {
  const { id } = useParams<{ id: string }>();
  const { isLoading, error, fetchCompanyInfoFromVacancy } = useProfileStore();

  useEffect(() => {
    if (id && !isNaN(Number(id))) {
      fetchCompanyInfoFromVacancy(Number(id)).catch((err) => {
        console.error('Failed to fetch company info:', err);
      });
    } else {
      console.error('Invalid company ID:', id);
    }
  }, [id, fetchCompanyInfoFromVacancy]);


  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка{}: {error }</div>;

  return (
    <>
      <CompanyProfileInfoComponent />
    </>
  )
}

export default CompanyProfilePage;