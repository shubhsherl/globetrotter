import { createTheme } from '@mui/material/styles';

// Kahoot-inspired color palette
const theme = createTheme({
  palette: {
    primary: {
      main: '#46178f', // Kahoot purple
      light: '#7b44c2',
      dark: '#2e0f5e',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#ff3355', // Kahoot pink/red
      light: '#ff6b7e',
      dark: '#c4002e',
      contrastText: '#ffffff',
    },
    success: {
      main: '#26890c', // Kahoot green
      light: '#5cb82f',
      dark: '#005c00',
      contrastText: '#ffffff',
    },
    error: {
      main: '#e21b3c', // Kahoot red
      light: '#ff5c66',
      dark: '#a80016',
      contrastText: '#ffffff',
    },
    info: {
      main: '#1368ce', // Kahoot blue
      light: '#5a96ff',
      dark: '#003e9c',
      contrastText: '#ffffff',
    },
    warning: {
      main: '#ffa602', // Kahoot yellow
      light: '#ffd74c',
      dark: '#c67700',
      contrastText: '#000000',
    },
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
    },
  },
  typography: {
    fontFamily: '"Montserrat", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
    },
    h2: {
      fontWeight: 700,
    },
    h3: {
      fontWeight: 600,
    },
    h4: {
      fontWeight: 600,
    },
    h5: {
      fontWeight: 500,
    },
    h6: {
      fontWeight: 500,
    },
    button: {
      fontWeight: 600,
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 30,
          padding: '10px 24px',
          boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
        },
        containedPrimary: {
          '&:hover': {
            boxShadow: '0 6px 10px rgba(0,0,0,0.2)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        rounded: {
          borderRadius: 16,
        },
        elevation1: {
          boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
        },
      },
    },
  },
});

export default theme; 