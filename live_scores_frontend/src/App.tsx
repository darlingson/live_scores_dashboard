import React, { useEffect, useState } from 'react';
import './App.css';

// Types
interface Scorer {
  playerName: string;
  team: string;
  minute: number;
}

type GameStatus = 'pending' | 'active' | 'finished';

interface Game {
  id: string;
  homeTeam: string;
  awayTeam: string;
  homeScore: number;
  awayScore: number;
  scorers: Scorer[];
  status: GameStatus;
  lastUpdate: string;
}

interface WebSocketMessage {
  type: 'initialGames' | 'gameUpdate';
  data: Game | Game[];
}

// Components
const GameCard: React.FC<{ game: Game }> = ({ game }) => {
  return (
    <div className={`game-card ${game.status}`}>
      <div className="team-scores">
        <span className="team-name">{game.homeTeam}</span>
        <span className="score">{game.homeScore}</span>
        <span className="separator">-</span>
        <span className="score">{game.awayScore}</span>
        <span className="team-name">{game.awayTeam}</span>
      </div>
      <div className="game-info">
        <span className="status">{game.status.toUpperCase()}</span>
        {game.status === 'active' && <span className="live-indicator">● LIVE</span>}
      </div>
      {game.scorers && game.scorers.length > 0 && (
        <div className="scorers-section">
          <h4>Goal Scorers:</h4>
          <ul>
            {game.scorers.map((scorer, index) => (
              <li key={index}>
                {scorer.playerName} ({scorer.team}) - {scorer.minute}'
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

// Main App
const App: React.FC = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [connectionStatus, setConnectionStatus] = useState<string>('Connecting...');

  useEffect(() => {
    let ws: WebSocket;

    const connect = () => {
      ws = new WebSocket('ws://localhost:8080/ws');

      ws.onopen = () => {
        console.log('Connected to WebSocket server');
        setConnectionStatus('Connected');
      };

      ws.onmessage = (event) => {
        const message: WebSocketMessage = JSON.parse(event.data);
        console.log('Received:', message);

        if (message.type === 'initialGames') {
          setGames(message.data as Game[]);
        } else if (message.type === 'gameUpdate') {
          const updatedGame = message.data as Game;
          setGames((prevGames) => {
            const index = prevGames.findIndex((g) => g.id === updatedGame.id);
            if (index !== -1) {
              const updated = [...prevGames];
              updated[index] = updatedGame;
              return updated;
            } else {
              return [...prevGames, updatedGame];
            }
          });
        }
      };

      ws.onclose = () => {
        console.warn('Disconnected from server. Retrying...');
        setConnectionStatus('Disconnected. Retrying in 5 seconds...');
        setTimeout(connect, 5000);
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('Error. Check server connection.');
      };
    };

    connect();

    return () => {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, []);

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Malawian Football Live Scores</h1>
        <p className="connection-status">Status: {connectionStatus}</p>
      </header>
      <main className="games-grid">
        {games.length === 0 && connectionStatus === 'Connected' ? (
          <p>No active games currently. Waiting for updates...</p>
        ) : (
          games
            .sort((a, b) => {
              const statusOrder: Record<GameStatus, number> = {
                active: 1,
                pending: 2,
                finished: 3,
              };
              return statusOrder[a.status] - statusOrder[b.status];
            })
            .map((game) => <GameCard key={game.id} game={game} />)
        )}
      </main>
      <footer className="app-footer">
        <p>&copy; {new Date().getFullYear()} Malawian Football Updates</p>
      </footer>
    </div>
  );
};

export default App;
