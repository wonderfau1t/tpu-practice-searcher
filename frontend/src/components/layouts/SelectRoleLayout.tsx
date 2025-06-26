import { Outlet } from 'react-router-dom';

export const SelectRoleLayout = () => {
  return (
    <div className="main-layout">
      <main>
        <Outlet />
      </main>
    </div>
  );
};