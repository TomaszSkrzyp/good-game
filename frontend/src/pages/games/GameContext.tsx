import { createContext, useContext, createSignal, createEffect, onMount, ParentComponent } from "solid-js";
import { auth } from "../../data/store";

interface GameContextType {
  hideScores: () => boolean;
  toggleHideScores: () => void;
}

const GameContext = createContext<GameContextType | undefined>(undefined);

export const GameProvider: ParentComponent = (props) => {
  const [hideScores, setHideScores] = createSignal<boolean>(false);

  onMount(async () => {
    try {
      if (auth.isLoggedIn && auth.token) {
        console.log("Fetching user settings for hideScores beacuse logged in");
        const res = await fetch("/api/userSettings", {
          headers: { Authorization: `Bearer ${auth.token}` },
        });
        if (res.ok) {
          const json = await res.json();
          if (typeof json.hideScores === "boolean") {
            setHideScores(json.hideScores);
            console.log("hideScores value from DB:", json.hideScores);
          }
          return; 
        }
      }else{
        
        console.log("fetching from local");
      }

      // Fall back to localStorage if not logged in or fetch failed
      const stored = typeof window !== "undefined" ? localStorage.getItem("hideScores") : null;
      if (stored !== null) setHideScores(stored === "1");
    } catch {
      /* ignore */
    }
  });

  const toggleHideScores = async () => {
    const next = !hideScores();
    setHideScores(next);

    // Persist to DB if logged in
    if (auth.isLoggedIn && auth.token) {
      try {
        await fetch("/api/userSettings", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${auth.token}`,
          },
          body: JSON.stringify({ hideScores: next }),
        });
      } catch {
        /* ignore network errors */
      }
    } else {
      // Fall back to localStorage if not logged in
      try {
        if (typeof window !== "undefined") localStorage.setItem("hideScores", next ? "1" : "0");
      } catch {
        /* ignore */
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