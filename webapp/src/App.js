import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { GameProvider } from './context/GameContext';
import Home from './pages/Home';
import Game from './pages/Game';
import Challenge from './pages/Challenge';
import theme from './theme';
import './App.css';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <GameProvider>
        <Router>
          <div className="App">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/game" element={<Game />} />
              <Route path="/challenge/:username" element={<Challenge />} />
              <Route path="/challenge/:username/:gameId" element={<Challenge />} />
            </Routes>
          </div>
        </Router>
      </GameProvider>
    </ThemeProvider>
  );
}

export default App; 