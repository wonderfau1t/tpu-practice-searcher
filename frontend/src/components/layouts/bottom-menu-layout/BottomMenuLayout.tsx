import {Outlet} from "react-router-dom";
import BottomMenuComponent from "../../shared/bottom-menu/BottomMenuComponent.tsx";
import './BottomMenuLayout.scss';
const BottomMenuLayout = () => {
  return (
    <div className="main-layout">
      <div className="main-layout-content">
        <Outlet />
      </div>
      <div className="bottom-menu">
        <BottomMenuComponent />
      </div>
    </div>
  );
};

export default BottomMenuLayout;