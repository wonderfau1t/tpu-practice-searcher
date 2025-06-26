import {useAuthStore} from "../../stores/authStore.ts";
import './CompanyUnderReviewPage.scss';

const CompanyUnderReviewPage = () => {
  const { errorMessage } = useAuthStore();
  return (
    <div className="company-under-review">
      <div>
        { errorMessage }
      </div>
    </div>
  )
}

export default CompanyUnderReviewPage;