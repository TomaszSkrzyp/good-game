import { Component, createSignal } from "solid-js";

interface TunerProps {
  initialConfig: any;
  onSave: (newConfig: any) => Promise<void>;
  title?: string;
}

const AlgorithmTuner: Component<TunerProps> = (props) => {
  // local state for the form so we don't save on every keystroke
  const [form, setForm] = createSignal({ ...props.initialConfig });
  const [loading, setLoading] = createSignal(false);
  const [saved, setSaved] = createSignal(false);
  const updateField = (path: string, value: any) => {
    setForm({ ...form(), [path]: value });
  };

  const handleSave = async (e: Event) => {
  e.preventDefault();
  setLoading(true);
  try {
    await props.onSave(form());
    setSaved(true);
    setTimeout(() => setSaved(false), 2000); // Reset after 2 seconds
  } finally {
    setLoading(false);
  }
};

  return (
    <div class="bg-gray-50 border border-gray-200 rounded-xl p-6 shadow-inner">
      <h3 class="text-lg font-bold mb-4 flex items-center gap-2 uppercase tracking-tight text-gray-800">
        🛠️ {props.title || "Algorithm Tuner"}
      </h3>
      
      <form onSubmit={handleSave} class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* bonus values */}
          <div class="space-y-4">
            <h4 class="text-xs font-bold text-gray-400 uppercase tracking-widest">Bonuses</h4>
            <div class="space-y-3">
              <TunerInput label="Clutch Bonus" value={form().clutchBonus} onInput={(v) => updateField('clutchBonus', v)} />
              <TunerInput label="Star Duel Bonus" value={form().starDuelBonus} onInput={(v) => updateField('starDuelBonus', v)} />
              <TunerInput label="Huge Swing" value={form().hugeSwingBonus} onInput={(v) => updateField('hugeSwingBonus', v)} />
            </div>
          </div>

          {/* logic thresholds */}
          <div class="space-y-4">
            <h4 class="text-xs font-bold text-gray-400 uppercase tracking-widest">Thresholds</h4>
            <div class="space-y-3">
              <TunerInput label="Shootout Total" value={form().shootoutThreshold} onInput={(v) => updateField('shootoutThreshold', v)} />
              <TunerInput label="Gritty Total" value={form().grittyThreshold} onInput={(v) => updateField('grittyThreshold', v)} />
              <TunerInput label="Star Performance" value={form().starPointsThreshold} onInput={(v) => updateField('starPointsThreshold', v)} />
            </div>
          </div>
        </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* playoff stakes */}
            <div class="space-y-4">
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-widest">Playoff Stakes</h4>
              <div class="space-y-3">
                <TunerInput label="Game 7 Bonus" value={form().game7Bonus} onInput={(v) => updateField('game7Bonus', v)} />
                <TunerInput label="Elimination Bonus" value={form().eliminationBonus} onInput={(v) => updateField('eliminationBonus', v)} />
                <TunerInput label="Base Playoff" value={form().playoffBonus} onInput={(v) => updateField('playoffBonus', v)} />
              </div>
            </div>

            {/* drama weights*/}
            <div class="space-y-4">
              <h4 class="text-xs font-bold text-gray-400 uppercase tracking-widest">Win Prob Drama</h4>
              <div class="space-y-3">
                <TunerInput label="Volatility (Jitter)" value={form().volatilityWeight} onInput={(v) => updateField('volatilityWeight', v)} />
                <TunerInput label="Max Swing (Comeback)" value={form().swingWeight} onInput={(v) => updateField('swingWeight', v)} />
                <TunerInput label="Comeback % (0.80)" value={form().comebackThreshold} onInput={(v) => updateField('comebackThreshold', v)} />
                <TunerInput label="Play-In Bonus" value={form().playInBonus} onInput={(v) => updateField('playInBonus', v)} />
                <TunerInput label="Season Series Grudge Bonus" value={form().seasonSeriesTiedBonus} onInput={(v) => updateField('seasonSeriesTiedBonus', v)} />
              </div>
            </div>
          </div>

        <button 
          type="submit" 
          disabled={loading()}
          class="w-full bg-black text-white font-bold py-2 rounded-lg hover:bg-gray-800 transition disabled:opacity-50 uppercase text-sm"
        >
          {loading() ? "Saving..." : saved() ? "Configuration Applied" : "Apply Global Configuration"}
        </button>
      </form>
      
      <p class="mt-4 text-[10px] text-gray-400 italic lowercase">
        note: current implementation updates the global algorithm. future presets will be saved to user profiles.
      </p>
    </div>
  );
};

// internal helper for clean inputs
const TunerInput: Component<{ label: string, value: number, onInput: (v: number) => void }> = (props) => (
  <div class="flex flex-col gap-1">
    <label class="text-xs font-bold text-gray-600 uppercase">{props.label}</label>
    <input 
      type="number" 
      value={props.value} 
      onInput={(e) => props.onInput(parseFloat(e.currentTarget.value) || 0)}
      class="p-2 border rounded bg-white font-mono text-sm focus:ring-2 focus:ring-blue-500 outline-none" 
    />
  </div>
);

export default AlgorithmTuner;