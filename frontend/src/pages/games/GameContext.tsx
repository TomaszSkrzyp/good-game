import { createContext, useContext, createSignal, createEffect, ParentComponent } from "solid-js";
import { auth } from "../../data/store";
import { api } from "../../utils/api";

interface GameContextType {
  hideScores: () => boolean;
  toggleHideScores: () => void;
}

interface UserSettings {
  hideScores: boolean;
}

const GameContext = createContext<GameContextType | undefined>(undefined);

export const GameProvider: ParentComponent = (props) => {
  const [hideScores, setHideScores] = createSignal<boolean>(false);

  createEffect(async () => {
    
    if (auth.isLoggedIn) {
      try {
        console.log("User logged in. Fetching settings from API...");
        const settings = await api.get("userSettings").json<UserSettings>();
        
        if (typeof settings.hideScores === "boolean") {
          setHideScores(settings.hideScores);
        }
      } catch (error) {
        console.warn("Failed to fetch remote settings, falling back to local.", error);
      }
    } else {
      if (typeof window !== "undefined") {
        const stored = localStorage.getItem("hideScores");
        if (stored !== null) {
          setHideScores(stored === "1");
        } else {
          setHideScores(false); 
        }
      }
    }
  });

  const toggleHideScores = async () => {
    const next = !hideScores();
    setHideScores(next);

    if (auth.isLoggedIn) {
      try {
        await api.post("userSettings", {
          json: { hideScores: next },
        });
      } catch (error) {
        console.error("Failed to persist setting to database:", error);
      }
    } else {
      if (typeof window !== "undefined") {
        localStorage.setItem("hideScores", next ? "1" : "0");
      }
    }
  };

  return (
    <GameContext.Provider value={{ hideScores, toggleHideScores }}>
      {props.children}
    </GameContext.Provider>
  );
};

export const useGameContext = () => {
  const ctx = useContext(GameContext);
  if (!ctx) throw new Error("useGameContext must be used within a GameProvider");
  return ctx;
};

export default GameContext;