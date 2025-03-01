import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getUser, createUser, getGameSummary } from '../services/api';
import { useGame } from '../context/GameContext';
import {
  Container,
  Typography,
  Box,
  Button,
  TextField,
  Paper,
  Avatar,
  Fade,
  Grow,
  CircularProgress,
  Divider,
  Card,
  CardMedia
} from '@mui/material';
import { styled } from '@mui/material/styles';
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents';
import PersonIcon from '@mui/icons-material/Person';
import ErrorIcon from '@mui/icons-material/Error';
import PublicIcon from '@mui/icons-material/Public';

// Styled components
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

const ScoreAvatar = styled(Avatar)(({ theme }) => ({
  backgroundColor: theme.palette.secondary.main,
  width: 80,
  height: 80,
  margin: '0 auto 16px',
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

// Random travel images for preview
const travelImages = [
  'https://images.unsplash.com/photo-1530521954074-e64f6810b32d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1350&q=80',
  'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1350&q=80',
  'https://images.unsplash.com/photo-1488646953014-85cb44e25828?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1350&q=80',
  'https://images.unsplash.com/photo-1507608616759-54f48f0af0ee?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1350&q=80',
  'https://images.unsplash.com/photo-1519922639192-e73293ca430e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1350&q=80'
];

function Challenge() {
  const { username, gameId } = useParams();
  const [challengerInfo, setChallengerInfo] = useState(null);
  const [gameSummary, setGameSummary] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [playerName, setPlayerName] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();
  const { setUser } = useGame();
  
  // Select a random image for the preview
  const previewImage = travelImages[Math.floor(Math.random() * travelImages.length)];
  
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        
        // Fetch user info
        const user = await getUser(username);
        setChallengerInfo(user);
        
        // If gameId is provided, fetch game summary
        if (gameId) {
          try {
            const summary = await getGameSummary(gameId);
            setGameSummary(summary);
          } catch (summaryErr) {
            console.error('Failed to fetch game summary:', summaryErr);
            // Continue even if summary fetch fails
          }
        }
        
        setLoading(false);
      } catch (err) {
        console.error('Error fetching challenge data:', err);
        setError('Could not find the challenger. The link might be invalid.');
        setLoading(false);
      }
    };
    
    fetchData();
  }, [username, gameId]);
  
  const handleStartChallenge = async (e) => {
    e.preventDefault();
    
    if (!playerName.trim()) {
      setError('Please enter your name to start the challenge');
      return;
    }
    
    setSubmitting(true);
    try {
      const user = await createUser(playerName);
      setUser(user);
      navigate('/game');
    } catch (err) {
      setError('Failed to create user. Please try again.');
      setSubmitting(false);
    }
  };
  
  // Note: Meta tags are now handled server-side for better social media sharing
  
  if (loading) {
    return (
      <Container maxWidth="md" sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress size={60} />
        <Typography variant="h5" sx={{ mt: 2 }}>
          Loading challenge...
        </Typography>
      </Container>
    );
  }
  
  if (error && !challengerInfo) {
    return (
      <Container maxWidth="md" sx={{ py: 8 }}>
        <Fade in={true} timeout={800}>
          <StyledPaper>
            <ErrorIcon color="error" sx={{ fontSize: 80, mb: 2 }} />
            <Typography variant="h4" gutterBottom>
              Challenge Not Found
            </Typography>
            <Typography variant="body1" paragraph align="center">
              {error}
            </Typography>
            <Button 
              variant="contained" 
              color="primary" 
              onClick={() => navigate('/')}
              sx={{ mt: 2 }}
            >
              Go Home
            </Button>
          </StyledPaper>
        </Fade>
      </Container>
    );
  }
  
  return (
    <Container maxWidth="md" sx={{ py: 8 }}>
      {/* Preview image for social media sharing */}
      <Box sx={{ display: 'none' }}>
        <img src={previewImage} alt="Globetrotter Challenge" />
      </Box>
      
      <Fade in={true} timeout={800}>
        <Box sx={{ textAlign: 'center', mb: 6 }}>
          <Typography 
            variant="h2" 
            component="h1" 
            gutterBottom 
            sx={{ fontWeight: 'bold', color: 'primary.main' }}
          >
            You've Been Challenged!
          </Typography>
        </Box>
      </Fade>
      
      <Grow in={true} timeout={1000}>
        <StyledPaper>
          {/* Preview image visible on the page */}
          <Card sx={{ width: '100%', mb: 3, borderRadius: 2, overflow: 'hidden' }}>
            <CardMedia
              component="img"
              height="200"
              image={previewImage}
              alt="Travel destination"
            />
          </Card>
          
          <ScoreAvatar>
            <EmojiEventsIcon fontSize="large" />
          </ScoreAvatar>
          
          <Typography variant="h5" gutterBottom>
            {challengerInfo.username} has challenged you!
          </Typography>
          
          <Box sx={{ my: 3, textAlign: 'center' }}>
            <Typography variant="body1" color="text.secondary" gutterBottom>
              Their current score:
            </Typography>
            <Typography variant="h3" color="secondary" sx={{ fontWeight: 'bold' }}>
              {gameSummary ? `${gameSummary.total_correct}/${gameSummary.total_questions}` : 
                `${challengerInfo.total_correct || 0}/${challengerInfo.total_count || 0}`}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Can you beat it?
            </Typography>
          </Box>
          
          <Box component="form" onSubmit={handleStartChallenge} sx={{ width: '100%' }}>
            <Typography variant="h6" gutterBottom>
              Enter your name to accept the challenge:
            </Typography>
            
            <TextField
              fullWidth
              label="Your Name"
              variant="outlined"
              value={playerName}
              onChange={(e) => setPlayerName(e.target.value)}
              error={!!error}
              helperText={error}
              sx={{ mb: 3 }}
            />
            
            <ColorButton
              type="submit"
              fullWidth
              disabled={submitting}
              startIcon={submitting ? <CircularProgress size={20} color="inherit" /> : <PersonIcon />}
            >
              {submitting ? 'Joining...' : 'Accept Challenge'}
            </ColorButton>
          </Box>
        </StyledPaper>
      </Grow>
    </Container>
  );
}

export default Challenge; 