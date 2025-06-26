import {useEffect} from "react";
import {useVacancyCardStore} from "../../../stores/vacanciesStore.ts";
import VacancyCardComponent from "../../../components/vacancy-card-component/VacancyCardComponent.tsx";
import {Link} from "react-router-dom";
import Loader from "../../../components/shared/loader/Loader.tsx";
import './StudentRepliesPage.scss';
const StudentRepliesPage = () => {
  const { repliedVacancies, totalCount, isLoading, error, fetchRepliedVacancies } = useVacancyCardStore();
  useEffect(() => {
    fetchRepliedVacancies();
  }, [fetchRepliedVacancies]);

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка: {error}</div>
  if (totalCount === 0) return <div className="no-replies">Вы не откликнулись ни на одну вакансию</div>

  return (
    <div className="company-vacancies">
      <div className="company-vacancies__list">
        { repliedVacancies.map((vacancy) => (
          <Link className="vacancy-card__wrapper" key={vacancy.id} to={`/vacancy/${vacancy.id}`}>
            <VacancyCardComponent
              key={vacancy.id}
              vacancy={{
                name: vacancy.name,
                companyName: vacancy.companyName,
              }}
            />
          </Link>
        )) }
      </div>
    </div>
  );
}

export default StudentRepliesPage;