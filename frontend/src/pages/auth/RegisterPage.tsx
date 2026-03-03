import { createSignal, Component, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";

const RegisterPage: Component = () => {
  const [username, setUsername] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const navigate = useNavigate();

  const handleRegister = async (e: Event) => {
    e.preventDefault();
    setError("");
    if (username().length < 3) { setError("Username must be at least 3 characters"); return; }
    if (password().length < 6) { setError("Password must be at least 6 characters"); return; }
    if (password() !== confirmPassword()) { setError("Passwords do not match!"); return; }

    setLoading(true);
    try {
      const response = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userName: username(), email: email(), password: password() }),
      });
      const data = await response.json();
      if (response.ok) {
        navigate("/");
      } else {
        setError(data.error || "Registration failed.");
      }
    } catch (err) {
      setError("Connection error. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="max-w-md mx-auto mt-12 p-6">
      <div class="bg-white shadow-md rounded-xl p-8 border border-gray-100">
        <h1 class="text-2xl font-bold text-gray-800 mb-6 text-center">Create Account</h1>
        <Show when={error()}>
          <div class="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm border border-red-100">
            {error()}
          </div>
        </Show>
        <form onSubmit={handleRegister} class="space-y-4">
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
            <label class="block text-sm font-medium text-gray-600 mb-1">Email</label>
            <input 
              type="email" 
              onInput={(e) => setEmail(e.currentTarget.value)} 
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
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Confirm Password</label>
            <input 
              type="password" 
              onInput={(e) => setConfirmPassword(e.currentTarget.value)} 
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
            {loading() ? "Creating..." : "Register"}
          </button>
          <button 
            type="button" 
            onClick={() => navigate("/")} 
            disabled={loading()}
            class="w-full text-sm text-gray-500 hover:text-blue-600 transition-colors pt-2"
          >
            Already have an account? Back to Login
          </button>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;