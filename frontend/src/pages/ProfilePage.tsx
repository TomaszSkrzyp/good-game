import { Component, createResource, Show, createEffect, For } from "solid-js";
import { auth} from "../data/store";
import { api } from "../utils/api";
import { useNavigate } from "@solidjs/router";
interface UserData {
  userName: string;
  email: string;
  userReactions: UserReactions[];
}
interface UserReactions{
  id: number;
  gameId: number;
  rating: number;
  createdAt: string;
}
const fetchUserData = async (): Promise<UserData> => {
  const profile = await api.get('profile').json<UserData>();
  //id gets taken from ctx
  const userReactions = await api.get('userReactions').json<UserReactions[]>();

  return { ...profile, userReactions };
};

const ProfilePage: Component = () => {
  const navigate = useNavigate();
  const [data] = createResource<UserData>(fetchUserData);

  createEffect(() => {
    console.log(auth);
    if (!auth.token) {
      navigate("/login");
      return;
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
          <Show when={data.loading}>
            <div class="animate-pulse space-y-4">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </Show>
          <Show when={data.error}>
            <div class="bg-red-50 text-red-600 p-4 rounded-lg border border-red-100">
              {data.error.message}
            </div>
          </Show>
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
            </div>
          </Show>
          <For each={data()?.userReactions}>{(reaction) => (
            <div class="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
              <p class="text-sm text-gray-500">Reaction ID: {reaction.id}</p>
              <p class="text-sm text-gray-500">Game ID: {reaction.gameId}</p>
              <p class="text-sm text-gray-500">Rating: {reaction.rating}</p>
              <p class="text-sm text-gray-500">Created At: {new Date(reaction.createdAt).toLocaleString()}</p>

            </div>
          )}</For>
      </div>
    </div>
    </div>
    );
};
  

export default ProfilePage;