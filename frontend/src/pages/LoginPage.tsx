import { createSignal, Component, onMount, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { loginUser, auth } from "../data/store";

const LoginPage: Component = () => {
  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const navigate = useNavigate();

  onMount(() => {
    if (auth.isLoggedIn) navigate("/profile");
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    if (username().trim().length < 3) {
      setError("Username is too short");
      return;
    }
    setLoading(true);
    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userName: username(), password: password() }),
      });
      const data = await response.json();
      if (response.ok) {
        loginUser(data); 
        navigate("/games"); 
      } else {
        setError(data.error || "Login failed");
      }
    } catch (err) {
      setError("Connection error. Is the server running?");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="max-w-md mx-auto mt-12 p-6">
      <div class="bg-white shadow-md rounded-xl p-8 border border-gray-100">
        <h1 class="text-2xl font-bold text-gray-800 mb-6 text-center">Login</h1>
        <Show when={error()}>
          <div class="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm border border-red-100">
            {error()}
          </div>
        </Show>
        <form onSubmit={handleSubmit} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Username</label>
            <input 
              type="text" 
              onInput={(e) => setUsername(e.currentTarget.value)} 
              class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition-all disabled:opacity-50"
              required 
              disabled={loading()}
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Password</label>
            <input 
              type="password" 
              onInput={(e) => setPassword(e.currentTarget.value)} 
              class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none transition-all disabled:opacity-50"
              required 
              disabled={loading()}
            />
          </div>
          <button 
            type="submit" 
            disabled={loading()}
            class="w-full bg-blue-600 text-white py-2 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:bg-blue-300"
          >
            {loading() ? "Logging in..." : "Login"}
          </button>
          <div class="relative py-4">
            <div class="absolute inset-0 flex items-center"><span class="w-full border-t border-gray-200"></span></div>
            <div class="relative flex justify-center text-xs uppercase"><span class="bg-white px-2 text-gray-500">Or</span></div>
          </div>
          <button 
            type="button" 
            onClick={() => navigate("/register")} 
            disabled={loading()}
            class="w-full border border-gray-300 text-gray-600 py-2 rounded-lg font-semibold hover:bg-gray-50 transition-colors"
          >
            Create Account
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;