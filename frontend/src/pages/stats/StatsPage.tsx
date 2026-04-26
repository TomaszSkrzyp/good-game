import { ParentComponent, For } from "solid-js";
import { A, useLocation } from "@solidjs/router";

const StatsPage: ParentComponent = (props) => {
  const location = useLocation();

  const tabs = [
    { path: "/stats/team-quality", label: "Team Quality Rankings" },
  ];

  const isActive = (path: string) => location.pathname === path;

  return (
    <div class="max-w-4xl mx-auto p-4">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-800 mb-6">Statistics Dashboard</h1>
        
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-8" aria-label="Tabs">
            <For each={tabs}>
              {(tab) => (
                <A
                  href={tab.path}
                  class={`
                    whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors
                    ${isActive(tab.path)
                      ? "border-blue-500 text-blue-600" 
                      : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                    }
                  `}
                >
                  {tab.label}
                </A>
              )}
            </For>
          </nav>
        </div>
      </div>

      <div class="mt-6">
        {props.children}
      </div>
    </div>
  );
};

export default StatsPage;