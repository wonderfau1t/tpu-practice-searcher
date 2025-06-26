import Input from '../../components/shared/input/Input.tsx';
import SelectMenuComponent from '../../components/shared/select-menu/SelectMenuComponent.tsx';
import ButtonComponent from '../../components/shared/button/ButtonComponent.tsx';
import TextareaComponent from '../../components/shared/textarea/TextareaComponent.tsx';
import { useVacancyStore } from '../../stores/createVacancyStore.ts';
import {useEffect, useState} from 'react';
import './CreateVacancyPage.scss';
import {toast} from "react-toastify";
import handleRequestContact from "../../utils/requestContact.ts";
import KeyWordsComponent from "../../components/shared/key-words-component/KeyWordsComponent.tsx";
import {useAuthStore} from "../../stores/authStore.ts";

const CreateVacancyPage = () => {
  const [formErrors, setFormErrors] = useState<{ name: boolean }>({ name: false });
  const {
    form,
    formats,
    courses,
    farePaymentMethods,
    accommodationPaymentMethods,
    isLoading,
    hasFetchedOptions,
    setFormField,
    setDescriptionField,
    fetchFormOptions,
    submitVacancy,
    submitVacancyWithoutCompany,
    resetForm,
  } = useVacancyStore();

  useEffect(() => {
    if (!hasFetchedOptions) {
      fetchFormOptions();
    }
    return () => resetForm();
  }, [hasFetchedOptions, fetchFormOptions, resetForm]);

  const { role } = useAuthStore();

  const navigateToVacancies = async () => {
    if (!form.vacancyName.trim()) {
      toast.error('Пожалуйста, введите название вакансии');
      setFormErrors({ name: true })
      return;
    }

    if (!form.formatId) {
      toast.error('Пожалуйста, выберите формат практики');
      setFormErrors({ name: true })
      return;
    }

    if (!form.courseIds.length) {
      toast.error('Пожалуйста, выберите хотя бы одно направление');
      setFormErrors({ name: true })
      return;
    }

    if (!form.farePaymentId) {
      toast.error('Пожалуйста, выберите оплату проезда');
      setFormErrors({ name: true })
      return;
    }
    if (farePaymentOtherId && form.farePaymentId === farePaymentOtherId && !form.farePaymentDetails.trim()) {
      toast.error('Пожалуйста, укажите детали оплаты проезда');
      setFormErrors({ name: true })
      return;
    }
    if (!form.paymentForAccommodationId) {
      toast.error('Пожалуйста, выберите оплату проживание');
      setFormErrors({ name: true })
      return;
    }
    if (accommodationPaymentOtherId && form.paymentForAccommodationId === accommodationPaymentOtherId && !form.paymentForAccommodationDetails.trim()) {
      toast.error('Пожалуйста, укажите детали оплаты проживания');
      setFormErrors({ name: true })
      return;
    }

    if (!form.deadlineAt) {
      toast.error('Пожалуйста, выберите конечную дату приёма заявок');
      setFormErrors({ name: true })
      return;
    }

    setFormErrors({ name: false });

    try {
      if (role === 'moderator' || role === 'admin') {
        await submitVacancyWithoutCompany();
      } else {
        await submitVacancy();
      }
      toast.success('Вакансия успешно создана!');
    } catch (error: any) {
      if (error.response && error.response?.status === 403) {
        handleRequestContact();
      }
      toast.error('Ошибка при создании вакансии');
    }
  };

  const farePaymentOtherOption = farePaymentMethods.find((option) => option.name === 'Другое');
  const farePaymentOtherId = farePaymentOtherOption ? farePaymentOtherOption.id : null;

  const accommodationPaymentOtherOption = accommodationPaymentMethods.find(
    (option) => option.name === 'Другое'
  );
  const accommodationPaymentOtherId = accommodationPaymentOtherOption
    ? accommodationPaymentOtherOption.id
    : null;

  const farePaymentOptions = farePaymentMethods.map((option) => ({
    value: option.id.toString(),
    label: option.name,
  }));

  const accommodationPaymentOptions = accommodationPaymentMethods.map((option) => ({
    value: option.id.toString(),
    label: option.name,
  }));

  return (
    <div className="form-container">
      <div className="company-reg__header">Заполните форму создания вакансии:</div>
      <div className="company-reg__form">
        { (role === 'moderator' || role === 'admin') && (
          <Input
            label="Название компании *"
            value={form.companyName}
            onChange={(e) => setFormField('companyName', e.target.value)}
            required={true}
            error={formErrors?.name}
            errorMessage="Укажите название компании"
          />
        )}

        <Input
          label="Название вакансии *"
          value={form.vacancyName}
          onChange={(e) => setFormField('vacancyName', e.target.value)}
          required={true}
          error={formErrors?.name}
          errorMessage="Укажите название вакансии"
        />
        <SelectMenuComponent
          label="Формат практики *"
          options={formats.map((format) => ({ value: format.id.toString(), label: format.name }))}
          value={form.formatId?.toString() || ''}
          onChange={(value) => setFormField('formatId', parseInt(value as string))}
          error={formErrors?.name}
          errorMessage="Укажите формат практики"
        />
        <SelectMenuComponent
          label="Направления *"
          searchable
          multiple
          options={courses.map((option) => ({ value: option.id.toString(), label: option.name }))}
          value={form.courseIds.map((id) => id.toString())}
          onChange={(values) => setFormField('courseIds', (values as string[]).map((v) => parseInt(v)))}
          error={formErrors?.name}
          errorMessage="Укажите хотя бы одно направление"
        />
        <SelectMenuComponent
          label="Оплата проезда *"
          options={farePaymentOptions}
          value={form.farePaymentId?.toString() || ''}
          onChange={(value) => {
            const selectedId = parseInt(value as string);
            setFormField('farePaymentId', selectedId);
            if (farePaymentOtherId && selectedId !== farePaymentOtherId) {
              setFormField('farePaymentDetails', '');
            }
          }}
          error={formErrors?.name}
          errorMessage="Укажите оплату проезда"
        />
        {farePaymentOtherId && form.farePaymentId === farePaymentOtherId && (
          <TextareaComponent
            label="Детали оплаты проезда *"
            value={form.farePaymentDetails}
            onChange={(e) => setFormField('farePaymentDetails', e.target.value)}
            required={false}
            error={formErrors?.name}
            errorMessage="Укажите детали оплаты проезда"
          />
        )}
        <SelectMenuComponent
          label="Оплата проживания *"
          options={accommodationPaymentOptions}
          value={form.paymentForAccommodationId?.toString() || ''}
          onChange={(value) => {
            const selectedId = parseInt(value as string);
            setFormField('paymentForAccommodationId', selectedId);
            if (accommodationPaymentOtherId && selectedId !== accommodationPaymentOtherId) {
              setFormField('paymentForAccommodationDetails', '');
            }
          }}
          error={formErrors?.name}
          errorMessage="Укажите оплату проживания"
        />
        {accommodationPaymentOtherId && form.paymentForAccommodationId === accommodationPaymentOtherId && (
          <TextareaComponent
            label="Детали оплаты проживания *"
            value={form.paymentForAccommodationDetails}
            onChange={(e) => setFormField('paymentForAccommodationDetails', e.target.value)}
            required={false}
            error={formErrors?.name}
            errorMessage="Укажите детали оплаты проживания"
          />
        )}
        <Input
          label="Конечная дата приема заявок *"
          type="date"
          value={form.deadlineAt ? form.deadlineAt.slice(0, 10) : ''}
          onChange={(e) => {
            const selectedDate = e.target.value;
            // Преобразуем в ISO с учетом времени (можно добавить 00:00 или текущее)
            const dateISO = new Date(selectedDate).toISOString();
            setFormField('deadlineAt', dateISO);
          }}
          required={true}
          error={formErrors?.name}
          errorMessage="Укажите конечную дату приема заявок"
        />
        <KeyWordsComponent
          keywords={form.keywords}
          setKeywords={(keywords) => setFormField('keywords', keywords)}
        />

        <TextareaComponent
          label="Место работы"
          value={form.description.workplace}
          onChange={(e) => setDescriptionField('workplace', e.target.value)}
          required={false}
        />
        <TextareaComponent
          label="Должность"
          value={form.description.position}
          onChange={(e) => setDescriptionField('position', e.target.value)}
          required={false}
        />
        <Input
          label="Зарплата"
          value={form.description.salary}
          onChange={(e) => setDescriptionField('salary', e.target.value)}
          required={false}
        />
        <TextareaComponent
          label="Питание"
          value={form.description.food}
          onChange={(e) => setDescriptionField('food', e.target.value)}
          required={false}
        />
        <TextareaComponent
          label="Требования"
          value={form.description.requirements}
          onChange={(e) => setDescriptionField('requirements', e.target.value)}
          required={false}
        />
        <TextareaComponent
          label="Условия"
          value={form.description.conditions}
          onChange={(e) => setDescriptionField('conditions', e.target.value)}
          required={false}
        />
        <TextareaComponent
          label="Дополнительная информация"
          value={form.description.additionalInfo}
          onChange={(e) => setDescriptionField('additionalInfo', e.target.value)}
          required={false}
        />
        <div className="company-reg__submit-button">
          <ButtonComponent
            variant="default"
            clickFunction={navigateToVacancies}
            text={isLoading ? 'Создание...' : 'Создать'}
            disabled={isLoading}
          />
        </div>
      </div>
    </div>
  );
};

export default CreateVacancyPage;