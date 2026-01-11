import { createSignal, Component, onMount } from "solid-js";
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
    
    // Walidacja wejściowa (Frontend)
    if (username().trim().length < 3) {
      setError("Username is too short");
      return;
    }

    setLoading(true);
    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
          userName: username(), 
          password: password() 
        }),
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
    <div class="auth-container">
      <h1>Login</h1>
      {error() && <p style="color: red; background: #fee; padding: 10px; border-radius: 4px;">{error()}</p>}
      
      <form onSubmit={handleSubmit}>
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username"
          onInput={(e) => setUsername(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <label for="password">Password</label>
        <input 
          type="password" 
          id="password"
          onInput={(e) => setPassword(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <div class="button-group">
          <button type="submit" disabled={loading()}>
            {loading() ? "Logging in..." : "Login"}
          </button>
          <button type="button" onClick={() => navigate("/register")} disabled={loading()}>
            Create Account
          </button>
        </div>
      </form>
    </div>
  );
};

export default LoginPage;