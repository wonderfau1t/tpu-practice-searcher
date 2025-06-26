// import './EditCompanyInfoComponent.scss';
// import ButtonComponent from "../shared/button/ButtonComponent.tsx";
// import TextareaComponent from "../shared/textarea/TextareaComponent.tsx";
// import Input from "../shared/input/Input.tsx";
// import React, {useState} from "react";
// import {useNavigate} from "react-router-dom";
// import {toast} from "react-toastify";
// import {tg} from "../../lib/telegram.ts";
//
// const EditCompanyInfoComponent = () => {
//   const [formErrors, setFormErrors] = useState<{ name: boolean }>({ name: false });
//   const [formData, setFormData] = useState<FormData>({
//     name: "",
//     description: "",
//     link: "",
//   });
//   const navigate = useNavigate();
//
//   const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
//     const { name, value } = e.target;
//     setFormData((prev) => ({
//       ...prev,
//       [name]: value,
//     }));
//   };
//
//   const handleSubmit = async () => {
//     if (!formData.name.trim()) {
//       toast.error("Пожалуйста, укажите название компании");
//       setFormErrors({ name: true })
//       return;
//     }
//
//     setFormErrors({ name: false });
//
//     try {
//       await registerCompany({
//         name: formData.name,
//         description: formData.description,
//         link: formData.link,
//       });
//       await verifyUser(tg.initData);
//       navigate("/company/vacanciesList");
//     } catch (err) {
//       console.log('ошибка регистрации')
//     }
//   };
//
//   return (
//     <div className="edit-company-info__wrapper">
//       <div className="edit-company-info__wrapper__button">
//         <ButtonComponent text="Редактировать данные" variant="outline"/>
//       </div>
//       <div className="edit-company-info__wrapper__form">
//         <Input
//           label="Название компании *"
//           value={formData.name}
//           onChange={(e) => handleInputChange({...e, target: {...e.target, name: "name"}})}
//           required
//           error={formErrors?.name}
//           errorMessage="Укажите название компании"
//         />
//         <Input
//           label="Ссылка на сайт компании"
//           value={formData.link}
//           onChange={(e) => handleInputChange({...e, target: {...e.target, name: "link"}})}
//         />
//         <TextareaComponent
//           label="Описание компании"
//           value={formData.description}
//           onChange={(e) => handleInputChange({...e, target: {...e.target, name: "description"}})}
//         />
//         <div className="edit-company-info__wrapper__form__submit-button">
//           <ButtonComponent
//             variant="default"
//             text={isLoading ? "Отправка..." : "Отправить"}
//             clickFunction={handleSubmit}
//             disabled={isLoading}
//           />
//         </div>
//       </div>
//     </div>
//   );
// };
//
// export default EditCompanyInfoComponent;