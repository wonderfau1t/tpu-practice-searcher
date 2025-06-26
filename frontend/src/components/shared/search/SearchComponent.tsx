import React, {useEffect, useMemo, useState} from 'react';
import './SearchComponent.scss';
import {useVacancyCardStore} from "../../../stores/vacanciesStore.ts";
import {debounce} from "../../../utils/debounce.ts";

interface SearchComponentsProps {
  placeholder?: string;
}

const SearchComponent: React.FC<SearchComponentsProps> = ({ placeholder = 'Поиск' }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const searchVacancies = useVacancyCardStore(state => state.searchVacancies);

  const debouncedSearch = useMemo(() => debounce(searchVacancies, 500), [searchVacancies]);

  useEffect(() => {
    if (searchQuery.trim()) {
      debouncedSearch(searchQuery);
    }
  }, [searchQuery]);

  return (
    <div className="search-component">
      <input
        className="search-component__input"
        placeholder={placeholder}
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
      />
    </div>
  );
};

export default SearchComponent;
