import { Router, Route } from "@solidjs/router";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import ProtectedRoute from "./components/ProtectedRoute";
import ProfilePage from "./pages/profile/ProfilePage";
import GamesPage from "./pages/games/GamesPage";
import Layout from "./components/Layout";
import { todayStr } from "./utils/dateUtils";
import { Navigate } from "@solidjs/router";

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
      <Route
    path="/games/:date"
    component={() => (
      <Layout>
        <GamesPage />
      </Layout>
    )}
  />
      <Route path="/games" component={() => <Navigate href={`/games/${todayStr()}`} />} /> 
    </Router>
  );
};

export default App;
