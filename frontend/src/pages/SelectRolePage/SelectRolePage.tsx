import { useNavigate } from "react-router-dom";
import { useAuthStore } from "../../stores/authStore.ts";
import axios from "axios";
import { API_URL } from "../../api/instance.ts";
import { tg } from "../../lib/telegram.ts";
import './SelectRolePage.scss';

const SelectRolePage = () => {
  const navigate = useNavigate();
  const { setRole } = useAuthStore();
  tg.backgroundColor
  const handleRegStudent = async () => {
    try {
      const initData = tg.initData;
      const response = await axios.get(`${API_URL}/register`, {
        headers: {
          Authorization: `tma ${initData}`
        },
      });
      localStorage.setItem('accessToken', response.data.result.accessToken);
      setRole(response.data.result.role);
      navigate('/test');
    } catch (error) {
      console.error(error);
    }
  }

  const handleRegCompany = () => {
    navigate('/companyRegistration');
  }
  return (
    <>
      <main className="select-role">
        <div className="select-role__header">Как вы хотите войти?</div>
        <div className="select-role__roles">
          <button className="select-role__button" onClick={handleRegStudent}>Студент</button>
          <button className="select-role__button" onClick={handleRegCompany}>Компания</button>
        </div>
      </main>
    </>
  );
}

export default SelectRolePage;