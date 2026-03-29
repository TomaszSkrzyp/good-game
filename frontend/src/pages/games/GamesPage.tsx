import { Component, createResource, For, Show, createEffect } from "solid-js";
import { useParams, useNavigate } from "@solidjs/router";
import { fetchWithAuth } from "../../data/store";
import GameItem, { Game } from "./GameItem";
import { useGameContext } from "./GameContext";
import { todayStr } from "../../utils/dateUtils";

const GamesPage: Component = () => {
  const params = useParams();
  const navigate = useNavigate();
  const { hideScores, toggleHideScores } = useGameContext();

  

  const day = () => params.date || todayStr();

  // Redirect to today's date if no date in URL
  createEffect(() => {
    if (!params.date) {
      navigate(`/games/${todayStr()}`, { replace: true });
    }
  });

  const fetchGames = async (date: string): Promise<Game[]> => {
    const res = await fetchWithAuth(`/api/games?date=${encodeURIComponent(date)}`);
    if (!res.ok) throw new Error("failed to fetch");
    return res.json();
  };

  const [games, { mutate }] = createResource(day, fetchGames);

  const shiftDay = (delta: number) => {
    const d = new Date(day() + 'T12:00:00'); // Noon local to avoid timezone issues
    d.setDate(d.getDate() + delta);
    const newDate = d.toISOString().slice(0, 10);
    navigate(`/games/${newDate}`);
  };

  const submitRating = async (gameId: number, rating: string) => {
    const val = parseInt(rating);
    try {
      const res = await fetchWithAuth(`/api/userReactions`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ gameId, rating: val }),
      });
      if (res.ok) {
        mutate((prev: Game[] | undefined) => 
          prev?.map(g => {
            if (g.id === gameId) {
              const isFirstVote = !g.rating || g.rating === 0;
              const newCount = isFirstVote ? (g.ratingCount || 0) + 1 : (g.ratingCount || 0);
              return { ...g, rating: val, ratingCount: newCount };
            }
            return g;
          })
        );
      }
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div class="max-w-4xl mx-auto p-4">
      <div class="flex flex-col md:flex-row justify-between items-center mb-6 gap-4">
        <h1 class="text-3xl font-bold text-gray-800">NBA Schedule</h1>
        <label class="ml-4 flex items-center gap-2 text-sm select-none cursor-pointer">
          <span class="text-xs text-gray-600">Hide scores</span>
          <div class="relative">
            <input
              type="checkbox"
              checked={hideScores()}
              onInput={(e) => toggleHideScores()}
              class="sr-only"
            />
            <div class={`w-10 h-6 rounded-full transition-colors ${hideScores() ? 'bg-blue-600' : 'bg-gray-300'}`}></div>
            <div class={`absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transform transition-transform ${hideScores() ? 'translate-x-4' : ''}`}></div>
          </div>
        </label>
        <div class="flex items-center gap-2">
          <button onClick={() => shiftDay(-1)} class="p-2 hover:bg-gray-100 rounded-full border cursor-pointer">←</button>
          <input
            type="date"
            value={day()}
            onInput={(e) => navigate(`/games/${e.currentTarget.value}`)}
            class="px-3 py-2 border rounded-lg outline-none"
          />
          <button onClick={() => shiftDay(1)} class="p-2 hover:bg-gray-100 rounded-full border cursor-pointer">→</button>
          <button onClick={() => navigate(`/games/${todayStr()}`)} class="ml-2 px-3 py-2 bg-gray-200 rounded-lg text-sm cursor-pointer">Today</button>
        </div>
      </div>

      <Show when={games.loading}>
        <div class="flex justify-center p-10 text-gray-500 italic">loading...</div>
      </Show>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <For each={games()}>
          {(game) => <GameItem game={game} onRate={submitRating} hideScores={hideScores()} />}
        </For>
      </div>

      <Show when={!games.loading && games()?.length === 0}>
        <div class="text-center py-20 text-gray-400">no games found</div>
      </Show>
    </div>
  );
};

export default GamesPage;