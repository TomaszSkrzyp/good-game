import { Component, createSignal, Show, For } from "solid-js";
import { auth } from "../../data/store";

interface Team {
  abbreviation: string;
  teamName: string;
}

interface GamePlayerStats {
  homeTopScorer: string;
  homeTopScorerPts: number;
  homeTopAssister: string;
  homeTopAssists: number;
  homeTopRebounder: string;
  homeTopRebounds: number;
  awayTopScorer: string;
  awayTopScorerPts: number;
  awayTopAssister: string;
  awayTopAssists: number;
  awayTopRebounder: string;
  awayTopRebounds: number;
}

export interface GameQuality {
  qualityScore: number;
  isBigScoring: boolean;
  isBigGame: boolean;
  isClutch: boolean;
  isStarDuel: boolean;
  isHugeSwing: boolean;
  isShootout: boolean;
  isGritty: boolean;
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
  hideScores?: boolean;
}

const GameItem: Component<GameItemProps> = (props) => {
  const [isSubmitting, setIsSubmitting] = createSignal(false);
  const [showStats, setShowStats] = createSignal(false);
  const [stats, setStats] = createSignal<GamePlayerStats | null>(null);
  const [isLoadingStats, setIsLoadingStats] = createSignal(false);
  const q = () => props.game.gameQuality;

  const handleRating = async (val: number) => {
    setIsSubmitting(true);
    await props.onRate(props.game.id, val.toString());
    setIsSubmitting(false);
  };

  const fetchGameStats = async () => {
    setIsLoadingStats(true);
    try {
      const response = await fetch(`/api/game/stats?gameId=${props.game.id}`);
      if (!response.ok) throw new Error("Failed to fetch stats");
      const data = await response.json();
      setStats(data);
      setShowStats(true);
    } catch (error) {
      console.error("Error fetching game stats:", error);
    } finally {
      setIsLoadingStats(false);
    }
  };

  const closeStats = () => {
    setShowStats(false);
    document.body.style.overflow = 'auto';
  };

  const openStats = async () => {
    await fetchGameStats();
    document.body.style.overflow = 'hidden';
  };

  return (
    <div class="bg-white border border-gray-200 rounded-xl shadow-sm p-5 relative overflow-hidden flex flex-col justify-between h-full">
      <Show when={isSubmitting()}>
        <div class="absolute inset-0 bg-white/60 z-10 flex items-center justify-center backdrop-blur-sm">
          <div class="animate-spin h-6 w-6 border-2 border-blue-500 border-t-transparent rounded-full" />
        </div>
      </Show>

      <Show when={showStats()}>
        <div 
          class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center" 
          onClick={closeStats}
        >
          <div 
            class="bg-white rounded-lg shadow-xl p-6 w-96 max-h-96 overflow-y-auto relative" 
            onClick={(e) => e.stopPropagation()}
          >
            <button 
              onClick={closeStats}
              class="sticky top-0 right-0 absolute top-4 right-4 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full p-1 transition-colors"
            >
              ✕
            </button>
            
            <h3 class="font-bold text-lg mb-4 pr-6">Game Stats</h3>
            
            <Show when={isLoadingStats()}>
              <div class="flex justify-center py-8">
                <div class="animate-spin h-6 w-6 border-2 border-blue-500 border-t-transparent rounded-full" />
              </div>
            </Show>

            <Show when={stats() && !isLoadingStats()}>
              <div class="space-y-4">
                <div>
                  <h4 class="font-bold text-sm text-gray-700 mb-2">{props.game.awayTeam.teamName}</h4>
                  <div class="space-y-1 text-sm">
                    <div class="flex justify-between"><span class="text-gray-600">Top Scorer:</span><span class="font-semibold">{stats()?.awayTopScorer} <span class="text-gray-400">({stats()?.awayTopScorerPts})</span></span></div>
                    <div class="flex justify-between"><span class="text-gray-600">Top Assister:</span><span class="font-semibold">{stats()?.awayTopAssister} <span class="text-gray-400">({stats()?.awayTopAssists})</span></span></div>
                    <div class="flex justify-between"><span class="text-gray-600">Top Rebounder:</span><span class="font-semibold">{stats()?.awayTopRebounder} <span class="text-gray-400">({stats()?.awayTopRebounds})</span></span></div>
                  </div>
                </div>

                <div class="border-t pt-4">
                  <h4 class="font-bold text-sm text-gray-700 mb-2">{props.game.homeTeam.teamName}</h4>
                  <div class="space-y-1 text-sm">
                    <div class="flex justify-between"><span class="text-gray-600">Top Scorer:</span><span class="font-semibold">{stats()?.homeTopScorer} <span class="text-gray-400">({stats()?.homeTopScorerPts})</span></span></div>
                    <div class="flex justify-between"><span class="text-gray-600">Top Assister:</span><span class="font-semibold">{stats()?.homeTopAssister} <span class="text-gray-400">({stats()?.homeTopAssists})</span></span></div>
                    <div class="flex justify-between"><span class="text-gray-600">Top Rebounder:</span><span class="font-semibold">{stats()?.homeTopRebounder} <span class="text-gray-400">({stats()?.homeTopRebounds})</span></span></div>
                  </div>
                </div>
              </div>
            </Show>
          </div>
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
            <TeamRow team={props.game.awayTeam} score={props.game.awayTeamPoints} hide={props.hideScores } />
            <TeamRow team={props.game.homeTeam} score={props.game.homeTeamPoints} hide={props.hideScores} />
          </div>

          <button
            onClick={openStats}
            disabled={isLoadingStats()}
            class="w-full px-3 py-2 mb-4 text-xs font-bold text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors disabled:opacity-50"
          >
            {isLoadingStats() ? "Loading..." : "View Stats"}
          </button>

          <div class="flex flex-wrap gap-1 mb-4">
            
            <Badge class="bg-emerald-50 text-emerald-600" title="Game finished">Finished</Badge>
            <Show when={q().isClutch}><Badge class="bg-red-50 text-red-600" title="4th quarter was close or game decided by 3 points or less">Clutch</Badge></Show>
            <Show when={q().isBigGame}><Badge class="bg-yellow-50 text-yellow-700" title="Player scored 35+ points with leadership in 3+ categories (points, rebounds, assists), or scored 50+ points">Big Game</Badge></Show>
            <Show when={q().isBigScoring}><Badge class="bg-purple-50 text-purple-600" title="A player scored 35+ points">Scoring</Badge></Show>
            <Show when={q().isStarDuel}><Badge class="bg-blue-50 text-blue-600" title="Both teams had a player scoring 35+ points">Star Duel</Badge></Show>
            <Show when={q().isHugeSwing}><Badge class="bg-green-50 text-green-600" title="Big lead after 3rd quarter that was closed in the 4th">Huge Swing</Badge></Show>
            <Show when={q().isShootout}><Badge class="bg-pink-50 text-pink-600" title="High scoring game (235+ total points)">Shootout</Badge></Show>
            <Show when={q().isGritty}><Badge class="bg-gray-50 text-gray-600" title="Low scoring game (200 or less total points)">Gritty</Badge></Show>
          </div>

          <div class={`px-2 py-1 rounded text-[11px] font-black ${q().qualityScore >= 70 ? 'bg-orange-500 text-white' : 'bg-gray-100 text-gray-600'}`}>
            {q().qualityScore} PTS
          </div>
        </Show>

        <Show when={props.game.status !== "STATUS_FINAL"}>
          <div class="space-y-3 mb-4">
            <TeamRow team={props.game.awayTeam} score={0} hide={true} />
            <TeamRow team={props.game.homeTeam} score={0} hide={true} />
          </div>
        </Show>
      </div>

      <Show when={auth.isLoggedIn && props.game.status === "STATUS_FINAL"}>
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

const TeamRow = (p: { team: Team, score: number,  hide?: boolean }) => (
  <div class="flex justify-between items-center">
    <div class="flex items-center gap-2">
      <span class="text-xs font-black text-gray-400 w-8">{p.team.abbreviation}</span>
      <span class='font-black text-gray-800' >{p.team.teamName}</span>
    </div>
    <span class='font-mono text-xl font-black text-gray-800'>
      {p.hide ? '—' : p.score}
    </span>
  </div>
);

const Badge = (p: { children: any, class: string, title?: string }) => (
  <span class={`text-[9px] font-bold px-2 py-0.5 rounded border border-current ${p.class}`} title={p.title}>{p.children}</span>
);

export default GameItem;