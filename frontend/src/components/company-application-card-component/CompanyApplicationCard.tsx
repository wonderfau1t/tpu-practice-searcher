import './CompanyApplicationCard.scss';

interface CompanyApplicationCardProps {
  application: {
    companyName: string;
    registeredAt: string;
  };
}

const VacancyCardComponent = ({ application }: CompanyApplicationCardProps) => {
  return (
    <div className="company-application-card">
      <div className="company-application-card__company-name">
        { application.companyName }
      </div>
      <div className="company-application-card__registered-date">
        { application.registeredAt }
      </div>
    </div>
  );
};

export default VacancyCardComponent;