import './ButtonComponent.scss';
import React, {JSX} from 'react';

interface ButtonProps {
  text: string | JSX.Element;
  clickFunction?: () => void;
  disabled?: boolean;
  variant:
    'default' |
    'decline' |
    'outline' |
    'filter' |
    'delete' |
    'delete-mini' |
    'tags-add-button' |
    'tags-remove-all-button' |
    'x-mini-button' |
    'close-button'
  ;
}

const ButtonComponent: React.FC<ButtonProps> = ({ text, clickFunction, disabled, variant = 'default' }) => {
  return (
    <button
      className={`button button--${variant}`}
      onClick={clickFunction}
      disabled={disabled}
    >
      {text}
    </button>
  );
};

export default ButtonComponent;