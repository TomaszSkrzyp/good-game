import { Router, Route } from "@solidjs/router";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import ProtectedRoute from "./components/ProtectedRoute";
import ProfilePage from "./pages/ProfilePage";
import GamesPage from "./pages/games/GamesPage";
import Layout from "./components/Layout";

const App = () => {
  return (
    <Router>
      <Route path="/" component={() => <Layout><LoginPage /></Layout>} />
      <Route path="/register" component={() => <Layout><RegisterPage /></Layout>} />
      <Route
        path="/profile"
        component={() => (
          <Layout>
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          </Layout>
        )}
      />
      <Route path="/games" component={() => <Layout><GamesPage /></Layout>} />
    </Router>
  );
};

export default App;
