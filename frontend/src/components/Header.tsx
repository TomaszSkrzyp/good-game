import { Component } from "solid-js"
import LogoutButton from "./Logout"

const Header: Component = () => {
  const goToProfile = () => (window.location.href = "/profile");

  return (
    <header class="flex items-center justify-between p-4 bg-gray-800 text-white">
      <div class="text-lg font-bold">Good Game</div>
      <div class="flex items-center gap-3">
        <button
          type="button"
          onClick={goToProfile}
          class="p-2 rounded-full bg-gray-700 hover:bg-gray-600 flex items-center justify-center"
          aria-label="Profile"
          title="Profile"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"  aria-hidden="true">
            <path {...{ "stroke-linecap": "round", "stroke-linejoin": "round", d: "M5.121 17.804A9 9 0 1118.879 6.196 9 9 0 015.121 17.804z" }} />
            <path {...{ "stroke-linecap": "round", "stroke-linejoin": "round", d: "M15 11a3 3 0 11-6 0 3 3 0 016 0z" }} />
          </svg>
        </button>
        <LogoutButton />
      </div>
    </header>
  );
};

export default Header;