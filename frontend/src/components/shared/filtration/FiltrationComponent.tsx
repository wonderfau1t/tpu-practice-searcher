import ButtonComponent from "../button/ButtonComponent.tsx";
import './FiltrationComponent.scss';
import SelectMenuComponent from "../select-menu/SelectMenuComponent.tsx";
import {useEffect, useState} from "react";
import {useVacancyStore} from "../../../stores/createVacancyStore.ts";
import {useVacancyCardStore} from "../../../stores/vacanciesStore.ts";
import SearchComponent from "../search/SearchComponent.tsx";

const FiltrationComponent = () => {
  const [isFiltersVisible, setIsFiltersVisible] = useState(false);
  const { courses, form, fetchFilters, setFormField } = useVacancyStore();
  const { filterVacancies, fetchAllVacancies } = useVacancyCardStore();
  useEffect(() => {
    fetchFilters();
  }, [fetchFilters]);

  const toggleFilters = () => {
    setIsFiltersVisible(!isFiltersVisible);
  };

  const handleGetAllVacancies = () => {
    fetchAllVacancies();
  }

  const applyFilters = async () => {
    try {
      await filterVacancies({
        courseIDs: form.courseIds.length > 0 ? form.courseIds : [],
      });
      setIsFiltersVisible(false);
    } catch (error) {
      console.error('Ошибка при применении фильтров:', error);
    }
  };

  return (
    <div className="filtration-component">
      <div className="filtration-component__wrapper">
        <div className="filtration-component__filtration__search">
          <SearchComponent placeholder="Поиск вакансий" />
        </div>
        <div className="filtration-component__filtration__buttons">
          <ButtonComponent
            text="Фильтры"
            variant="filter"
            clickFunction={toggleFilters}
          />
          <ButtonComponent
            text="Сбросить"
            variant="filter"
            clickFunction={handleGetAllVacancies}
          />
        </div>
      </div>
      {isFiltersVisible && (
        <div className={`filtration-component__filters-content ${isFiltersVisible ? 'visible' : ''}`}>
          <SelectMenuComponent
            label="Направления"
            searchable
            multiple
            options={courses.map((option) => ({ value: option.id.toString(), label: option.name }))}
            value={form.courseIds.map((id) => id.toString())}
            onChange={(values) => setFormField('courseIds', (values as string[]).map((v) => parseInt(v)))}
          />
          <ButtonComponent text="Применить" variant="default" clickFunction={applyFilters} />
        </div>
      )}
    </div>
  )
}

export default FiltrationComponent;