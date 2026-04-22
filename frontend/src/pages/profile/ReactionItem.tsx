import { Component } from "solid-js";

// helper interfaces for clarity
interface Team {
  teamName: string;
}

interface Game {
  homeTeam: Team;
  awayTeam: Team;
}

interface ReactionProps {
  reaction: {
    id: number;
    rating: number;
    createdAt: string;
    game: Game;    
  };
}

const ReactionItem: Component<ReactionProps> = (props) => {
  return (
    <div class="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200 flex justify-between items-center shadow-sm">
      <div>
        <h4 class="text-xs font-bold text-blue-600 uppercase tracking-widest mb-1">Matchup</h4>
        <p class="text-md font-semibold text-gray-800">
          {props.reaction.game.homeTeam.teamName} 
          <span class="mx-2 text-gray-400 font-normal text-sm">vs</span> 
          {props.reaction.game.awayTeam.teamName}
        </p>
        <p class="text-xs text-gray-400 mt-2">
          {new Date(props.reaction.createdAt).toLocaleDateString()}
        </p>
      </div>
      
      <div class="text-right">
        <span class="text-xs text-gray-400 block mb-1">Your Rating</span>
        <div class="bg-white border px-3 py-1 rounded-full font-bold text-blue-700 shadow-sm">
          {props.reaction.rating} / 5
        </div>
      </div>
    </div>
  );
};

export default ReactionItem;