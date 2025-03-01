import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { startGame, getNextQuestion, submitAnswer, getGameResults, resetUserScore } from '../services/api';
import { useGame } from '../context/GameContext';
import ShareButton from '../components/ShareButton';
import {
  Container,
  Typography,
  Box,
  Button,
  Grid,
  Paper,
  LinearProgress,
  Fade,
  Grow,
  CircularProgress,
  Chip,
  Divider
} from '@mui/material';
import { styled } from '@mui/material/styles';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CancelIcon from '@mui/icons-material/Cancel';
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents';
import PublicIcon from '@mui/icons-material/Public';
import SentimentVeryDissatisfiedIcon from '@mui/icons-material/SentimentVeryDissatisfied';
import Confetti from 'react-confetti';
import { useWindowSize } from 'react-use';

// Styled components
const OptionButton = styled(Button)(({ theme, color, selected, correct, incorrect, isCorrectOption }) => ({
  width: '100%',
  padding: '16px',
  borderRadius: 12,
  fontSize: '1.1rem',
  fontWeight: 'bold',
  marginBottom: theme.spacing(2),
  backgroundColor: selected 
    ? (correct ? theme.palette.success.main : incorrect ? theme.palette.error.main : color)
    : isCorrectOption ? theme.palette.success.main : color,
  color: 'white',
  '&:hover': {
    backgroundColor: selected 
      ? (correct ? theme.palette.success.dark : incorrect ? theme.palette.error.dark : color)
      : isCorrectOption ? theme.palette.success.dark : color,
    opacity: 0.9,
  },
  transition: 'all 0.3s ease',
}));

const QuestionCard = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  borderRadius: 16,
  boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
  marginBottom: theme.spacing(4),
}));

const ResultCard = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  borderRadius: 16,
  boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
  textAlign: 'center',
}));

const ScoreChip = styled(Chip)(({ theme }) => ({
  fontWeight: 'bold',
  fontSize: '1rem',
  padding: theme.spacing(1),
}));

function Game() {
  const [loading, setLoading] = useState(true);
  const [answering, setAnswering] = useState(false);
  const [selectedOption, setSelectedOption] = useState(null);
  const [feedback, setFeedback] = useState(null);
  const [showConfetti, setShowConfetti] = useState(false);
  const { width, height } = useWindowSize();
  const navigate = useNavigate();
  // Use refs to track state and prevent infinite loops
  const gameInitializedRef = useRef(false);
  const apiCallInProgressRef = useRef(false);
  const scoreResetRef = useRef(false);
  
  const { 
    user, 
    gameId, 
    setGameId, 
    gameState, 
    setGameState,
    currentQuestion,
    setCurrentQuestion,
    results,
    setResults,
    setUser,
    resetScore
  } = useGame();
  
  // Colors for option buttons (Kahoot-style)
  const optionColors = ['#e21b3c', '#1368ce', '#26890c', '#ffa602'];
  
  // Define loadNextQuestion outside useEffect to avoid dependency issues
  const loadNextQuestion = async (id) => {
    if (!id) {
      console.error('Cannot load next question: gameId is undefined');
      return;
    }
    
    setLoading(true);
    setSelectedOption(null);
    setFeedback(null);
    setShowConfetti(false);
    
    try {
      console.log(`Loading next question for game ID: ${id}`);
      const questionData = await getNextQuestion(id);
      console.log('Question data received:', questionData);
      
      // Check if the game is finished
      if (questionData.game_finished || questionData.has_next === false) {
        console.log('Game is finished, loading results');
        try {
          const resultsData = await getGameResults(id);
          console.log('Game results received:', resultsData);
          setResults(resultsData);
          setGameState('finished');
        } catch (resultErr) {
          console.error('Failed to load game results:', resultErr);
          // Even if results fail to load, mark the game as finished
          setGameState('finished');
        }
        return;
      }
      
      // Verify that questionData has the expected structure after transformation
      if (!questionData.clues || !questionData.options) {
        console.error('Question data is missing expected properties after transformation:', questionData);
        // Try to recover by setting a default structure
        const recoveredData = {
          ...questionData,
          clues: questionData.clues || [questionData.question || 'Question not available'],
          options: questionData.options || []
        };
        console.log('Attempting to recover with:', recoveredData);
        setCurrentQuestion(recoveredData);
      } else {
        console.log('Setting current question with valid data');
        setCurrentQuestion(questionData);
      }
    } catch (err) {
      console.error('Failed to load question:', err);
      setCurrentQuestion(null); // Reset current question on error
    } finally {
      setLoading(false);
    }
  };
  
  // Separate useEffect for handling game initialization
  useEffect(() => {
    if (!user) {
      navigate('/');
      return;
    }
    
    const initGame = async () => {
      // Prevent multiple API calls
      if (gameInitializedRef.current || apiCallInProgressRef.current) return;
      
      apiCallInProgressRef.current = true;
      gameInitializedRef.current = true;
      
      setLoading(true);
      try {
        console.log('Starting game for user:', user.username);
        
        // Reset the score reset flag
        scoreResetRef.current = true;
        
        // Reset the score counters in the UI before starting a new game
        if (user && (user.correct_count > 0 || user.total_count > 0)) {
          // Reset score locally
          const resetUser = await resetUserScore(user.username);
          setUser(resetUser);
          console.log('Reset user score locally (init game):', resetUser);
        }
        
        // Reset score in context
        resetScore();
        
        const gameData = await startGame(user.username);
        console.log('Game data received:', gameData);
        
        if (!gameData || !gameData.game_id) {
          console.error('Invalid game data received:', gameData);
          setLoading(false);
          gameInitializedRef.current = false;
          apiCallInProgressRef.current = false;
          return;
        }
        
        setGameId(gameData.game_id);
        setGameState('playing');
        await loadNextQuestion(gameData.game_id);
      } catch (err) {
        console.error('Failed to start game:', err);
        setLoading(false);
        gameInitializedRef.current = false;
      } finally {
        apiCallInProgressRef.current = false;
      }
    };
    
    if (gameState === 'initial') {
      initGame();
    }
  }, [user, navigate, gameState, setGameId, setGameState]); // Remove currentQuestion, gameId, and loadNextQuestion
  
  // Separate useEffect for loading next question
  useEffect(() => {
    if (gameState === 'playing' && gameId && !currentQuestion && !apiCallInProgressRef.current) {
      loadNextQuestion(gameId);
    }
  }, [gameState, gameId, currentQuestion]);
  
  // Reset the ref when the game state changes back to initial
  useEffect(() => {
    if (gameState === 'initial') {
      gameInitializedRef.current = false;
      
      // Reset the user's score when the game state changes to initial, but only once
      if (user && (user.correct_count > 0 || user.total_count > 0) && !scoreResetRef.current) {
        scoreResetRef.current = true;
        
        // Use an async function inside useEffect
        const resetUserScoreAsync = async () => {
          // Reset score locally
          const resetUser = await resetUserScore(user.username);
          setUser(resetUser);
          console.log('Reset user score locally (state change):', resetUser);
        };
        
        resetUserScoreAsync();
        
        // Reset score in context
        resetScore();
      }
    } else {
      // Reset the flag when game state changes from initial
      scoreResetRef.current = false;
    }
  }, [gameState, user, setUser, resetScore]);
  
  const handleAnswer = async (optionId, optionText) => {
    if (answering) return;
    
    // If this is the correct answer after an incorrect answer, just update the selection
    if (feedback && !feedback.correct && optionId === feedback.correct_option_id) {
      setSelectedOption(optionText);
      return;
    }
    
    setSelectedOption(optionText);
    setAnswering(true);
    
    try {
      console.log(`Submitting answer: gameId=${gameId}, questionId=${currentQuestion.question_id}, optionId=${optionId}`);
      
      // Make sure we have all required parameters
      if (!gameId || !currentQuestion || !currentQuestion.question_id) {
        console.error('Missing required parameters for submitting answer:', {
          gameId,
          questionId: currentQuestion?.question_id,
          optionId
        });
        throw new Error('Missing required parameters for submitting answer');
      }
      
      const result = await submitAnswer(gameId, currentQuestion.question_id, optionId);
      setFeedback(result);
      
      // Update the user object with the latest score information
      if (user && result) {
        const updatedUser = {
          ...user,
          correct_count: user.correct_count ? (result.correct ? user.correct_count + 1 : user.correct_count) : (result.correct ? 1 : 0),
          total_count: user.total_count ? user.total_count + 1 : 1
        };
        console.log('Updating user with new score:', updatedUser);
        setUser(updatedUser);
      }
      
      // Only show confetti for correct answers
      setShowConfetti(result.correct);
    } catch (err) {
      console.error('Failed to submit answer:', err);
    } finally {
      setAnswering(false);
    }
  };
  
  const handleNextQuestion = () => {
    loadNextQuestion(gameId);
  };
  
  const handlePlayAgain = async () => {
    // Reset the score reset flag
    scoreResetRef.current = false;
    
    // Reset the user's score when starting a new game
    if (user) {
      // Reset score locally
      const resetUser = await resetUserScore(user.username);
      setUser(resetUser);
      console.log('Reset user score locally:', resetUser);
    }
    
    // Reset score in context
    resetScore();
    
    // Reset game state
    setGameState('initial');
    setCurrentQuestion(null);
    setResults(null);
  };
  
  if (loading && !currentQuestion) {
    return (
      <Container maxWidth="md" sx={{ py: 8, textAlign: 'center' }}>
        <CircularProgress size={60} />
        <Typography variant="h5" sx={{ mt: 2 }}>
          Loading your adventure...
        </Typography>
      </Container>
    );
  }
  
  if (gameState === 'finished') {
    // If we have results, show them
    if (results) {
      return (
        <Container maxWidth="md" sx={{ py: 8 }}>
          {results.score_percentage >= 70 && <Confetti width={width} height={height} recycle={false} numberOfPieces={200} />}
          
          <Fade in={true} timeout={800}>
            <Box sx={{ textAlign: 'center', mb: 6 }}>
              <Typography 
                variant="h2" 
                component="h1" 
                gutterBottom 
                sx={{ fontWeight: 'bold', color: 'primary.main' }}
              >
                Game Results
              </Typography>
            </Box>
          </Fade>
          
          <Grow in={true} timeout={1000}>
            <ResultCard>
              <Box sx={{ mb: 3 }}>
                <EmojiEventsIcon sx={{ fontSize: 80, color: 'primary.main', mb: 2 }} />
                <Typography variant="h4" gutterBottom>
                  {results.score_percentage >= 70 
                    ? 'Congratulations, World Explorer!' 
                    : 'Nice Try, Adventurer!'}
                </Typography>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  {results.score_percentage >= 70 
                    ? 'You really know your way around the globe!' 
                    : 'Keep exploring to improve your knowledge!'}
                </Typography>
              </Box>
              
              <Box sx={{ mb: 4 }}>
                <Typography variant="h3" color="primary" sx={{ fontWeight: 'bold' }}>
                  {((results.total_correct / results.total_questions) * 100).toFixed(2)}%
                </Typography>
                <Typography variant="body1" color="text.secondary">
                  You answered {results.total_correct} out of {results.total_questions} questions correctly
                </Typography>
              </Box>
              
              <Divider sx={{ my: 3 }} />
              
              <Grid container spacing={2} sx={{ mb: 4 }}>
                <Grid item xs={12} sm={6}>
                  <Button 
                    variant="contained" 
                    color="primary" 
                    fullWidth 
                    size="large"
                    onClick={handlePlayAgain}
                  >
                    Play Again
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <ShareButton 
                    username={user.username} 
                    gameId={gameId}
                    score={{
                      correct: results.total_correct,
                      total: results.total_questions
                    }}
                  />
                </Grid>
              </Grid>
            </ResultCard>
          </Grow>
        </Container>
      );
    } else {
      // If game is finished but we don't have results, show a fallback
      return (
        <Container maxWidth="md" sx={{ py: 8, textAlign: 'center' }}>
          <Typography variant="h4" gutterBottom>
            Game Complete!
          </Typography>
          <Typography variant="body1" sx={{ mb: 4 }}>
            Your game has been completed, but we couldn't load the results.
          </Typography>
          <Button 
            variant="contained" 
            color="primary" 
            size="large"
            onClick={handlePlayAgain}
          >
            Play Again
          </Button>
        </Container>
      );
    }
  }
  
  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      {showConfetti && (
        <Confetti 
          width={width} 
          height={height} 
          recycle={false} 
          numberOfPieces={200}
        />
      )}
      
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <ScoreChip 
          icon={<PublicIcon />} 
          label={`${user?.username || 'Player'}`} 
          color="primary" 
          variant="outlined" 
        />
        <ScoreChip 
          icon={<EmojiEventsIcon />} 
          label={`Score: ${user?.correct_count || 0}/${user?.total_count || 0}`} 
          color="secondary" 
        />
      </Box>
      
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', my: 8 }}>
          <CircularProgress size={60} />
        </Box>
      ) : currentQuestion ? (
        <>
          <LinearProgress 
            variant="determinate" 
            value={feedback ? 100 : 0} 
            color="primary" 
            sx={{ height: 8, borderRadius: 4, mb: 4 }} 
          />
          
          <Fade in={true} timeout={500}>
            <QuestionCard>
              <Typography variant="h4" gutterBottom sx={{ fontWeight: 'bold' }}>
                Where am I?
              </Typography>
              {currentQuestion && currentQuestion.clues && currentQuestion.clues.map((clue, index) => (
                <Typography key={index} variant="h6" sx={{ mb: 2 }}>
                  {clue}
                </Typography>
              ))}
            </QuestionCard>
          </Fade>
          
          {currentQuestion && currentQuestion.options && currentQuestion.options.length > 0 ? (
            <Grid container spacing={2}>
              {currentQuestion.options.map((option, index) => (
                <Grid item xs={12} sm={6} key={index}>
                  <Grow in={true} timeout={500 + (index * 100)}>
                    <Box>
                      <OptionButton
                        color={optionColors[index % optionColors.length]}
                        onClick={() => {
                          // Allow clicking if no option selected yet OR
                          // if this is the correct option after an incorrect answer
                          if (!selectedOption || (feedback && !feedback.correct && feedback.correct_option_id === option.id)) {
                            handleAnswer(option.id, option.text);
                          }
                        }}
                        disabled={!!selectedOption && !(feedback && feedback.correct_option_id === option.id)}
                        selected={selectedOption === option.text}
                        correct={feedback && selectedOption === option.text && feedback.correct}
                        incorrect={feedback && selectedOption === option.text && !feedback.correct}
                        isCorrectOption={feedback && feedback.correct_option_id === option.id}
                      >
                        {option.text}
                        {feedback && (
                          (selectedOption === option.text && feedback.correct && <CheckCircleIcon sx={{ ml: 1 }} />) || 
                          (selectedOption === option.text && !feedback.correct && <CancelIcon sx={{ ml: 1 }} />) ||
                          (feedback.correct_option_id === option.id && selectedOption !== option.text && <CheckCircleIcon sx={{ ml: 1 }} />)
                        )}
                      </OptionButton>
                    </Box>
                  </Grow>
                </Grid>
              ))}
            </Grid>
          ) : (
            <Box sx={{ textAlign: 'center', my: 4 }}>
              <Typography variant="h6" color="text.secondary">
                No options available for this question.
              </Typography>
              <Button 
                variant="contained" 
                color="primary" 
                sx={{ mt: 3 }}
                onClick={handleNextQuestion}
              >
                Try Next Question
              </Button>
            </Box>
          )}
          
          {feedback && (
            <Fade in={true} timeout={500}>
              <Box sx={{ mt: 4, textAlign: 'center' }}>
                <Typography variant="h4" sx={{ mb: 2, color: feedback.correct ? 'success.main' : 'error.main' }}>
                  {feedback.correct ? 'Correct!' : 'Incorrect!'}
                </Typography>
                {!feedback.correct && (
                  <SentimentVeryDissatisfiedIcon 
                    className="shake"
                    sx={{ 
                      fontSize: 80, 
                      color: 'error.main', 
                      mb: 2
                    }} 
                  />
                )}
                <Typography variant="h6" sx={{ 
                  mb: 3, 
                  color: feedback.correct ? 'text.primary' : 'text.secondary',
                  fontStyle: feedback.correct ? 'normal' : 'italic'
                }}>
                  {feedback.correct ? feedback.fun_fact : feedback.trivia}
                </Typography>
                <Button 
                  variant="contained" 
                  color="primary" 
                  size="large"
                  onClick={handleNextQuestion}
                >
                  Next Question
                </Button>
              </Box>
            </Fade>
          )}
        </>
      ) : (
        <Box sx={{ textAlign: 'center', my: 8 }}>
          <Typography variant="h5" color="text.secondary">
            No question available. Please try again.
          </Typography>
          <Button 
            variant="contained" 
            color="primary" 
            sx={{ mt: 3 }}
            onClick={() => {
              localStorage.removeItem('username');
              navigate('/');
            }}
          >
            Restart Game
          </Button>
        </Box>
      )}
    </Container>
  );
}

export default Game; 