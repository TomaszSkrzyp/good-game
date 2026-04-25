import { Component, Show } from "solid-js";
import { A } from "@solidjs/router";
import LogoutButton from "../pages/auth/Logout";
import { auth } from "../data/auth";

const Header: Component = () => {
  return (
    <header class="flex items-center justify-between p-4 bg-gray-800 text-white shadow-md">
      <A href="/games" class="text-xl font-extrabold tracking-tight hover:text-blue-400 transition-colors">
        Good Game
      </A>
      <A href="/how-it-works" class="text-sm font-medium hover:text-blue-500">
        Algorithm
      </A>
      <div class="flex items-center gap-4">
        <Show 
          when={auth.isLoggedIn} 
          fallback={
            <div class="flex items-center gap-2">
              <A 
                href="/" 
                class="px-4 py-2 text-sm font-medium hover:text-blue-400 transition-colors"
              >
                Sign In
              </A>
              <A 
                href="/register" 
                class="px-4 py-2 text-sm font-bold bg-blue-600 hover:bg-blue-700 rounded-lg transition-all shadow-sm"
              >
                Sign Up
              </A>
            </div>
          }
        >
          <div class="flex items-center gap-3">
            <span class="hidden md:block text-sm text-gray-400 mr-2">
              Hi, <span class="text-white font-medium">{auth.userName}</span>
            </span>
            <A
              href="/profile"
              class="p-2 rounded-full bg-gray-700 hover:bg-gray-600 flex items-center justify-center transition-colors border border-gray-600"
              aria-label="Profile"
              title="Profile"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </A>
            <LogoutButton />
          </div>
        </Show>
      </div>
    </header>
  );
};

export default Header;