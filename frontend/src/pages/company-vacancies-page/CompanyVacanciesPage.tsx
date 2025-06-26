import './CompanyVacanciesPage.scss';
import VacancyCardComponent from "../../components/vacancy-card-component/VacancyCardComponent.tsx";
import { useVacancyCardStore } from "../../stores/vacanciesStore.ts";
import { useEffect } from "react";
import { Link } from "react-router-dom";
import Loader from "../../components/shared/loader/Loader.tsx";
const CompanyVacanciesPage = () => {
  const { vacancies, totalCount, isLoading, error, fetchCompanyVacancies } = useVacancyCardStore();
  useEffect(() => {
    fetchCompanyVacancies();
  }, [fetchCompanyVacancies]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: {error}</div>
  if (totalCount === 0) return <div className="no-vacancies">У вас пока нет созданных вакансий</div>

  return (
    <div className="company-vacancies">
      <div className="company-vacancies__list">
        { vacancies.map((vacancy) => (
          <Link className="vacancy-card__wrapper" key={vacancy.id} to={`/vacancy/${vacancy.id}`}>
            <VacancyCardComponent
              key={vacancy.id}
              vacancy={{
                name: vacancy.name,
                countOfReplies: vacancy.countOfReplies,
              }}
            />
          </Link>
        )) }
      </div>
    </div>
  );
}

export default CompanyVacanciesPage;