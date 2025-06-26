import './VacancyCardComponent.scss';
import {ContactRound} from "lucide-react";

interface VacancyCardProps {
  vacancy: {
    name: string;
    courses?: string;
    countOfReplies?: number;
    companyName?: string;
    isCreatedByUser?: boolean;
  };
}

const VacancyCardComponent = ({ vacancy }: VacancyCardProps) => {
  return (
    <div className="vacancy-card">
      { vacancy.isCreatedByUser && (
        <div className="vacancy-card__who-create-label">
          <ContactRound stroke-width="1.7" size="24" color="white" className="vacancy-card__who-create-label__icon" />
        </div>
      ) }
      <div className="vacancy-card__vacancy-name">
        { vacancy.name }
      </div>
      { vacancy.companyName && (
        <div className="vacancy-card__course">
          <span className="course-label">Компания:</span>
          <span className="course-text">
          {vacancy.companyName}
        </span>
        </div>
      )}
      {typeof vacancy.countOfReplies !== 'undefined' && (
        <div className="vacancy-card__responses">
          <div className="vacancy-card__responses__content">
            <div className="text">Откликов:</div>
            <div className="number">{vacancy.countOfReplies}</div>
          </div>
        </div>
      )}
    </div>
  );
};

export default VacancyCardComponent;