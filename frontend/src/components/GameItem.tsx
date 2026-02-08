import { Component, createSignal, Show, For } from "solid-js";
import { auth } from "../data/store";

interface Team {
  abbreviation: string;
  teamName: string;
}
export interface GameQuality {
  qualityScore: number;
  isBigScoring: boolean;
  isBigGame: boolean;
  //TODO ADD INFO ABOUR BIG GAMES ON BOTH SIDES
  isClutch: boolean;
}

export interface Game {
  id: number;
  gameTime: string;
  homeTeamPoints: number;
  awayTeamPoints: number;
  homeTeam: Team;
  awayTeam: Team;
  status: string;
  gameQuality: GameQuality; 
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
  const q = () => props.game.gameQuality;

  const handleRating = async (val: number) => {
    setIsSubmitting(true);
    await props.onRate(props.game.id, val.toString());
    setIsSubmitting(false);
  };

  return (
    <div class="bg-white border border-gray-200 rounded-xl shadow-sm p-5 relative overflow-hidden flex flex-col justify-between h-full">
      <Show when={isSubmitting()}>
        <div class="absolute inset-0 bg-white/60 z-10 flex items-center justify-center backdrop-blur-sm">
          <div class="animate-spin h-6 w-6 border-2 border-blue-500 border-t-transparent rounded-full" />
        </div>
      </Show>

      <div>
        <div class="flex justify-between items-start mb-4">
          <div class="text-[10px] font-bold text-gray-400 uppercase">
            {new Date(props.game.gameTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
          </div>
          
        </div>
        <Show when={props.game.status === "STATUS_FINAL"}>
        <div class="space-y-3 mb-4">
          <TeamRow team={props.game.awayTeam} score={props.game.awayTeamPoints} isWinner={props.game.awayTeamPoints > props.game.homeTeamPoints} />
          <TeamRow team={props.game.homeTeam} score={props.game.homeTeamPoints} isWinner={props.game.homeTeamPoints > props.game.awayTeamPoints} />
        </div>

        <div class="flex flex-wrap gap-1 mb-4">
          <Show when={q().isClutch}><Badge class="bg-red-50 text-red-600">🔥 Clutch</Badge></Show>
          <Show when={q().isBigGame}><Badge class="bg-yellow-50 text-yellow-700">⭐ Big Game</Badge></Show>
          <Show when={q().isBigScoring}><Badge class="bg-purple-50 text-purple-600">🏀 Scoring</Badge></Show>
          
        </div>
        <div class={`px-2 py-1 rounded text-[11px] font-black ${q().qualityScore >= 70 ? 'bg-orange-500 text-white' : 'bg-gray-100 text-gray-600'}`}>
            {q().qualityScore} PTS
        </div>
        </Show>
      </div>
      

      <Show when={props.game.status === "STATUS_FINAL" && auth.isLoggedIn}>
        <div class="pt-4 border-t border-gray-100 flex items-center justify-between">
          <span class="text-[10px] font-black text-gray-400 uppercase">Rate:</span>
          <div class="flex gap-1">
            <For each={[1, 2, 3, 4, 5]}>{(v) => (
              <button 
                onClick={() => handleRating(v)}
                class={`w-7 h-7 rounded text-xs font-bold transition-colors ${props.game.rating === v ? 'bg-blue-600 text-white' : 'bg-gray-50 text-gray-400 hover:bg-gray-200'}`}
              >{v}</button>
            )}</For>
          </div>
        </div>
      </Show>
    </div>
  );
};

const TeamRow = (p: { team: Team, score: number, isWinner: boolean }) => (
  <div class="flex justify-between items-center">
    <div class="flex items-center gap-2">
      <span class="text-xs font-black text-gray-400 w-8">{p.team.abbreviation}</span>
      <span class={`text-sm ${p.isWinner ? 'font-bold text-gray-900' : 'text-gray-500'}`}>{p.team.teamName}</span>
    </div>
    <span class={`font-mono text-xl ${p.isWinner ? 'font-black text-gray-900' : 'text-gray-300'}`}>{p.score}</span>
  </div>
);

const Badge = (p: { children: any, class: string }) => (
  <span class={`text-[9px] font-bold px-2 py-0.5 rounded border border-current ${p.class}`}>{p.children}</span>
);

export default GameItem;