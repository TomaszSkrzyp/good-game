import { Component, createSignal, Show, For } from "solid-js";
import { auth } from "../data/store";

interface Team {
  abbreviation: string;
  teamName: string;
}

export interface Game {
  id: number;
  gameTime: string;
  homeTeamPoints: number;
  awayTeamPoints: number;
  homeTeam: Team;
  awayTeam: Team;
  rating?: number;
  avgRating?: number;
  ratingCount?: number;
}

interface GameItemProps {
  game: Game;
  onRate: (gameId: number, rating: string) => Promise<void>;
}

const GameItem: Component<GameItemProps> = (props) => {
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  
  const isFinished = () => props.game.homeTeamPoints > 0 || props.game.awayTeamPoints > 0;

  const handleRatingChange = async (val: string) => {
    setIsSubmitting(true);
    await props.onRate(props.game.id, val);
    setIsSubmitting(false);
  };

  return (
    <div class="bg-white border border-gray-200 rounded-xl shadow-sm hover:shadow-md transition-shadow p-5 flex flex-col justify-between relative">
      <Show when={isSubmitting()}>
        <div class="absolute inset-0 bg-white/50 z-10 flex items-center justify-center rounded-xl">
          <div class="animate-spin h-5 w-5 border-2 border-blue-500 border-t-transparent rounded-full"></div>
        </div>
      </Show>

      <div>
        <div class="flex justify-between items-center mb-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">
          <span>{new Date(props.game.gameTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
          <Show when={isFinished()}>
            <div class="flex items-center gap-2">
              <Show when={props.game.avgRating}>
                <span class="text-blue-600 bg-blue-50 px-2 py-0.5 rounded">
                  ★ {props.game.avgRating?.toFixed(1)} ({props.game.ratingCount})
                </span>
              </Show>
              <span class="bg-green-100 text-green-700 px-2 py-0.5 rounded">Final</span>
            </div>
          </Show>
        </div>

        <div class="space-y-3 mb-4">
          {/* away team */}
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <span class="w-8 font-bold text-gray-500 text-sm">{props.game.awayTeam?.abbreviation}</span>
              <span class="font-medium text-gray-800">{props.game.awayTeam?.teamName}</span>
            </div>
            <span class={`text-xl font-bold ${props.game.awayTeamPoints > props.game.homeTeamPoints ? 'text-black' : 'text-gray-400'}`}>
              {props.game.awayTeamPoints}
            </span>
          </div>

          {/* home team */}
          <div class="flex justify-between items-center">
            <div class="flex items-center gap-3">
              <span class="w-8 font-bold text-gray-500 text-sm">{props.game.homeTeam?.abbreviation}</span>
              <span class="font-medium text-gray-800">{props.game.homeTeam?.teamName}</span>
            </div>
            <span class={`text-xl font-bold ${props.game.homeTeamPoints > props.game.awayTeamPoints ? 'text-black' : 'text-gray-400'}`}>
              {props.game.homeTeamPoints}
            </span>
          </div>
        </div>
      </div>

      <Show when={isFinished() && auth.isLoggedIn}>
        <div class="mt-4 pt-4 border-t flex items-center justify-between">
          <span class="text-sm font-semibold text-gray-600">Your Rating:</span>
          <div class="flex gap-2 items-center">
            <select 
              disabled={isSubmitting()}
              class="text-sm border rounded px-1 py-1 bg-gray-50 focus:ring-2 focus:ring-blue-500 outline-none"
              onChange={(e) => handleRatingChange(e.currentTarget.value)}
              value={props.game.rating || ""}
            >
              <option value="" disabled>Rate</option>
              <For each={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}>{(val) => 
                <option value={val}>{val}</option>
              }</For>
            </select>
            <Show when={props.game.rating}>
              <span class="text-blue-600 font-bold text-sm">★ {props.game.rating}</span>
            </Show>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default GameItem;