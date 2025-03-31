import { useState, useEffect, useCallback } from 'react';
import { DataGrid } from '@mui/x-data-grid';
import Box from '@mui/material/Box';

const columns = [
  {
    field: 'name',
    headerName: 'Name',
    width: 150,
    editable: false,
  },
  {
    field: 'traitor',
    headerName: 'Is Traitor',
    width: 110,
    editable: false,
  },
  {
    field: 'decision1',
    headerName: 'Round 1 Decision',
    width: 150,
    editable: false,
  },
  {
    field: 'data1',
    headerName: 'Round 1 Data',
    width: 200,
    editable: false,
    renderCell: (params) => {
      const rows = [];
      if (!params.row.data1) {
        return <div></div>;
      }

      for (let i = 0; i < params.row.data1.length; i++) {
        rows.push(<div key={i}>{params.row.data1[i].from}: {params.row.data1[i].command}</div>);
      }

      return (
        <div>
          {rows}
        </div>
      );
    },
  },
  {
    field: 'decision2',
    headerName: 'Round 2 Decision',
    width: 150,
    editable: false,
  },
  {
    field: 'data2',
    headerName: 'Round 2 Data',
    width: 200,
    editable: false,
    renderCell: (params) => {
      const rows = [];
      if (!params.row.data2) {
        return <div></div>;
      }

      for (let i = 0; i < params.row.data2.length; i++) {
        rows.push(<div key={i}>{params.row.data2[i].from}: {params.row.data2[i].command}</div>);
      }

      return (
        <div>
          {rows}
        </div>
      );
    },
  },
];

function Table() {
  let [rows, setRows] = useState([]);

  async function getData() {
    // fetch the node data
    const nodeUrl = 'http://localhost:8080/proto.server.v1.CommanderService/Nodes';
    const nodesResponse = await fetch(nodeUrl, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})});
    const nodesData = await nodesResponse.json();
   
    // fetch the decision data
    const decisionUrl = 'http://localhost:8080/proto.server.v1.CommanderService/Decisions';
    const decisionsResponse = await fetch(decisionUrl, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({allData: true})});
    const decisionData = await decisionsResponse.json();

    rows = [];
    nodesData?.nodes.forEach((node) => {
      const decision1 = decisionData.decisions?.find((d) => d.from === node.id && d.round === 1);
      // round 1 data
      const data1 = decision1?.commands
      for (let i = 0; i < data1?.length; i++) {
        const fromNode = nodesData.nodes.find((n) => n.id === data1[i].from);
        if (fromNode) {
          data1[i].from = fromNode.data.label;
        }
      }
      
      const decision2 = decisionData.decisions?.find((d) => d.from === node.id && d.round === 2);
      // round 1 data
      const data2 = decision1?.commands
      for (let i = 0; i < data1?.length; i++) {
        const fromNode = nodesData.nodes.find((n) => n.id === data2[i].from);
        if (fromNode) {
          data2[i].from = fromNode.data.label;
        }
      }

      const row = {
        id: node.id,
        name: node.data.label,
        traitor: node.isTraitor ? 'yes' : 'no',
        decision1: decision1 ? decision1.decision : 'undecided',
        data1: data1,
        decision2: decision2 ? decision2.decision : 'undecided',
        data2: data2,
      };

      rows.push(row);
    });

    setRows(rows);
  }

  useEffect(() => {
    const dataInterval = setInterval(getData, 5000);

    return () => clearInterval(dataInterval);
  }, []);

  return (
    <Box sx={{ width: '100%' }}>
      <DataGrid
        sx={{
          '&.MuiDataGrid-root--densityCompact .MuiDataGrid-cell': { py: '8px' },
          '&.MuiDataGrid-root--densityStandard .MuiDataGrid-cell': { py: '15px' },
          '&.MuiDataGrid-root--densityComfortable .MuiDataGrid-cell': { py: '22px' },
        }}
        rows={rows}
        columns={columns}
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
  )
}

export default Table;