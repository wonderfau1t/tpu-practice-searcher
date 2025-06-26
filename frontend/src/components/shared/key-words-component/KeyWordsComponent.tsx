import ButtonComponent from "../button/ButtonComponent.tsx";
import {Plus, Trash2, X} from "lucide-react";
import './KeyWordsComponent.scss';
import React, {useRef, useState} from "react";

interface KeywordsComponentProps {
  keywords: string[];
  setKeywords: (keywords: string[]) => void;
}

const KeyWordsComponent: React.FC<KeywordsComponentProps> = ({ keywords, setKeywords }) => {
  const [inputValue, setInputValue] = useState<string>('');
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  }

  const handleAddKeyword = () => {
    if (inputValue.trim()) {
      setKeywords([...keywords, inputValue.trim()]);
      setInputValue('');
      inputRef.current?.focus();
    }
  }

  const handleClearKeywords = () => {
    setKeywords([]);
  }

  const handleRemoveKeyword = (index: number) => {
    keywords.splice(index, 1);
    setKeywords([...keywords]);
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && inputValue.trim()) {
      handleAddKeyword();
    }
  }

  return (
    <div className="key-words-component">
      <div className="key-words-component__label">
        Ключевые слова
      </div>
      <div className="key-words-component__keywords-box">
        <ul className="key-words-component__keywords-box__list">
          { keywords.map((keyword, index) => (
            <li key={index} >
              { keyword }
              <ButtonComponent
                text={<X className="x-icon" size="14" color="gray"/>}
                variant={"x-mini-button"}
                clickFunction={() => handleRemoveKeyword(index)}
              />
            </li>
          ))}
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            placeholder="Введите ключевое слово"
            ref={inputRef}
          />
        </ul>
      </div>
      <div className="key-words-component__buttons">
        <ButtonComponent
          text={<Plus size="24" color="green" />}
          variant="tags-add-button"
          clickFunction={handleAddKeyword}
        />
        <ButtonComponent
          text={<Trash2 size="24" stroke-width="1.6" color="green"/>}
          variant="tags-remove-all-button"
          clickFunction={handleClearKeywords}
        />
      </div>
    </div>
  )
}

export default KeyWordsComponent;