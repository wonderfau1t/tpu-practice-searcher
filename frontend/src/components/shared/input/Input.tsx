import './Input.scss';
import React from 'react';

interface InputParams {
  label: string;
  value: string | number | undefined; // Теперь принимаем string | number
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  type?: string;
  required?: boolean;
  error?: boolean;
  onBlur?: () => void;
  errorMessage?: string;
}

const Input: React.FC<InputParams> = ({
                                        label,
                                        value,
                                        type = 'text',
                                        onChange,
                                        required = false,
                                        error = false,
                                        errorMessage,
                                      }) => {

  const stringValue = value != null ? String(value) : '';

  return (
    <div className="custom-input__wrapper">
      <div className="custom-input__wrapper__label">
        {label}
      </div>
      <label className="custom-input">
        <input
          type={type}
          required={required}
          value={stringValue}
          onChange={onChange}
          placeholder={label}
        />
      </label>
      {error && <div className="error-text">{errorMessage || 'Поле обязательно'}</div>}
    </div>
  );
};

export default Input;