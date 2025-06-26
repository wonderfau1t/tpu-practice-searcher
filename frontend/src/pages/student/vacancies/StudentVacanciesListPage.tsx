import {useVacancyCardStore} from "../../../stores/vacanciesStore.ts";
import {useEffect} from "react";
import './StudentVacanciesListPage.scss';
import FiltrationComponent from "../../../components/shared/filtration/FiltrationComponent.tsx";
import VacanciesListComponent from "../../../components/vacancies-list/VacanciesListComponent.tsx";
const StudentVacanciesListPage = () => {
  const { fetchAllVacancies } = useVacancyCardStore();
  useEffect(() => {
    fetchAllVacancies();
  }, [fetchAllVacancies]);

  return (
    <div className="company-vacancies">
      <div className="company-vacancies__list">
        <div className="company-vacancies__list__filtration-block">
          <FiltrationComponent />
        </div>
        <VacanciesListComponent />
      </div>
    </div>
  );
}

export default StudentVacanciesListPage;