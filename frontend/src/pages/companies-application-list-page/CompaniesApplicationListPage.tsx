import './CompaniesApplicationListPage.scss';
import { useEffect } from "react";
import { Link } from "react-router-dom";
import Loader from "../../components/shared/loader/Loader.tsx";
import CompanyApplicationCard from "../../components/company-application-card-component/CompanyApplicationCard.tsx";
import {useApplicationsStore} from "../../stores/applicationsStore.ts";
const CompaniesApplicationListPage = () => {
  const { companies, totalCount, isLoading, error, fetchCompaniesApplications } = useApplicationsStore();
  useEffect(() => {
    fetchCompaniesApplications();
  }, [fetchCompaniesApplications]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: {error}</div>
  if (totalCount === 0) return <div className="no-applications">Заявок пока нет</div>

  return (
    <div className="companies-applications-list">
      <div className="companies-applications-list__list">
        {Array.isArray(companies) && companies.length > 0 ? (
          companies.map((application) => (
            <Link
              className="application-card__wrapper"
              key={application.companyID}
              to={`/currentCompanyApplication/${application.companyID}`}
            >
              <CompanyApplicationCard
                application={{
                  companyName: application.name,
                  registeredAt: application.registeredAt,
                }}
              />
            </Link>
          ))
        ) : (
          <div className="no-applications">Заявок пока нет</div>
        )}
      </div>
    </div>
  );
}

export default CompaniesApplicationListPage;