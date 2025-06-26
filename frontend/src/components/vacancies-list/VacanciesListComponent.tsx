
import { Link } from "react-router-dom";
import {useVacancyCardStore} from "../../stores/vacanciesStore.ts";
import Loader from "../shared/loader/Loader.tsx";
import VacancyCardComponent from "../vacancy-card-component/VacancyCardComponent.tsx";
import './VacanciesListComponent.scss';

const VacanciesList = () => {
  const { allVacancies, totalCount, isLoading, error } = useVacancyCardStore();

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />;
  if (error) return <div>Ошибка: {error}</div>;
  if (totalCount === 0) return <div className="no-vacancies">Пока нет выложенных вакансий</div>;

  return (
    <>
      {allVacancies.map((vacancy) => (
        <Link className="vacancy-card__wrapper" key={vacancy.id} to={`/vacancy/${vacancy.id}`}>
          <VacancyCardComponent
            vacancy={{
              name: vacancy.name,
              companyName: vacancy.companyName,
              countOfReplies: vacancy.countOfReplies,
              isCreatedByUser: vacancy.isCreatedByUser,
            }}
          />
        </Link>
      ))}
    </>
  );
};

export default VacanciesList;
