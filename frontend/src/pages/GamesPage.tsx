import { Component, createSignal, createResource, For, Show } from "solid-js";
import { auth } from "../data/store";
import GameItem, { Game } from "../components/GameItem";

const GamesPage: Component = () => {
  const todayStr = () => {
  const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  };
  const [day, setDay] = createSignal(todayStr());

  const fetchGames = async (date: string): Promise<Game[]> => {
    const headers: Record<string, string> = {};
    if (auth.token) headers["Authorization"] = `Bearer ${auth.token}`;

    const res = await fetch(`/api/games?date=${encodeURIComponent(date)}`, { headers });
    if (!res.ok) throw new Error("failed to fetch");
    return res.json();
  };

  const [games, { mutate }] = createResource(day, fetchGames);

  const shiftDay = (delta: number) => {
    const d = new Date(day());
    d.setDate(d.getDate() + delta);
    setDay(d.toISOString().slice(0, 10));
  };

  const submitRating = async (gameId: number, rating: string) => {
    const token = auth.token;
    if (!token) return;

    const val = parseInt(rating);

    try {
      const res = await fetch(`/api/userReactions`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
        body: JSON.stringify({ gameId, rating: val }),
      });

      if (res.ok) {
        mutate((prev: Game[] | undefined) => 
          prev?.map(g => {
            if (g.id === gameId) {
              // optimistic update of rating count
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
        <div class="flex items-center gap-2">
          <button onClick={() => shiftDay(-1)} class="p-2 hover:bg-gray-100 rounded-full border">←</button>
          <input
            type="date"
            value={day()}
            onInput={(e) => setDay(e.currentTarget.value)}
            class="px-3 py-2 border rounded-lg outline-none"
          />
          <button onClick={() => shiftDay(1)} class="p-2 hover:bg-gray-100 rounded-full border">→</button>
          <button onClick={() => setDay(todayStr())} class="ml-2 px-3 py-2 bg-gray-200 rounded-lg text-sm">Today</button>
        </div>
      </div>

      <Show when={games.loading}>
        <div class="flex justify-center p-10 text-gray-500 italic">loading...</div>
      </Show>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <For each={games()}>
          {(game) => <GameItem game={game} onRate={submitRating} />}
        </For>
      </div>

      <Show when={!games.loading && games()?.length === 0}>
        <div class="text-center py-20 text-gray-400">no games found</div>
      </Show>
    </div>
  );
};

export default GamesPage;