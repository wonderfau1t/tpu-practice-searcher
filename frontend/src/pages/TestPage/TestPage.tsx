import {useNavigate} from "react-router-dom";
import {useAuthStore} from "../../stores/authStore.ts";

const TestPage = () => {
  const navigate = useNavigate();
  const { logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/');
  };
  return (
    <div>
      <h1>Отлично</h1>
      <div>
        <button onClick={handleLogout}>Выйти</button>
      </div>
    </div>
  );
}

export default TestPage;