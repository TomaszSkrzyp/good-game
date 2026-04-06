import { Component, createSignal, Show, For } from "solid-js";
import { auth } from "../../data/store";

interface Team {
  abbreviation: string;
  teamName: string;
}

interface GamePlayerStats {
  homeTopScorer: string; homeTopScorerPts: number;
  homeTopAssister: string; homeTopAssists: number;
  homeTopRebounder: string; homeTopRebounds: number;
  awayTopScorer: string; awayTopScorerPts: number;
  awayTopAssister: string; awayTopAssists: number;
  awayTopRebounder: string; awayTopRebounds: number;
}

export interface GameQuality {
  qualityScore: number;
  isBigScoring: boolean; isBigGame: boolean;
  isClutch: boolean; isStarDuel: boolean;
  isHugeSwing: boolean; isShootout: boolean;
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
  const [localToggle, setLocalToggle] = createSignal(false); 
  
  // Rating states
  const [showConfirmModal, setShowConfirmModal] = createSignal(false);
  const [pendingValue, setPendingValue] = createSignal<number>(0);

  const q = () => props.game.gameQuality;

  const isCurrentlyHidden = () => {
    if (props.game.status !== "STATUS_FINAL") return true;
    return (props.hideScores ?? false) !== localToggle();
  };

  const handleStarClick = async (val: number) => {
    const currentRating = props.game.rating ?? 0;
    if (currentRating === 0) {
      await executeRate(val);
    } else if (val !== currentRating) {
      setPendingValue(val);
      setShowConfirmModal(true);
    }
  };

  const executeRate = async (val: number) => {
    setIsSubmitting(true);
    await props.onRate(props.game.id, val.toString());
    setShowConfirmModal(false);
    setIsSubmitting(false);
  };

  const fetchGameStats = async () => {
    try {
      const response = await fetch(`/api/game/stats?gameId=${props.game.id}`);
      if (!response.ok) throw new Error("failed to fetch stats");
      const data = await response.json();
      setStats(data);
      setShowStats(true);
      document.body.style.overflow = 'hidden';
    } catch (error) {
      console.error("stats error:", error);
    }
  };

  const closeStats = () => {
    setShowStats(false);
    document.body.style.overflow = 'auto';
  };

  return (
    <div class="bg-white border border-gray-200 rounded-xl p-5 relative flex flex-col justify-between h-full group hover:border-blue-200 transition-all">
      
      <Show when={isSubmitting()}>
        <div class="absolute inset-0 bg-white/60 z-20 flex items-center justify-center backdrop-blur-sm">
          <div class="animate-spin h-6 w-6 border-2 border-blue-500 border-t-transparent rounded-full" />
        </div>
      </Show>

      {/* Stats Modal */}
      <Show when={showStats()}>
        <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" onClick={closeStats}>
          <div class="bg-white rounded-lg shadow-xl p-6 w-full max-w-md max-h-[80vh] overflow-y-auto relative" onClick={(e) => e.stopPropagation()}>
            <button onClick={closeStats} class="absolute top-4 right-4 text-gray-400 hover:text-gray-600 p-2 cursor-pointer">✕</button>
            <h3 class="font-bold text-lg mb-6 text-center text-gray-800">Game Stats</h3>
            <Show when={stats()}>
               <div class="space-y-6">
                 <TeamStatsBlock teamName={props.game.awayTeam.teamName} stats={{
                   scorer: stats()!.awayTopScorer, scorerPts: stats()!.awayTopScorerPts,
                   assister: stats()!.awayTopAssister, assists: stats()!.awayTopAssists,
                   rebounder: stats()!.awayTopRebounder, rebounds: stats()!.awayTopRebounds
                 }} />
                 <TeamStatsBlock teamName={props.game.homeTeam.teamName} stats={{
                   scorer: stats()!.homeTopScorer, scorerPts: stats()!.homeTopScorerPts,
                   assister: stats()!.homeTopAssister, assists: stats()!.homeTopAssists,
                   rebounder: stats()!.homeTopRebounder, rebounds: stats()!.homeTopRebounds
                 }} />
               </div>
            </Show>
          </div>
        </div>
      </Show>

      {/* Rating Confirm Modal */}
      <Show when={showConfirmModal()}>
        <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-lg p-6 max-w-xs w-full shadow-xl text-center">
            <h3 class="font-bold text-gray-800 mb-2 text-sm">Update Rating?</h3>
            <p class="text-xs text-gray-500 mb-6">Change from {props.game.rating} to {pendingValue()} stars?</p>
            <div class="flex gap-2">
              <button onClick={() => executeRate(pendingValue())} class="flex-1 py-2 bg-blue-600 text-white text-[10px] font-bold rounded uppercase cursor-pointer">Confirm</button>
              <button onClick={() => setShowConfirmModal(false)} class="flex-1 py-2 bg-gray-100 text-gray-500 text-[10px] font-bold rounded uppercase cursor-pointer">Cancel</button>
            </div>
          </div>
        </div>
      </Show>

      <div>
        <div class="flex justify-between items-start mb-4">
          <div class="text-[10px] font-bold text-gray-400 uppercase">
            {new Date(props.game.gameTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
          </div>
          <Show when={props.game.status === "STATUS_FINAL"}>
            <div class="flex gap-2">
              <button onClick={() => setLocalToggle(!localToggle())} class={`px-3 py-1 text-[10px] font-bold rounded transition-colors cursor-pointer ${isCurrentlyHidden() ? 'text-green-600 bg-green-50' : 'text-gray-500 bg-gray-50'}`}>
                {isCurrentlyHidden() ? "Show Score" : "Hide Score"}
              </button>
              <button onClick={fetchGameStats} class="px-3 py-1 text-[10px] font-bold text-blue-600 bg-blue-50 rounded hover:bg-blue-100 cursor-pointer">Stats</button>
            </div>
          </Show>
        </div>

        <div class="space-y-3 mb-4">
          <TeamRow team={props.game.awayTeam} score={props.game.awayTeamPoints} hide={isCurrentlyHidden()} />
          <TeamRow team={props.game.homeTeam} score={props.game.homeTeamPoints} hide={isCurrentlyHidden()} />
        </div>

        <Show when={props.game.status === "STATUS_FINAL"}>
          <div class="flex flex-wrap gap-1 mb-4">
            <Badge class="bg-emerald-50 text-emerald-600" title="Game finished">Finished</Badge>
            <Show when={q().isClutch}><Badge class="bg-red-50 text-red-600" title="4th quarter was close or game decided by 3 points or less">Clutch</Badge></Show>
            <Show when={q().isBigGame}><Badge class="bg-yellow-50 text-yellow-700" title="High impact player performance">Big Game</Badge></Show>
            <Show when={q().isBigScoring}><Badge class="bg-purple-50 text-purple-600" title="A player scored 35+ points">Scoring</Badge></Show>
            <Show when={q().isStarDuel}><Badge class="bg-blue-50 text-blue-600" title="Star players on both teams performed">Star Duel</Badge></Show>
            <Show when={q().isHugeSwing}><Badge class="bg-green-50 text-green-600" title="Big lead closed in 4th quarter">Huge Swing</Badge></Show>
            <Show when={q().isShootout}><Badge class="bg-pink-50 text-pink-600" title="High total points shootout">Shootout</Badge></Show>
            <Show when={q().isGritty}><Badge class="bg-gray-50 text-gray-600" title="Low scoring defensive game">Gritty</Badge></Show>
          </div>
          <div class={`inline-block px-2 py-1 rounded text-[11px] font-black ${q().qualityScore >= 70 ? 'bg-orange-500 text-white' : 'bg-gray-100 text-gray-600'}`}>
            {q().qualityScore} PTS
          </div>
        </Show>
      </div>

      <div class="pt-4 border-t border-gray-100 flex flex-col gap-3 mt-4">
        <div class="flex justify-between items-center">
          <div class="flex flex-col">
            <span class="text-[9px] font-black text-gray-400 uppercase tracking-tighter">Avg Rating</span>
            <span class="text-sm font-bold text-blue-600">
              {props.game.avgRating?.toFixed(1) || "0.0"} 
              <span class="text-gray-400 font-normal text-[10px] ml-1">({props.game.ratingCount || 0})</span>
            </span>
          </div>

          <Show when={auth.isLoggedIn} fallback={<span class="text-[9px] text-gray-400 italic">Sign in to rate</span>}>
            <div class="flex gap-1">
              <For each={[1, 2, 3, 4, 5]}>{(v) => {
                const isSelected = () => (props.game.rating ?? 0) === v;
                return (
                  <button onClick={() => handleStarClick(v)} class={`w-8 h-8 rounded text-xs font-bold border cursor-pointer transition-colors ${isSelected() ? 'bg-blue-600 text-white border-blue-600' : 'bg-white text-gray-400 border-gray-200 hover:border-blue-300'}`}>{v}</button>
                )
              }}</For>
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
};

// Fixed Badge Component
const Badge = (p: { children: any, class: string, title?: string }) => (
  <span title={p.title} class={`text-[9px] font-bold px-2 py-0.5 rounded border border-current cursor-help ${p.class}`}>
    {p.children}
  </span>
);

const TeamRow = (p: { team: Team, score: number, hide: boolean }) => (
  <div class="flex justify-between items-center">
    <div class="flex items-center gap-2">
      <span class="text-xs font-black text-gray-400 w-8">{p.team.abbreviation}</span>
      <span class="font-bold text-gray-800">{p.team.teamName}</span>
    </div>
    <span class="font-mono text-xl font-black text-gray-800">{p.hide ? '—' : p.score}</span>
  </div>
);

const TeamStatsBlock = (p: { teamName: string, stats: any }) => (
  <div class="bg-gray-50 rounded-lg p-4 text-sm w-full">
    <h4 class="font-bold text-xs text-gray-500 uppercase mb-3 text-center border-b pb-1">{p.teamName}</h4>
    <div class="flex justify-between mb-1"><span>Points</span><span class="font-bold">{p.stats.scorer} ({p.stats.scorerPts})</span></div>
    <div class="flex justify-between mb-1"><span>Assists</span><span class="font-bold">{p.stats.assister} ({p.stats.assists})</span></div>
    <div class="flex justify-between"><span>Rebounds</span><span class="font-bold">{p.stats.rebounder} ({p.stats.rebounds})</span></div>
  </div>
);

export default GameItem;