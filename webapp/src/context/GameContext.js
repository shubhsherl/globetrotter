import React, { createContext, useState, useContext, useEffect } from 'react';

const GameContext = createContext();

// Helper function to check if stored data is expired
const isExpired = (timestamp) => {
  const TWO_HOURS = 2 * 60 * 60 * 1000; // 2 hours in milliseconds
  return Date.now() - timestamp > TWO_HOURS;
};

export const GameProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [score, setScore] = useState({ correct: 0, total: 0 });
  const [gameId, setGameId] = useState(null);
  const [gameState, setGameState] = useState('initial'); // initial, playing, finished
  const [currentQuestion, setCurrentQuestion] = useState(null);
  const [results, setResults] = useState(null);
  
  // Check localStorage for existing user on mount
  useEffect(() => {
    const storedUser = localStorage.getItem('globetrotter_user');
    const storedTimestamp = localStorage.getItem('globetrotter_timestamp');
    
    if (storedUser && storedTimestamp && !isExpired(parseInt(storedTimestamp))) {
      setUser(JSON.parse(storedUser));
    } else {
      // Clear expired data
      localStorage.removeItem('globetrotter_user');
      localStorage.removeItem('globetrotter_timestamp');
    }
  }, []);
  
  // Save user to localStorage with timestamp
  const saveUser = (userData) => {
    setUser(userData);
    localStorage.setItem('globetrotter_user', JSON.stringify(userData));
    localStorage.setItem('globetrotter_timestamp', Date.now().toString());
  };
  
  const updateScore = (correct) => {
    setScore(prev => ({
      correct: correct ? prev.correct + 1 : prev.correct,
      total: prev.total + 1
    }));
  };
  
  const resetScore = () => {
    setScore({ correct: 0, total: 0 });
  };
  
  const resetGame = () => {
    setGameId(null);
    setGameState('initial');
    setCurrentQuestion(null);
    setResults(null);
    resetScore(); // Reset score when resetting the game
  };
  
  return (
    <GameContext.Provider value={{ 
      user, 
      setUser: saveUser, 
      score, 
      updateScore,
      resetScore,
      gameId,
      setGameId,
      gameState,
      setGameState,
      currentQuestion,
      setCurrentQuestion,
      results,
      setResults,
      resetGame
    }}>
      {children}
    </GameContext.Provider>
  );
};

export const useGame = () => useContext(GameContext); 