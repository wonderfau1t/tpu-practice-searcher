import {BookDashed, Building2, Link2} from "lucide-react";
interface CompanyInfoProps {
  companyInfo: {
    name: string,
    description: string,
    website: string,
  } | null;
}
const CompanyProfileComponent = ({ companyInfo }: CompanyInfoProps) => {
  return (
    <div className="company-profile__company-info">
      <div className="company-profile__company-info__name">
        <div className="company-profile__company-info__name__label">
          <Building2 size={22}/>
          Название компании:
        </div>
        <div className="company-profile__company-info__name__text">
          {companyInfo?.name}
        </div>
      </div>
      <div className="company-profile__company-info__description">
        <div className="company-profile__company-info__description__label">
          <BookDashed size={22}/>
          Описание:
        </div>
        <div className="company-profile__company-info__description__text">
          {companyInfo?.description}
        </div>
      </div>
      <div className="company-profile__company-info__link">
        <div className="company-profile__company-info__link__label">
          <Link2 size={22}/>
          Ссылка на сайт компании:
        </div>
        <div className="company-profile__company-info__link__text">
          {companyInfo?.website}
        </div>
      </div>
    </div>
  );
}

export default CompanyProfileComponent;