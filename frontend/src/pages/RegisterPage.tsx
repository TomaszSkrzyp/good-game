import { createSignal, Component } from "solid-js";
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
    
    // Walidacja biznesowa (Frontend)
    if (username().length < 3) {
      setError("Username must be at least 3 characters");
      return;
    }
    if (password().length < 6) {
      setError("Password must be at least 6 characters");
      return;
    }
    if (password() !== confirmPassword()) {
      setError("Passwords do not match!");
      return;
    }

    setLoading(true);
    try {
      const response = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ 
          userName: username(), // Upewnij się, że pasuje do struktury w Go (userName)
          email: email(),
          password: password() 
        }),
      });

      const data = await response.json();

      if (response.ok) {
        alert("Account created! Please log in.");
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
    <div class="auth-container">
      <h1>Create Account</h1>
      {error() && <p style="color: red; background: #fee; padding: 10px; border-radius: 4px;">{error()}</p>}
      
      <form onSubmit={handleRegister}>
        <label for="reg-username">Username</label>
        <input 
          type="text" 
          id="reg-username"
          onInput={(e) => setUsername(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <label for="reg-email">Email</label>
        <input 
          type="email" 
          id="reg-email"
          onInput={(e) => setEmail(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <label for="reg-password">Password</label>
        <input 
          type="password" 
          id="reg-password"
          onInput={(e) => setPassword(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <label for="confirm-password">Confirm Password</label>
        <input 
          type="password" 
          id="confirm-password"
          onInput={(e) => setConfirmPassword(e.currentTarget.value)} 
          required 
          disabled={loading()}
        />

        <div class="button-group">
          <button type="submit" disabled={loading()}>
            {loading() ? "Creating..." : "Register"}
          </button>
          <button type="button" onClick={() => navigate("/")} disabled={loading()}>
            Back to Login
          </button>
        </div>
      </form>
    </div>
  );
};

export default RegisterPage;