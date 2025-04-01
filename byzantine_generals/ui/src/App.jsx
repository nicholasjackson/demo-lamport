import { useState } from 'react';
import PropTypes from 'prop-types';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import CssBaseline from '@mui/material/CssBaseline';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';

import Grid from './Grid';

import '@xyflow/react/dist/style.css';
function CustomTabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3, width: '100vw', height: '100vh'}}>{children}</Box>}
    </div>
  );
}

CustomTabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.number.isRequired,
  value: PropTypes.number.isRequired,
};

function Flow() {

  const [value, setValue] =useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const voteClick = () => {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/IssueCommand';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
      });
  }
  
  const resetClick = () => {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Reset';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
      });
  }

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
              <Button key="vote" sx={{ color: '#fff' }} onClick={voteClick}>
                Vote
              </Button>
              <Button key="reset" sx={{ color: '#fff' }} onClick={resetClick}>
                Reset 
              </Button>
          </Box>
        </Toolbar>
      </AppBar>
      <Box component="main" sx={{ p: 3, width: '100vw', height: '95vh', marginTop: '48px' }}>
        <Grid/>
      </Box>
    </Box>
  );
}
 
export default Flow;