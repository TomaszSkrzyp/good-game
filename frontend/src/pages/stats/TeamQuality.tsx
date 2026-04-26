import { Component, createResource, For, Show } from "solid-js";
import { api } from "../../utils/api";

export type TeamQualityStat = {
  teamId: number;
  teamName: string;
  avgGameQuality: number;
  gamesRated: number;
};

const TeamQualityTable: Component = () => {
  const [teamQualityResource] = createResource(async () => {
    return await api.get('teams/quality').json<TeamQualityStat[]>();
  });

  return (
    <>
      <Show when={teamQualityResource.loading}>
        <div class="flex justify-center p-10 text-gray-500 italic">loading stats...</div>
      </Show>

      <Show when={!teamQualityResource.loading && teamQualityResource()?.length === 0}>
        <div class="text-center py-10 text-gray-400">No stats available yet. Rate some games!</div>
      </Show>

      <Show when={teamQualityResource() && teamQualityResource()!.length > 0}>
        <div class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-200">
                <th class="p-4 font-semibold text-gray-600 text-center w-16">Rank</th>
                <th class="p-4 font-semibold text-gray-600">Team</th>
                <th class="p-4 font-semibold text-gray-600 text-right">Avg Quality</th>
                <th class="p-4 font-semibold text-gray-600 text-right">Games Played</th>
              </tr>
            </thead>
            <tbody>
              <For each={teamQualityResource()}>
                {(stat, index) => (
                  <tr class="border-b border-gray-100 hover:bg-gray-50 transition-colors last:border-0">
                    <td class="p-4 font-bold text-gray-400 text-center">{index() + 1}</td>
                    <td class="p-4 font-medium text-gray-800">{stat.teamName}</td>
                    <td class="p-4 text-right font-semibold text-blue-600">
                      {stat.avgGameQuality.toFixed(2)}
                    </td>
                    <td class="p-4 text-right text-gray-500">
                      {stat.gamesRated}
                    </td>
                  </tr>
                )}
              </For>
            </tbody>
          </table>
        </div>
      </Show>
    </>
  );
};

export default TeamQualityTable;