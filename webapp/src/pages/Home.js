import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { createUser } from '../services/api';
import { useGame } from '../context/GameContext';
import { 
  Container, 
  Typography, 
  TextField, 
  Button, 
  Box, 
  Paper, 
  Avatar,
  Fade,
  CircularProgress
} from '@mui/material';
import { styled } from '@mui/material/styles';
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents';
import PublicIcon from '@mui/icons-material/Public';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';

// Styled components
const GlobeAvatar = styled(Avatar)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  width: 80,
  height: 80,
  margin: '0 auto 16px',
}));

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  borderRadius: 16,
  maxWidth: 500,
  margin: '0 auto',
  boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
}));

const ColorButton = styled(Button)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  color: 'white',
  fontWeight: 'bold',
  padding: '12px 32px',
  fontSize: '1.1rem',
  borderRadius: 30,
  '&:hover': {
    backgroundColor: theme.palette.primary.dark,
  },
}));

function Home() {
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showWelcome, setShowWelcome] = useState(false);
  const navigate = useNavigate();
  const { user, setUser, resetGame } = useGame();
  const resetGameCalledRef = useRef(false);
  
  useEffect(() => {
    // Reset game state when returning to home, but only once
    if (!resetGameCalledRef.current) {
      resetGame();
      resetGameCalledRef.current = true;
    }
    
    // If user already exists in context, show welcome screen
    if (user) {
      setUsername(user.username);
      setShowWelcome(true);
    }
  }, [user, resetGame]);
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!username.trim()) {
      setError('Username is required');
      return;
    }
    
    setLoading(true);
    try {
      const userData = await createUser(username);
      setUser(userData);
      setLoading(false);
      setShowWelcome(true);
    } catch (err) {
      setError(`Failed to create user. Please try again. ${err.message}`);
      setLoading(false);
    }
  };
  
  const handlePlayGame = () => {
    navigate('/game');
  };
  
  return (
    <Container maxWidth="md" sx={{ py: 8 }}>
      <Fade in={true} timeout={800}>
        <Box sx={{ textAlign: 'center', mb: 6 }}>
          <Typography 
            variant="h2" 
            component="h1" 
            gutterBottom 
            sx={{ 
              fontWeight: 'bold',
              color: 'primary.main',
              letterSpacing: '-0.5px'
            }}
          >
            Globetrotter Challenge
          </Typography>
          <Typography variant="h5" color="text.secondary" sx={{ mb: 4 }}>
            The Ultimate Travel Guessing Game!
          </Typography>
        </Box>
      </Fade>
      
      <Fade in={true} timeout={1000}>
        <StyledPaper elevation={3}>
          {!showWelcome ? (
            <>
              <GlobeAvatar>
                <PublicIcon fontSize="large" />
              </GlobeAvatar>
              <Typography variant="h4" gutterBottom>
                Join the Adventure
              </Typography>
              <Typography variant="body1" color="text.secondary" sx={{ mb: 3, textAlign: 'center' }}>
                Enter your username to start exploring the world!
              </Typography>
              
              <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%', mt: 1 }}>
                <TextField
                  fullWidth
                  label="Username"
                  variant="outlined"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  error={!!error}
                  helperText={error}
                  sx={{ mb: 3 }}
                />
                <ColorButton
                  type="submit"
                  fullWidth
                  disabled={loading}
                  startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
                >
                  {loading ? 'Joining...' : 'Join Game'}
                </ColorButton>
              </Box>
            </>
          ) : (
            <>
              <GlobeAvatar>
                <EmojiEventsIcon fontSize="large" />
              </GlobeAvatar>
              <Typography variant="h4" gutterBottom>
                Welcome, {username}!
              </Typography>
              <Typography variant="body1" color="text.secondary" sx={{ mb: 4, textAlign: 'center' }}>
                Ready to test your knowledge of famous destinations around the world?
              </Typography>
              <ColorButton
                onClick={handlePlayGame}
                startIcon={<PlayArrowIcon />}
                size="large"
              >
                Play Game
              </ColorButton>
            </>
          )}
        </StyledPaper>
      </Fade>
    </Container>
  );
}

export default Home; 