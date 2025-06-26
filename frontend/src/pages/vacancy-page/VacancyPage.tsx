import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useVacancyCardStore } from '../../stores/vacanciesStore.ts';
import './VacancyPage.scss';
import ButtonComponent from '../../components/shared/button/ButtonComponent.tsx';
import { useAuthStore } from '../../stores/authStore.ts';
import Loader from '../../components/shared/loader/Loader.tsx';
import {Link2, Pencil, X} from 'lucide-react';
import { toast } from 'react-toastify';
import handleRequestContact from '../../utils/requestContact.ts';
import { useVacancyStore } from '../../stores/createVacancyStore.ts';
import Input from '../../components/shared/input/Input.tsx';
import SelectMenuComponent from '../../components/shared/select-menu/SelectMenuComponent.tsx';
import TextareaComponent from '../../components/shared/textarea/TextareaComponent.tsx';
import { Course} from "../../stores/vacanciesStore.ts";
import { VacancyInfo } from "../../stores/vacanciesStore.ts";
import KeyWordsComponent from "../../components/shared/key-words-component/KeyWordsComponent.tsx";

const VacancyPage = () => {
  const { id } = useParams<{ id: string }>();
  const {
    currentVacancy,
    isLoading,
    fetchVacancyById,
    replyToVacancy,
    declineVacancy,
    fetchVacancyByModerator,
    deleteVacancy,
    editVacancy,
  } = useVacancyCardStore();
  const { role } = useAuthStore();
  const navigate = useNavigate();
  const isReplied = useVacancyCardStore((state) => state.currentVacancy?.isReplied);

  const {
    form,
    formats,
    courses,
    farePaymentMethods,
    accommodationPaymentMethods,
    hasFetchedOptions,
    setFormField,
    fetchFormOptions,
  } = useVacancyStore();

  const [isEditing, setEditing] = useState(false);
  const [editedInfo, setEditedInfo] = useState({
    name: '',
    companyName: '',
    keywords: [] as string[],
    deadlineAt: '',
    description: {
      workplace: '',
      position: '',
      salary: '',
      requirements: '',
      food: '',
      conditions: '',
      additionalInfo: '',
    },
  });

  useEffect(() => {
    if (!hasFetchedOptions) {
      fetchFormOptions();
    }
  }, [hasFetchedOptions, fetchFormOptions]);

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

  const handleEditMode = () => {
    setEditing(!isEditing);
  };

  useEffect(() => {
    if (id) {
      if (role === 'headHR' || role === 'HR' || role === 'student') {
        fetchVacancyById(Number(id)).catch((err) => {
          console.error('Failed to fetch vacancy:', err);
        });
      } else if (role === 'moderator' || role === 'admin') {
        fetchVacancyByModerator(Number(id)).catch((err) => {
          console.error('Failed to fetch vacancy:', err);
        });
      }
    }
  }, [id, fetchVacancyById, fetchVacancyByModerator]);

  useEffect(() => {
    if (currentVacancy) {
      setEditedInfo({
        name: currentVacancy.vacancyName || '',
        companyName: currentVacancy.companyName || '',
        keywords: currentVacancy.keywords || [],
        deadlineAt: currentVacancy.deadlineAt || '',
        description: {
          workplace: currentVacancy.description.workplace || '',
          position: currentVacancy.description.position || '',
          salary: currentVacancy.description.salary || '',
          requirements: currentVacancy.description.requirements || '',
          food: currentVacancy.description.food || '',
          conditions: currentVacancy.description.conditions || '',
          additionalInfo: currentVacancy.description.additionalInfo || '',
        },
      });
      setFormField('formatId', currentVacancy.formatID);
      setFormField(
        'courseIds',
        currentVacancy.courses && Array.isArray(currentVacancy.courses)
          ? currentVacancy.courses.map((c: any) => c.courseId)
          : []
      );
      console.log('Initial courseIds:', currentVacancy);
      setFormField('farePaymentId', currentVacancy.farePaymentID);
      setFormField('farePaymentDetails', currentVacancy.farePaymentDetails || '');
      setFormField('paymentForAccommodationId', currentVacancy.paymentForAccommodationID);
      setFormField(
        'paymentForAccommodationDetails',
        currentVacancy.paymentForAccommodationDetails || ''
      );
    }
  }, [currentVacancy, setFormField]);

  const moveToCompany = async () => {
    if (!currentVacancy || !currentVacancy.companyID) {
      return;
    }
    try {
      navigate(`/currentCompany/${currentVacancy.companyID}`);
    } catch (error) {
      console.error('Failed to fetch company info:', error);
    }
  };

  const handleReply = async (vacancyId: number) => {
    try {
      await replyToVacancy(vacancyId);
      toast.success('Вы успешно откликнулись на вакансию!');
    } catch (error: any) {
      if (error.response && error.response?.status === 403) {
        handleRequestContact();
      } else {
        toast.error('Вы уже откликнулись или произошла ошибка');
      }
    }
  };

  const handleDecline = async (vacancyId: number) => {
    try {
      await declineVacancy(vacancyId);
      toast.success('Отклик на вакансию успешно отменен!');
    } catch (e) {
      toast.error('Ошибка');
    }
  };

  const handleDelete = async (vacancyId: number) => {
    try {
      await deleteVacancy(vacancyId);
      if (role === 'admin') {
        navigate('/admin/vacanciesList');
      } else {
        navigate('/company/vacanciesList');
      }
    } catch (error) {
      toast.error('Ошибка');
    }
  };

  const handleEditVacancy = async () => {
    if (!id) return;

    try {
      // Преобразуем courseIds в courses
      const selectedCourses = form.courseIds
        .map((courseId) => {
          const course = courses.find((c) => c.id === Number(courseId));
          return course ? { courseId: course.id, name: course.name } : null;
        })
        .filter((course): course is Course => course !== null);

      console.log('Selected courses:', selectedCourses);

      // Формируем updatedData
      const updatedData: Partial<VacancyInfo> = {
        vacancyName: editedInfo.name,
        keywords: editedInfo.keywords,
        deadlineAt: editedInfo.deadlineAt,
        formatID: form.formatId ?? 0,
        courses: selectedCourses,
        farePaymentID: form.farePaymentId ?? 0,
        farePaymentDetails: form.farePaymentDetails,
        paymentForAccommodationID: form.paymentForAccommodationId ?? 0,
        paymentForAccommodationDetails: form.paymentForAccommodationDetails,
        description: editedInfo.description,
      };

      // Добавляем companyName только для moderator и admin, если isCreatedByUser === true
      if ((role === 'moderator' || role === 'admin') && currentVacancy?.isCreatedByUser) {
        updatedData.companyName = editedInfo.companyName;
      }

      await editVacancy(Number(id), updatedData);
      toast.success('Вакансия успешно обновлена!');
      setEditing(false);
      await fetchVacancyById(Number(id)); // Перезагружаем данные
    } catch (error) {
      toast.error('Ошибка при сохранении изменений');
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    field: string
  ) => {
    const { value } = e.target;
    // Проверяем, является ли field вложенным полем в description
    const descriptionFields = [
      'workplace',
      'position',
      'salary',
      'requirements',
      'food',
      'conditions',
      'additionalInfo',
    ];
    if (descriptionFields.includes(field)) {
      setEditedInfo((prev) => ({
        ...prev,
        description: {
          ...prev.description,
          [field]: value,
        },
      }));
    } else {
      setEditedInfo((prev) => ({
        ...prev,
        [field]: value,
      }));
    }
  };

  if (isLoading) return <Loader isLoading={isLoading} delay={300} />;
  if (!currentVacancy) return <div>Вакансия не найдена</div>;

  return (
    <div className="vacancy-page__wrapper">
      {isEditing ? (
        <div className="vacancy-page__wrapper__edit-mode">
          <div className="vacancy-page__wrapper__edit-mode__close-buttom">
            <ButtonComponent
              text={<X size="22" color="green" />}
              variant="close-button"
              clickFunction={handleEditMode}
            />
          </div>
          <div className="vacancy-page__wrapper__edit-mode__header">
            Редактировать вакансию
          </div>
          <Input
            label="Название вакансии *"
            value={editedInfo.name}
            onChange={(e) => handleInputChange(e, 'name')}
            required={true}
            errorMessage="Укажите название вакансии"
          />
          {(role === 'moderator' || role === 'admin') && currentVacancy?.isCreatedByUser && (
            <Input
              label="Название компании"
              value={editedInfo.companyName}
              onChange={(e) => handleInputChange(e, 'companyName')}
              required={false}
              errorMessage="Укажите название компании"
            />
          )}
          <SelectMenuComponent
            label="Формат практики *"
            options={formats.map((format) => ({
              value: format.id.toString(),
              label: format.name,
            }))}
            value={form.formatId?.toString() || ''}
            onChange={(value) => setFormField('formatId', parseInt(value as string))}
            errorMessage="Укажите формат практики"
          />
          <SelectMenuComponent
            label="Направления *"
            searchable
            multiple
            options={courses.map((option) => ({
              value: option.id.toString(),
              label: option.name,
            }))}
            value={form.courseIds.map((id) => id.toString())}
            onChange={(values) =>
              setFormField('courseIds', (values as string[]).map((v) => parseInt(v)))
            }
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
            errorMessage="Укажите оплату проезда"
          />
          {farePaymentOtherId && form.farePaymentId === farePaymentOtherId && (
            <TextareaComponent
              label="Детали оплаты проезда *"
              value={form.farePaymentDetails}
              onChange={(e) => setFormField('farePaymentDetails', e.target.value)}
              required={false}
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
            errorMessage="Укажите оплату проживания"
          />
          {accommodationPaymentOtherId &&
            form.paymentForAccommodationId === accommodationPaymentOtherId && (
              <TextareaComponent
                label="Детали оплаты проживания *"
                value={form.paymentForAccommodationDetails}
                onChange={(e) => setFormField('paymentForAccommodationDetails', e.target.value)}
                required={false}
                errorMessage="Укажите детали оплаты проживания"
              />
            )}
          <Input
            label="Конечная дата приема заявок *"
            type="date"
            value={editedInfo.deadlineAt ? editedInfo.deadlineAt.slice(0, 10) : ''}
            onChange={(e) => handleInputChange(e, 'deadlineAt')} // Передаем field
            required={true}
            errorMessage="Укажите конечную дату приема заявок"
          />
          <KeyWordsComponent
            keywords={editedInfo.keywords}
            setKeywords={(newKeywords: string[]) => setEditedInfo((prev) => ({ ...prev, keywords: newKeywords }))}
          />
          <TextareaComponent
            label="Место работы"
            value={editedInfo.description.workplace}
            onChange={(e) => handleInputChange(e, 'workplace')}
            required={false}
          />
          <TextareaComponent
            label="Должность"
            value={editedInfo.description.position}
            onChange={(e) => handleInputChange(e, 'position')}
            required={false}
          />
          <Input
            label="Зарплата"
            value={editedInfo.description.salary}
            onChange={(e) => handleInputChange(e, 'salary')}
            required={false}
          />
          <TextareaComponent
            label="Питание"
            value={editedInfo.description.food}
            onChange={(e) => handleInputChange(e, 'food')}
            required={false}
          />
          <TextareaComponent
            label="Требования"
            value={editedInfo.description.requirements}
            onChange={(e) => handleInputChange(e, 'requirements')}
            required={false}
          />
          <TextareaComponent
            label="Условия"
            value={editedInfo.description.conditions}
            onChange={(e) => handleInputChange(e, 'conditions')}
            required={false}
          />
          <TextareaComponent
            label="Дополнительная информация"
            value={editedInfo.description.additionalInfo}
            onChange={(e) => handleInputChange(e, 'additionalInfo')}
            required={false}
          />
          <ButtonComponent
            text="Сохранить"
            variant="default"
            clickFunction={handleEditVacancy}
          />
        </div>
      ) : (
        // Остальная часть компонента без изменений
        <>
          <div className="vacancy-page">
            <h1 className="vacancy-page__title">{currentVacancy.vacancyName || '-'}</h1>
            {currentVacancy?.isCreatedByUser && currentVacancy?.companyName && (
              <h3 className="vacancy-page__company-name">{currentVacancy.companyName}</h3>
            )}
            <div className="vacancy-page__content-wrapper">
              <div className="vacancy-page__content">
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Формат проведения:</strong>{' '}
                  {currentVacancy.format || '-'}
                </div>
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Направления:</strong>{' '}
                  {Array.isArray(currentVacancy.courses) && currentVacancy.courses.length > 0
                    ? currentVacancy.courses.map((c) => c.name).join(', ')
                    : '-'}
                </div>
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Ключевые слова:</strong>{' '}
                  {Array.isArray(currentVacancy.keywords) && currentVacancy.keywords.length > 0
                    ? currentVacancy.keywords.join(', ')
                    : '-'}
                </div>
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Конец приема заявок:</strong>{' '}
                  {currentVacancy.deadlineAt.slice(0, 10) || '-'}
                </div>
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Проживание:</strong>{' '}
                  { currentVacancy.paymentForAccommodation === 'Другое'
                    ? currentVacancy.paymentForAccommodationDetails
                    : currentVacancy.paymentForAccommodation
                  }
                </div>
                <div className="vacancy-page__item">
                  <strong className="vacancy-page__label">Проезд:</strong>{' '}
                  { currentVacancy.farePayment === 'Другое'
                    ? currentVacancy.farePaymentDetails
                    : currentVacancy.farePayment
                  }
                </div>
              </div>

              {(role === 'admin' ||
                (role === 'headHR' && currentVacancy?.isCreatedByUser) ||
                (role === 'HR' && currentVacancy?.isCreatedByUser) ||
                (role === 'moderator' && currentVacancy?.isCreatedByUser)) && (
                <div className="vacancy-page__to-company-button">
                  <ButtonComponent
                    text={<Pencil size="24" color="green" />}
                    variant="default"
                    clickFunction={handleEditMode}
                  />
                  <ButtonComponent
                    variant="delete"
                    text="Удалить"
                    clickFunction={() => handleDelete(Number(id))}
                    disabled={isLoading}
                  />
                </div>
              )}

              {((role === 'student' && currentVacancy?.hasCompanyProfile) || currentVacancy?.isCreatedByUser === false) && (
                <div className="vacancy-page__to-company-button">
                  <ButtonComponent
                    variant="outline"
                    text="Компания"
                    clickFunction={moveToCompany}
                    disabled={isLoading}
                  />
                </div>
              )}

              {role === 'student' && (
                <div className="vacancy-page__student-buttons">
                  {!isReplied ? (
                    <ButtonComponent
                      variant="default"
                      text="Откликнуться"
                      clickFunction={() => handleReply(Number(id))}
                      disabled={isLoading}
                    />
                  ) : (
                    <ButtonComponent
                      variant="decline"
                      text="Отменить отклик"
                      clickFunction={() => handleDecline(Number(id))}
                      disabled={isLoading}
                    />
                  )}
                </div>
              )}

              <div className="vacancy-page__additional-info">
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Место работы:</strong>{' '}
                  {currentVacancy.description.workplace || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Должность:</strong>{' '}
                  {currentVacancy.description.position || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Зарплата:</strong>{' '}
                  {currentVacancy.description.salary || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Требования:</strong>{' '}
                  {currentVacancy.description.requirements || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Питание:</strong>{' '}
                  {currentVacancy.description.food || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Условия труда:</strong>{' '}
                  {currentVacancy.description.conditions || '-'}
                </div>
                <div className="vacancy-page__additional-info__item">
                  <strong className="vacancy-page__label">Дополнительная информация:</strong>{' '}
                  {currentVacancy.description.additionalInfo || '-'}
                </div>
              </div>
            </div>
          </div>
          { (role === 'moderator' || role === 'admin') && (
            <>
              <div className="vacancy-page__hr-block">
                <div className="vacancy-page__hr-block__label">Опубликовал</div>
                <div className="vacancy-page__hr-block__info">
                  <div className="vacancy-page__hr-block__info__username">
                    <Link2 size={22} />
                    @{currentVacancy.hrInfo?.username}
                  </div>
                  <div className="vacancy-page__hr-block__info__phone-number">
                    {currentVacancy.hrInfo?.phoneNumber}
                  </div>
                </div>
              </div>

              <div className="vacancy-page__replies-block">
                <div className="vacancy-page__replies-block__label">Откликнувшиеся студенты:</div>
                {currentVacancy.repliedStudents && currentVacancy.repliedStudents.length > 0 ? (
                  <ul className="vacancy-page__replies-block__list">
                    {currentVacancy.repliedStudents.map((student, index) => (
                      <li key={index} className="vacancy-page__replies-block__list__list-item">
                        @{student.username}
                        {student.phoneNumber ? ` (${student.phoneNumber})` : '-'}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <div className="vacancy-page__replies-block__empty">Нет откликнувшихся студентов</div>
                )}
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
};

export default VacancyPage;