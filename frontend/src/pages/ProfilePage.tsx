import { Component, createResource } from "solid-js";
import { auth } from "../data/store";
const fetchUserData = async (token: string) => {
  const response = await fetch("/api/profile", {
    headers: {
      "Authorization": `Bearer ${token}`
    }
  });
  if (!response.ok) throw new Error("Failed to fetch profile");
  return response.json();
};

const ProfilePage: Component = () => {
  // createResource automatically handles the async loading state
  const [data] = createResource(() => auth.token, fetchUserData);

  return (
    <div>
      <h1>User Profile</h1>
      {data.loading && <p>Loading your data...</p>}
      {data.error && <p style="color: red">Error: {data.error.message}</p>}
      
      {data() && (
        <div class="profile-card">
          <p><strong>Username:</strong> {data().userName}</p>
          <p><strong>Email:</strong> {data().email}</p>
          <p><strong>User ID:</strong> {data().id}</p>
        </div>
      )}
    </div>
  );
};

export default ProfilePage;