import './SelectMenuComponent.scss';
import {ChevronDown, Plus, Search} from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';
import ButtonComponent from "../button/ButtonComponent.tsx";

interface SelectOption {
  value: string;
  label: string;
}

interface SelectMenuProps {
  label: string;
  options: SelectOption[];
  value: string | string[];
  onChange: (value: string | string[]) => void;
  searchable?: boolean;
  multiple?: boolean;
  addable?: boolean;
  error?: boolean;
  errorMessage?: string;

}

const SelectMenuComponent: React.FC<SelectMenuProps> = ({
                                                          label,
                                                          options,
                                                          value,
                                                          onChange,
                                                          searchable = false,
                                                          multiple = false,
                                                          addable = false,
                                                          error = false,
                                                          errorMessage,
                                                        }) => {
  const [isActive, setIsActive] = useState<boolean>(false);
  const [searchTerm, setSearchTerm] = useState<string>('');
  const selectRef = useRef<HTMLDivElement>(null);
  const optionsRef = useRef<HTMLUListElement>(null);


  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (selectRef.current && !selectRef.current.contains(event.target as Node)) {
        setIsActive(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const toggleSelectMenu = () => {
    setIsActive(!isActive);
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const filteredOptions = options.filter((option) =>
    option.label.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleOptionClick = (optionValue: string) => {
    if (multiple) {
      const currentValues = Array.isArray(value) ? value : [];
      const newValues = currentValues.includes(optionValue)
        ? currentValues.filter((v) => v !== optionValue)
        : [...currentValues, optionValue];
      onChange(newValues);
    } else {
      onChange(optionValue);
      setIsActive(false);
      setSearchTerm('');
    }
  };

  const selectedLabel = multiple
    ? Array.isArray(value) && value.length > 0
      ? `Выбрано: ${value.length}`
      : ''
    : options.find((opt) => opt.value === value)?.label || '';

  return (
    <div className="select__wrapper">
      <div className="select__wrapper__label">
        { label }
      </div>
      <div className="wrapper" ref={selectRef}>
        <div className="select-wrapper">
        <span
          className={`floating-placeholder ${
            isActive || (Array.isArray(value) ? value.length > 0 : value) ? 'active' : ''
          }`}
        >
          {label}
        </span>
          <div
            className={`select-button ${isActive ? 'active' : ''} ${
              (Array.isArray(value) && value.length > 0) || (!Array.isArray(value) && value) ? 'has-value' : ''
            }`}
            onClick={toggleSelectMenu}
          >
            <span>{selectedLabel}</span>
            <ChevronDown
              size="22"
              color="#8aeb9d"
              className={`chevron ${isActive ? 'rotated' : ''}`}
            />
          </div>
        </div>
        {isActive && (
          <div className={`content ${isActive ? 'active' : ''}`}>
            {searchable && (
              <div className="search">
                <Search className="search-icon" size="20" color="green"/>
                <input
                  type="text"
                  placeholder="Поиск"
                  value={searchTerm}
                  onChange={handleSearch}
                />
              </div>
            )}

            {addable && (
              <ButtonComponent
                text={<Plus size="20" color="green"/>}
                variant="default"
              />

            )}

            <ul className="options" ref={optionsRef}>
              {filteredOptions.map((option) => (
                <li
                  key={option.value}
                  className={
                    multiple
                      ? Array.isArray(value) && value.includes(option.value)
                        ? 'selected'
                        : ''
                      : value === option.value
                        ? 'selected'
                        : ''
                  }
                  onClick={() => handleOptionClick(option.value)}
                >
                  {option.label}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
      {error && <div className="error-text">{errorMessage || 'Поле обязательно'}</div>}

    </div>
  );
};

export default SelectMenuComponent;