import { useState, useEffect, useCallback } from 'react';
import { ReactFlow, Controls, Background, ConnectionMode, MarkerType, applyNodeChanges, applyEdgeChanges } from '@xyflow/react';

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
  let [edges, setEdges] = useState(defaultEdges);

  function getDecisions() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Decisions';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {

        if(data.decisions === undefined) {
          nodes.forEach((node) => {
            // remove the decision from the node
            node.data.decision = undefined;
            if (node.style !== undefined) {
              node.style.backgroundColor = undefined;
            }
          });

          setNodes(nodes);
          return;
        }

        data.decisions.forEach((decision) => {
          // add the decision to the node
          const index = nodes.findIndex((item) => item.id === decision.from);
          if (index > -1 && nodes[index].data.decision === undefined) {
            console.log("Adding decision to node", decision.from, decision.decision);

            let bgColor = '#fd4848'; // retreat color (red)
            if(decision.decision === 'attack') {
              bgColor = '#36ff40'; //attack color (green)
            }

            nodes[index] = {
              ...nodes[index],
              style: {...nodes[index].style, backgroundColor: bgColor},
              data: {...nodes[index].data, decision: decision.decision}
            };
          }
        });

        // update the nodes
        setNodes(nodes);
    });
  }

  function getNodes() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Nodes';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {
        let newNodes = [];

        data.nodes.forEach((node, index) => {
          let copyNode = {...node};
          
          // does the item exist
          const found = nodes.find((item) => item.id === node.id);
          if (found) {
            copyNode = {...found};
          }

          let borderColor = '1px solid #EDEDED'; // default color
          // if a traitor mark the node border
          if (node.isTraitor) {
            borderColor = '2px dashed #404040'; // gray color
          }

          copyNode = {
            ...copyNode,
            style: {...copyNode.style, border: borderColor},
          }
          
          newNodes.push(copyNode);
        });

        // set the node positions
        let generalNodes = 0;
        for(let i = 0; i < newNodes.length; i++) { 
          if (newNodes[i].id === '0') {
            newNodes[i].position = { x: 0, y: 100 };
          } else {
            newNodes[i].position = { x: 0 + 200 * generalNodes, y: 500 };
            generalNodes++;
          }
        }

        setNodes(newNodes);
        nodes = newNodes;
      });
  }
  
  function getEdges() {
    const url = 'http://localhost:8080/proto.server.v1.CommanderService/Edges';
    fetch(url, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({})})
      .then((response) => response.json())
      .then((data) => {

        // if we have no edges, then we need to set the default edges
        if(data.edges === undefined) {
          setEdges(defaultEdges);
          edges = defaultEdges;
          return;
        }
        
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
  
  useEffect(() => {
    const decisionsInterval = setInterval(getDecisions, 5000);

    return () => clearInterval(decisionsInterval);
  }, []);

  const onNodesChange = useCallback(
    (changes) => {
      setNodes((oldNodes) => applyNodeChanges(changes, oldNodes));
    },
    [setNodes],
  );
  
  const onEdgesChange = useCallback(
    (changes) => {
      setEdges((oldEdges) => applyEdgeChanges(changes, oldEdges));
    },
    [setEdges],
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