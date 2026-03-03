import { clearSession } from "../../data/store";

const LogoutButton = () => {
  const handleLogout = () => {
    clearSession();
    window.location.href = "/";
  };

  return <button onClick={handleLogout}>Logout</button>;
};
export default LogoutButton