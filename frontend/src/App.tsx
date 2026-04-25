import { Router, Route } from "@solidjs/router";
import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import ProtectedRoute from "./components/ProtectedRoute";
import ProfilePage from "./pages/profile/ProfilePage";
import GamesPage from "./pages/games/GamesPage";
import Layout from "./components/Layout";
import { todayStr } from "./utils/dateUtils";
import { Navigate } from "@solidjs/router";
import HowItWorks from "./pages/algorithm/HowItWorks";

const App = () => {
  return (
    <Router>
      <Route component={Layout}> 
        <Route path="/" component={LoginPage} />
        <Route path="/register" component={RegisterPage} />
        <Route path="/how-it-works" component={HowItWorks} />
        <Route path="/profile" component={() => (
          <ProtectedRoute>
            <ProfilePage />
          </ProtectedRoute>
        )} />
        
        <Route path="/games/:date" component={GamesPage} />
        <Route path="/games" component={() => <Navigate href={`/games/${todayStr()}`} />} />
      </Route>
    </Router>
  );
};

export default App;
