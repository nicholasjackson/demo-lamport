import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import CssBaseline from '@mui/material/CssBaseline';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';

import Grid from './Grid';

import '@xyflow/react/dist/style.css';

function Flow() {

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar component="nav">
        <Toolbar>
          <Typography
            variant="h6"
            component="div"
            sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block' } }}
          >
            Byzantine Generals
          </Typography>
          <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
              <Button key="vote" sx={{ color: '#fff' }}>
                Vote
              </Button>
          </Box>
        </Toolbar>
      </AppBar>
      <Box component="main" sx={{ p: 3, width: '100vw', height: '100vh' }}>
        <Grid/>
      </Box>
    </Box>
  );
}
 
export default Flow;
