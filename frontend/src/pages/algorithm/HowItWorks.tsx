import { Component, createResource, createSignal, For, Show } from "solid-js";
import { api } from "../../utils/api";
import AlgorithmTuner from "./AlgorithmTuner";
import { isAdmin } from "../../data/jwt";

interface MarginWeight { maxMargin: number; points: number; }
interface GameQualityConfig {
  margins: MarginWeight[];
  hugeSwingBonus: number;
  clutchBonus: number;
  overtimeBonus: number;
  shootoutBonus: number;
  shootoutThreshold: number;
  grittyThreshold: number;
  starDuelBonus: number;
  starPointsThreshold: number;
  bigGameBonus: number;
}

const HowItWorks: Component = () => {
  const [showTuner, setShowTuner] = createSignal(false);
  // fetch live config
  const [config, { mutate }] = createResource(async () => {
    return await api.get('config').json<GameQualityConfig>();
  });

  // the callback passed to the tuner
  const handleAdminUpdate = async (newConfig: GameQualityConfig) => {
    try {
      await api.post('config', { json: newConfig });
      mutate(newConfig); // reactive sync for the rest of the page
    } catch (e) {
      console.error("update failed", e);
    }
  };

  return (
    <div class="max-w-3xl mx-auto p-6 bg-white shadow-sm rounded-xl my-10 border border-gray-100">
      <header class="mb-10 border-b pb-6">
        <h1 class="text-4xl font-extrabold text-gray-900 mb-2">The Quality Algorithm</h1>
        <p class="text-gray-500 uppercase text-xs tracking-widest font-bold">Automatic "Good Game" calculation breakdown.</p>
      </header>
      <Show when={isAdmin()}>
      <div class="mb-8 flex justify-end">
        <button 
          onClick={() => setShowTuner(!showTuner())}
          class="text-[10px] text-gray-300 hover:text-blue-500 underline italic lowercase transition"
        >
          {showTuner() ? "close experimental settings" : "open algorithm tuner"}
        </button>
      </div>

      <Show when={showTuner() && config() }>
        <div class="mb-12 animate-in fade-in slide-in-from-top-4 duration-300">
          <AlgorithmTuner 
            initialConfig={config()!} 
            onSave={handleAdminUpdate} 
            title="Global Admin Controls"
          />
        </div>
      </Show>
      </Show>

      <Show when={!config.loading} fallback={<p class="text-center py-10 italic lowercase text-gray-400">loading logic...</p>}>
        <section class="space-y-12">
          <div>
            <h2 class="text-2xl font-bold text-gray-800 mb-4 flex items-center gap-2 text-transform-none">
              <span class="bg-blue-100 text-blue-600 px-3 py-1 rounded-full text-sm font-mono font-bold">01</span>
              Game Tension
            </h2>
            <p class="text-gray-600 mb-4 italic lowercase text-sm">points added based on the final score difference.</p>
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
              <For each={config()?.margins}>
                {(m) => (
                  <div class="p-4 bg-gray-50 rounded-lg border border-gray-100 text-center">
                    <div class="text-[10px] uppercase tracking-tighter text-gray-400 font-bold">Under {m.maxMargin} Pts</div>
                    <div class="text-2xl font-bold text-blue-600">+{m.points}</div>
                  </div>
                )}
              </For>
            </div>
          </div>

          <div>
          <h2 class="text-2xl font-bold text-gray-800 mb-6 flex items-center gap-2 text-transform-none">
          <span class="bg-purple-100 text-purple-600 px-3 py-1 rounded-full text-sm font-mono font-bold">02</span>
            Badge Legend
          </h2>
          <div class="grid grid-cols-1 gap-3">
              <BadgeDesc label="Clutch" color="bg-red-50 text-red-700 border-red-100"
                desc={`Final margin ≤ 3 pts, or game was within 8 pts entering the 4th quarter. (+${config()?.clutchBonus} Pts)`} />
    
              <BadgeDesc label="Star Duel" color="bg-purple-50 text-purple-700 border-purple-100"
                desc={`Opposing stars both scored over ${config()?.starPointsThreshold} points. (+${config()?.starDuelBonus} Pts)`} />
    
              <BadgeDesc label="Comeback" color="bg-blue-50 text-blue-700 border-blue-100"
                desc={`A 15+ point deficit at the end of the 3rd was cut to 7 pts or less by the finish. (+${config()?.hugeSwingBonus} Pts)`} />
              
              <BadgeDesc label="Shootout" color="bg-orange-50 text-orange-700 border-orange-100"
                desc={`High-octane game with a total combined score over ${config()?.shootoutThreshold}. (+${config()?.shootoutBonus} Pts)`} />

              <BadgeDesc label="Overtime" color="bg-emerald-50 text-emerald-700 border-emerald-100"
                desc={`Extra periods were required to decide the winner. (+${config()?.overtimeBonus} Pts)`} />

              <BadgeDesc label="Gritty" color="bg-slate-50 text-slate-700 border-slate-100"
                desc={`Defensive battle with a total score under ${config()?.grittyThreshold} points in regulation. (+10 Pts)`} />

              <BadgeDesc label="Big Game" color="bg-amber-50 text-amber-700 border-amber-100"
                desc={`Triggered by historic individual play: a 50pt performance or a 30pt Triple-Double. (+${config()?.bigGameBonus} Pts)`} />
            </div>
          </div>
        </section>
      </Show>

      <footer class="mt-16 text-center text-gray-300 text-[10px] uppercase tracking-widest font-bold">
        &copy; {new Date().getFullYear()} Good-Game NBA Algorithm
      </footer>
    </div>
  );
};

// reusable badge description row
const BadgeDesc: Component<{ label: string, color: string, desc: string }> = (props) => (
  <div class="flex items-center gap-4 p-4 rounded-lg border border-gray-50 hover:border-gray-200 transition">
    <span class={`px-2 py-1 rounded text-[10px] font-bold uppercase border w-24 text-center shrink-0 ${props.color}`}>
      {props.label}
    </span>
    <p class="text-sm text-gray-600 italic lowercase">{props.desc}</p>
  </div>
);

export default HowItWorks;