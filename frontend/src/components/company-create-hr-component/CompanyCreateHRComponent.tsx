import Input from "../shared/input/Input.tsx";
import './CompanyCreateHRComponent.scss';
import ButtonComponent from "../shared/button/ButtonComponent.tsx";
import React, {useState} from "react";
import {useProfileStore} from "../../stores/profileStore.ts";
import {toast} from "react-toastify";
const CompanyCreateHRComponent = () => {
  const [username, setUsername] = useState<string>("");
  const { createHR, isLoading, error } = useProfileStore();

  const handleSubmit = async () => {
    if (username) {
      try {
        await createHR(username);
        toast.success("Новый HR успешно добавлен!");
        setUsername("");
      } catch (error) {
        toast.error("Ошибка при добавлении HR")
      }
    }
  };

  return (
    <div className="create-hr">
      <div className="create-hr__form">
        <div className="create-hr__form__header">Добавить HR</div>
        <Input
          label="Telegram @username"
          value={username}
          required={true}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)}
        />
        <ButtonComponent variant="default" text="Добавить" clickFunction={handleSubmit} disabled={isLoading} />
        {error && <div className="error">{error}</div>}
      </div>
    </div>
  );
};

export default CompanyCreateHRComponent;