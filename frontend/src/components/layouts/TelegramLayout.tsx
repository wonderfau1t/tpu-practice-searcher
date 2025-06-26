import { useEffect } from "react";
import {Outlet, useNavigate} from "react-router-dom";
import { tg } from "../../lib/telegram.ts";

const TelegramLayout = () => {
  const navigate = useNavigate();

  useEffect(() => {
    if (location.pathname === "/" || location.pathname === "/moderator/vacanciesList" || location.pathname === "/company/vacanciesList" || location.pathname === "/student/vacanciesList") {
      tg.BackButton.hide();
    } else {
      tg.BackButton.show();
    }
    const handleBack = () => navigate(-1);

    tg.BackButton.onClick(handleBack);

    return () => {
      tg.BackButton.offClick(handleBack);
      tg.BackButton.hide();
    };
  }, [navigate, location.pathname]);

  return (
      <Outlet />
  );
};

export default TelegramLayout;
