import { logout } from "../../data/auth";

const LogoutButton = () => {
  const handleLogout = async () => {
    await logout(); 
  };

  return (
    <button 
      onClick={handleLogout}
      class="text-sm text-gray-500 hover:text-red-600 transition"
    >
      Logout
    </button>
  );
};
export default LogoutButton