import React, { useState } from 'react';
import { 
  Button, 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  TextField, 
  IconButton, 
  Box, 
  Typography,
  Tooltip,
  Snackbar,
  Alert
} from '@mui/material';
import { styled } from '@mui/material/styles';
import ShareIcon from '@mui/icons-material/Share';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import WhatsAppIcon from '@mui/icons-material/WhatsApp';
import CloseIcon from '@mui/icons-material/Close';

const ShareIconButton = styled(IconButton)(({ theme }) => ({
  backgroundColor: '#25D366', // WhatsApp green
  color: 'white',
  padding: 12,
  '&:hover': {
    backgroundColor: '#128C7E',
  },
}));

const CopyButton = styled(Button)(({ theme }) => ({
  marginLeft: theme.spacing(1),
}));

function ShareButton({ username, gameId, score }) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  
  const scoreText = score ? `${score.correct}/${score.total}` : '';
  const shareUrl = `${window.location.origin}/challenge/${username}/${gameId}`;
  const shareTitle = `I scored ${scoreText} in the Globetrotter Challenge! Can you beat me?`;
  const whatsappUrl = `https://wa.me/?text=${encodeURIComponent(`${shareTitle} ${shareUrl}`)}`;
  
  const handleOpen = () => {
    setOpen(true);
  };
  
  const handleClose = () => {
    setOpen(false);
  };
  
  const handleCopy = () => {
    navigator.clipboard.writeText(shareUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 3000);
  };
  
  return (
    <>
      <Button
        variant="contained"
        color="secondary"
        fullWidth
        size="large"
        startIcon={<ShareIcon />}
        onClick={handleOpen}
      >
        Challenge a Friend
      </Button>
      
      <Dialog 
        open={open} 
        onClose={handleClose}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          Challenge a Friend
          <IconButton
            aria-label="close"
            onClick={handleClose}
            sx={{
              position: 'absolute',
              right: 8,
              top: 8,
            }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        
        <DialogContent>
          <Typography variant="body1" gutterBottom>
            Share this challenge with your friends and see if they can beat your score!
          </Typography>
          
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 3 }}>
            <Tooltip title="Share on WhatsApp">
              <a href={whatsappUrl} target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none' }}>
                <ShareIconButton size="large">
                  <WhatsAppIcon fontSize="large" />
                </ShareIconButton>
              </a>
            </Tooltip>
          </Box>
          
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
            <TextField
              fullWidth
              variant="outlined"
              value={shareUrl}
              InputProps={{
                readOnly: true,
              }}
              onClick={(e) => e.target.select()}
            />
            <CopyButton 
              variant="contained" 
              color="primary"
              onClick={handleCopy}
              startIcon={<ContentCopyIcon />}
            >
              Copy
            </CopyButton>
          </Box>
        </DialogContent>
        
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            Close
          </Button>
        </DialogActions>
      </Dialog>
      
      <Snackbar 
        open={copied} 
        autoHideDuration={3000} 
        onClose={() => setCopied(false)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert severity="success" variant="filled">
          Link copied to clipboard!
        </Alert>
      </Snackbar>
    </>
  );
}

export default ShareButton; 