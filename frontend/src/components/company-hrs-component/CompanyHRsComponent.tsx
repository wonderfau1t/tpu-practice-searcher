import './CompanyHRsComponent.scss';

interface HRsProps {
  hrs: {
    id: number;
    username: string;
    countOfVacancies: number;
  } [];
}

const CompanyHRsComponent = ({ hrs }: HRsProps) => {
  return (
    <div className="company-hrs">
      <div className="company-hrs__title">
        <div>Username</div>
        <div>Кол-во вакансий</div>
      </div>
      <div className="company-hrs__line"></div>
      <div className="company-hrs__list">
        { hrs.map((hr) => (
          <div key={hr.id} className="company-hrs__list__item">
            <div className="company-hrs__list__item__username">@{ hr.username }</div>
            <div className="company-hrs__list__item__vacancies-count">{ hr.countOfVacancies }</div>
          </div>
        )) }
      </div>
    </div>
  );
};

export default CompanyHRsComponent;