import { useState, useEffect, useCallback } from 'react';
import { ReactFlow, Controls, Background, Position, ConnectionMode, MarkerType, applyEdgeChanges, applyNodeChanges } from '@xyflow/react';

import BiDirectionalEdge from './BiDirectionalEdge';
import BiDirectionalNode from './BiDirectionalNode';

const edgeTypes = {
  bidirectional: BiDirectionalEdge,
};
 
const nodeTypes = {
  bidirectional: BiDirectionalNode,
};

const defaultNodes = [];
const defaultEdges = [];

function Grid() {
  let [nodes, setNodes] = useState(defaultNodes);
  let [edges, setEdges, onEdgesChange] = useState(defaultEdges);

  function getNodes() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Nodes';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        let nodesToAdd = [];

        data.nodes.forEach((node, index) => {
          // does the item exist
          const found = nodes.find((item) => item.id === node.id);
          if (found) {
            return;
          }

          // set the positions of the nodes
          if (node.id === '0') {
            node.position = { x: 100, y: 100 };
          } else {
            node.position = { x: 0 + 200 * index, y: 500 };
          }

          nodesToAdd.push(node);
        });

        if (nodesToAdd.length === 0) {
          return;
        }

        const newNodes = [...nodes, ...nodesToAdd];
        setNodes(newNodes);
        nodes = newNodes;
      });
  }
  
  function getEdges() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Edges';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        let edgesToAdd = [];

        // loop over the nodes and find the nodes for each edge
        // we need this to calculate the source and target
        data.edges.forEach((edge) => {
          const source = nodes.find((node) => node.id === edge.source);
          const target = nodes.find((node) => node.id === edge.target);
         
          // do we have these nodes
          if (!source || !target) {
            return;
          }
          
          // does the item exist
          const found = edges.find((item) => item.id === edge.id);
          if (found) {
            return;
          }
          
          console.log(`Adding edge from ${source.data.label} to ${target.data.label}`);

          // work out which node on the target to connect to
          if(source.position.y < target.position.y) {
            edge.targetHandle = 'top';
          }else if (source.position.x > target.position.x) {
            edge.targetHandle = 'right';
            edge.sourceHandle = 'left';
          } else {
            edge.targetHandle = 'left';
            edge.sourceHandle = 'right';
          }

          edge.markerEnd = {type: MarkerType.Arrow};
          edge.type = 'bidirectional';
          edgesToAdd.push(edge);
        });

        if (edgesToAdd.length > 0) {
          edges = [...edges, ...edgesToAdd];
          setEdges(edgesToAdd);
        }

      });
  }

  useEffect(() => {
    const nodesInterval = setInterval(getNodes, 5000);

    return () => clearInterval(nodesInterval);
  }, []);
  
  useEffect(() => {
    const edgesInterval = setInterval(getEdges, 5000);

    return () => clearInterval(edgesInterval);
  }, []);


  const onNodesChange = useCallback(
    (changes) => {
      setNodes((nds) => applyNodeChanges(changes, nds))
    },
    [],
  );

  return (
    <ReactFlow 
      nodes={nodes} 
      edges={edges} 
      onNodesChange={onNodesChange} 
      onEdgesChange={onEdgesChange} 
      connectionMode={ConnectionMode.Loose}
      nodeTypes={nodeTypes}
      edgeTypes={edgeTypes}
      fitView
    >
      <Background />
      <Controls />
    </ReactFlow>
  );
}
 
export default Grid;