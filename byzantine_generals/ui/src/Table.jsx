import { useState, useEffect, useCallback } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import Box from '@mui/material/Box';

const baseColumns = [
  {
    field: 'col1',
    headerName: '',
    width: 150,
    editable: false,
    cellClassName: (params) => {
      if (params.row.col1 === '') {
        return 'white_column';
      }

      return 'general_column';
    },
  },
];

function Table(props) {
  let [rows, setRows] = useState([]);
  let [cols, setCols] = useState([]);

  async function getData() {
    // fetch the node data
    const nodeUrl = 'http://localhost:8080/proto.server.v1.CommanderService/Nodes';
    const nodesResponse = await fetch(nodeUrl, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})});
    const nodesData = await nodesResponse.json();
   
    // fetch the decision data
    const decisionUrl = 'http://localhost:8080/proto.server.v1.CommanderService/Decisions';
    const decisionsResponse = await fetch(decisionUrl, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({allData: true})});
    const decisionData = await decisionsResponse.json();

    // build the columns
    cols = [...baseColumns];
    nodesData?.nodes.forEach((node) => {
      if (node.data.label === "Commander") {
        return;
      }

      cols.push({
        field: node.id,
        headerName: node.data.label,
        width: 100,
        editable: false,
        headerClassName: 'table_header',
      });
      
      cols.push({
        field: node.id + '_decision',
        headerName: '',
        width: 100,
        editable: false,
        headerClassName: 'table_header',
      });
    });

    cols.push({
      field: 'decision',
      headerName: 'Decision',
      width: 120,
      editable: false,
      headerClassName: 'table_header',
    });

    console.log(cols);
    setCols(cols);

    // find the commander
    const commander = nodesData.nodes.find((node) => node.data.label === 'Commander');

    if (props.round === 1) {
      rows = buildRound1Rows(decisionData);
    }

    if( props.round === 2) {
      rows = buildRound2Rows(decisionData,nodesData);
    }

    setRows(rows);
  }

  function buildRound1Rows(decisionData) {
    let rows = [];
    
    let row = {
      id: 0,
      col1: 'Commander',
    };

    decisionData.decisions?.forEach((decision) => {
      if (decision.round === 1) {
        row[decision.from] = decision.decision;
        row[decision.from + '_decision'] = '';

      }
    });


    rows.push(row);
    
    row = {
      id: 1,
      col1: '',
    };

    var attackCount = 0;
    var retreatCount = 0;
    var totalCount = 0;
    decisionData.decisions?.forEach((decision) => {
      if (decision.round === 1) {
        row[decision.from] = decision.decision;
        row[decision.from + '_decision'] = '';
        
        if (decision.decision === 'attack') {
          attackCount++;
        }
        
        if (decision.decision === 'retreat') {
          retreatCount++;
        }

        totalCount++;
      }
    });
    
    row.decision = 'Inconclusive';

    if (attackCount === totalCount) {
      row.decision = 'attack';
    }

    if (retreatCount === totalCount) {
      row.decision = 'retreat';
    }

    rows.push(row);

    return rows;
  }
  
  function buildRound2Rows(decisionData, nodesData) {
    let rows = [];
  
    let rowID = 0;

    nodesData.nodes.forEach((node) => {

    if (node.data.label === "Commander") {
      return;
    }

    let row = {
      id: rowID,
      col1: node.data.label,
    };

    rowID++;

    decisionData.decisions?.forEach((decision) => {
      if (decision.round === 1) {
        row[decision.from] = decision.decision;
        row[decision.from + '_decision'] = '';
      }
    });


    rows.push(row);
    
    row = {
      id: 1,
      col1: '',
    };
  });

    //var attackCount = 0;
    //var retreatCount = 0;
    //var totalCount = 0;
    //decisionData.decisions?.forEach((decision) => {
    //  if (decision.round === 1) {
    //    row[decision.from] = decision.decision;
    //    row[decision.from + '_decision'] = '';
    //    
    //    if (decision.decision === 'attack') {
    //      attackCount++;
    //    }
    //    
    //    if (decision.decision === 'retreat') {
    //      retreatCount++;
    //    }

    //    totalCount++;
    //  }
    //});
    
    //row.decision = 'Inconclusive';

    //if (attackCount === totalCount) {
    //  row.decision = 'attack';
    //}

    //if (retreatCount === totalCount) {
    //  row.decision = 'retreat';
    //}


    return rows;
  }

  let getRowClass = (params) => {
    if (params.row.col1 === '') {
      return 'decision_column';
    }

    return '';
  }

  useEffect(() => {
    const dataInterval = setInterval(getData, 5000);

    return () => clearInterval(dataInterval);
  }, []);

  return (
    <>
    <h2>Round {props.round}</h2>
    <Box sx={{ width: '100%' }}>
      <DataGrid
        sx={{
          '&.MuiDataGrid-root--densityCompact .MuiDataGrid-cell': { py: '8px' },
          '&.MuiDataGrid-root--densityStandard .MuiDataGrid-cell': { py: '15px' },
          '&.MuiDataGrid-root--densityComfortable .MuiDataGrid-cell': { py: '22px' },
        }}
        getRowClassName={getRowClass}
        rows={rows}
        columns={cols}
        initialState={{
          pagination: {
            paginationModel: {
              pageSize: 10,
            },
          },
        }}
        pageSizeOptions={[10]}
        getRowHeight={() => 'auto'}
        disableRowSelectionOnClick
      />
    </Box> 
    </>
  )
}

export default Table;