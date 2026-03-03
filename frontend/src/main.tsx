
import './index.css';
import { render } from "solid-js/web";
import App from "./App";
import { GameProvider } from "./pages/games/GameContext";

render(() => (
  <GameProvider>
    <App />
  </GameProvider>
), document.getElementById("app") as HTMLElement);