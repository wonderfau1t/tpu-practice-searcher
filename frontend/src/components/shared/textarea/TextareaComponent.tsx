import './TextareaComponent.scss';
import React, { useRef } from 'react';

interface TextareaProps {
  label: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  required?: boolean;
  error?: boolean;
  errorMessage?: string;
}

const TextareaComponent: React.FC<TextareaProps> = ({
                                                      label,
                                                      value = '',
                                                      onChange,
                                                      required = false,
                                                      error = false,
                                                      errorMessage,
                                                    }) => {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleInput = () => {
    const element = textareaRef.current;
    if (element) {
      element.style.height = 'auto';
      const scrollHeight = element.scrollHeight;
      element.style.height = `${Math.min(scrollHeight, 200)}px`;
    }
  };

  return (
    <div className="custom-textarea__wrapper">
      <div className="custom-textarea__wrapper__label">{label}</div>
      <label className="custom-textarea">
        <textarea
          ref={textareaRef}
          value={value}
          onChange={onChange}
          onInput={handleInput}
          required={required}
          placeholder={label}
        />
      </label>
      {error && <div className="error-text">{errorMessage || 'Поле обязательно'}</div>}
    </div>
  );
};

export default TextareaComponent;