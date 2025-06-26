import {Link2, Pencil, X} from "lucide-react";
import './CompanyProfileInfoComponent.scss';
// import ButtonComponent from "../shared/button/ButtonComponent.tsx";
import React, {useEffect, useState} from "react";
import ButtonComponent from "../shared/button/ButtonComponent.tsx";
import Input from "../shared/input/Input.tsx";
import TextareaComponent from "../shared/textarea/TextareaComponent.tsx";
import {useCompanyStore} from "../../stores/companyStore.ts";
import {useProfileStore} from "../../stores/profileStore.ts";
import {useAuthStore} from "../../stores/authStore.ts";
import {toast} from "react-toastify";

const CompanyProfileInfoComponent: React.FC = () => {
  const [isEditing, setIsEditing] = useState(false);
  const { companyInfo, fetchCompanyInfo } = useProfileStore();
  const { editCompany, isLoading } = useCompanyStore();
  const [editedInfo, setEditedInfo] = useState({
    name: companyInfo?.name,
    description: companyInfo?.description,
    website: companyInfo?.link,
  });
  const [formErrors, setFormErrors] = useState<{ name: boolean }>({ name: false });

  const { role } = useAuthStore();

  const handleEditMode = () => {
    if (isEditing) {
      setIsEditing(false);
    } else {
      setIsEditing(true);
    }
  }

  useEffect(() => {
    if (companyInfo) {
      setEditedInfo({
        name: companyInfo.name || '',
        description: companyInfo.description || '',
        website: companyInfo.link || '',
      });
    }
  }, [companyInfo]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setEditedInfo((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmitEditing = async () => {
    if (!editedInfo?.name?.trim()) {
      toast.error("Пожалуйста, укажите название компании");
      setFormErrors({ name: true })
      return;
    }

    setFormErrors({ name: false });

    try {
      await editCompany({
        name: editedInfo?.name || '',
        description: editedInfo?.description || '',
        link: editedInfo.website || '',
      });
      await fetchCompanyInfo();
      toast.success("Изменения успешно сохранены!");
      setIsEditing(false);
    } catch (err) {
      console.log('ошибка регистрации')
    }
  };

  return (
    <div className="company-profile__company-info">
      { isEditing ? (
        <div className="company-profile__company-info__edit-block">
          <div className="company-profile__company-info__edit-block__edit-button">
            <ButtonComponent
              text={<X size="24" color="green" />}
              variant="close-button"
              clickFunction={handleEditMode}
            />
          </div>
          <Input
            label="Название компании *"
            value={editedInfo?.name}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "name" } })}
            required
            error={formErrors?.name}
            errorMessage="Укажите название компании"
          />
          <TextareaComponent
            label="Описание компании"
            value={editedInfo?.description}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "description" } })}
          />
          <Input
            label="Ссылка на сайт компании"
            value={editedInfo?.website}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "website" } })}
          />

          <ButtonComponent
            variant="default"
            text={isLoading ? "Отправка..." : "Отправить"}
            clickFunction={handleSubmitEditing}
            disabled={isLoading}
          />
        </div>
      ) : (
        <>
          { (role === 'headHR') && (
            <div className="company-profile__company-info__edit-button">
              <ButtonComponent
                text={<Pencil size="24" color="green" />}
                variant="tags-add-button"
                clickFunction={handleEditMode}
              />
            </div>
          )}
          <div className="company-profile__company-info__name">
            <div className="company-profile__company-info__name__text">
              {companyInfo?.name}
            </div>
          </div>
          <div className="company-profile__company-info__description">
            <div className="company-profile__company-info__description__label">
              Описание:
            </div>
            <div className="company-profile__company-info__description__text">
              { companyInfo?.description }
            </div>
          </div>
          <div className="company-profile__company-info__link">
            <div className="company-profile__company-info__link__label">
              Ссылки:
            </div>
            <div className="company-profile__company-info__link__text">
              <Link2 size={26}/>
              { companyInfo?.link }
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default CompanyProfileInfoComponent;