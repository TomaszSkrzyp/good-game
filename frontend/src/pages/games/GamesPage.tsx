import { Component, createResource, For, Show, createEffect } from "solid-js";
import { createStore, reconcile } from "solid-js/store";
import { useParams, useNavigate } from "@solidjs/router";
import { api } from "../../utils/api"; // Your Ky instance
import GameItem, { Game } from "./GameItem";
import { useGameContext } from "./GameContext";
import { todayStr } from "../../utils/dateUtils";

const GamesPage: Component = () => {
  const params = useParams();
  const navigate = useNavigate();
  const { hideScores, toggleHideScores } = useGameContext();
  const day = () => params.date || todayStr();

  const [gamesResource] = createResource(day, async (date) => {
    return await api.get('games', { 
      searchParams: { date } 
    }).json<Game[]>();
  });

  const [gamesStore, setGamesStore] = createStore<{ list: Game[] }>({ list: [] });

  createEffect(() => {
    const data = gamesResource();
    if (data) {
      setGamesStore("list", reconcile(data));
    }
  });

  createEffect(() => {
    if (!params.date) navigate(`/games/${todayStr()}`, { replace: true });
  });

  const shiftDay = (delta: number) => {
    const d = new Date(day() + 'T12:00:00');
    d.setDate(d.getDate() + delta);
    navigate(`/games/${d.toISOString().slice(0, 10)}`);
  };

  const submitRating = async (gameId: number, rating: string) => {
    const val = parseInt(rating);
    try {
      await api.post('userReactions', {
        json: { gameId, rating: val },
      });
      setGamesStore("list", (g) => g.id === gameId, { rating: val });
    } catch (err) {
      console.error("rating failed:", err);
    }
  };

  return (
    <div class="max-w-4xl mx-auto p-4">
      <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
        <h1 class="text-3xl font-bold text-gray-800">NBA Schedule</h1>
        
        <div class="flex items-center gap-4">
          <label class="flex items-center gap-2 text-sm select-none cursor-pointer">
            <span class="text-xs text-gray-600">Hide scores</span>
            <div class="relative">
              <input 
                type="checkbox" 
                checked={hideScores()} 
                onInput={() => toggleHideScores()} 
                class="sr-only" 
              />
              <div class={`w-10 h-6 rounded-full transition-colors ${hideScores() ? 'bg-blue-600' : 'bg-gray-300'}`}></div>
              <div class={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transform transition-transform ${hideScores() ? 'translate-x-4' : ''}`}></div>
            </div>
          </label>

          <div class="flex items-center gap-2">
            <button onClick={() => shiftDay(-1)} class="p-2 hover:bg-gray-100 rounded-full border border-gray-200 cursor-pointer">←</button>
            <input 
              type="date" 
              value={day()} 
              onInput={(e) => navigate(`/games/${e.currentTarget.value}`)} 
              class="px-3 py-2 border rounded-lg outline-none" 
            />
            <button onClick={() => shiftDay(1)} class="p-2 hover:bg-gray-100 rounded-full border border-gray-200 cursor-pointer">→</button>
            <button onClick={() => navigate(`/games/${todayStr()}`)} class="ml-2 px-3 py-2 bg-gray-200 rounded-lg text-sm cursor-pointer">Today</button>
          </div>
        </div>
      </div>

      <Show when={gamesResource.loading}>
        <div class="flex justify-center p-10 text-gray-500 italic">loading...</div>
      </Show>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <For each={gamesStore.list} fallback={<div class="col-span-full text-center py-10 text-gray-400">No games scheduled for this date.</div>}>
          {(game) => <GameItem game={game} onRate={submitRating} hideScores={hideScores()} />}
        </For>
      </div>
    </div>
  );
};

export default GamesPage;