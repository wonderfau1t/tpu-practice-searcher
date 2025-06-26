import {useNavigate, useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import Loader from "../../components/shared/loader/Loader.tsx";
import {useApplicationsStore} from "../../stores/applicationsStore.ts";
import './CompanyApplicationPage.scss';
import ButtonComponent from "../../components/shared/button/ButtonComponent.tsx";
import TextareaComponent from "../../components/shared/textarea/TextareaComponent.tsx";
import {toast} from "react-toastify";

const CompanyApplicationPage = () => {
  const { id } = useParams<{ id: string }>();
  const [declineMessage, setDeclineMessage] = useState<string>('');
  const [isDeclineDetailsShown, setIsDeclineDetailsShown] = useState(false);
  const navigate = useNavigate();
  const {
    selectedApplication,
    approveCompanyApplication,
    rejectCompanyApplication,
    isLoading,
    error,
    fetchCompanyApplicationByID
  } = useApplicationsStore();
  const [formErrors, setFormErrors] = useState<{ name: boolean }>({ name: false });

  useEffect(() => {
    if (id && !isNaN(Number(id))) {
      fetchCompanyApplicationByID(Number(id)).catch((err) => {
        console.error('Failed to fetch company info:', err);
      });
    } else {
      console.error('Invalid company ID:', id);
    }
  }, [id, fetchCompanyApplicationByID]);


  if (isLoading) return <Loader isLoading={isLoading} delay={300} />
  if (error) return <div>Ошибка{}: {error }</div>;

  const handleShowDeclineDetails = () => {
    setIsDeclineDetailsShown(!isDeclineDetailsShown);
  }

  const handleRejectApplication = () => {
    if (!declineMessage.trim()) {
      toast.error("Пожалуйста, опишите причины отклонения!");
      setFormErrors({ name: true })
      return;
    }

    setFormErrors({ name: false });

    if (id && !isNaN(Number(id)) && declineMessage.trim().length > 0) {
      rejectCompanyApplication(Number(id), declineMessage).catch((err) => {
        console.log('Failed to reject application:', err);
      });
      navigate('/admin/companiesApplicationsList');
    }
  }

  const handleApproveApplication = () => {
    approveCompanyApplication(Number(id)).catch((err) => {
      console.log('Failed to approve application:', err);
    })
    navigate('/admin/companiesApplicationsList');
  }


  return (
    <div className="company-application">
      <div className="company-application__header">
        Заявка компании
      </div>
      <div className="company-application__company-info">
        <div className="company-application__company-info__item">
          <div className="company-application__company-info__item__label">
            Название:
          </div>
          <div>
            {selectedApplication?.name}
          </div>
        </div>

        <div className="company-application__company-info__item">
          <div className="company-application__company-info__item__label">
            Описание:
          </div>
          <div>
            {selectedApplication?.description}
          </div>
        </div>

        <div className="company-application__company-info__item">
          <div className="company-application__company-info__item__label">
            Ссылка на сайт:
          </div>
          <div>
            {selectedApplication?.link}
          </div>
        </div>

        <div className="company-application__company-info__item">
          <div className="company-application__company-info__item__label">
            Дата регистрации:
          </div>
          <div>
            {selectedApplication?.registeredAt}
          </div>
        </div>

        <div className="company-application__company-info__item">
          <div className="company-application__company-info__item__label">
            Отправивший HR:
          </div>
          <div>
            {selectedApplication?.HRUsername}
          </div>
        </div>
      </div>
      <div className="company-application__approving-block">
        <ButtonComponent
          text="Подтвердить"
          clickFunction={handleApproveApplication}
          variant="default"
        />
        <ButtonComponent
          text="Отклонить"
          variant="decline"
          clickFunction={handleShowDeclineDetails}
        />
      </div>
      { isDeclineDetailsShown ? (
        <div className="company-application__details">
          <TextareaComponent
            label="Причины отклонения"
            value={declineMessage}
            onChange={(e) => setDeclineMessage(e.target.value)}
            required={false}
            error={formErrors?.name}
            errorMessage="Укажите причины отклонения"
          />
          <ButtonComponent
            text="Отклонить вакансию"
            variant="decline"
            clickFunction={() => handleRejectApplication()}
          />
        </div>
      ) : null}
    </div>
  )
}

export default CompanyApplicationPage;