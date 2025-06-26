import {useVacancyCardStore} from "../../../stores/vacanciesStore.ts";
import {useEffect} from "react";
import VacancyCardComponent from "../../../components/vacancy-card-component/VacancyCardComponent.tsx";
import {Link} from "react-router-dom";
import Loader from "../../../components/shared/loader/Loader.tsx";

const ModeratorVacanciesListPage = () => {
  const { moderatorVacancies, totalCount, isLoading, error, fetchModeratorVacancies } = useVacancyCardStore();
  useEffect(() => {
    fetchModeratorVacancies();
  }, [fetchModeratorVacancies]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: {error}</div>
  if (totalCount === 0) return <div className="no-vacancies">Пока нет выложенных вакансий</div>

  return (
    <div className="company-vacancies">
      <div className="company-vacancies__list">
        { moderatorVacancies.map((vacancy) => (
          <Link className="vacancy-card__wrapper" key={vacancy.id} to={`/vacancy/${vacancy.id}`}>
            <VacancyCardComponent
              key={vacancy.id}
              vacancy={{
                name: vacancy.name,
                companyName: vacancy.companyName,
                countOfReplies: vacancy.countOfReplies,
                isCreatedByUser: vacancy.isCreatedByUser,
              }}
            />
          </Link>
        )) }
      </div>
    </div>
  );
}

export default ModeratorVacanciesListPage;