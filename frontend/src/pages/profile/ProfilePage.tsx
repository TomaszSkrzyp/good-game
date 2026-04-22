import { Component, createResource, Show, createEffect, For } from "solid-js";
import { auth, logout } from "../../data/store";
import { api } from "../../utils/api";
import { useNavigate } from "@solidjs/router";
import ReactionItem from "./ReactionItem";

interface UserReaction {
  id: number;
  gameId: number;
  rating: number;
  createdAt: string;
  game: {
    homeTeam: { teamName: string };
    awayTeam: { teamName: string };
  };
}

interface UserData {
  userName: string;
  email: string;
  userReactions: UserReaction[];
}

const fetchUserData = async (): Promise<UserData> => {
  const profile = await api.get('profile').json<any>();
  const userReactions = await api.get('userReactions').json<UserReaction[]>();
  return { 
    userName: profile.userName, 
    email: profile.email, 
    userReactions 
  };
};

const ProfilePage: Component = () => {
  const navigate = useNavigate();
  const [data] = createResource<UserData>(fetchUserData);

  createEffect(() => {
    if (!auth.token) navigate("/");
    if (data.error) {
      logout();
      navigate("/");
    }
  });

  return (
    <div class="max-w-2xl mx-auto mt-12 p-6">
      <div class="bg-white shadow-lg rounded-xl overflow-hidden border border-gray-100">
        <div class="bg-blue-600 h-24 w-full"></div>
        <div class="px-8 pb-8">
          <div class="relative -mt-12 mb-6">
            <div class="w-24 h-24 bg-gray-200 border-4 border-white rounded-full flex items-center justify-center text-3xl shadow-sm text-gray-500 font-bold uppercase">
              {data()?.userName?.slice(0, 2) || "?"}
            </div>
          </div>

          <h1 class="text-3xl font-bold text-gray-800 mb-6">User Profile</h1>

          <Show when={data()}>
            <div class="space-y-6">
              <div class="border-b border-gray-100 pb-4">
                <label class="text-xs uppercase tracking-wider text-gray-400 font-bold">Username</label>
                <p class="text-lg text-gray-700 font-medium">{data()?.userName}</p>
              </div>
              <div class="border-b border-gray-100 pb-4">
                <label class="text-xs uppercase tracking-wider text-gray-400 font-bold">Email Address</label>
                <p class="text-lg text-gray-700 font-medium">{data()?.email}</p>
              </div>

              <div class="mt-8">
                <h2 class="text-xl font-bold text-gray-800 mb-4">Activity History</h2>
                <For each={data()?.userReactions} fallback={<p class="text-gray-500 italic">No ratings yet.</p>}>
                  {(reaction) => (
                    <ReactionItem reaction={reaction} />
                  )}
                </For>
              </div>
            </div>
          </Show>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;