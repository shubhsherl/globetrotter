import axios from 'axios';

const API_URL = '/api';

// Simple debounce implementation to prevent duplicate API calls
let lastCallTimestamp = 0;
const DEBOUNCE_DELAY = 1000; // 1 second

export const getRandomDestination = async () => {
  const response = await axios.get(`${API_URL}/destinations/random`);
  return response.data;
};

export const createUser = async (username) => {
  const response = await axios.post(`${API_URL}/users`, { username });
  return response.data;
};

export const getUser = async (username) => {
  const response = await axios.get(`${API_URL}/users/${username}`);
  return response.data;
};

export const updateUserScore = async (username, correct) => {
  const response = await axios.post(`${API_URL}/users/${username}/score`, { correct });
  return response.data;
};

export const resetUserScore = async (username) => {
  try {
    console.log(`Resetting score for user: ${username}`);
    const response = await axios.post(`${API_URL}/users/${username}/reset-score`);
    console.log('Reset score response:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error resetting user score:', error);
    // Even if the API call fails, return a default user object with reset scores
    return {
      username,
      correct_count: 0,
      total_count: 0
    };
  }
};

// Track in-flight requests to prevent duplicates
let startGameInProgress = false;

export const startGame = async (username) => {
  console.log(`API call: startGame for user ${username}`);
  
  // Debounce mechanism
  const now = Date.now();
  if (now - lastCallTimestamp < DEBOUNCE_DELAY) {
    console.log('Debouncing startGame call - too soon after previous call');
    throw new Error('Please wait before starting a new game');
  }
  
  // Prevent concurrent calls
  if (startGameInProgress) {
    console.log('startGame call already in progress, skipping duplicate call');
    throw new Error('Game start already in progress');
  }
  
  try {
    startGameInProgress = true;
    lastCallTimestamp = now;
    
    const response = await axios.post(`${API_URL}/game/play`, { username });
    console.log('startGame response:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error in startGame:', error);
    throw error;
  } finally {
    startGameInProgress = false;
  }
};

export const getNextQuestion = async (gameId) => {
  try {
    const response = await axios.get(`${API_URL}/game/${gameId}/next-question`);
    console.log('Raw API response from next-question:', response.data);
    
    // Validate the response data structure
    const data = response.data;
    if (!data) {
      console.error('Empty response from next-question API');
      throw new Error('Invalid response from server');
    }
    
    // Check if game is finished - either by game_finished flag or has_next being false
    if (data.game_finished || data.has_next === false) {
      console.log('Game is finished based on API response:', data);
      return {
        ...data,
        game_finished: true // Ensure game_finished is set to true
      };
    }
    
    // Store the question ID globally so it can be accessed by submitAnswer
    window.currentQuestionId = data.question_id;
    console.log('Stored question ID:', window.currentQuestionId);
    
    // Transform the data structure to match what the frontend expects
    // The API returns options_display as an object with id:text pairs
    // We need to transform it to an array of {id, text} objects
    if (data.options_display) {
      const transformedData = {
        ...data,
        question_id: data.question_id, // Keep the question_id in the transformed data
        clues: [data.question], // Use the question as a clue
        options: Object.entries(data.options_display).map(([id, text]) => ({
          id: parseInt(id),
          text
        }))
      };
      
      console.log('Transformed question data:', transformedData);
      return transformedData;
    }
    
    // If we can't transform the data, throw an error
    console.error('Invalid question data structure:', data);
    throw new Error('Invalid question data received');
  } catch (error) {
    console.error('Error fetching next question:', error);
    throw error;
  }
};

export const submitAnswer = async (gameId, questionId, answerId) => {
  try {
    console.log(`Submitting answer for game ${gameId}, question ${questionId}, option ID: ${answerId}`);
    
    const response = await axios.post(`${API_URL}/game/${gameId}/submit-answer`, { 
      game_id: gameId,
      question_id: questionId,
      selected_destination: answerId
    });
    
    console.log('Submit answer response:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error submitting answer:', error);
    throw error;
  }
};

export const getGameResults = async (gameId) => {
  const response = await axios.get(`${API_URL}/game/${gameId}/result`);
  return response.data;
};

export const getGameSummary = async (gameId) => {
  try {
    console.log(`Fetching game summary for game ID: ${gameId}`);
    const response = await axios.get(`${API_URL}/game/${gameId}/summary`);
    console.log('Game summary response:', response.data);
    return response.data;
  } catch (error) {
    console.error('Error fetching game summary:', error);
    throw error;
  }
}; 