import TextareaComponent from "../../components/shared/textarea/TextareaComponent.tsx";
import ButtonComponent from "../../components/shared/button/ButtonComponent.tsx";
import Input from "../../components/shared/input/Input.tsx";
import React, {useState} from "react";
import {useCompanyStore} from "../../stores/companyStore.ts";
import {useNavigate} from "react-router-dom";
import {tg} from "../../lib/telegram.ts";
import {toast} from "react-toastify";
import {useAuthStore} from "../../stores/authStore.ts";

interface FormData {
  name: string;
  description: string;
  link: string;
  phoneNumber?: string;
}

const CompanyRegistrationPage = () => {
  const [formErrors, setFormErrors] = useState<{ name: boolean }>({ name: false });
  const { verifyUser } = useAuthStore();
  const [formData, setFormData] = useState<FormData>({
    name: "",
    description: "",
    link: "",
  });
  const { registerCompany, isLoading, error } = useCompanyStore();

  const navigate = useNavigate();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async () => {
    if (!formData.name.trim()) {
      toast.error("Пожалуйста, укажите название компании");
      setFormErrors({ name: true })
      return;
    }

    setFormErrors({ name: false });

    try {
      await registerCompany({
        name: formData.name,
        description: formData.description,
        link: formData.link,
      });
      await verifyUser(tg.initData);
      navigate("/company/vacanciesList");
    } catch (err) {
      console.log('ошибка регистрации')
    }
  };

  return (
    <div className="form-container">
      <div className="company-reg__header">
        Заполните форму подачи заявки на регистрацию компании:
      </div>
      {error && <div className="company-reg__error">{error}</div>}
      <div className="company-reg__form">
        <Input
          label="Название компании *"
          value={formData.name}
          onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "name" } })}
          required
          error={formErrors?.name}
          errorMessage="Укажите название компании"
        />
        <Input
          label="Ссылка на сайт компании"
          value={formData.link}
          onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "link" } })}
        />
        <TextareaComponent
          label="Описание компании"
          value={formData.description}
          onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: "description" } })}
        />
        <div className="company-reg__submit-button">
          <ButtonComponent
            variant="default"
            text={isLoading ? "Отправка..." : "Отправить"}
            clickFunction={handleSubmit}
            disabled={isLoading}
          />
        </div>
      </div>
    </div>
  );
};

export default CompanyRegistrationPage;